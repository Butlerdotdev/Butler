package bootstrap

import (
	"butler/internal/models"
	service "butler/internal/services/bootstrap"
	"context"
	"fmt"

	"github.com/spf13/viper"
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

// HandleProvisionCluster loads config and calls the bootstrap service.
func (h *BootstrapHandler) HandleProvisionCluster() error {
	h.logger.Info("Handling cluster provisioning request...")

	// Load config into the BootstrapConfig model
	var config models.BootstrapConfig
	if err := viper.Unmarshal(&config); err != nil {
		h.logger.Error("Failed to load config", zap.Error(err))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate critical fields
	if err := validateBootstrapConfig(&config); err != nil {
		h.logger.Error("Configuration validation failed", zap.Error(err))
		return fmt.Errorf("configuration invalid: %w", err)
	}

	// Initialize bootstrap service
	bootstrapService, err := service.NewBootstrapService(h.ctx, &config, h.logger)
	if err != nil {
		h.logger.Error("Failed to initialize bootstrap service", zap.Error(err))
		return err
	}

	// Run provisioning process
	err = bootstrapService.ProvisionManagementCluster()
	if err != nil {
		h.logger.Error("Cluster provisioning failed", zap.Error(err))
		return err
	}

	h.logger.Info("Cluster provisioning completed successfully.")
	return nil
}

// validateBootstrapConfig handles validation of required fields in the config.
func validateBootstrapConfig(cfg *models.BootstrapConfig) error {
	if cfg.ManagementCluster.Name == "" {
		return fmt.Errorf("managementcluster.name is required")
	}
	if cfg.ManagementCluster.Provider == "" {
		return fmt.Errorf("managementcluster.provider is required")
	}
	return nil
}
