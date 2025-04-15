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

// BootstrapConfig defines the structure for provisioning the Nutanix management cluster
type BootstrapConfig struct {
	ManagementCluster struct {
		Name string `yaml:"name"`

		Nutanix struct {
			Endpoint    string `yaml:"endpoint"`
			Username    string `yaml:"username"`
			Password    string `yaml:"password"`
			ClusterUUID string `yaml:"clusterUUID"`
			SubnetUUID  string `yaml:"subnetUUID"`
		} `yaml:"nutanix"`

		Nodes []struct {
			Role       string   `yaml:"role"`
			CPU        int      `yaml:"cpu"`
			RAM        string   `yaml:"ram"`
			Disk       string   `yaml:"disk"`
			ExtraDisks []string `yaml:"extraDisks"`
			Count      int      `yaml:"count"`
			IsoUUID    string   `yaml:"isoUUID"`
		} `yaml:"nodes"`

		Talos struct {
			Version              string `yaml:"version"`
			ControlPlaneEndpoint string `yaml:"controlPlaneEndpoint"`
			ClusterName          string `yaml:"clusterName"`
			CIDR                 string `yaml:"cidr"`
			Gateway              string `yaml:"gateway"`
		} `yaml:"talos"`

		ClusterAPI struct {
			Version              string `yaml:"version"`
			Provider             string `yaml:"provider"`
			BootstrapProvider    string `yaml:"bootstrapProvider"`
			ControlPlaneProvider string `yaml:"controlPlaneProvider"`
		} `yaml:"clusterAPI"`
	} `yaml:"managementCluster"`
}

// NutanixVMConfig represents the Nutanix API request format for creating a VM.
type NutanixVMConfig struct {
	Metadata Metadata `json:"metadata"`
	// Name     string   `json:"name"`
	Spec Spec `json:"spec"`
}

// Metadata defines the VM kind.
type Metadata struct {
	Kind string `json:"kind"`
}

// Spec represents the VM configuration.
type Spec struct {
	Name             string           `json:"name"`
	Resources        Resources        `json:"resources"`
	ClusterReference ClusterReference `json:"cluster_reference"`
}

// Resources defines the VM hardware settings.
type Resources struct {
	PowerState        string     `json:"power_state"`
	NumSockets        int        `json:"num_sockets"`
	NumVCPUsPerSocket int        `json:"num_vcpus_per_socket"`
	MemorySizeMiB     int        `json:"memory_size_mib"`
	BootConfig        BootConfig `json:"boot_config"`
	DiskList          []Disk     `json:"disk_list"`
	NicList           []Nic      `json:"nic_list"`
}

// BootConfig defines boot order.
type BootConfig struct {
	BootDeviceOrderList []string `json:"boot_device_order_list"`
}

// Disk defines VM disk settings.
type Disk struct {
	DeviceProperties    DeviceProperties     `json:"device_properties"`
	DiskSizeMiB         int                  `json:"disk_size_mib,omitempty"`
	DataSourceReference *DataSourceReference `json:"data_source_reference,omitempty"`
}

// DeviceProperties defines a disk device.
type DeviceProperties struct {
	DeviceType  string      `json:"device_type"`
	DiskAddress DiskAddress `json:"disk_address"`
}

// DiskAddress defines the disk attachment type.
type DiskAddress struct {
	AdapterType string `json:"adapter_type"`
	DeviceIndex int    `json:"device_index"`
}

// DataSourceReference references an ISO image for booting.
type DataSourceReference struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
}

// Nic defines network settings.
type Nic struct {
	SubnetReference SubnetReference `json:"subnet_reference"`
}

// SubnetReference references the subnet.
type SubnetReference struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
}

// ClusterReference references the Nutanix cluster.
type ClusterReference struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
}

type NutanixClusterList struct {
	Entities []NutanixClusterEntities `json:"entities"`
}

type NutanixClusterEntities struct {
	Metadata NutanixClusterMetadata `json:"metadata"`
	Spec     NutanixClusterSpec     `json:"spec"`
}

type NutanixClusterSpec struct {
	Name string `json:"name"`
}

type NutanixClusterMetadata struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type NutanixSubnetList struct {
	Entities []NutanixSubnetEntities `json:"entities"`
}

type NutanixSubnetEntities struct {
	Metadata NutanixSubnetMetadata `json:"metadata"`
	Spec     NutanixSubnetSpec     `json:"spec"`
}

type NutanixSubnetSpec struct {
	Name             string                  `json:"name"`
	ClusterReference NutanixClusterReference `json:"cluster_reference"`
}

type NutanixClusterReference struct {
	UUID string `json:"uuid"`
}

type NutanixSubnetMetadata struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
