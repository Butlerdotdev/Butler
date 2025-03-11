// Package models defines data structures for Talos and its config .
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

package models

// TalosConfig holds cluster bootstrap information.
type TalosConfig struct {
	ClusterName       string   `json:"cluster_name"`
	ControlPlaneIP    string   `json:"control_plane_ip"`
	OutputDir         string   `json:"output_dir"`
	ControlPlaneNodes []string `json:"control_plane_nodes"`
	WorkerNodes       []string `json:"worker_nodes"`
}
