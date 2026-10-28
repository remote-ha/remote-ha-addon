package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Options struct {
	APIURL            string `json:"api_url"`
	PairingToken      string `json:"pairing_token"`
	CoreHost          string `json:"core_host"`
	CorePort          int    `json:"core_port"`
	CoreScheme        string `json:"core_scheme"`
	CoreServerName    string `json:"core_server_name"`
	CertDir           string `json:"certificate_dir"`
	Email             string `json:"acme_email"`
	Terms             bool   `json:"accept_acme_terms"`
	ManageCertificate bool   `json:"manage_certificate"`
	ACMECaFile        string `json:"acme_ca_file"`
	CoreCaFile        string `json:"core_ca_file"`
	APICaFile         string `json:"api_ca_file"`
}

type Identity struct {
	DeviceID      string  `json:"device_id"`
	Domain        string  `json:"domain"`
	Certificate   string  `json:"certificate"`
	CA            string  `json:"ca"`
	Expires       float64 `json:"expires"`
	RelayHost     string  `json:"relay_host"`
	RelayPort     int     `json:"relay_port"`
	RelayName     string  `json:"relay_server_name"`
	ACMEDirectory string  `json:"acme_directory"`
	ProxyType     string  `json:"proxy_type"`
}

type Agent struct {
	mu              sync.RWMutex
	options         Options
	state           string
	identity        Identity
	key             *ecdsa.PrivateKey
	client          *http.Client
	running         bool
	restartRequired bool
	certExpires     int64
	certError       bool
	configChanged   chan struct{}
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".pending-")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if err = f.Chmod(mode); err == nil {
		_, err = f.Write(data)
	}
	if err == nil {
		err = f.Sync()
	}
	if e := f.Close(); err == nil {
		err = e
	}
	if err != nil {
		return err
	}
	if err = os.Rename(name, path); err != nil {
		return err
	}
	directory, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}

func switchLink(directory, name, target string) error {
	temporary := filepath.Join(directory, "."+name+"-pending")
	_ = os.Remove(temporary)
	if err := os.Symlink(target, temporary); err != nil {
		return err
	}
	if err := os.Rename(temporary, filepath.Join(directory, name)); err != nil {
		return err
	}
	f, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Sync()
}

func keyPEM(key *ecdsa.PrivateKey) []byte {
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		panic(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func loadKey(path string) (*ecdsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("invalid key")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := parsed.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("P256 key required")
	}
	return key, nil
}

func newCSR(key *ecdsa.PrivateKey) (string, error) {
	b, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{Subject: pkix.Name{CommonName: "device"}}, key)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: b})), nil
}

func roots(file string) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if !pool.AppendCertsFromPEM(b) {
			return nil, errors.New("invalid CA file")
		}
	}
	return pool, nil
}

func (a *Agent) request(path string, body any, signed bool) ([]byte, error) {
	return a.requestWithKey(path, body, signed, nil)
}

type managementStatusError struct{ code int }

func (e *managementStatusError) Error() string {
	return fmt.Sprintf("management request rejected (%d)", e.code)
}

func (a *Agent) requestWithKey(path string, body any, signed bool, override *ecdsa.PrivateKey) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", strings.TrimRight(a.options.APIURL, "/")+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if signed {
		a.mu.RLock()
		id, key := a.identity.DeviceID, a.key
		a.mu.RUnlock()
		if override != nil {
			key = override
		}
		nonceBytes := make([]byte, 16)
		if _, err := rand.Read(nonceBytes); err != nil {
			return nil, err
		}
		nonce := hex.EncodeToString(nonceBytes)
		stamp := fmt.Sprint(time.Now().Unix())
		digest := sha256.Sum256(data)
		canonical := "POST\n" + path + "\n" + stamp + "\n" + nonce + "\n" + hex.EncodeToString(digest[:])
		hash := sha256.Sum256([]byte(canonical))
		signature, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Device", id)
		req.Header.Set("X-Timestamp", stamp)
		req.Header.Set("X-Nonce", nonce)
		req.Header.Set("X-Signature", base64.StdEncoding.EncodeToString(signature))
	}
	response, err := a.client.Do(req)
	if err != nil {
		return nil, errors.New("management connection failed")
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, &managementStatusError{code: response.StatusCode}
	}
	return io.ReadAll(io.LimitReader(response.Body, 32768))
}

