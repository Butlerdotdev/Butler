// Package models defines data structures for Butler's cluster provisioning.
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

type ProxmoxSessionTokenResponse struct {
	Data ProxmoxTokenData `json:"data"`
}

type ProxmoxTokenData struct {
	Ticket string `json:"ticket"`
	CSRF   string `json:"CSRFPreventionToken"`
}

type ProxmoxVMConfig struct {
	VMId    int    `json:"vmid"`
	Name    string `json:"name"`
	OSType  string `json:"ostype"`
	Memory  int    `json:"memory"`
	Cores   int    `json:"cores"`
	Sockets int    `json:"sockets"`
	Start   bool   `json:"start"`
	OnBoot  bool   `json:"onboot"`
	Agent   string `json:"agent"`
	Ide2    string `json:"ide2"`
	Scsihw  string `json:"scsihw"`
	Scsi0   string `json:"scsi0"`
	Numa    bool   `json:"numa"`
	Cpu     string `json:"cpu"`
	Net0    string `json:"net0"`
}

type ProxmoxAllVMRequest struct {
	Type string `json:"type"`
}

type ProxmoxAllVMResponse struct {
	Data []ProxmoxVMResponse `json:"data"`
}

type ProxmoxVMResponse struct {
	Name    string `json:"name"`
	MaxCpu  int    `json:"maxcpu"`
	Uptime  int    `json:"uptime"`
	Node    string `json:"node"`
	Status  string `json:"status"`
	VMId    int    `json:"vmid"`
	MaxMem  int    `json:"maxmem"`
	MaxDisk int    `json:"maxdisk"`
}

type ProxmoxNetworkInterfaceResponse struct {
	Data ProxmoxNetworkInterfacesData `json:"data"`
}

type ProxmoxNetworkInterfacesData struct {
	Result []ProxmoxNetworkInterfaces `json:"result"`
}

type ProxmoxNetworkInterfaces struct {
	Name            string             `json:"name"`
	HardwareAddress string             `json:"hardware-address"`
	IpAddresses     []ProxmoxIpAddress `json:"ip-addresses"`
}

type ProxmoxIpAddress struct {
	IpAddressType string `json:"ip-address-type"`
	IpAddress     string `json:"ip-address"`
	Prefix        int    `json:"prefix"`
}
