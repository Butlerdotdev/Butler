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
	"os"
	"path/filepath"
	"time"
	// "path/filepath"
)

// TalosClient interacts with Talos using the ExecAdapter.
type TalosClient struct {
	execAdapter exec.ExecAdapter
}

// NewTalosClient creates a new Talos client.
func NewTalosClient(execAdapter exec.ExecAdapter) *TalosClient {
	return &TalosClient{execAdapter: execAdapter}
}

func (c *TalosClient) GenerateConfig(ctx context.Context, config models.TalosConfig) error {
	_, err := c.execAdapter.RunCommand(ctx, "talosctl", "gen", "config",
		config.ClusterName, fmt.Sprintf("https://%s", config.ControlPlaneEndpoint),
		"--output", config.OutputDir,
		"--config-patch", `[{"op": "replace", "path": "/machine/time", "value": {"disabled": true}}]`,
	)
	return err
}

// ApplyConfig applies the Talos configuration file to a node
func (c *TalosClient) ApplyConfig(ctx context.Context, node, configDir, role string, insecure bool) error {
	// Map role to match Talos file naming convention
	roleMapping := map[string]string{
		"control-plane": "controlplane",
		"worker":        "worker",
	}

	mappedRole, exists := roleMapping[role]
	if !exists {
		return fmt.Errorf("invalid role %s for Talos configuration", role)
	}

	// Get absolute paths for config files
	absConfigDir, err := filepath.Abs(configDir)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path of config directory: %w", err)
	}

	configFile := filepath.Join(absConfigDir, mappedRole+".yaml")
	talosConfigPath := filepath.Join(absConfigDir, "talosconfig")

	// Ensure Talos configuration exists
	if _, err := os.Stat(talosConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("talos configuration file missing: %s", talosConfigPath)
	}

	// Ensure Talos knows about the endpoint (fix: explicitly set --talosconfig)
	_, err = c.execAdapter.RunCommand(ctx, "talosctl", "config", "endpoint", node, "--talosconfig", talosConfigPath)
	if err != nil {
		return fmt.Errorf("failed to configure Talos endpoint: %w", err)
	}

	// Apply the configuration
	args := []string{
		"apply-config",
		"--nodes", node,
		"--file", configFile,
		"--talosconfig", talosConfigPath,
	}
	if insecure {
		args = append(args, "--insecure")
	}

	_, err = c.execAdapter.RunCommand(ctx, "talosctl", args...)
	if err != nil {
		return fmt.Errorf("failed to apply Talos config: %w", err)
	}

	return nil
}

func (c *TalosClient) BootstrapControlPlane(ctx context.Context, node, configDir string) error {

	// Configure short sleep for talos nodes to register before bootstrapping.
	waitTime := 180 * time.Second
	time.Sleep(waitTime)

	// Get absolute paths for config files
	absConfigDir, err := filepath.Abs(configDir)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path of config directory: %w", err)
	}

	talosConfigPath := filepath.Join(absConfigDir, "talosconfig")

	// Ensure Talos configuration exists
	if _, err := os.Stat(talosConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("talos configuration file missing: %s", talosConfigPath)
	}

	// Ensure Talos knows about the endpoint
	_, err = c.execAdapter.RunCommand(ctx, "talosctl", "config", "endpoint", node, "--talosconfig", talosConfigPath)
	if err != nil {
		return fmt.Errorf("failed to configure Talos endpoint: %w", err)
	}

	// Bootstrap Talos on the control plane node
	args := []string{
		"bootstrap",
		"--nodes", node,
		"--talosconfig", talosConfigPath,
	}

	_, err = c.execAdapter.RunCommand(ctx, "talosctl", args...)
	if err != nil {
		return fmt.Errorf("failed to bootstrap Talos: %w", err)
	}

	return nil
}

// GetStatus fetches Talos health status.
func (c *TalosClient) GetStatus(ctx context.Context) (models.PlatformStatus, error) {
	result, err := c.execAdapter.RunCommand(ctx, "talosctl", "health")
	if err != nil {
		return models.PlatformStatus{}, err
	}

	// Extract the Stdout to pass as a string
	return parseTalosHealth(result.Stdout), nil
}
