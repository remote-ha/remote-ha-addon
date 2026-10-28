package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// The default Supervisor role can only read /info endpoints. Its token is
// never sent to Core, FRP, the management platform, or diagnostics.
func discoverCore(options *Options) {
	token := os.Getenv("SUPERVISOR_TOKEN")
	if token == "" {
		return
	}
	request, _ := http.NewRequest("GET", "http://supervisor/core/info", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("redirect denied") }}
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return
	}
	var info struct {
		Result string `json:"result"`
		Data   struct {
			IP   string `json:"ip_address"`
			Port int    `json:"port"`
			SSL  bool   `json:"ssl"`
		} `json:"data"`
	}
	if json.NewDecoder(io.LimitReader(response.Body, 16384)).Decode(&info) != nil || info.Result != "ok" {
		return
	}
	if options.CoreHost == "homeassistant" {
		if net.ParseIP(info.Data.IP) != nil {
			options.CoreHost = info.Data.IP
		}
		if info.Data.Port > 0 && info.Data.Port < 65536 {
			options.CorePort = info.Data.Port
		}
		if info.Data.SSL {
			options.CoreScheme = "https"
		}
	}
}

func (a *Agent) coreURL() *url.URL {
	scheme := a.options.CoreScheme
	if scheme == "" {
		scheme = "http"
	}
	return &url.URL{Scheme: scheme, Host: net.JoinHostPort(a.options.CoreHost, strconv.Itoa(a.options.CorePort)), Path: "/"}
}

func (a *Agent) coreTransport() (*http.Transport, error) {
	if a.options.CoreScheme != "" && a.options.CoreScheme != "http" && a.options.CoreScheme != "https" {
		return nil, errors.New("invalid Core scheme")
	}
	pool, err := roots(a.options.CoreCaFile)
	if err != nil {
		return nil, err
	}
	return &http.Transport{DialContext: (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: pool, ServerName: a.options.CoreServerName},
		ResponseHeaderTimeout: 60 * time.Second, IdleConnTimeout: 90 * time.Second}, nil
}

func (a *Agent) gateway(transport http.RoundTripper) http.Handler {
	proxy := &httputil.ReverseProxy{Transport: transport, FlushInterval: -1,
		Rewrite: func(request *httputil.ProxyRequest) {
			request.SetURL(a.coreURL())
			request.Out.Host = request.In.Host
			// No forwarded identity is trusted from the cloud or the Internet. Core
			// sees the local app address, so it needs no trusted_proxies configuration.
			for key := range request.Out.Header {
				lower := strings.ToLower(key)
				if strings.HasPrefix(lower, "x-forwarded-") || strings.HasPrefix(lower, "cf-") || lower == "forwarded" || lower == "true-client-ip" {
					request.Out.Header.Del(key)
				}
			}
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, _ error) {
			http.Error(w, "Home Assistant temporarily unavailable", http.StatusBadGateway)
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.mu.RLock()
		domain := a.identity.Domain
		a.mu.RUnlock()
		host := strings.ToLower(r.Host)
		if name, _, err := net.SplitHostPort(host); err == nil {
			host = name
		}
		if domain == "" || host != domain {
			http.Error(w, "Installation not found", http.StatusNotFound)
			return
		}
		proxy.ServeHTTP(w, r)
	})
}
