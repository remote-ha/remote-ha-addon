package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type testRoundTripper func(*http.Request) (*http.Response, error)

func (f testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestIdentityRenewalRecoversCommittedButLostResponse(t *testing.T) {
	caKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ca := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "test device CA"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(365 * 24 * time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign}
	caDER, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatal(err)
	}
	ca, _ = x509.ParseCertificate(caDER)
	id := "12345678901234567890123456789012"
	serial := int64(1)
	identityFor := func(public *ecdsa.PublicKey, expiry time.Time) (Identity, error) {
		serial++
		uri, _ := url.Parse("urn:rha:device:" + id)
		leaf := &x509.Certificate{SerialNumber: big.NewInt(serial), Subject: pkix.Name{CommonName: id}, NotBefore: time.Now().Add(-time.Minute), NotAfter: expiry, BasicConstraintsValid: true, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, URIs: []*url.URL{uri}}
		der, err := x509.CreateCertificate(rand.Reader, leaf, ca, public, caKey)
		return Identity{DeviceID: id, Domain: "installation.test", Expires: float64(expiry.Unix()), Certificate: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})), CA: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER}))}, err
	}
	oldKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	oldIdentity, err := identityFor(&oldKey.PublicKey, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	active := &oldKey.PublicKey
	lost := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.URL.Path != "/v1/device/renew" {
			http.NotFound(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		hash := sha256.Sum256(body)
		canonical := strings.Join([]string{r.Method, r.URL.Path, r.Header.Get("X-Timestamp"), r.Header.Get("X-Nonce"), hex.EncodeToString(hash[:])}, "\n")
		digest := sha256.Sum256([]byte(canonical))
		signature, _ := base64.StdEncoding.DecodeString(r.Header.Get("X-Signature"))
		if !ecdsa.VerifyASN1(active, digest[:], signature) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var payload map[string]string
		if json.Unmarshal(body, &payload) != nil {
			w.WriteHeader(400)
			return
		}
		block, _ := pem.Decode([]byte(payload["csr"]))
		if block == nil {
			w.WriteHeader(400)
			return
		}
		csr, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil || csr.CheckSignature() != nil {
			w.WriteHeader(400)
			return
		}
		active = csr.PublicKey.(*ecdsa.PublicKey)
		// Model an API commit followed by a failed proxy response, before the client
		// has persisted its new identity. No private key crosses this test boundary.
		if !lost {
			lost = true
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		identity, err := identityFor(active, time.Now().Add(7*24*time.Hour))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(identity)
	})
	client := &http.Client{Transport: testRoundTripper(func(r *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, r)
		return recorder.Result(), nil
	})}
	directory := t.TempDir()
	a := &Agent{state: directory, options: Options{APIURL: "https://api.test"}, client: client}
	if err = a.saveIdentity(oldIdentity, oldKey); err != nil {
		t.Fatal(err)
	}
	if err = a.renewIdentity(); err == nil {
		t.Fatal("lost response must be reported")
	}
	if _, err = os.Stat(filepath.Join(directory, "renewal.key")); err != nil {
		t.Fatal("pending key was not persisted")
	}
	restarted := &Agent{state: directory, options: a.options, client: client}
	if err = restarted.pair(); err != nil {
		t.Fatal(err)
	}
	if err = restarted.renewIdentity(); err != nil {
		t.Fatal("renewal did not recover after restart:", err)
	}
	mu.Lock()
	matches := restarted.key.PublicKey.Equal(active)
	mu.Unlock()
	if !matches || restarted.key.PublicKey.Equal(&oldKey.PublicKey) {
		t.Fatal("wrong identity key activated")
	}
	if _, err = os.Stat(filepath.Join(directory, "renewal.key")); !os.IsNotExist(err) {
		t.Fatal("completed pending key not removed")
	}
	// A paired device learns policy changes at token refresh, without rotating
	// its private key or waiting five days for the next identity renewal.
	production := "https://acme-v02.api.letsencrypt.org/directory"
	restarted.configChanged = make(chan struct{}, 1)
	configurationDomain := oldIdentity.Domain
	restarted.client = &http.Client{Transport: testRoundTripper(func(r *http.Request) (*http.Response, error) {
		payload, _ := json.Marshal(map[string]any{"access_token": "test-only", "configuration": map[string]string{
			"domain": configurationDomain, "acme_directory": production}})
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(string(payload)))}, nil
	})}
	keyBefore, certificateBefore := restarted.key, restarted.identity.Certificate
	if err = restarted.refreshToken(); err != nil {
		t.Fatal(err)
	}
	if restarted.identity.ACMEDirectory != production || restarted.key != keyBefore || restarted.identity.Certificate != certificateBefore {
		t.Fatal("configuration refresh changed device identity or failed to update policy")
	}
	select {
	case <-restarted.configChanged:
	default:
		t.Fatal("certificate worker was not notified")
	}
	reloaded := &Agent{state: directory}
	if err = reloaded.pair(); err != nil || reloaded.identity.ACMEDirectory != production {
		t.Fatal("configuration was not persisted across restart", err)
	}
	configurationDomain = "foreign.test"
	if err = restarted.refreshToken(); err == nil {
		t.Fatal("foreign installation configuration accepted")
	}
	configurationDomain = oldIdentity.Domain
	production = "http://insecure.test/directory"
	if err = restarted.refreshToken(); err == nil {
		t.Fatal("HTTP ACME downgrade accepted")
	}
}
