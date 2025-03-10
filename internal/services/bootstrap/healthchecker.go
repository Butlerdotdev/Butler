package bootstrap

import (
	"butler/internal/adapters/providers"
	"butler/internal/models"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// HealthChecker waits for VMs to report healthy and have allocated IPs.
type HealthChecker struct {
	provider providers.ProviderInterface
	logger   *zap.Logger
}

// NewHealthChecker initializes a HealthChecker.
func NewHealthChecker(provider providers.ProviderInterface, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		provider: provider,
		logger:   logger,
	}
}

// WaitForVMsToBeReady waits for all VMs to become healthy and have allocated IPs.
func (h *HealthChecker) WaitForVMsToBeReady(config *models.BootstrapConfig, timeout time.Duration) error {
	h.logger.Info("Waiting for VMs to report healthy and have allocated IPs")

	deadline := time.Now().Add(timeout)

	for _, node := range config.ManagementCluster.Nodes {
		for i := 1; i <= node.Count; i++ {
			vmName := fmt.Sprintf("%s-%s-%d", config.ManagementCluster.Name, node.Role, i)

			for time.Now().Before(deadline) {
				status, err := h.provider.GetVMStatus(vmName)
				if err != nil {
					h.logger.Warn("Failed to get VM status", zap.String("vm_name", vmName), zap.Error(err))
					time.Sleep(10 * time.Second)
					continue
				}

				if status.Healthy && status.IP != "" {
					h.logger.Info("VM is healthy and has an allocated IP", zap.String("vm_name", vmName), zap.String("ip", status.IP))
					break
				}

				h.logger.Info("Waiting for VM to be ready", zap.String("vm_name", vmName))
				time.Sleep(10 * time.Second)
			}

			// If VM is still not ready after timeout, return error
			status, err := h.provider.GetVMStatus(vmName)
			if err != nil || !status.Healthy || status.IP == "" {
				return fmt.Errorf("timeout: VM %s did not become healthy with an allocated IP", vmName)
			}
		}
	}

	h.logger.Info("All VMs are healthy and have allocated IPs")
	return nil
}
