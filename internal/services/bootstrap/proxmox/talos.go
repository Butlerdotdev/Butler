// Package bootstrap provides services for provisioning the Butler management cluster.
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

package bootstrap

import (
	"butler/internal/adapters/platforms/talos"
	"butler/internal/models"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// TalosInitializer handles configuring Talos on provisioned VMs.
type TalosInitializer struct {
	talosAdapter *talos.TalosAdapter
	logger       *zap.Logger
}

// NewTalosInitializer creates a new Talos initializer.
func NewTalosInitializer(talosAdapter *talos.TalosAdapter, logger *zap.Logger) *TalosInitializer {
	return &TalosInitializer{talosAdapter: talosAdapter, logger: logger}
}

// ConfigureTalos sets up Talos on the cluster nodes.
func (t *TalosInitializer) ConfigureTalos(ctx context.Context, config *models.TalosConfig, insecure bool) error {
	t.logger.Info("Starting Talos setup", zap.String("cluster", config.ClusterName))

	// Generate Talos Configuration
	if err := t.GenerateConfig(ctx, config); err != nil {
		return fmt.Errorf("failed to generate Talos config: %w", err)
	}

	// Apply Config to Control Plane Nodes
	for _, node := range config.ControlPlaneNodes {
		if err := t.ApplyConfig(ctx, node, config.OutputDir, "controlplane.yaml", insecure); err != nil {
			return fmt.Errorf("failed to apply Talos config to control plane node %s: %w", node, err)
		}
	}

	// Apply Config to Worker Nodes
	for _, node := range config.WorkerNodes {
		if err := t.ApplyConfig(ctx, node, config.OutputDir, "worker.yaml", insecure); err != nil {
			return fmt.Errorf("failed to apply Talos config to worker node %s: %w", node, err)
		}
	}

	// Wait for Nodes to Register
	t.WaitForNodesToRegister()

	// Set Talos Endpoint
	controlPlaneNode := config.ControlPlaneNodes[0]
	if err := t.SetEndpoint(ctx, controlPlaneNode); err != nil {
		return fmt.Errorf("failed to configure Talos endpoint: %w", err)
	}

	// Bootstrap Talos on the First Control Plane Node
	if err := t.BootstrapControlPlane(ctx, controlPlaneNode); err != nil {
		return fmt.Errorf("failed to bootstrap Talos on control plane: %w", err)
	}

	// Retrieve KubeConfig
	if err := t.RetrieveKubeConfig(ctx, controlPlaneNode); err != nil {
		return fmt.Errorf("failed to retrieve kubeconfig: %w", err)
	}

	t.logger.Info("Talos setup complete")
	return nil
}

// GenerateConfig generates Talos configuration files.
func (t *TalosInitializer) GenerateConfig(ctx context.Context, config *models.TalosConfig) error {
	t.logger.Info("Generating Talos configuration",
		zap.String("cluster", config.ClusterName),
		zap.String("endpoint", config.ControlPlaneEndpoint),
	)

	_, err := t.talosAdapter.ExecuteCommand(ctx,
		"gen", "config", config.ClusterName, fmt.Sprintf("https://%s", config.ControlPlaneEndpoint),
		"--output", config.OutputDir,
		"--config-patch", `[
	{
		"op": "replace",
		"path": "/machine/time",
		"value": { "disabled": true }
	},
	{
		"op": "add",
		"path": "/machine/kernel",
		"value": {
			"modules": [
				{ "name": "openvswitch" }
			]
		}
	},
	{
		"op": "add",
		"path": "/cluster/network/cni",
		"value": { "name": "none" }
	},
	{
		"op": "replace",
		"path": "/cluster/network/podSubnets",
		"value": [ "10.16.0.0/16" ]
	},
	{
		"op": "add",
		"path": "/machine/kubelet/extraMounts",
		"value": [
			{
				"source": "/run/openvswitch",
				"destination": "/run/openvswitch",
				"type": "bind",
				"options": ["rbind", "rw"]
			},
			{
				"source": "/run/ovn",
				"destination": "/run/ovn",
				"type": "bind",
				"options": ["rbind", "rw"]
			},
			{
				"source": "/var/log/openvswitch",
				"destination": "/var/log/openvswitch",
				"type": "bind",
				"options": ["rbind", "rw"]
			},
			{
				"source": "/var/log/ovn",
				"destination": "/var/log/ovn",
				"type": "bind",
				"options": ["rbind", "rw"]
			}
		]
	}
]`,
	)
	return err
}

// ApplyConfig applies the Talos configuration to a node.
func (t *TalosInitializer) ApplyConfig(ctx context.Context, node, configDir, configFile string, insecure bool) error {
	t.logger.Info("Applying Talos config", zap.String("node", node), zap.String("file", configFile))

	args := []string{
		"apply-config",
		"--nodes", node,
		"--file", fmt.Sprintf("%s/%s", configDir, configFile),
		"--talosconfig", "talosconfig/talosconfig",
	}
	if insecure {
		args = append(args, "--insecure")
	}

	_, err := t.talosAdapter.ExecuteCommand(ctx, args...)
	return err
}

// WaitForNodesToRegister adds a delay to allow nodes to register.
func (t *TalosInitializer) WaitForNodesToRegister() {
	waitTime := 180 * time.Second
	t.logger.Info("Waiting for Talos nodes to register before bootstrapping", zap.Duration("waitTime", waitTime))
	time.Sleep(waitTime)
}

// SetEndpoint configures the Talos endpoint.
func (t *TalosInitializer) SetEndpoint(ctx context.Context, node string) error {
	t.logger.Info("Configuring Talos endpoint", zap.String("node", node))

	_, err := t.talosAdapter.ExecuteCommand(ctx,
		"config", "endpoint", node, "--talosconfig", "talosconfig/talosconfig",
	)
	return err
}

// BootstrapControlPlane bootstraps Talos on a control-plane node.
func (t *TalosInitializer) BootstrapControlPlane(ctx context.Context, node string) error {
	t.logger.Info("Bootstrapping Talos control plane", zap.String("node", node))

	_, err := t.talosAdapter.ExecuteCommand(ctx,
		"bootstrap", "--nodes", node, "--talosconfig", "talosconfig/talosconfig",
	)
	return err
}

// RetrieveKubeConfig fetches and stores the Kubernetes kubeconfig.
func (t *TalosInitializer) RetrieveKubeConfig(ctx context.Context, node string) error {
	t.logger.Info("Retrieving kubeconfig from Talos", zap.String("node", node))

	_, err := t.talosAdapter.ExecuteCommand(ctx,
		"kubeconfig", "talosconfig/kubeconfig",
		"--nodes", node,
		"--talosconfig", "talosconfig/talosconfig",
		"--force",
		"--merge",
	)
	return err
}
