// Package talos defines an adapter for Talos and bootstrapping the OS.
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

package talos

import (
	"butler/internal/adapters/exec"
	"butler/internal/models"
	"context"
	"fmt"
)

// TalosAdapter provides methods for managing Talos.
type TalosAdapter struct {
	client *TalosClient
}

// NewTalosAdapter initializes a new Talos adapter.
func NewTalosAdapter(execAdapter exec.ExecAdapter) *TalosAdapter {
	return &TalosAdapter{client: NewTalosClient(execAdapter)}
}

// Configure is required to implement PlatformAdapter.
func (t *TalosAdapter) Configure(ctx context.Context) error {
	// This can be a placeholder if not needed yet.
	return nil
}

// GenerateConfig delegates to TalosClient.
func (t *TalosAdapter) GenerateConfig(ctx context.Context, config models.TalosConfig) error {
	return t.client.GenerateConfig(ctx, config)
}

// ApplyConfig applies Talos configuration to a node.
func (t *TalosAdapter) ApplyConfig(ctx context.Context, node, configDir, role string, insecure bool) error {
	// Map "control-plane" to "controlplane" (Talos naming convention)
	roleMapping := map[string]string{
		"control-plane": "controlplane",
		"worker":        "worker",
	}

	mappedRole, exists := roleMapping[role]
	if !exists {
		return fmt.Errorf("invalid role %s for Talos configuration", role)
	}

	fmt.Printf("[DEBUG] Applying Talos config: %s/%s.yaml -> Node: %s\n", configDir, mappedRole, node)

	// Apply Talos configuration using the updated function
	return t.client.ApplyConfig(ctx, node, configDir, role, insecure)
}

// BootstrapControlPlane runs the Talos bootstrap process.
func (t *TalosAdapter) BootstrapControlPlane(ctx context.Context, node, configDir string) error {
	return t.client.BootstrapControlPlane(ctx, node, configDir)
}

// GetStatus fetches the Talos node status.
func (t *TalosAdapter) GetStatus(ctx context.Context) (models.PlatformStatus, error) {
	return t.client.GetStatus(ctx)
}
