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
	ManagementCluster ManagementClusterConfig `mapstructure:"managementCluster" yaml:"managementCluster"`
}

// ManagementClusterConfig holds the cluster configuration.
type ManagementClusterConfig struct {
	Name       string        `mapstructure:"name" yaml:"name"`
	Provider   string        `mapstructure:"provider" yaml:"provider"`
	Nutanix    NutanixConfig `mapstructure:"nutanix" yaml:"nutanix"`
	Proxmox    ProxmoxConfig `mapstructure:"proxmox" yaml:"proxmox"`
	Nodes      []NodeConfig  `mapstructure:"nodes" yaml:"nodes"`
	Talos      TalosConfig   `mapstructure:"talos" yaml:"talos"`
	ClusterAPI ClusterAPI    `mapstructure:"clusterAPI" yaml:"clusterAPI"`
	Flux       FluxConfig    `mapstructure:"flux" yaml:"flux"`
}

// FluxConfig holds Flux GitOps settings.
type FluxConfig struct {
	GitOwner      string `mapstructure:"gitOwner" yaml:"gitOwner"`
	GitRepository string `mapstructure:"gitRepository" yaml:"gitRepository"`
	GitBranch     string `mapstructure:"gitBranch" yaml:"gitBranch"`
	GitPath       string `mapstructure:"gitPath" yaml:"gitPath"`
	GitHostname   string `mapstructure:"gitHostname" yaml:"gitHostname"`
	GitPAT        string `mapstructure:"gitPAT" yaml:"gitPAT"`
}

// NutanixConfig defines the Nutanix API connection and cluster details.
type NutanixConfig struct {
	Endpoint    string `mapstructure:"endpoint" yaml:"endpoint"`
	Username    string `mapstructure:"username" yaml:"username"`
	Password    string `mapstructure:"password" yaml:"password"`
	ClusterUUID string `mapstructure:"clusterUUID" yaml:"clusterUUID"`
	SubnetUUID  string `mapstructure:"subnetUUID" yaml:"subnetUUID"`
}

// ProxmoxConfig defines the Proxmox API connection and cluster details.
type ProxmoxConfig struct {
	Endpoint           string   `mapstructure:"endpoint" yaml:"endpoint"`
	Username           string   `mapstructure:"username" yaml:"username"`
	Password           string   `mapstructure:"password" yaml:"password"`
	StorageLocation    string   `mapstructure:"storageLocation" yaml:"storageLocation"`
	AvailableVMIdStart int      `mapstructure:"availableVMIdStart" yaml:"availableVMIdStart"`
	AvailableVMIdEnd   int      `mapstructure:"availableVMIdEnd" yaml:"availableVMIdEnd"`
	Nodes              []string `mapstructure:"nodes" yaml:"nodes"`
}

// NodeConfig represents a single VM configuration.
type NodeConfig struct {
	Role       string   `mapstructure:"role" yaml:"role"`
	Count      int      `mapstructure:"count" yaml:"count"`
	CPU        int      `mapstructure:"cpu" yaml:"cpu"`
	RAM        string   `mapstructure:"ram" yaml:"ram"`
	Disk       string   `mapstructure:"disk" yaml:"disk"`
	IsoUUID    string   `mapstructure:"isoUUID" yaml:"isoUUID"`
	ExtraDisks []string `mapstructure:"extraDisks" yaml:"extraDisks"`
}

// TalosConfig holds Talos Linux bootstrapping details.
type TalosConfig struct {
	Version              string   `mapstructure:"version" yaml:"version"`
	ControlPlaneEndpoint string   `mapstructure:"controlPlaneEndpoint" yaml:"controlPlaneEndpoint"`
	ControlPlaneVIP      string   `mapstructure:"controlPlaneVIP" yaml:"controlPlaneVIP"`
	BoundNodeIP          string   `mapstructure:"boundNodeIP" yaml:"boundNodeIP"`
	ClusterName          string   `mapstructure:"clusterName" yaml:"clusterName"`
	CIDR                 string   `mapstructure:"cidr" yaml:"cidr"`
	Gateway              string   `mapstructure:"gateway" yaml:"gateway"`
	OutputDir            string   `mapstructure:"outputDir" yaml:"outputDir"`
	ControlPlaneNodes    []string `mapstructure:"controlPlaneNodes" yaml:"controlPlaneNodes"`
	WorkerNodes          []string `mapstructure:"workerNodes" yaml:"workerNodes"`
}

// ClusterAPI represents the Cluster API provider settings.
type ClusterAPI struct {
	Version              string `mapstructure:"version" yaml:"version"`
	Provider             string `mapstructure:"provider" yaml:"provider"`
	BootstrapProvider    string `mapstructure:"bootstrapProvider" yaml:"bootstrapProvider"`
	ControlPlaneProvider string `mapstructure:"controlPlaneProvider" yaml:"controlPlaneProvider"`
}
