package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Supervisor apps accept only the pairing token. Standalone operator fixtures
// retain endpoint overrides without exposing them in the HA configuration UI.
func loadOptions(data []byte, supervised bool) (Options, error) {
	options := Options{APIURL: "https://remote-ha.com", CoreHost: "homeassistant", CorePort: 8123, CoreScheme: "http"}
	if supervised {
		var user struct {
			PairingToken string `json:"pairing_token"`
		}
		if err := json.Unmarshal(data, &user); err != nil {
			return Options{}, errors.New("invalid options JSON")
		}
		options.PairingToken = user.PairingToken
	} else if err := json.Unmarshal(data, &options); err != nil {
		return Options{}, errors.New("invalid options JSON")
	}
	return options, nil
}

// Replace only this app's options through the supported Supervisor API. Read
// the stored value to retain !secret references; never touch /data/rha identity.
func migrateSupervisorOptions(client *http.Client, token string) error {
	if token == "" {
		return nil
	}
	call := func(method, path string, body []byte) ([]byte, error) {
		request, err := http.NewRequest(method, "http://supervisor/addons/self/"+path, bytes.NewReader(body))
		if err != nil {
			return nil, errors.New("cannot prepare app options migration")
		}
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			return nil, errors.New("app options migration unavailable")
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.New("app options migration rejected")
		}
		data, err := io.ReadAll(io.LimitReader(response.Body, 65536))
		if err != nil {
			return nil, errors.New("cannot read app options migration response")
		}
		return data, nil
	}
	data, err := call(http.MethodGet, "info", nil)
	if err != nil {
		return err
	}
	var info struct {
		Result string `json:"result"`
		Data   struct {
			Options map[string]json.RawMessage `json:"options"`
		} `json:"data"`
	}
	if json.Unmarshal(data, &info) != nil || info.Result != "ok" || info.Data.Options == nil {
		return errors.New("invalid app options migration response")
	}
	pairingToken, exists := info.Data.Options["pairing_token"]
	if exists && len(info.Data.Options) == 1 {
		return nil
	}
	if !exists {
		pairingToken = json.RawMessage(`""`)
	}
	body, err := json.Marshal(map[string]any{"options": map[string]json.RawMessage{"pairing_token": pairingToken}})
	if err != nil {
		return errors.New("cannot encode app options migration")
	}
	data, err = call(http.MethodPost, "options", body)
	if err != nil {
		return err
	}
	var result struct {
		Result string `json:"result"`
	}
	if json.Unmarshal(data, &result) != nil || result.Result != "ok" {
		return errors.New("app options migration failed")
	}
	return nil
}
