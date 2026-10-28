package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSupervisedOptionsOnlyAcceptPairingToken(t *testing.T) {
	legacy := []byte(`{"pairing_token":"pair-me","api_url":"https://foreign.test","core_host":"foreign-core","core_port":443,"core_scheme":"https","manage_certificate":true,"core_ca_file":"/ssl/old.pem"}`)
	options, err := loadOptions(legacy, true)
	if err != nil || options.PairingToken != "pair-me" || options.APIURL != "https://remote-ha.com" || options.CoreHost != "homeassistant" || options.CorePort != 8123 || options.CoreScheme != "http" || options.ManageCertificate || options.CoreCaFile != "" {
		t.Fatal("legacy settings altered the supervised app defaults")
	}
	standalone, err := loadOptions(legacy, false)
	if err != nil || standalone.CoreHost != "foreign-core" || standalone.APIURL != "https://foreign.test" {
		t.Fatal("standalone operator overrides lost")
	}
	if _, err := loadOptions([]byte(`{"pairing_token":false}`), true); err == nil {
		t.Fatal("invalid pairing token accepted")
	}
}

func TestSupervisorMigrationOnlyReplacesOwnOptionsPreservingSecretReference(t *testing.T) {
	posts := 0
	client := &http.Client{Transport: testRoundTripper(func(request *http.Request) (*http.Response, error) {
		if request.URL.Host != "supervisor" || request.Header.Get("Authorization") != "Bearer local-supervisor-test-token" {
			t.Fatal("unexpected Supervisor target or authentication")
		}
		body := `{"result":"ok"}`
		switch request.URL.Path {
		case "/addons/self/info":
			if request.Method != http.MethodGet {
				t.Fatal("expected read-only info request")
			}
			body = `{"result":"ok","data":{"options":{"pairing_token":"!secret rha_pairing","api_url":"https://old.test","manage_certificate":true}}}`
		case "/addons/self/options":
			if request.Method != http.MethodPost {
				t.Fatal("expected options replacement")
			}
			var payload map[string]map[string]string
			if json.NewDecoder(request.Body).Decode(&payload) != nil || len(payload) != 1 || len(payload["options"]) != 1 || payload["options"]["pairing_token"] != "!secret rha_pairing" {
				t.Fatal("migration changed secret reference or retained legacy settings")
			}
			posts++
		default:
			t.Fatal("migration reached an endpoint outside its own settings")
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
	})}
	if err := migrateSupervisorOptions(client, "local-supervisor-test-token"); err != nil || posts != 1 {
		t.Fatal("migration did not complete")
	}
}

func TestSupervisorMigrationLeavesTokenOnlyOptionsUntouched(t *testing.T) {
	requests := 0
	client := &http.Client{Transport: testRoundTripper(func(request *http.Request) (*http.Response, error) {
		requests++
		if request.Method != http.MethodGet {
			t.Fatal("token-only options rewritten")
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"result":"ok","data":{"options":{"pairing_token":""}}}`))}, nil
	})}
	if err := migrateSupervisorOptions(client, "test"); err != nil || requests != 1 {
		t.Fatal("token-only migration failed")
	}
	if err := migrateSupervisorOptions(client, ""); err != nil || requests != 1 {
		t.Fatal("standalone contacted Supervisor")
	}
}

func TestSupervisorMigrationDoesNotLeakFailureResponse(t *testing.T) {
	client := &http.Client{Transport: testRoundTripper(func(request *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 403, Body: io.NopCloser(strings.NewReader("private response content"))}, nil
	})}
	err := migrateSupervisorOptions(client, "secret-token")
	if err == nil || strings.Contains(err.Error(), "secret-token") || strings.Contains(err.Error(), "private response") {
		t.Fatal("Supervisor failure was hidden or leaked credentials")
	}
}
