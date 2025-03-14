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

// BootstrapConfig defines the cluster provisioning configuration.
type BootstrapConfig struct {
	ManagementCluster ManagementClusterConfig `yaml:"managementCluster"`
}

// ManagementClusterConfig holds the cluster configuration.
type ManagementClusterConfig struct {
	Name       string        `yaml:"name"`
	Provider   string        `yaml:"provider"`
	Nutanix    NutanixConfig `yaml:"nutanix,omitempty"`
	Nodes      []NodeConfig  `yaml:"nodes"`
	Talos      TalosConfig   `yaml:"talos"`
	ClusterAPI ClusterAPI    `yaml:"clusterAPI"`
	Flux       FluxConfig    `yaml:"flux"`
}

type FluxConfig struct {
	GitOwner      string `yaml:"gitOwner"`
	GitRepository string `yaml:"gitRepository"`
	GitBranch     string `yaml:"gitBranch"`
	GitPath       string `yaml:"gitPath"`
	GitHostname   string `yaml:"gitHostname"`
	GitPAT        string `yaml:"gitPAT"`
}

// NutanixConfig defines the Nutanix API connection and cluster details.
type NutanixConfig struct {
	Endpoint    string `yaml:"endpoint"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	ClusterUUID string `yaml:"clusterUUID,omitempty"`
	SubnetUUID  string `yaml:"subnetUUID,omitempty"`
}

// NodeConfig represents a single VM configuration.
type NodeConfig struct {
	Role    string `yaml:"role"`
	Count   int    `yaml:"count"`
	CPU     int    `yaml:"cpu"`
	RAM     string `yaml:"ram"`
	Disk    string `yaml:"disk"`
	IsoUUID string `yaml:"isoUUID"`
}

// TalosConfig holds Talos Linux bootstrapping details.
type TalosConfig struct {
	Version              string   `yaml:"version"`
	ControlPlaneEndpoint string   `yaml:"controlPlaneEndpoint"`
	ControlPlaneVIP      string   `yaml:"controlPlaneVIP"`
	ClusterName          string   `yaml:"clusterName"`
	CIDR                 string   `yaml:"cidr"`
	Gateway              string   `yaml:"gateway"`
	ControlPlaneIP       string   `json:"control_plane_ip"`
	OutputDir            string   `json:"output_dir"`
	ControlPlaneNodes    []string `json:"control_plane_nodes"`
	WorkerNodes          []string `json:"worker_nodes"`
}

// ClusterAPI represents the Cluster API provider settings.
type ClusterAPI struct {
	Version              string `yaml:"version"`
	Provider             string `yaml:"provider"`
	BootstrapProvider    string `yaml:"bootstrapProvider"`
	ControlPlaneProvider string `yaml:"controlPlaneProvider"`
}

// ToMap converts NutanixConfig to a map[string]string for provider factory.
func (n NutanixConfig) ToMap() map[string]string {
	return map[string]string{
		"endpoint":    n.Endpoint,
		"username":    n.Username,
		"password":    n.Password,
		"clusterUUID": n.ClusterUUID,
		"subnetUUID":  n.SubnetUUID,
	}
}
