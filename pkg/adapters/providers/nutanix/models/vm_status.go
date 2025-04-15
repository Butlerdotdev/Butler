// Package models defines data structures for Butler's cluster provisioning.
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

package models

// NutanixVMStatus represents the structure of the response when querying VM status.
type NutanixVMStatus struct {
	Entities []VMEntity `json:"entities"`
}

// VMEntity represents an individual VM entity returned from Nutanix API.
type VMEntity struct {
	Status VMStatus `json:"status"`
}

// VMStatus represents the status details of a Nutanix VM.
type VMStatus struct {
	State     string      `json:"state"`
	Resources VMResources `json:"resources"`
}

// VMResources represents VM resource details including power state and NICs.
type VMResources struct {
	PowerState string  `json:"power_state"`
	NICs       []VMNic `json:"nic_list"` // Renamed from Nic
}

// IPEndpoint represents an IP address assigned to a NIC.
type IPEndpoint struct {
	IP string `json:"ip"`
}

// VMNic represents a network interface card (NIC) of the VM.
type VMNic struct { // Renamed from Nic to avoid conflict
	IpEndpointList []IPEndpoint `json:"ip_endpoint_list"`
}
