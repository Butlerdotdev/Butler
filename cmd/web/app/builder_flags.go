// Package app
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

package app

import (
	"flag"
	"github.com/butdotdev/butler/pkg/config/tlscfg"
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/viper"
	"time"
)

const (
	webHTTPHostPort                = "web-server"
	webGRPCHostPort                = "grpc-server.host-port"
	webGRPCMaxReceiveMessageLength = "grpc-server.max-message-size"
	webMaxConnectionAge            = "grpc-server.max-connection-age"
	webMaxConnectionAgeGrace       = "grpc-server.max-connection-age-grace"
)

var tlsHTTPFlagsConfig = tlscfg.ServerFlagsConfig{
	Prefix: "web.http",
}

var tlsGRPCFlagsConfig = tlscfg.ServerFlagsConfig{
	Prefix: "web.grpc",
}

// WebOptions struct holds config for web
type WebOptions struct {
	// WebHTTPHostPort is the host:port address that the web service listens in on for http requests
	WebHTTPHostPort string

	//WebGRPCHostPort is the host:port addres that the grpc server listens on for RPC requests
	WebGRPCHostPort string

	// TLSHTTP configures secure transport for HTTP endpoint
	TLSHTTP tlscfg.Options

	//TLSGRPC configures secure transport for gRPC
	TLSGRPC tlscfg.Options

	//WebGRPCMaxReceiveMessageLength is the max message size receivable by the gRPC server
	WebGRPCMaxReceiveMessageLength int

	//WebGRPCMaxConnectionAge is the duration for the max amount of time a connection can exist.
	WebGRPCMaxConnectionAge time.Duration

	//WebGRPCMaxConnectionAgeGrace is an additive period of MaxConnectionAge after which teh connection will forceibly close
	WebGRPCMaxConnectionAgeGrace time.Duration
}

// AddFlags adds flags for WebOptions
func AddFlags(flags *flag.FlagSet) {
	flags.String(webHTTPHostPort, ports.PortToHostPort(ports.WebHTTP), "The host:port (e.g. 127.0.0.1:3000 ) of Butler's HTTP server")
	flags.String(webGRPCHostPort, ports.PortToHostPort(ports.WebGRPC), "The host:port (e.g. 127.0.0.1:9998) of the GRPC server")
	flags.Duration(webMaxConnectionAge, 0, "The max amount of time a connection may exist. Set this value to a few secs or minutes on highly elastic envs")
	flags.Duration(webMaxConnectionAgeGrace, 0, "The additive period after MaxConnectionAge after which the connection will forcibly close. See https://pkg.go.dev/google.golang.org/grpc/keepalive#ServerParameters")

	tlsHTTPFlagsConfig.AddFlags(flags)
	tlsGRPCFlagsConfig.AddFlags(flags)
}

// InitFromViper initializes WebOptions with props from viper
func (wOpts *WebOptions) InitFromViper(v *viper.Viper) *WebOptions {
	wOpts.WebHTTPHostPort = ports.FormatHostPort(v.GetString(webHTTPHostPort))
	wOpts.WebGRPCHostPort = ports.FormatHostPort(v.GetString(webGRPCHostPort))
	wOpts.WebGRPCMaxConnectionAgeGrace = v.GetDuration(webMaxConnectionAgeGrace)
	wOpts.WebGRPCMaxConnectionAge = v.GetDuration(webMaxConnectionAge)
	wOpts.WebGRPCMaxReceiveMessageLength = v.GetInt(webGRPCMaxReceiveMessageLength)
	wOpts.TLSHTTP = tlsHTTPFlagsConfig.InitFromViper(v)
	return wOpts
}
