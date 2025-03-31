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
	"butler/internal/adapters/platforms/kubectl"
	"butler/internal/models"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// KubeVipInitializer handles generating and deploying the Kube-Vip manifest.
type KubeVipInitializer struct {
	docker  *docker.DockerAdapter
	kubectl *kubectl.KubectlAdapter
	logger  *zap.Logger
}

// NewKubeVipInitializer creates a new KubeVip initializer.
func NewKubeVipInitializer(docker *docker.DockerAdapter, kubectl *kubectl.KubectlAdapter, logger *zap.Logger) *KubeVipInitializer {
	return &KubeVipInitializer{docker: docker, kubectl: kubectl, logger: logger}
}

// ConfigureKubeVip performs all steps to set up Kube-Vip: generating the manifest, applying RBAC, and deploying the DaemonSet.
func (k *KubeVipInitializer) ConfigureKubeVip(ctx context.Context, config *models.BootstrapConfig, server string) error {
	k.logger.Info("Starting Kube-Vip configuration")

	// Generate Kube-Vip manifest
	if err := k.GenerateManifest(ctx, config); err != nil {
		return fmt.Errorf("failed to generate Kube-Vip manifest: %w", err)
	}

	// Apply Kube-Vip RBAC
	if err := k.ApplyRBAC(ctx, server); err != nil {
		return fmt.Errorf("failed to apply Kube-Vip RBAC: %w", err)
	}

	// Apply Kube-Vip DaemonSet
	if err := k.ApplyDaemonSet(ctx, server); err != nil {
		return fmt.Errorf("failed to apply Kube-Vip DaemonSet: %w", err)
	}

	k.logger.Info("Kube-Vip setup completed successfully")
	return nil
}

// GenerateManifest creates the Kube-Vip DaemonSet manifest.
func (k *KubeVipInitializer) GenerateManifest(ctx context.Context, config *models.BootstrapConfig) error {
	k.logger.Info("Generating Kube-Vip manifest...")

	vip := config.ManagementCluster.Talos.ControlPlaneVIP
	interfaceName := "ens3" // TODO: Make dynamic
	version := config.ManagementCluster.Talos.Version
	outputDir := "talosconfig"
	outputFile := filepath.Join(outputDir, "kube-vip-ds.yaml")

	args := []string{
		"run", "--network", "host", "--rm",
		fmt.Sprintf("ghcr.io/kube-vip/kube-vip:%s", version),
		"manifest", "daemonset",
		"--interface", interfaceName,
		"--address", vip,
		"--inCluster",
		"--taint",
		"--controlplane",
		"--services",
		"--arp",
		"--leaderElection",
	}
	result, err := k.docker.ExecuteCommand(ctx, args...)
	if err != nil {
		return fmt.Errorf("error generating Kube-Vip manifest: %w", err)
	}

	if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
		return fmt.Errorf("failed to write Kube-Vip manifest to file: %w", err)
	}

	k.logger.Info("Kube-Vip manifest saved", zap.String("file", "talosconfig/kube-vip-ds.yaml"))
	return nil
}

// ApplyRBAC applies the Kube-Vip RBAC manifest.
func (k *KubeVipInitializer) ApplyRBAC(ctx context.Context, server string) error {
	rbacManifestURL := "https://kube-vip.io/manifests/rbac.yaml"

	k.logger.Info("Applying Kube-Vip RBAC configuration",
		zap.String("server", server),
		zap.String("manifest", rbacManifestURL),
	)

	_, err := k.kubectl.ExecuteCommand(ctx,
		"--server", server,
		"--kubeconfig", "talosconfig/kubeconfig",
		"apply", "-f", rbacManifestURL,
		"--insecure-skip-tls-verify=true",
	)
	if err != nil {
		return fmt.Errorf("failed to apply Kube-Vip RBAC: %w", err)
	}

	k.logger.Info("Kube-Vip RBAC successfully applied")
	return nil
}

// ApplyDaemonSet applies the Kube-Vip DaemonSet.
func (k *KubeVipInitializer) ApplyDaemonSet(ctx context.Context, server string) error {
	manifestPath := "talosconfig/kube-vip-ds.yaml"

	k.logger.Info("Applying Kube-Vip DaemonSet using kubectl",
		zap.String("server", server),
		zap.String("manifest", manifestPath),
	)

	_, err := k.kubectl.ExecuteCommand(ctx,
		"--server", server,
		"--kubeconfig", "talosconfig/kubeconfig",
		"apply", "-f", manifestPath,
		"--insecure-skip-tls-verify=true",
	)
	if err != nil {
		return fmt.Errorf("failed to apply Kube-Vip DaemonSet: %w", err)
	}

	k.logger.Info("Kube-Vip successfully applied")
	return nil
}
