// Package bootstrap provides handlers for provisioning the Butler management cluster.
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
	"butler/internal/models"
	service "butler/internal/services/bootstrap"
	"context"
	"fmt"

	"go.uber.org/zap"
)

// BootstrapHandler handles requests for provisioning clusters.
type BootstrapHandler struct {
	ctx    context.Context
	logger *zap.Logger
}

// NewBootstrapHandler initializes a new BootstrapHandler.
func NewBootstrapHandler(ctx context.Context, logger *zap.Logger) *BootstrapHandler {
	return &BootstrapHandler{
		ctx:    ctx,
		logger: logger,
	}
}

// HandleProvisionCluster validates input and calls the bootstrap service.
func (h *BootstrapHandler) HandleProvisionCluster(config *models.BootstrapConfig) error {
	h.logger.Info("Handling cluster provisioning request...")

	// validate Required Fields
	if config.ManagementCluster.Name == "" {
		return fmt.Errorf("cluster name is required")
	}
	if config.ManagementCluster.Provider == "" {
		return fmt.Errorf("provider is required (e.g., 'nutanix', 'aws', 'azure')")
	}
	if len(config.ManagementCluster.Nodes) == 0 {
		return fmt.Errorf("no nodes defined in configuration")
	}

	// Call the Service
	bootstrapService, err := service.NewBootstrapService(h.ctx, config, h.logger)
	if err != nil {
		h.logger.Error("Failed to initialize bootstrap service", zap.Error(err))
		return err
	}

	// Proceed with provisioning
	err = bootstrapService.ProvisionManagementCluster(config)
	if err != nil {
		h.logger.Error("Cluster provisioning failed", zap.Error(err))
		return err
	}

	h.logger.Info("Cluster provisioning completed successfully.")
	return nil
}
