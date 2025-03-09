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
	"butler/internal/adapters/providers"
	"butler/internal/models"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// BootstrapService orchestrates provisioning the management cluster.
type BootstrapService struct {
	logger   *zap.Logger
	provider providers.ProviderInterface
}

// NewBootstrapService initializes a new BootstrapService.
func NewBootstrapService(ctx context.Context, config *models.BootstrapConfig, logger *zap.Logger) (*BootstrapService, error) {
	// Initialize the correct cloud provider adapter
	provider, err := providers.NewProviderFactory(ctx, config.ManagementCluster.Provider, config.ManagementCluster.Nutanix.ToMap(), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	return &BootstrapService{
		logger:   logger,
		provider: provider,
	}, nil
}

// ProvisionManagementCluster provisions the management cluster VMs.
func (b *BootstrapService) ProvisionManagementCluster(config *models.BootstrapConfig) error {
	b.logger.Info("Starting provisioning of management cluster", zap.String("cluster_name", config.ManagementCluster.Name))

	for _, node := range config.ManagementCluster.Nodes {
		for i := 1; i <= node.Count; i++ {
			vmName := fmt.Sprintf("%s-%s-%d", config.ManagementCluster.Name, node.Role, i)
			b.logger.Info("Creating VM", zap.String("name", vmName), zap.String("role", node.Role))

			vmConfig := models.VMConfig{
				Name:        vmName,
				CPU:         node.CPU,
				RAM:         node.RAM,
				Disk:        node.Disk,
				IsoUUID:     node.IsoUUID,
				SubnetUUID:  config.ManagementCluster.Nutanix.SubnetUUID,
				ClusterUUID: config.ManagementCluster.Nutanix.ClusterUUID,
			}

			// Call the provider's CreateVM function
			_, err := b.provider.CreateVM(vmConfig)
			if err != nil {
				b.logger.Error("Failed to create VM", zap.String("vm_name", vmName), zap.Error(err))
				return err
			}
		}
	}

	b.logger.Info("Management cluster provisioned successfully!")
	return nil
}