func (a *Agent) saveIdentity(identity Identity, key *ecdsa.PrivateKey) error {
	// Validate CA chain, client EKU, public key and device URI before activating.
	pair, err := tls.X509KeyPair([]byte(identity.Certificate), keyPEM(key))
	if err != nil {
		return err
	}
	cert, err := x509.ParseCertificate(pair.Certificate[0])
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(identity.CA)) {
		return errors.New("invalid tunnel CA")
	}
	if _, err = cert.Verify(x509.VerifyOptions{Roots: pool, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}}); err != nil {
		return err
	}
	if len(cert.URIs) != 1 || cert.URIs[0].String() != "urn:rha:device:"+identity.DeviceID {
		return errors.New("wrong device certificate")
	}
	parsed, err := url.Parse("https://" + identity.Domain)
	if err != nil || parsed.Hostname() != identity.Domain || strings.Contains(identity.Domain, "/") {
		return errors.New("invalid domain")
	}
	directory := filepath.Join(a.state, "identity-"+fmt.Sprint(cert.SerialNumber))
	if err = os.MkdirAll(directory, 0700); err != nil {
		return err
	}
	data, _ := json.Marshal(identity)
	for name, b := range map[string][]byte{"identity.json": data, "device.key": keyPEM(key), "device.pem": []byte(identity.Certificate), "ca.pem": []byte(identity.CA)} {
		if err = atomicWrite(filepath.Join(directory, name), b, 0600); err != nil {
			return err
		}
	}
	if err = switchLink(a.state, "identity", filepath.Base(directory)); err != nil {
		return err
	}
	a.mu.Lock()
	a.identity, a.key = identity, key
	a.mu.Unlock()
	return nil
}

func (a *Agent) pair() error {
	directory := filepath.Join(a.state, "identity")
	cached := false
	if data, err := os.ReadFile(filepath.Join(directory, "identity.json")); err == nil {
		if err = json.Unmarshal(data, &a.identity); err != nil {
			return err
		}
		a.key, err = loadKey(filepath.Join(directory, "device.key"))
		if err != nil {
			return err
		}
		cached = true
		if pending, e := loadKey(filepath.Join(a.state, "renewal.key")); e == nil && pending.PublicKey.Equal(&a.key.PublicKey) {
			_ = os.Remove(filepath.Join(a.state, "renewal.key"))
		}
	}
	fingerprint := sha256.Sum256([]byte(a.options.PairingToken))
	fingerprintText := hex.EncodeToString(fingerprint[:])
	marker := filepath.Join(a.state, "pairing-token.sha256")
	previous, _ := os.ReadFile(marker)
	if cached && (a.options.PairingToken == "" || string(previous) == fingerprintText) {
		return nil
	}
	if a.options.PairingToken == "" {
		return errors.New("pairing token required")
	}
	key, err := loadKey(filepath.Join(a.state, "pairing.key"))
	if os.IsNotExist(err) {
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err == nil {
			err = atomicWrite(filepath.Join(a.state, "pairing.key"), keyPEM(key), 0600)
		}
	}
	if err != nil {
		return err
	}
	csr, err := newCSR(key)
	if err != nil {
		return err
	}
	data, err := a.request("/v1/pair", map[string]string{"token": a.options.PairingToken, "csr": csr}, false)
	if err != nil {
		var status *managementStatusError
		if cached && errors.As(err, &status) && (status.code == 401 || status.code == 409) {
			log.Print("pairing token rejected or already used; retaining the saved identity")
			return nil
		}
		return err
	}
	var identity Identity
	if err = json.Unmarshal(data, &identity); err != nil {
		return err
	}
	if err = a.saveIdentity(identity, key); err != nil {
		return err
	}
	if err = atomicWrite(marker, []byte(fingerprintText), 0600); err != nil {
		log.Print("pairing succeeded; token fingerprint could not be persisted")
	}
	a.options.PairingToken = ""
	_ = os.Remove(filepath.Join(a.state, "pairing.key"))
	log.Print("pairing completed; new device identity saved")
	return nil
}
func (a *Agent) refreshToken() error {
	data, err := a.request("/v1/device/token", map[string]string{}, true)
	if err != nil {
		return err
	}
	var token struct {
		Token         string `json:"access_token"`
		Configuration *struct {
			Domain        string `json:"domain"`
			ACMEDirectory string `json:"acme_directory"`
			ProxyType     string `json:"proxy_type"`
		} `json:"configuration"`
	}
	if err = json.Unmarshal(data, &token); err != nil || token.Token == "" {
		return errors.New("missing access token")
	}
	if token.Configuration != nil {
		a.mu.RLock()
		identity, key := a.identity, a.key
		a.mu.RUnlock()
		config := token.Configuration
		u, parseErr := url.Parse(config.ACMEDirectory)
		if config.Domain != identity.Domain || parseErr != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Fragment != "" {
			return errors.New("invalid device configuration")
		}
		if config.ProxyType != "" && config.ProxyType != "http" {
			return errors.New("unsupported tunnel protocol; update the platform")
		}
		if config.ACMEDirectory != identity.ACMEDirectory || config.ProxyType != identity.ProxyType {
			identity.ACMEDirectory = config.ACMEDirectory
			identity.ProxyType = config.ProxyType
			if err = a.saveIdentity(identity, key); err != nil {
				return err
			}
			if a.configChanged != nil {
				select {
				case a.configChanged <- struct{}{}:
				default:
				}
			}
		}
	}
	return atomicWrite(filepath.Join(a.state, "access.jwt"), []byte(token.Token), 0600)
}

