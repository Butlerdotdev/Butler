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
	"butler/internal/adapters/platforms"
	"butler/internal/models"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// TalosInitializer handles configuring Talos on provisioned VMs.
type TalosInitializer struct {
	talos  platforms.PlatformAdapter
	logger *zap.Logger
}

// NewTalosInitializer creates a new Talos initializer.
func NewTalosInitializer(talos platforms.PlatformAdapter, logger *zap.Logger) *TalosInitializer {
	return &TalosInitializer{talos: talos, logger: logger}
}

// ConfigureTalos sets up Talos on the cluster nodes.
func (t *TalosInitializer) ConfigureTalos(ctx context.Context, config *models.TalosConfig, insecure bool) error {
	t.logger.Info("Generating Talos configuration")

	// Generate config
	err := t.talos.GenerateConfig(ctx, *config)
	if err != nil {
		return fmt.Errorf("failed to generate Talos config: %w", err)
	}

	// Apply Talos config to each control plane node
	for _, node := range config.ControlPlaneNodes {
		t.logger.Info("Applying Talos config to control plane", zap.String("node", node))
		err := t.talos.ApplyConfig(ctx, node, config.OutputDir, "control-plane", insecure)
		if err != nil {
			return fmt.Errorf("failed to apply Talos config to control plane: %w", err)
		}
	}

	// Apply Talos config to each worker node
	for _, node := range config.WorkerNodes {
		t.logger.Info("Applying Talos config to worker node", zap.String("node", node))
		err := t.talos.ApplyConfig(ctx, node, config.OutputDir, "worker", insecure)
		if err != nil {
			return fmt.Errorf("failed to apply Talos config to worker: %w", err)
		}
	}

	// Bootstrap Talos on a control plane node (run once)
	controlPlaneNode := config.ControlPlaneNodes[0]
	t.logger.Info("Bootstrapping Talos on control plane", zap.String("node", controlPlaneNode))
	err = t.talos.BootstrapControlPlane(ctx, controlPlaneNode, config.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to bootstrap Talos: %w", err)
	}

	t.logger.Info("Talos setup complete")
	return nil
}
