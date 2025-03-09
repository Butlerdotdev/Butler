// Package nutanix provides a Nutanix API client for handling HTTP communication.
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

package nutanix

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// NutanixClient handles raw HTTP communication with the Nutanix API.
type NutanixClient struct {
	ctx      context.Context
	logger   *zap.Logger
	endpoint string
	username string
	password string
	client   *http.Client
}

// NewNutanixClient initializes a NutanixClient.
func NewNutanixClient(ctx context.Context, endpoint, username, password string, logger *zap.Logger) *NutanixClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disables TLS verification (TODO: Improve security)
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &NutanixClient{
		ctx:      ctx,
		logger:   logger,
		endpoint: endpoint,
		username: username,
		password: password,
		client:   client,
	}
}

// DoRequest makes an authenticated Nutanix API request.
func (n *NutanixClient) DoRequest(method, path string, payload interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", n.endpoint, path)

	// Encode credentials
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(n.username+":"+n.password))

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

	// Set headers
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}
