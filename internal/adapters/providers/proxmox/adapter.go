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
	"encoding/json"
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

	vmId, err := n.GetNextVMId(vm.AvailableVMIdStart, vm.AvailableVMIdEnd)
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
		Scsi0:   strings.Join([]string{vm.StorageLocation, ":", parseDisk(vm.Disk), ",iothread=on"}, ""),
		Numa:    false,
		Agent:   true,
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
	// We dont know what node the VM is on, so we need to get all VMs and find the one with the right name
	allVms, err := n.GetAllVms()
	if err != nil {
		n.logger.Error("Failed to get all VMs", zap.Error(err))
		return sharedModels.VMStatus{}, err
	}

	for _, vm := range allVms.Data {
		if vm.Name == vmName {
			n.logger.Info("Found VM", zap.String("name", vm.Name), zap.Int("id", vm.VMId), zap.String("status", vm.Status))
			var vmStatus sharedModels.VMStatus
			if vm.Status == "running" {
				vmStatus.Healthy = true
			} else {
				vmStatus.Healthy = false
			}
			vmStatus.IP = "" // TODO: Get IP address from Proxmox API
			return vmStatus, nil
		}
	}

	n.logger.Error("VM not found", zap.String("name", vmName))
	return sharedModels.VMStatus{}, fmt.Errorf("VM %s not found", vmName)
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

func (n *ProxmoxAdapter) GetNextVMId(vmRangeStart int, vmRangeEnd int) (int, error) {
	allVms, err := n.GetAllVms()
	if err != nil {
		n.logger.Error("Failed to get all VMs", zap.Error(err))
		return -1, err
	}

	ids := map[int]bool{}
	for _, vm := range allVms.Data {
		ids[vm.VMId] = true
	}

	for i := vmRangeStart; i <= vmRangeEnd; i++ {
		if !ids[i] {
			return i, nil
		}
	}

	// If no IDs are available in the range, return an error
	return -1, fmt.Errorf("no available VM IDs in the range %d-%d", vmRangeStart, vmRangeEnd)
}

func (n *ProxmoxAdapter) GetAllVms() (models.ProxmoxAllVMResponse, error) {
	n.logger.Info("Getting VM List")

	// Get the list of all VMs from Proxmox API
	path := "/api2/json/cluster/resources?type=vm"
	resp, err := n.client.DoRequest("GET", path, nil)
	if err != nil {
		n.logger.Error("Failed to send request to retrieve all VMs", zap.Error(err))
		return models.ProxmoxAllVMResponse{}, err
	}
	defer resp.Body.Close()

	// Read and check response
	body, _ := io.ReadAll(resp.Body)
	n.logger.Info("Response body", zap.ByteString("body", body))
	if resp.StatusCode >= 300 {
		n.logger.Error("Failed to get all VMs", zap.Int("status", resp.StatusCode), zap.ByteString("response", body))
		return models.ProxmoxAllVMResponse{}, fmt.Errorf("Failed to get all VMs: %s", body)
	}

	var vms models.ProxmoxAllVMResponse
	if err := json.Unmarshal(body, &vms); err != nil {
		n.logger.Error("Failed to unmarshal VM response", zap.Error(err))
		return models.ProxmoxAllVMResponse{}, fmt.Errorf("failed to unmarshal VM response: %w", err)
	}

	n.logger.Info("VMs retrieved successfully")
	return vms, nil
}
