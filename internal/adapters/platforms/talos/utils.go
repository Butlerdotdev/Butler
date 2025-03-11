// Package talos defines an adapter for Talos and bootstrapping the OS.
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

package talos

import (
	"butler/internal/models"
	"strings"
)

// parseTalosHealth converts `talosctl health` output into a structured PlatformStatus.
func parseTalosHealth(output string) models.PlatformStatus {
	healthy := strings.Contains(output, "health: OK")
	version := "unknown"

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Talos version:") {
			version = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return models.PlatformStatus{
		Healthy: healthy,
		Details: map[string]string{"version": version},
	}
}
