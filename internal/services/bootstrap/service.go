// // Package bootstrap provides services for provisioning the Butler management cluster.
// //
// // Copyright (c) 2025, The Butler Authors
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package bootstrap

import (
	"butler/internal/adapters/providers"
	"butler/internal/models"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// BootstrapService orchestrates provisioning the management cluster.
type BootstrapService struct {
	logger      *zap.Logger
	provider    providers.ProviderInterface
	provisioner *Provisioner
	healthCheck *HealthChecker
}

// NewBootstrapService initializes a new BootstrapService with dependencies.
func NewBootstrapService(ctx context.Context, config *models.BootstrapConfig, logger *zap.Logger) (*BootstrapService, error) {
	provider, err := providers.NewProviderFactory(ctx, config.ManagementCluster.Provider, config.ManagementCluster.Nutanix.ToMap(), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	return &BootstrapService{
		logger:      logger,
		provider:    provider,
		provisioner: NewProvisioner(provider, logger),
		healthCheck: NewHealthChecker(provider, logger),
	}, nil
}

// ProvisionManagementCluster provisions the management cluster VMs and waits for them to be ready.
func (b *BootstrapService) ProvisionManagementCluster(config *models.BootstrapConfig) error {
	b.logger.Info("Starting provisioning of management cluster", zap.String("cluster_name", config.ManagementCluster.Name))

	// Provision VMs
	err := b.provisioner.ProvisionVMs(config)
	if err != nil {
		return err
	}

	// Wait for VMs to be healthy
	err = b.healthCheck.WaitForVMsToBeReady(config, 10*time.Minute)
	if err != nil {
		return err
	}

	b.logger.Info("Management cluster provisioned successfully!")
	return nil
}
