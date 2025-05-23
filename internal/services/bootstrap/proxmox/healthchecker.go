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
	"fmt"
	"time"

	"github.com/butlerdotdev/butler/pkg/adapters/providers"
	"github.com/butlerdotdev/butler/pkg/models"
	"go.uber.org/zap"
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

// WaitForVMsToBeReady waits for all VMs to become healthy and collects their assigned IPs.
func (h *HealthChecker) WaitForVMsToBeReady(config *models.BootstrapConfig, timeout time.Duration) (map[string]string, error) {
	h.logger.Info("Waiting for VMs to report healthy and have allocated IPs")

	deadline := time.Now().Add(timeout)
	nodeIPs := make(map[string]string)

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
					nodeIPs[vmName] = status.IP
					break
				}

				h.logger.Info("Waiting for VM to be ready", zap.String("vm_name", vmName))
				time.Sleep(10 * time.Second)
			}

			// If VM is still not ready after timeout, return an error
			if _, exists := nodeIPs[vmName]; !exists {
				return nil, fmt.Errorf("timeout: VM %s did not become healthy with an allocated IP", vmName)
			}
		}
	}

	h.logger.Info("All VMs are healthy and have allocated IPs")
	return nodeIPs, nil
}
