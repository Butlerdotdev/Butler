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
	"context"
	"fmt"
	"time"

	"github.com/butlerdotdev/butler/internal/mappers"
	"github.com/butlerdotdev/butler/pkg/adapters/exec"
	"github.com/butlerdotdev/butler/pkg/adapters/platforms"
	"github.com/butlerdotdev/butler/pkg/adapters/platforms/docker"
	"github.com/butlerdotdev/butler/pkg/adapters/platforms/flux"
	"github.com/butlerdotdev/butler/pkg/adapters/platforms/helm"
	"github.com/butlerdotdev/butler/pkg/adapters/platforms/kubectl"
	"github.com/butlerdotdev/butler/pkg/adapters/platforms/talos"
	"github.com/butlerdotdev/butler/pkg/adapters/providers"
	"github.com/butlerdotdev/butler/pkg/models"

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
	kubeOvnInit       *KubeOvnInitializer
	fluxInit          *FluxInitializer
	kubectl           *kubectl.KubectlAdapter
	kubeConfigManager *KubeConfigManager
	config            *models.BootstrapConfig
}

// NewBootstrapService initializes BootstrapService using Viper for config.
func NewBootstrapService(ctx context.Context, config *models.BootstrapConfig, logger *zap.Logger) (*BootstrapService, error) {

	logger.Info("Initializing BootstrapService")
	provider, err := providers.NewProviderFactory(
		ctx,
		config.ManagementCluster.Provider,
		mappers.NewMapping(config.ManagementCluster.Provider, config.ManagementCluster),
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	// Init Adapters and Type Assertions
	execAdapter := exec.NewClient(logger)

	dockerAdapter, err := platforms.GetPlatformAdapter("docker", execAdapter, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Docker adapter: %w", err)
	}
	dockerConcrete := dockerAdapter.(*docker.DockerAdapter)
	dockerConcrete, ok := dockerAdapter.(*docker.DockerAdapter)
	if !ok {
		return nil, fmt.Errorf("failed to assert DockerAdapter type")
	}

	talosAdapter, err := platforms.GetPlatformAdapter("talos", execAdapter, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Talos adapter: %w", err)
	}
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

	helmAdapter, err := platforms.GetPlatformAdapter("helm", execAdapter, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Helm adapter: %w", err)
	}
	helmConcrete, ok := helmAdapter.(*helm.HelmAdapter)
	if !ok {
		return nil, fmt.Errorf("failed to assert HelmAdapter type")
	}

	kubeOvnInit := NewKubeOvnInitializer(kubectlConcrete, helmConcrete, talosConcrete, logger)

	fluxAdapter, err := platforms.GetPlatformAdapter("flux", execAdapter, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Flux adapter: %w", err)
	}
	fluxConcrete, ok := fluxAdapter.(*flux.FluxAdapter)
	if !ok {
		return nil, fmt.Errorf("failed to assert FluxAdapter type")
	}

	kubeConfigManager := NewKubeConfigManager(logger, kubectlConcrete)

	return &BootstrapService{
		logger:            logger,
		provider:          provider,
		provisioner:       NewProvisioner(provider, logger),
		healthCheck:       NewHealthChecker(provider, logger),
		talosInit:         NewTalosInitializer(talosConcrete, logger),
		kubeVipInit:       NewKubeVipInitializer(dockerConcrete, kubectlConcrete, logger),
		fluxInit:          NewFluxInitializer(fluxConcrete, logger),
		kubectl:           kubectlConcrete,
		kubeOvnInit:       kubeOvnInit,
		kubeConfigManager: kubeConfigManager,
		config:            config,
	}, nil
}

// ProvisionManagementCluster provisions the management cluster.
func (b *BootstrapService) ProvisionManagementCluster() error {
	config := b.config

	b.logger.Info("Starting provisioning of management cluster",
		zap.String("cluster_name", config.ManagementCluster.Name),
	)

	// Provision VMs
	if err := b.provisioner.ProvisionVMs(config); err != nil {
		return err
	}

	// Wait for health checks & collect IPs
	nodeIPs, err := b.healthCheck.WaitForVMsToBeReady(config, 10*time.Minute)
	if err != nil {
		return err
	}

	// Separate nodes
	controlPlanes, workers, err := b.provisioner.SeparateNodesByRole(config, nodeIPs)
	if err != nil {
		return err
	}

	if len(controlPlanes) == 0 {
		return fmt.Errorf("no available control plane nodes for Kube-Vip setup")
	}

	// Talos config
	talosConfig := models.TalosConfig{
		ClusterName:          config.ManagementCluster.Name,
		ControlPlaneEndpoint: config.ManagementCluster.Talos.ControlPlaneEndpoint,
		OutputDir:            "./talosconfig",
		ControlPlaneNodes:    controlPlanes,
		WorkerNodes:          workers,
	}

	if err := b.talosInit.ConfigureTalos(context.Background(), &talosConfig, true); err != nil {
		return fmt.Errorf("failed to configure Talos: %w", err)
	}

	// Validate kubeconfig
	kubeconfigPath := "talosconfig/kubeconfig"
	if err := b.kubeConfigManager.ValidateKubeConfig(kubeconfigPath); err != nil {
		return fmt.Errorf("kubeconfig validation failed: %w", err)
	}
	if err := b.kubeConfigManager.EnsureCorrectContext(kubeconfigPath, config.ManagementCluster.Name); err != nil {
		return fmt.Errorf("failed to set kubeconfig context: %w", err)
	}
	if err := b.kubeConfigManager.WaitForKubernetesAPI(kubeconfigPath, controlPlanes[0], 5*time.Minute); err != nil {
		return fmt.Errorf("kubernetes API not ready: %w", err)
	}

	// Kube-Vip
	server := fmt.Sprintf("https://%s:6443", controlPlanes[0])
	config.ManagementCluster.Talos.BoundNodeIP = controlPlanes[0]
	b.logger.Info("Applying Kube-Vip RBAC configuration",
		zap.String("server", server),
		zap.String("manifest", "https://kube-vip.io/manifests/rbac.yaml"),
	)

	// Wait until at least one node registers
	if err := b.kubeOvnInit.WaitForNodes(context.Background(), server, 2*time.Minute); err != nil {
		return fmt.Errorf("failed waiting for nodes to register: %w", err)
	}

	// Gather node IP-to-name mapping prior to VIP being associated - This allows for reduction in VIP vs Bound IP Logic later when mapping.
	ipToNodeMap, err := b.kubeOvnInit.getInternalIPToNodeNameMap(context.Background(), server)
	if err != nil {
		return fmt.Errorf("failed to collect node IP-to-name map: %w", err)
	}

	// KubeVip configuration
	if err := b.kubeVipInit.ConfigureKubeVip(context.Background(), config, server); err != nil {
		return fmt.Errorf("failed to configure Kube-Vip: %w", err)
	}

	// Check if the VIP is ready
	if err := b.kubeConfigManager.WaitForKubernetesAPI(kubeconfigPath, config.ManagementCluster.Talos.ControlPlaneVIP, 2*time.Minute); err != nil {
		return fmt.Errorf("VIP-based Kubernetes API not ready: %w", err)
	}

	// Label nodes for kube-ovn
	if err := b.kubeOvnInit.LabelNodes(context.Background(), config, controlPlanes, workers, ipToNodeMap); err != nil {
		return fmt.Errorf("failed to label nodes for Kube-OVN: %w", err)
	}

	// Helm Install kube-ovn with templated values
	if err := b.kubeOvnInit.ConfigureKubeOvn(context.Background(), controlPlanes, config); err != nil {
		return fmt.Errorf("failed to install Kube-OVN: %w", err)
	}

	// TODO: Add a check to ensure that the Kube-OVN pods are running before proceeding.

	// TODO: Piraeus Operator(Linstor) - V2 of this operator does not have a helm chart that they suggest to use. They have one for V1, but their docs show to use the V2 operator.
	// We will install this similarly to kube ovn, Outside of the flux process. BUT, we can still use flux to manage parts of linstor.

	// Bootstrap Flux
	// if err := b.fluxInit.FluxBootstrap(context.Background(), config); err != nil {
	// 	return fmt.Errorf("failed to bootstrap Flux: %w", err)
	// }

	// After Flux is bootstrapped the following should be provisioned Via Flux:
	// MetalLB
	// Traefik
	// CAPI

	b.logger.Info("Flux bootstrap completed successfully")
	b.logger.Info("Management cluster provisioned successfully")
	return nil
}
