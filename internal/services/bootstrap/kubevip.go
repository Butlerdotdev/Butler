// Package bootstrap provides services for configuring Kube-Vip in Butler.
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

package bootstrap

import (
	"butler/internal/adapters/platforms/docker"
	"butler/internal/models"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// KubeVipInitializer handles generating and deploying the Kube-Vip manifest.
type KubeVipInitializer struct {
	docker docker.DockerAdapter
	logger *zap.Logger
}

// NewKubeVipInitializer creates a new KubeVip initializer.
func NewKubeVipInitializer(docker docker.DockerAdapter, logger *zap.Logger) *KubeVipInitializer {
	return &KubeVipInitializer{docker: docker, logger: logger}
}

// ConfigureKubeVip generates and logs the Kube-Vip manifest.
func (k *KubeVipInitializer) ConfigureKubeVip(ctx context.Context, config *models.BootstrapConfig) error {
	k.logger.Info("Generating Kube-Vip manifest...")

	// Use ControlPlaneEndpoint as the VIP
	vip := config.ManagementCluster.Talos.ControlPlaneEndpoint
	// TODO: Make interface assignment dynamic (This will definitely be a thing with the more providers we support)
	interfaceName := "ens3"

	err := k.docker.GenerateKubeVipManifest(ctx, config.ManagementCluster.Talos.Version, interfaceName, vip)
	if err != nil {
		fmt.Println("Error generating Kube-Vip manifest:", err)
	} else {
		fmt.Println("Kube-Vip manifest saved to kube-vip-ds.yaml")
	}

	// Save or apply the manifest (future: pass to Kubectl adapter)
	return nil
}
