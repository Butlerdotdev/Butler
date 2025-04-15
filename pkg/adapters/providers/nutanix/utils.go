// Package nutanix provides utility functions for parsing VM specifications.
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
	sharedModels "butler/internal/models"
	"butler/pkg/adapters/providers/nutanix/models"
	"fmt"
)

// parseRAM converts "8GB" to MiB
func parseRAM(ram string) int {
	var size int
	fmt.Sscanf(ram, "%dGB", &size)
	return size * 1024
}

// parseDisk converts "50GB" to MiB
func parseDisk(disk string) int {
	var size int
	fmt.Sscanf(disk, "%dGB", &size)
	return size * 1024
}

func buildDiskList(vm sharedModels.VMConfig) []models.Disk {
	disks := []models.Disk{
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
	}

	for i, extra := range vm.ExtraDisks {
		disks = append(disks, models.Disk{
			DeviceProperties: models.DeviceProperties{
				DeviceType: "DISK",
				DiskAddress: models.DiskAddress{
					AdapterType: "SCSI",
					DeviceIndex: i + 1,
				},
			},
			DiskSizeMiB: parseDisk(extra),
		})
	}
	return disks
}
