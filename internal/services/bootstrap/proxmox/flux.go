// Package bootstrap provides services for bootstrapping Flux in Butler.
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
package bootstrap

import (
	"butler/internal/models"
	"butler/pkg/adapters/platforms/flux"
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// FluxInitializer handles bootstrapping FluxCD into the management cluster.
type FluxInitializer struct {
	fluxAdapter *flux.FluxAdapter
	logger      *zap.Logger
}

// NewFluxInitializer creates a new Flux initializer.
func NewFluxInitializer(fluxAdapter *flux.FluxAdapter, logger *zap.Logger) *FluxInitializer {
	return &FluxInitializer{
		fluxAdapter: fluxAdapter,
		logger:      logger,
	}
}

// FluxBootstrap initializes FluxCD on the management cluster with retries.
func (f *FluxInitializer) FluxBootstrap(ctx context.Context, config *models.BootstrapConfig) error {
	f.logger.Info("Starting Flux bootstrap for the management cluster")

	// Get values from config
	clusterName := config.ManagementCluster.Name
	gitOwner := config.ManagementCluster.Flux.GitOwner
	gitRepo := config.ManagementCluster.Flux.GitRepository
	gitBranch := config.ManagementCluster.Flux.GitBranch
	gitPath := config.ManagementCluster.Flux.GitPath
	gitHostname := config.ManagementCluster.Flux.GitHostname

	// Get GitLab token from environment
	gitToken := os.Getenv("GITLAB_TOKEN")

	// If token is missing, prompt user
	if gitToken == "" {
		fmt.Print("Enter your GitLab token (will not be stored): ")
		fmt.Scanln(&gitToken) // Read input (not secure for automation, but fine for CLI users)
	}

	// Ensure GitLab token is set in the environment
	if err := os.Setenv("GITLAB_TOKEN", gitToken); err != nil {
		f.logger.Error("Failed to set GitLab token environment variable",
			zap.String("clusterName", clusterName),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set GitLab token environment variable: %w", err)
	}

	maxRetries := 3
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		f.logger.Info("Executing Flux bootstrap",
			zap.String("clusterName", clusterName),
			zap.Int("attempt", attempt),
		)

		_, err = f.fluxAdapter.ExecuteCommand(ctx,
			"bootstrap", "gitlab",
			"--owner", gitOwner,
			"--repository", gitRepo,
			"--branch", gitBranch,
			"--path", gitPath,
			"--token-auth=true",
			"--hostname", gitHostname,
			"--read-write-key=true",
			"--components-extra", "image-reflector-controller,image-automation-controller",
			"--insecure-skip-tls-verify=true",
			"--kubeconfig", "talosconfig/kubeconfig",
		)
		if err == nil {
			f.logger.Info("Flux bootstrap completed successfully",
				zap.String("clusterName", clusterName),
			)
			return nil
		}

		f.logger.Warn("Flux bootstrap failed, retrying...",
			zap.String("clusterName", clusterName),
			zap.Int("attempt", attempt),
			zap.Error(err),
		)

		// Exponential backoff
		backoff := attempt * 10
		f.logger.Info("Waiting before retrying...",
			zap.Int("seconds", backoff),
		)
		time.Sleep(time.Duration(backoff) * time.Second)
	}

	// If all retries fail, return an error
	f.logger.Error("Flux bootstrap failed after multiple attempts",
		zap.String("clusterName", clusterName),
		zap.Error(err),
	)
	return fmt.Errorf("failed to bootstrap Flux after %d attempts: %w", maxRetries, err)
}
