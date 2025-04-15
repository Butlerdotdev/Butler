// Package platforms defines generic interface contract and factory.
//
// # Copyright (c) 2025, The Butler Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package platforms

import (
	"butler/pkg/adapters/exec"
	"butler/pkg/adapters/platforms/docker"
	"butler/pkg/adapters/platforms/flux"
	"butler/pkg/adapters/platforms/helm"
	"butler/pkg/adapters/platforms/kubectl"
	"butler/pkg/adapters/platforms/talos"
	"fmt"

	"go.uber.org/zap"
)

// GetPlatformAdapter returns the correct platform adapter.
func GetPlatformAdapter(name string, execAdapter exec.ExecAdapter, logger *zap.Logger) (PlatformAdapter, error) {
	switch name {
	case "talos":
		return talos.NewTalosAdapter(execAdapter, logger), nil
	case "kubectl":
		return kubectl.NewKubectlAdapter(execAdapter, logger), nil
	case "docker":
		return docker.NewDockerAdapter(execAdapter, logger), nil
	case "flux":
		return flux.NewFluxAdapter(execAdapter, logger), nil
	case "helm":
		return helm.NewHelmAdapter(execAdapter, logger), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", name)
	}
}
