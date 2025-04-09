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
	"encoding/json"
	"fmt"
	"io"
	"strings"

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

	diskList := buildDiskList(vm)

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
				DiskList: diskList,
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

// GetVMStatus fetches the VM's health status and IP address from Nutanix Prism Central.
func (n *NutanixAdapter) GetVMStatus(vmName string) (sharedModels.VMStatus, error) {
	n.logger.Info("Fetching VM status", zap.String("vm_name", vmName))

	apiPath := "/api/nutanix/v3/vms/list"
	requestPayload := map[string]interface{}{
		"kind":   "vm",
		"filter": fmt.Sprintf("vm_name==%s", vmName),
	}

	// Send request using DoRequest
	resp, err := n.client.DoRequest("POST", apiPath, requestPayload)
	if err != nil {
		n.logger.Error("Failed to fetch VM status", zap.String("vm_name", vmName), zap.Error(err))
		return sharedModels.VMStatus{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return sharedModels.VMStatus{}, fmt.Errorf("nutanix API returned error: %d", resp.StatusCode)
	}

	// Read response
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Parse JSON response using the NutanixVMStatus model
	var responseData models.NutanixVMStatus
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return sharedModels.VMStatus{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if VM was found
	if len(responseData.Entities) == 0 {
		return sharedModels.VMStatus{}, fmt.Errorf("VM %s not found in Nutanix", vmName)
	}

	vm := responseData.Entities[0]

	// Extract power state & execution state
	isPoweredOn := strings.EqualFold(vm.Status.Resources.PowerState, "ON")
	isComplete := (vm.Status.State == "COMPLETE")

	// Extract IP from nic_list > ip_endpoint_list
	var assignedIP string
	for _, nic := range vm.Status.Resources.NICs {
		if len(nic.IpEndpointList) > 0 {
			assignedIP = nic.IpEndpointList[0].IP
			break // Stop at first found IP
		}
	}

	// Determine overall health:
	isHealthy := isPoweredOn && isComplete && assignedIP != ""

	// Log extracted values
	n.logger.Info("Fetched VM status successfully",
		zap.String("vm_name", vmName),
		zap.Bool("powered_on", isPoweredOn),
		zap.String("state", vm.Status.State),
		zap.Bool("healthy", isHealthy),
		zap.String("ip", assignedIP),
	)

	return sharedModels.VMStatus{
		Healthy: isHealthy,
		IP:      assignedIP,
	}, nil
}

func (n *NutanixAdapter) GetClusterUuids() ([]models.NutanixClusterEntities, error) {
	requestPayload := map[string]interface{}{
		"kind":   "cluster",
		"length": 1,
		"offset": 0,
	}
	resp, err := n.client.DoRequest("POST", "/api/nutanix/v3/clusters/list", requestPayload)
	if err != nil {
		return []models.NutanixClusterEntities{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return []models.NutanixClusterEntities{}, fmt.Errorf("failed to fetch clusters: %d", resp.StatusCode)
	}

	var clusters models.NutanixClusterList
	if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
		return []models.NutanixClusterEntities{}, err
	}

	return clusters.Entities, nil
}

func (n *NutanixAdapter) GetSubnetUuids(clusterUUID string) ([]models.NutanixSubnetEntities, error) {
	requestPayload := map[string]interface{}{
		"kind":   "subnet",
		"length": 1,
		"offset": 0,
	}
	resp, err := n.client.DoRequest("POST", "/api/nutanix/v3/subnets/list", requestPayload)
	if err != nil {
		return []models.NutanixSubnetEntities{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return []models.NutanixSubnetEntities{}, fmt.Errorf("failed to fetch subnets: %d", resp.StatusCode)
	}

	var subnets models.NutanixSubnetList
	if err := json.NewDecoder(resp.Body).Decode(&subnets); err != nil {
		return []models.NutanixSubnetEntities{}, err
	}

	return subnets.Entities, nil
}
