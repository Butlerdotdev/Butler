// Package bootstrap provides functionality for initializing and configuring Kube-OVN in Butler.
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
	"bufio"
	"butler/internal/adapters/platforms/helm"
	"butler/internal/adapters/platforms/kubectl"
	"butler/internal/models"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"go.uber.org/zap"
)

// Not sold on this approach of embedding the values.yaml file like this. Will need to re-evaluate this later.
//
//go:embed assets/kubeovn/values.yaml
var baseKubeOvnValues string

// KubeOvnInitializer provides functionality to bootstrap Kube-OVN,
// including Helm install and node labeling for control and worker planes.
type KubeOvnInitializer struct {
	kubectl *kubectl.KubectlAdapter
	helm    *helm.HelmAdapter
	logger  *zap.Logger
}

// NewKubeOvnInitializer constructs a new KubeOvnInitializer instance.
func NewKubeOvnInitializer(kubectl *kubectl.KubectlAdapter, helm *helm.HelmAdapter, logger *zap.Logger) *KubeOvnInitializer {
	return &KubeOvnInitializer{
		kubectl: kubectl,
		helm:    helm,
		logger:  logger,
	}
}

// LabelNodes applies Kube-OVN-specific labels to control plane and worker nodes.
// Control plane nodes receive 'kube-ovn/role=master'; workers get 'node-role.kubernetes.io/worker='.
func (k *KubeOvnInitializer) LabelNodes(
	ctx context.Context,
	config *models.BootstrapConfig,
	controlPlaneIPs, workerIPs []string,
	ipToNodeNameMap map[string]string,
) error {
	kubeconfig := "talosconfig/kubeconfig"
	server := fmt.Sprintf("https://%s:6443", config.ManagementCluster.Talos.ControlPlaneVIP)

	controlPlaneNames := resolveNodeNames(k.logger, ipToNodeNameMap, controlPlaneIPs)
	workerNames := resolveNodeNames(k.logger, ipToNodeNameMap, workerIPs)

	for _, nodeName := range controlPlaneNames {
		k.logger.Info("Labeling node", zap.String("node", nodeName), zap.String("label", "kube-ovn/role=master"))
		if _, err := k.kubectl.ExecuteCommand(
			ctx,
			"--server", server,
			"--kubeconfig", kubeconfig,
			"label", "node", nodeName, "kube-ovn/role=master", "--overwrite",
			"--insecure-skip-tls-verify=true",
		); err != nil {
			return fmt.Errorf("failed to label control plane node %s: %w", nodeName, err)
		}
	}

	for _, nodeName := range workerNames {
		k.logger.Info("Labeling node", zap.String("node", nodeName), zap.String("label", "node-role.kubernetes.io/worker="))
		if _, err := k.kubectl.ExecuteCommand(
			ctx,
			"--server", server,
			"--kubeconfig", kubeconfig,
			"label", "node", nodeName, "node-role.kubernetes.io/worker=", "--overwrite",
			"--insecure-skip-tls-verify=true",
		); err != nil {
			return fmt.Errorf("failed to label worker node %s: %w", nodeName, err)
		}
	}

	return nil
}