func (a *Agent) renewIdentity() error {
	a.mu.RLock()
	expiry := a.identity.Expires
	a.mu.RUnlock()
	if time.Unix(int64(expiry), 0).After(time.Now().Add(365 * 24 * time.Hour)) {
		return nil
	}
	// Persist the next key before sending a CSR. If the server commits but its
	// response is lost, that key can authenticate a retry after a process restart.
	pending := filepath.Join(a.state, "renewal.key")
	key, err := loadKey(pending)
	if os.IsNotExist(err) {
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err == nil {
			err = atomicWrite(pending, keyPEM(key), 0600)
		}
	}
	if err != nil {
		return err
	}
	csr, err := newCSR(key)
	if err != nil {
		return err
	}
	data, err := a.request("/v1/device/renew", map[string]string{"csr": csr}, true)
	var status *managementStatusError
	if errors.As(err, &status) && status.code == http.StatusUnauthorized {
		data, err = a.requestWithKey("/v1/device/renew", map[string]string{"csr": csr}, true, key)
	}
	if err != nil {
		return err
	}
	var identity Identity
	if err = json.Unmarshal(data, &identity); err != nil {
		return err
	}
	if err = a.saveIdentity(identity, key); err != nil {
		return err
	}
	return os.Remove(pending)
}

func (a *Agent) writeFRP() error {
	a.mu.RLock()
	identity := a.identity
	a.mu.RUnlock()
	dir := filepath.Join(a.state, "identity")
	config := map[string]any{
		"serverAddr": identity.RelayHost, "serverPort": identity.RelayPort, "loginFailExit": true,
		"auth": map[string]any{"method": "oidc", "additionalScopes": []string{"HeartBeats", "NewWorkConns"},
			"oidc": map[string]any{"tokenSource": map[string]any{"type": "file", "file": map[string]string{"path": filepath.Join(a.state, "access.jwt")}}}},
		"transport": map[string]any{"protocol": "tcp", "poolCount": 1, "heartbeatInterval": 10, "heartbeatTimeout": 30,
			"tls": map[string]any{"enable": true, "serverName": identity.RelayName, "certFile": filepath.Join(dir, "device.pem"),
				"keyFile": filepath.Join(dir, "device.key"), "trustedCaFile": filepath.Join(dir, "ca.pem")}},
		"proxies": []map[string]any{{"name": strings.Split(identity.Domain, ".")[0], "type": "http", "customDomains": []string{identity.Domain},
			"localIP": "127.0.0.1", "localPort": 8089}},
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return atomicWrite(filepath.Join(a.state, "frpc.json"), data, 0600)
}

func (a *Agent) probeCore() error {
	transport, err := a.coreTransport()
	if err != nil {
		return err
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: 3 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Get(a.coreURL().String())
	if err != nil {
		return errors.New("local Home Assistant unavailable; verify detected host/port and existing TLS trust")
	}
	response.Body.Close()
	if response.StatusCode >= 500 {
		return fmt.Errorf("local Home Assistant returned HTTP %d", response.StatusCode)
	}
	a.mu.Lock()
	a.restartRequired = false
	a.certExpires = int64(a.identity.Expires)
	a.mu.Unlock()
	return nil
}

func jitter(base time.Duration) time.Duration {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(base/2)+1))
	if err != nil {
		return base
	}
	return base + time.Duration(n.Int64())
}

func (a *Agent) health(w http.ResponseWriter, r *http.Request) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"paired": a.identity.DeviceID != "", "tunnel_process_running": a.running,
		"domain": a.identity.Domain, "cert_expires": a.certExpires, "restart_required": a.restartRequired, "certificate_renewal_error": a.certError})
}

