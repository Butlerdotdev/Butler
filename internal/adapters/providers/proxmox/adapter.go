// Package proxmox provides an adapter for provisioning and managing VMs on Proxmox VE.
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

package proxmox

import (
	sharedModels "butler/internal/models"

	"go.uber.org/zap"
)

// Proxmox Adapter implements ProviderInterface.
type ProxmoxAdapter struct {
	client *ProxmoxClient
	logger *zap.Logger
}

// NewProxmoxAdapter initializes the Proxmox adapter.
func NewProxmoxAdapter(client *ProxmoxClient, logger *zap.Logger) *ProxmoxAdapter {
	return &ProxmoxAdapter{
		client: client,
		logger: logger,
	}
}

// CreateVM provisions a VM in Proxmox VE.
func (n *ProxmoxAdapter) CreateVM(vm sharedModels.VMConfig) (string, error) {
	// TODO: Implement
	return vm.Name, nil
}

// DeleteVM removes a VM from Proxmox.
func (n *ProxmoxAdapter) DeleteVM(vmID string) error {
	// TODO: Implement
	return nil
}

// GetVMStatus fetches the VM's health status and IP address from Proxmox VE.
func (n *ProxmoxAdapter) GetVMStatus(vmName string) (sharedModels.VMStatus, error) {
	// TODO: Implement

	return sharedModels.VMStatus{}, nil
}
