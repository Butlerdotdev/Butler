// Package proxmox provides a Proxmox API client for handling HTTP communication.
//
// Copyright (c) 2025, The Butler Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxmox

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/butlerdotdev/butler/pkg/adapters/providers/proxmox/models"

	"go.uber.org/zap"
)

// ProxmoxClient handles raw HTTP communication with the Proxmox API.
type ProxmoxClient struct {
	ctx      context.Context
	logger   *zap.Logger
	endpoint string
	username string
	password string
	tokens   models.ProxmoxTokenData
	client   *http.Client
	nodes    []string
}

// NewProxmoxClient initializes a ProxmoxClient.
func NewProxmoxClient(ctx context.Context, endpoint, username, password, nodes string, logger *zap.Logger) *ProxmoxClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disables TLS verification (TODO: Improve security, make this optional)
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &ProxmoxClient{
		ctx:      ctx,
		logger:   logger,
		endpoint: endpoint,
		username: username,
		password: password,
		client:   client,
		nodes:    strings.Split(nodes, ","),
	}
}

// DoRequest makes an authenticated Proxmox API request.
func (n *ProxmoxClient) DoRequest(method, path string, payload interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", n.endpoint, path)

	// Obtain session tokens if not already set
	if n.tokens.Ticket == "" {
		var err error
		n.tokens, err = n.GetSessionToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get session tokens: %w", err)
		}
	}

	// Convert payload to JSON
	var jsonPayload []byte
	if payload != nil {
		var err error
		jsonPayload, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
		}
	}

	// Create request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers and cookies
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CSRFPreventionToken", n.tokens.CSRF)
	cookie := &http.Cookie{
		Name:  "PVEAuthCookie",
		Value: n.tokens.Ticket,
	}
	req.AddCookie(cookie)

	// Execute request
	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}

func (n *ProxmoxClient) GetSessionToken() (models.ProxmoxTokenData, error) {
	url := fmt.Sprintf("%s%s", n.endpoint, "/api2/json/access/ticket")

	payload := map[string]string{
		"username": n.username,
		"password": n.password,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return models.ProxmoxTokenData{}, fmt.Errorf("failed to create request for session token: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return models.ProxmoxTokenData{}, fmt.Errorf("HTTP request failed for session token: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ProxmoxTokenData{}, fmt.Errorf("failed to read session token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.ProxmoxTokenData{}, fmt.Errorf("session token request returned unhealthy status code: %s", body)
	}

	var tokenData models.ProxmoxSessionTokenResponse
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return models.ProxmoxTokenData{}, fmt.Errorf("failed to decode session token response: %w", err)
	}

	return tokenData.Data, nil
}
