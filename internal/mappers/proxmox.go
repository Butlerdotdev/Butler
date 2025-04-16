// Package mappers provides helper functions to transform and map internal
// configuration models into the formats required by Butler's platform adapters,
// providers, and other services. These pure functions are responsible for
// struct-to-map conversions, type reshaping, and general data preparation
// without side effects.
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

package mappers

import (
	"fmt"
	"strings"

	"github.com/butlerdotdev/butler/internal/models"
)

func ProxmoxToMap(cfg models.ProxmoxConfig) map[string]string {
	return map[string]string{
		"endpoint":           cfg.Endpoint,
		"username":           cfg.Username,
		"password":           cfg.Password,
		"storageLocation":    cfg.StorageLocation,
		"availableVMIdStart": fmt.Sprintf("%d", cfg.AvailableVMIdStart),
		"availableVMIdEnd":   fmt.Sprintf("%d", cfg.AvailableVMIdEnd),
		"nodes":              strings.Join(cfg.Nodes, ","),
	}
}
