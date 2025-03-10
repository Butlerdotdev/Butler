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

// VMConfig represents a generic VM configuration that works across all providers.
type VMConfig struct {
	Name        string
	Role        string
	CPU         int
	RAM         string
	Disk        string
	Count       int
	IsoUUID     string
	ClusterUUID string
	SubnetUUID  string
}

// ClusterConfig represents a generic cluster configuration across providers.
type ClusterConfig struct {
	Name     string
	Nodes    []VMConfig
	Provider string
}

// VMStatus represents the status of a VM in Nutanix.
type VMStatus struct {
	Healthy bool   `json:"healthy"`
	IP      string `json:"ip"`
}