func (a *Agent) run(ctx context.Context) {
	backoff := time.Second
	for ctx.Err() == nil {
		if err := a.pair(); err == nil {
			break
		}
		log.Print("pairing unavailable; check configuration and token expiry")
		select {
		case <-ctx.Done():
			return
		case <-time.After(jitter(backoff)):
		}
		if backoff < time.Minute {
			backoff *= 2
		}
	}
	if ctx.Err() != nil {
		return
	}
	log.Print("device identity loaded; requesting authorization")
	a.configChanged = make(chan struct{}, 1)
	log.Print("encrypted connection enabled; no additional Home Assistant setup is required")
	var child *exec.Cmd
	var done chan error
	stop := func() {
		if child != nil {
			_ = child.Process.Kill()
			<-done
			child = nil
			a.mu.Lock()
			a.running = false
			a.mu.Unlock()
		}
	}
	defer stop()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	tokenAt, renewAt, nextStart := time.Time{}, time.Time{}, time.Time{}
	backoff = time.Second
	lastIdentity := ""
	for ctx.Err() == nil {
		if child != nil {
			select {
			case <-done:
				child = nil
				a.mu.Lock()
				a.running = false
				a.mu.Unlock()
				nextStart = time.Now().Add(jitter(backoff))
				if backoff < time.Minute {
					backoff *= 2
				}
			default:
			}
		}
		if time.Since(renewAt) > time.Hour {
			if err := a.renewIdentity(); err != nil {
				log.Print("device identity renewal failed")
			} else {
				renewAt = time.Now()
			}
		}
		a.mu.RLock()
		identityVersion := a.identity.Certificate + a.identity.ProxyType
		a.mu.RUnlock()
		if identityVersion != lastIdentity {
			stop()
			lastIdentity = identityVersion
			tokenAt = time.Time{}
		}
		tokenOK := true
		if time.Since(tokenAt) > 20*time.Second {
			if err := a.refreshToken(); err != nil {
				tokenOK = false
				stop()
				log.Print("device authorization unavailable; tunnel stopped")
			} else {
				tokenAt = time.Now()
			}
		}
		if tokenOK && child == nil && time.Now().After(nextStart) {
			if err := a.probeCore(); err == nil {
				if err = a.writeFRP(); err == nil {
					child = exec.CommandContext(ctx, "frpc", "-c", filepath.Join(a.state, "frpc.json"))
					// FRP diagnostics may contain bearer tokens on error. Keep only agent status.
					child.Stdout = io.Discard
					child.Stderr = io.Discard
					if err = child.Start(); err == nil {
						done = make(chan error, 1)
						go func(p *exec.Cmd) { done <- p.Wait() }(child)
						a.mu.Lock()
						a.running = true
						a.mu.Unlock()
						backoff = time.Second
					} else {
						child = nil
						log.Print("frpc failed to start")
					}
				}
			} else {
				log.Printf("waiting for local Home Assistant: %v", err)
			}
			nextStart = time.Now().Add(jitter(backoff))
		}
		_ = a.probeCore()
		a.mu.RLock()
		expiry, restart, certError := a.certExpires, a.restartRequired, a.certError
		a.mu.RUnlock()
		_, _ = a.request("/v1/device/status", map[string]any{"cert_expires": expiry, "restart_required": restart, "certificate_error": certError}, true)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func main() {
	optionsFile := os.Getenv("OPTIONS_FILE")
	if optionsFile == "" {
		optionsFile = "/data/options.json"
	}
	state := os.Getenv("STATE_DIR")
	if state == "" {
		state = "/data/rha"
	}
	b, err := os.ReadFile(optionsFile)
	if err != nil {
		log.Fatal("cannot read options")
	}
	supervisorToken := os.Getenv("SUPERVISOR_TOKEN")
	options, err := loadOptions(b, supervisorToken != "")
	if err != nil {
		log.Fatal("invalid options JSON")
	}
	localClient := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("redirect denied") }}
	if err := migrateSupervisorOptions(localClient, supervisorToken); err != nil {
		log.Print("legacy app options cleanup deferred; retry on next start")
	}
	discoverCore(&options)
	u, err := url.Parse(options.APIURL)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.Path != "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		log.Fatal("api_url must be an HTTPS origin")
	}
	if options.CoreHost == "" || options.CorePort < 1 || options.CorePort > 65535 {
		log.Fatal("invalid Core endpoint")
	}
	if err = os.MkdirAll(state, 0700); err != nil {
		log.Fatal("cannot create state directory")
	}
	pool, err := roots(options.APICaFile)
	if err != nil {
		log.Fatal("invalid API trust file")
	}
	a := &Agent{options: options, state: state, client: &http.Client{Timeout: 20 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("redirect denied") },
		Transport:     &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: pool}}}}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", a.health)
	server := &http.Server{Addr: ":8088", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = server.ListenAndServe() }()
	transport, err := a.coreTransport()
	if err != nil {
		log.Fatal("invalid local Core HTTPS trust")
	}
	defer transport.CloseIdleConnections()
	gateway := &http.Server{Handler: a.gateway(transport), ReadHeaderTimeout: 10 * time.Second}
	listener, err := net.Listen("tcp", "127.0.0.1:8089")
	if err != nil {
		log.Fatal("cannot start local tunnel gateway")
	}
	go func() { _ = gateway.Serve(listener) }()
	a.run(ctx)
	_ = gateway.Close()
	_ = server.Close()
}
