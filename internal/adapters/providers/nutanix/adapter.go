// Package nutanix provides an adapter for provisioning and managing VMs on Nutanix AHV.
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

package nutanix

import (
	"butler/internal/adapters/providers/nutanix/models"
	sharedModels "butler/internal/models"
	"fmt"
	"io"

	"go.uber.org/zap"
)

// NutanixAdapter implements ProviderInterface.
type NutanixAdapter struct {
	client *NutanixClient
	logger *zap.Logger
}

// NewNutanixAdapter initializes the Nutanix adapter.
func NewNutanixAdapter(client *NutanixClient, logger *zap.Logger) *NutanixAdapter {
	return &NutanixAdapter{
		client: client,
		logger: logger,
	}
}

// CreateVM provisions a VM in Nutanix AHV.
func (n *NutanixAdapter) CreateVM(vm sharedModels.VMConfig) (string, error) {
	n.logger.Info("Creating VM", zap.String("name", vm.Name), zap.Int("CPU", vm.CPU), zap.String("RAM", vm.RAM), zap.String("Disk", vm.Disk))

	// Construct the Nutanix VM payload using structs
	payload := models.NutanixVMConfig{
		Metadata: models.Metadata{
			Kind: "vm",
		},
		Spec: models.Spec{
			Name: vm.Name,
			Resources: models.Resources{
				PowerState:        "ON",
				NumSockets:        vm.CPU,
				NumVCPUsPerSocket: 1,
				MemorySizeMiB:     parseRAM(vm.RAM),
				BootConfig: models.BootConfig{
					BootDeviceOrderList: []string{"CDROM", "DISK"},
				},
				DiskList: []models.Disk{
					{
						DeviceProperties: models.DeviceProperties{
							DeviceType: "DISK",
							DiskAddress: models.DiskAddress{
								AdapterType: "SCSI",
								DeviceIndex: 0,
							},
						},
						DiskSizeMiB: parseDisk(vm.Disk),
					},
					{
						DeviceProperties: models.DeviceProperties{
							DeviceType: "CDROM",
							DiskAddress: models.DiskAddress{
								AdapterType: "IDE",
								DeviceIndex: 1,
							},
						},
						DataSourceReference: &models.DataSourceReference{
							Kind: "image",
							UUID: vm.IsoUUID,
						},
					},
				},
				NicList: []models.Nic{
					{
						SubnetReference: models.SubnetReference{
							Kind: "subnet",
							UUID: vm.SubnetUUID,
						},
					},
				},
			},
			ClusterReference: models.ClusterReference{
				Kind: "cluster",
				UUID: vm.ClusterUUID,
			},
		},
	}

	// Send API request to Nutanix
	resp, err := n.client.DoRequest("POST", "/api/nutanix/v3/vms", payload)
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

// DeleteVM removes a VM from Nutanix.
func (n *NutanixAdapter) DeleteVM(vmID string) error {
	url := fmt.Sprintf("/api/nutanix/v3/vms/%s", vmID)

	resp, err := n.client.DoRequest("DELETE", url, nil)
	if err != nil {
		n.logger.Error("Failed to delete VM", zap.String("vmID", vmID), zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		n.logger.Error("Failed to delete VM", zap.String("vmID", vmID), zap.Int("status", resp.StatusCode))
		return fmt.Errorf("failed to delete VM %s", vmID)
	}

	n.logger.Info("VM deleted successfully", zap.String("vmID", vmID))
	return nil
}
