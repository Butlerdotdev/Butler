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
	"butler/internal/models"
	"butler/pkg/adapters/providers"
	"fmt"

	"go.uber.org/zap"
)

// Provisioner handles VM provisioning.
type Provisioner struct {
	provider providers.ProviderInterface
	logger   *zap.Logger
}

// NewProvisioner initializes a Provisioner.
func NewProvisioner(provider providers.ProviderInterface, logger *zap.Logger) *Provisioner {
	return &Provisioner{
		provider: provider,
		logger:   logger,
	}
}

// ProvisionVMs provisions the required VMs for the management cluster.
func (p *Provisioner) ProvisionVMs(config *models.BootstrapConfig) error {
	for _, node := range config.ManagementCluster.Nodes {
		for i := 1; i <= node.Count; i++ {
			vmName := fmt.Sprintf("%s-%s-%d", config.ManagementCluster.Name, node.Role, i)
			p.logger.Info("Creating VM", zap.String("name", vmName), zap.String("role", node.Role))

			vmConfig := models.VMConfig{
				Name:               vmName,
				CPU:                node.CPU,
				RAM:                node.RAM,
				Disk:               node.Disk,
				IsoUUID:            node.IsoUUID,
				SubnetUUID:         config.ManagementCluster.Nutanix.SubnetUUID,
				ClusterUUID:        config.ManagementCluster.Nutanix.ClusterUUID,
				StorageLocation:    config.ManagementCluster.Proxmox.StorageLocation,
				AvailableVMIdStart: config.ManagementCluster.Proxmox.AvailableVMIdStart,
				AvailableVMIdEnd:   config.ManagementCluster.Proxmox.AvailableVMIdEnd,
			}

			_, err := p.provider.CreateVM(vmConfig)
			if err != nil {
				p.logger.Error("Failed to create VM", zap.String("vm_name", vmName), zap.Error(err))
				return err
			}
		}
	}
	return nil
}

// SeparateNodesByRole classifies VMs into control planes and workers based on role.
func (p *Provisioner) SeparateNodesByRole(config *models.BootstrapConfig, nodeIPs map[string]string) ([]string, []string, error) {
	var controlPlanes, workers []string

	for _, node := range config.ManagementCluster.Nodes {
		for i := 1; i <= node.Count; i++ {
			vmName := fmt.Sprintf("%s-%s-%d", config.ManagementCluster.Name, node.Role, i)
			ip, exists := nodeIPs[vmName]
			if !exists {
				return nil, nil, fmt.Errorf("missing IP for VM %s", vmName)
			}

			if node.Role == "control-plane" {
				controlPlanes = append(controlPlanes, ip)
			} else if node.Role == "worker" {
				workers = append(workers, ip)
			}
		}
	}

	if len(controlPlanes) == 0 {
		return nil, nil, fmt.Errorf("no control plane nodes found for Talos")
	}

	return controlPlanes, workers, nil
}
