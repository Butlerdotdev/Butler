// Package docker defines an adapter for Docker.
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

package docker

import (
	"butler/internal/adapters/exec"
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

// DockerClient runs Docker commands.
type DockerClient struct {
	execAdapter exec.ExecAdapter
	logger      *zap.Logger
}

// NewDockerClient initializes the client.
func NewDockerClient(execAdapter exec.ExecAdapter, logger *zap.Logger) *DockerClient {
	return &DockerClient{
		execAdapter: execAdapter,
		logger:      logger,
	}
}

// GenerateKubeVipManifest runs the Docker command to generate the Kube-Vip manifest.
func (c *DockerClient) GenerateKubeVipManifest(ctx context.Context, version, interfaceName, vip string) error {
	if version == "" {
		version = "v0.8.9"
	}

	outputDir := "talosconfig"
	outputFile := filepath.Join(outputDir, "kube-vip-ds.yaml")

	c.logger.Info("Generating Kube-Vip manifest",
		zap.String("version", version),
		zap.String("interface", interfaceName),
		zap.String("vip", vip),
		zap.String("output_file", outputFile),
	)
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

	// Run Docker command
	result, err := c.execAdapter.RunCommand(ctx, "docker", args...)
	if err != nil {
		c.logger.Error("Failed to generate Kube-Vip manifest",
			zap.Error(err),
			zap.String("version", version),
			zap.String("interface", interfaceName),
			zap.String("vip", vip),
		)
		return err
	}

	// Write manifest to file
	if err := os.WriteFile(outputFile, []byte(result.Stdout), 0644); err != nil {
		c.logger.Error("Failed to write Kube-Vip manifest to file", zap.String("file", outputFile), zap.Error(err))
		return err
	}

	return nil
}
