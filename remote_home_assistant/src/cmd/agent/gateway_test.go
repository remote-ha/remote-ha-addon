package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGatewayForwardsHAAuthAndStripsUntrustedProxyHeaders(t *testing.T) {
	var seen bool
	transport := testRoundTripper(func(r *http.Request) (*http.Response, error) {
		seen = true
		if r.URL.Path == "/api/" && (r.URL.RawQuery != "q=test" || r.Header.Get("Authorization") != "Bearer home-assistant-user-token") {
			t.Error("HA request corrupted")
		}
		for key := range r.Header {
			lower := strings.ToLower(key)
			if strings.HasPrefix(lower, "x-forwarded-") || strings.HasPrefix(lower, "cf-") || lower == "forwarded" {
				t.Error("forwarded identity trusted:", key)
			}
		}
		if r.URL.Scheme != "http" || r.URL.Host != "homeassistant:8123" {
			t.Error("wrong Core endpoint")
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("local Core HTTP works"))}, nil
	})
	a := &Agent{options: Options{CoreHost: "homeassistant", CoreScheme: "http", CorePort: 8123}, identity: Identity{Domain: "installation.test"}}
	request := httptest.NewRequest("GET", "http://installation.test/api/?q=test", nil)
	request.Header.Set("Authorization", "Bearer home-assistant-user-token")
	for _, header := range []string{"X-Forwarded-For", "X-Forwarded-Proto", "X-Forwarded-Host", "CF-Connecting-IP", "Forwarded"} {
		request.Header.Set(header, "attacker")
	}
	response := httptest.NewRecorder()
	a.gateway(transport).ServeHTTP(response, request)
	if !seen || response.Code != 200 {
		t.Fatal("plain HTTP Core unavailable", response.Code)
	}
	seen = false
	request.Host = "foreign.test"
	a.gateway(transport).ServeHTTP(httptest.NewRecorder(), request)
	if seen {
		t.Fatal("foreign installation reached Core")
	}
}
