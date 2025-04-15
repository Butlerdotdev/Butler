// Package bootstrap provides services for KubeConfig in Butler.
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
	"butler/pkg/adapters/platforms"
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// KubeConfigManager handles kubeconfig validation and Kubernetes API readiness checks.
type KubeConfigManager struct {
	logger          *zap.Logger
	platformAdapter platforms.PlatformAdapter
}

// NewKubeConfigManager initializes a KubeConfigManager.
func NewKubeConfigManager(logger *zap.Logger, platformAdapter platforms.PlatformAdapter) *KubeConfigManager {
	return &KubeConfigManager{
		logger:          logger,
		platformAdapter: platformAdapter,
	}
}

// ValidateKubeConfig ensures the kubeconfig is valid.
func (k *KubeConfigManager) ValidateKubeConfig(kubeconfigPath string) error {
	k.logger.Info("Validating kubeconfig", zap.String("path", kubeconfigPath))

	// Ensure the kubeconfig file exists
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return fmt.Errorf("kubeconfig file missing: %s", kubeconfigPath)
	}

	// Validate kubeconfig by running `kubectl config view`
	_, err := k.platformAdapter.ExecuteCommand(context.Background(), "--kubeconfig", kubeconfigPath, "config", "view")
	if err != nil {
		return fmt.Errorf("kubectl config view failed, kubeconfig may be invalid: %w", err)
	}

	k.logger.Info("Kubeconfig is valid", zap.String("path", kubeconfigPath))
	return nil
}

// WaitForKubernetesAPI ensures that the Kubernetes API is accessible.
func (k *KubeConfigManager) WaitForKubernetesAPI(kubeconfigPath, controlPlaneNode string, timeout time.Duration) error {
	k.logger.Info("Waiting 60 seconds for Kubernetes API to initialize before checking readiness")
	time.Sleep(60 * time.Second)

	k.logger.Info("Waiting for Kubernetes API to be ready", zap.String("controlPlaneNode", controlPlaneNode))
	server := fmt.Sprintf("https://%s:6443", controlPlaneNode)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		_, err := k.platformAdapter.ExecuteCommand(context.Background(),
			"--server", server, "--kubeconfig", kubeconfigPath, "get", "nodes", "--request-timeout=15s",
		)
		if err == nil {
			k.logger.Info("Kubernetes API is ready", zap.String("server", server))
			return nil
		}

		k.logger.Warn("Kubernetes API not yet ready, retrying...", zap.String("server", server), zap.Error(err))

		// 🔄 Attempt to refresh kubeconfig if API is still unreachable after multiple tries
		if time.Now().Add(-30 * time.Second).After(deadline) {
			k.logger.Warn("Kubernetes API still unavailable, attempting to refresh kubeconfig")

			_, err = k.platformAdapter.ExecuteCommand(context.Background(),
				"--server", server, "--kubeconfig", kubeconfigPath, "config", "view",
			)
			if err != nil {
				k.logger.Warn("Kubeconfig may be invalid, attempting to refresh it...")
				_, _ = k.platformAdapter.ExecuteCommand(context.Background(),
					"--server", server, "--kubeconfig", kubeconfigPath, "config", "unset", "current-context",
				)
				_, _ = k.platformAdapter.ExecuteCommand(context.Background(),
					"--server", server, "--kubeconfig", kubeconfigPath, "config", "use-context", "admin@butler-r",
				)
			}
		}

		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("timed out waiting for Kubernetes API on node %s", controlPlaneNode)
}

// EnsureCorrectContext ensures the kubeconfig context is set correctly.
func (k *KubeConfigManager) EnsureCorrectContext(kubeconfigPath string, clusterName string) error {
	k.logger.Info("Ensuring kubeconfig context is set", zap.String("path", kubeconfigPath))

	contextName := fmt.Sprintf("admin@%s", clusterName)

	// Force reload of the kubeconfig by unsetting current context
	_, err := k.platformAdapter.ExecuteCommand(context.Background(), "--kubeconfig", kubeconfigPath, "config", "unset", "current-context")
	if err != nil {
		k.logger.Warn("Failed to unset current-context", zap.Error(err))
	}

	// Switch to the correct context
	_, err = k.platformAdapter.ExecuteCommand(context.Background(), "--kubeconfig", kubeconfigPath, "config", "use-context", contextName)
	if err != nil {
		return fmt.Errorf("failed to switch kubeconfig context: %w", err)
	}

	// Verify the active context
	result, err := k.platformAdapter.ExecuteCommand(context.Background(), "--kubeconfig", kubeconfigPath, "config", "current-context")
	if err != nil {
		return fmt.Errorf("failed to get current context: %w", err)
	}

	k.logger.Info("Kubeconfig context set", zap.String("context", result))
	return nil
}
