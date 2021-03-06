// Package main
// Copyright (c) 2022, The Butler Authors
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
//

package ports

import (
	"strconv"
	"strings"
)

const (

	// WebHTTP is the default web HTTP port
	WebHTTP      = 3000
	HTTPHostPort = 3001
	// Carbon is the default Carbon port
	Carbon = 2003
	// Carbonserver is the default Carbonserver port
	Carbonserver = 8001
	// WebGRPC is the default GRPC port
	WebGRPC = 9998
)

// FormatHostPort returns hostPort in a usable format (host:port)
func FormatHostPort(hostPort string) string {
	if hostPort == "" {
		return ""
	}

	if strings.Contains(hostPort, ":") {
		return hostPort
	}

	return ":" + hostPort
}

// PortToHostPort converts the port into a host:port address string
func PortToHostPort(port int) string {
	return ":" + strconv.Itoa(port)
}
