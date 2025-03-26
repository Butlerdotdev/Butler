// Package proxmox provides utility functions for parsing VM specifications.
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

import "fmt"

// parseRAM converts "8GB" to MiB
func parseRAM(ram string) int {
	var size int
	fmt.Sscanf(ram, "%dGB", &size)
	return size * 1024
}

// parseDisk converts "50GB" to just the value, excluding "GB"
func parseDisk(disk string) string {
	var size int
	fmt.Sscanf(disk, "%dGB", &size)
	return fmt.Sprintf("%d", size)
}