// ConfigureKubeOvn generates a rendered values.yaml file and performs a Helm install
// of the Kube-OVN chart using the VIP and remaining control plane IPs.
// The bound control plane node IP is excluded from the list.
func (k *KubeOvnInitializer) ConfigureKubeOvn(ctx context.Context, controlPlaneIPs []string, config *models.BootstrapConfig) error {
	k.logger.Info("Rendering and applying Kube-OVN values.yaml")

	vip := config.ManagementCluster.Talos.ControlPlaneVIP
	boundIP := config.ManagementCluster.Talos.BoundNodeIP

	k.logger.Info("Starting Kube-OVN control plane IP rendering",
		zap.String("vip", vip),
		zap.String("boundNodeIP", boundIP),
		zap.Strings("controlPlaneIPs", controlPlaneIPs),
	)

	renderedIPs := []string{vip}
	seen := map[string]bool{vip: true}

	for _, ip := range controlPlaneIPs {
		if ip == boundIP {
			k.logger.Info("Skipping bound control plane IP", zap.String("ip", ip))
			continue
		}
		if !seen[ip] {
			k.logger.Info("Adding control plane IP", zap.String("ip", ip))
			renderedIPs = append(renderedIPs, ip)
			seen[ip] = true
		}
	}

	k.logger.Info("Final rendered control plane IP list for Kube-OVN",
		zap.Strings("renderedIPs", renderedIPs))

	// Template values
	data := map[string]interface{}{
		"MASTER_NODES": strings.Join(renderedIPs, ","),
		"NODE_IPS":     strings.Join(renderedIPs, ","),
	}

	// Render template
	tmpl, err := template.New("kube-ovn-values").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(baseKubeOvnValues)
	if err != nil {
		return fmt.Errorf("failed to parse Kube-OVN values template: %w", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("failed to render Kube-OVN values template: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "kube-ovn-values-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(rendered.Bytes()); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	k.logger.Info("Installing Kube-OVN via Helm", zap.String("valuesFile", tmpFile.Name()))
	_, err = k.helm.ExecuteCommand(
		ctx,
		"install", "kube-ovn", "kube-ovn/kube-ovn",
		"-n", "kube-system",
		"-f", tmpFile.Name(),
		"--kubeconfig", "talosconfig/kubeconfig",
		"--insecure-skip-tls-verify",
	)
	if err != nil {
		return fmt.Errorf("failed to install Kube-OVN: %w", err)
	}

	k.logger.Info("Kube-OVN installed successfully")
	return nil
}

// resolveNodeNames converts a list of internal IPs to node names
// using the provided IP-to-node name map. Logs and skips missing entries.
func resolveNodeNames(logger *zap.Logger, ipToNode map[string]string, ips []string) []string {
	var names []string
	seen := make(map[string]bool)

	for _, ip := range ips {
		nodeName, found := ipToNode[ip]
		if !found {
			logger.Warn("IP not found in cluster nodes", zap.String("ip", ip))
			continue
		}
		if !seen[nodeName] {
			names = append(names, nodeName)
			seen[nodeName] = true
		}
	}
	return names
}

// getInternalIPToNodeNameMap retrieves the mapping of internal IPs to node names
// by parsing the output of `kubectl get nodes -o wide`.
func (k *KubeOvnInitializer) getInternalIPToNodeNameMap(ctx context.Context, server string) (map[string]string, error) {
	out, err := k.kubectl.ExecuteCommand(ctx,
		"--server", server,
		"--kubeconfig", "talosconfig/kubeconfig",
		"get", "nodes", "-o", "wide",
		"--insecure-skip-tls-verify=true",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader([]byte(out)))
	var headers []string

	if scanner.Scan() {
		headers = strings.Fields(scanner.Text())
	}

	var ipIndex, nameIndex int = -1, -1
	for i, header := range headers {
		if header == "INTERNAL-IP" {
			ipIndex = i
		}
		if header == "NAME" {
			nameIndex = i
		}
	}
	if ipIndex == -1 || nameIndex == -1 {
		return nil, fmt.Errorf("could not find required columns in kubectl output")
	}

	ipToNode := make(map[string]string)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) <= ipIndex || len(fields) <= nameIndex {
			continue
		}
		ip := fields[ipIndex]
		name := fields[nameIndex]
		ipToNode[ip] = name
	}

	return ipToNode, nil
}

// WaitForNodes waits until at least one node is registered and in Ready state,
// polling every 5 seconds up to the provided timeout.
func (k *KubeOvnInitializer) WaitForNodes(ctx context.Context, server string, timeout time.Duration) error {
	k.logger.Info("Waiting for nodes to be ready...")

	kubeconfig := "talosconfig/kubeconfig"
	start := time.Now()

	for {
		out, err := k.kubectl.ExecuteCommand(ctx,
			"--server", server,
			"--kubeconfig", kubeconfig,
			"get", "nodes", "--request-timeout=15s",
			"--insecure-skip-tls-verify=true",
		)
		if err == nil && strings.Contains(out, "Ready") {
			k.logger.Info("Nodes detected in cluster")
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for nodes to become ready")
		}
		time.Sleep(5 * time.Second)
	}
}
