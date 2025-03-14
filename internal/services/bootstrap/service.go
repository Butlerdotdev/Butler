// Package bootstrap provides services for provisioning the Butler management cluster.
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
	"butler/internal/adapters/exec"
	"butler/internal/adapters/platforms"
	"butler/internal/adapters/platforms/docker"
	"butler/internal/adapters/platforms/kubectl"
	"butler/internal/adapters/platforms/talos"
	"butler/internal/adapters/providers"
	"butler/internal/models"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// BootstrapService orchestrates provisioning the management cluster.
type BootstrapService struct {
	logger            *zap.Logger
	provider          providers.ProviderInterface
	provisioner       *Provisioner
	healthCheck       *HealthChecker
	talosInit         *TalosInitializer
	kubeVipInit       *KubeVipInitializer
	kubectl           *kubectl.KubectlAdapter
	kubeConfigManager *KubeConfigManager
}

// NewBootstrapService initializes a new BootstrapService with dependencies.
func NewBootstrapService(ctx context.Context, config *models.BootstrapConfig, logger *zap.Logger) (*BootstrapService, error) {
	provider, err := providers.NewProviderFactory(ctx, config.ManagementCluster.Provider, config.ManagementCluster.Nutanix.ToMap(), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	// Initialize Exec Adapter
	execAdapter := exec.NewClient(logger)

	// Initialize Docker Adapter
	dockerAdapter, err := platforms.GetPlatformAdapter("docker", execAdapter, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Docker adapter: %w", err)
	}
	dockerConcrete, ok := dockerAdapter.(*docker.DockerAdapter)
	if !ok {
		return nil, fmt.Errorf("failed to assert DockerAdapter type")
	}

	// Initialize Talos adapter
	talosAdapter, err := platforms.GetPlatformAdapter("talos", execAdapter, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Talos adapter: %w", err)
	}

	// Type assertion to *talos.TalosAdapter
	talosConcrete, ok := talosAdapter.(*talos.TalosAdapter)
	if !ok {
		return nil, fmt.Errorf("failed to assert TalosAdapter type")
	}

	kubectlAdapter, err := platforms.GetPlatformAdapter("kubectl", execAdapter, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubectl adapter: %w", err)
	}

	kubectlConcrete, ok := kubectlAdapter.(*kubectl.KubectlAdapter)
	if !ok {
		return nil, fmt.Errorf("failed to assert KubectlAdapter type")
	}

	// Pass the interface
	kubeConfigManager := NewKubeConfigManager(logger, kubectlAdapter)

	return &BootstrapService{
		logger:            logger,
		provider:          provider,
		provisioner:       NewProvisioner(provider, logger),
		healthCheck:       NewHealthChecker(provider, logger),
		talosInit:         NewTalosInitializer(talosConcrete, logger),
		kubeVipInit:       NewKubeVipInitializer(dockerConcrete, kubectlConcrete, logger),
		kubectl:           kubectlConcrete,
		kubeConfigManager: kubeConfigManager,
	}, nil
}

// ProvisionManagementCluster provisions the management cluster VMs and waits for them to be ready.
func (b *BootstrapService) ProvisionManagementCluster(config *models.BootstrapConfig) error {
	b.logger.Info("Starting provisioning of management cluster", zap.String("cluster_name", config.ManagementCluster.Name))

	// Provision VMs
	err := b.provisioner.ProvisionVMs(config)
	if err != nil {
		return err
	}

	// Wait for health checks & collect IPs
	nodeIPs, err := b.healthCheck.WaitForVMsToBeReady(config, 10*time.Minute)
	if err != nil {
		return err
	}

	// Separate Control Plane & Worker Nodes
	controlPlanes, workers, err := b.provisioner.SeparateNodesByRole(config, nodeIPs)
	if err != nil {
		return err
	}

	// Generate & Apply Talos Configuration
	talosConfig := models.TalosConfig{
		ClusterName:          config.ManagementCluster.Name,
		ControlPlaneEndpoint: config.ManagementCluster.Talos.ControlPlaneEndpoint,
		OutputDir:            "./talosconfig",
		ControlPlaneNodes:    controlPlanes,
		WorkerNodes:          workers,
	}
	insecure := true

	err = b.talosInit.ConfigureTalos(context.Background(), &talosConfig, insecure)
	if err != nil {
		return fmt.Errorf("failed to configure Talos: %w", err)
	}

	// Ensure At Least One Control Plane Node Exists
	if len(controlPlanes) == 0 {
		return fmt.Errorf("no available control plane nodes for Kube-Vip setup")
	}
	server := fmt.Sprintf("https://%s:6443", controlPlanes[0])

	// Validate KubeConfig Before Proceeding
	kubeconfigPath := "talosconfig/kubeconfig"
	err = b.kubeConfigManager.ValidateKubeConfig(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("kubeconfig validation failed: %w", err)
	}

	// Ensure kubeconfig context is correct
	err = b.kubeConfigManager.EnsureCorrectContext(kubeconfigPath, config.ManagementCluster.Name)
	if err != nil {
		return fmt.Errorf("failed to ensure correct kubeconfig context: %w", err)
	}

	// Ensure Kubernetes API is Ready Before Applying Kube-Vip
	err = b.kubeConfigManager.WaitForKubernetesAPI(kubeconfigPath, controlPlanes[0], 5*time.Minute)
	if err != nil {
		return fmt.Errorf("kubernetes API did not become available: %w", err)
	}

	// Apply Kube-Vip RBAC
	rbacManifestURL := "https://kube-vip.io/manifests/rbac.yaml"
	b.logger.Info("Applying Kube-Vip RBAC configuration",
		zap.String("server", server),
		zap.String("manifest", rbacManifestURL),
	)

	// Deploy Kube-Vip DaemonSet
	err = b.kubeVipInit.ConfigureKubeVip(context.Background(), config, server)
	if err != nil {
		return fmt.Errorf("failed to configure Kube-Vip: %w", err)
	}

	b.logger.Info("Management cluster provisioned successfully with Talos!")
	return nil
}
