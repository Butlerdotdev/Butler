// Package tlscfg
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

package tlscfg

import (
	"flag"
	"github.com/spf13/viper"
)

const (
	tlsPrefix         = ".tls"
	tlsEnabled        = tlsPrefix + ".enabled"
	tlsCA             = tlsPrefix + ".ca"
	tlsCert           = tlsPrefix + ".cert"
	tlsKey            = tlsPrefix + ".key"
	tlsServerName     = tlsPrefix + "server-name"
	tlsClientCA       = tlsPrefix + ".client-ca"
	tlsSkipHostVerify = tlsPrefix + ".skip-host-verify"
)

// ClientFlagsConfig describes which CLI flags for TLS client should be generated
type ClientFlagsConfig struct {
	Prefix string
}

// ServerFlagConfig describe which CLI flags for TLS server should be generated
type ServerFlagConfig struct {
	Prefix string
}

// AddFlags adds flags for TLS to the FlagSet.
func (c ClientFlagsConfig) AddFlags(flags *flag.FlagSet) {
	flags.Bool(c.Prefix+tlsEnabled, false, "Enable TLS when talking to the remote server")
	flags.String(c.Prefix+tlsCA, "", "Path to a TLS CA file used to verify the remote server")
	flags.String(c.Prefix+tlsCert, "", "Path to the TLS cert file, used to identify this process to the remote server")
	flags.String(c.Prefix+tlsKey, "", "Path to the TLS private key file, used to identify this process to the remote server")
	flags.String(c.Prefix+tlsServerName, "", "Override the TLS server name we expect in the cer of the remote server")
	flags.Bool(c.Prefix+tlsSkipHostVerify, false, "(insecure) Skip server's cert chain and host name verification")
}

// AddFlags adds flags for TLS to the FlatSet
func (c ServerFlagConfig) AddFlags(flags *flag.FlagSet) {
	flags.Bool(c.Prefix+tlsEnabled, false, "Enable TLS on the server")
	flags.String(c.Prefix+tlsCert, "", "Path to a TLS cert file, used to Identify this server to clients")
	flags.String(c.Prefix+tlsKey, "", "Path to a TLS Private Key file, used to identify this server to clients")
	flags.String(c.Prefix+tlsClientCA, "", "Path to a TLS CA file used to verify certs presented by clients (if unset, all clients are permitted")
}

// InitFromViper creates tls.Config populated with values retrieved from Viper.
func (c ClientFlagsConfig) InitFromViper(v *viper.Viper) Options {
	var p Options
	p.Enabled = v.GetBool(c.Prefix + tlsEnabled)
	p.CAPath = v.GetString(c.Prefix + tlsCA)
	p.CertPath = v.GetString(c.Prefix + tlsCert)
	p.KeyPath = v.GetString(c.Prefix + tlsKey)
	p.ServerName = v.GetString(c.Prefix + tlsServerName)
	p.SkipHostVerify = v.GetBool(c.Prefix + tlsSkipHostVerify)
	return p
}

func (c ServerFlagConfig) InitFromViper(v *viper.Viper) Options {
	var p Options
	p.Enabled = v.GetBool(c.Prefix + tlsEnabled)
	p.CertPath = v.GetString(c.Prefix + tlsCert)
	p.KeyPath = v.GetString(c.Prefix + tlsKey)
	p.ClientCAPath = v.GetString(c.Prefix + tlsClientCA)
	return p
}
