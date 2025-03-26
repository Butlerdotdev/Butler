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
	"butler/internal/adapters/providers/proxmox/models"
	sharedModels "butler/internal/models"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

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
	n.logger.Info("Creating VM", zap.String("name", vm.Name), zap.Int("CPU", vm.CPU), zap.String("RAM", vm.RAM), zap.String("Disk", vm.Disk))

	vmId, err := n.GetNextVMId()
	if err != nil {
		n.logger.Error("Failed to get next VM ID", zap.String("name", vm.Name), zap.Error(err))
		return "", err
	}
	// Construct the Proxmox VM payload using structs
	payload := models.ProxmoxVMConfig{
		VMId:    vmId,
		Name:    vm.Name,
		OSType:  "l26",
		Memory:  parseRAM(vm.RAM),
		Cores:   vm.CPU,
		Sockets: 1,
		Start:   true,
		OnBoot:  true,
		Ide2:    strings.Join([]string{vm.IsoUUID, "media=cdrom"}, ","),
		Scsihw:  "virtio-scsi-single",
		Scsi0:   strings.Join([]string{vm.StorageLocation, ":", parseDisk(vm.Disk), "iothread=on"}, ""),
		Numa:    false,
		Cpu:     "host",
		Net0:    "virtio,bridge=vmbr0,firewall=1",
	}

	// Assign the VM to a random node
	randomNode, err := n.GetRandomNode()
	if err != nil {
		return "", fmt.Errorf("failed to get random node: %w", err)
	}

	path := fmt.Sprintf("/api2/json/nodes/%s/qemu", randomNode)
	resp, err := n.client.DoRequest("POST", path, payload)
	if err != nil {
		n.logger.Error("Failed to send VM creation request", zap.String("name", vm.Name), zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	// Read and check response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		n.logger.Error("Failed to create VM", zap.String("name", vm.Name), zap.Int("status", resp.StatusCode), zap.ByteString("response", body))
		return "", fmt.Errorf("failed to create VM %s: %s", vm.Name, body)
	}

	n.logger.Info("VM created successfully", zap.String("name", vm.Name))
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

func (n *ProxmoxAdapter) GetRandomNode() (string, error) {
	// return random node from n.client.nodes
	if len(n.client.nodes) == 0 {
		return "", fmt.Errorf("no nodes available")
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(n.client.nodes))
	return n.client.nodes[randomIndex], nil
}

func (n *ProxmoxAdapter) GetNextVMId() (int, error) {
	// TODO: Implement with range defined in the config

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(900)
	return randomIndex + 100, nil
}
