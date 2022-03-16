// Package scylla
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

package scylla

import (
	"flag"
	"github.com/butdotdev/butler/pkg/scylla/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"time"
)

const (
	// session settings
	suffixEnabled            = ".enabled"
	suffixConnPerHost        = ".connections-per-host"
	suffixMaxRetryAttempts   = ".max-retry-attempts"
	suffixTimeout            = ".timeout"
	suffixConnectTimeout     = ".connect-timeout"
	suffixReconnectInterval  = ".reconnect-interval"
	suffixServers            = ".servers"
	suffixPort               = ".port"
	suffixKeyspace           = ".keyspace"
	suffixDC                 = ".local-dc"
	suffixConsistency        = ".consistency"
	suffixDisableCompression = ".disable-compression"
	suffixProtoVer           = ".proto-version"
	suffixSocketKeepAlive    = ".socket-keep-alive"
	suffixUsername           = ".username"
	suffixPassword           = ".password"

	// common storage settings
	suffixRuleStoreWriteCacheTTL = ".rule-store-write-cache-ttl"
	suffixIndexTagsBlacklist     = ".index.tag-blacklist"
	suffixIndexTagsWhitelist     = ".index.tag-whitelist"
	suffixIndexLogs              = ".index.logs"
	suffixIndexTags              = ".index.tags"
	suffixIndexProcessTags       = ".index.process-tags"
)

type Options struct {
	Primary                namespaceConfig `mapstructure:",squash"`
	others                 map[string]*namespaceConfig
	RuleStoreWriteCacheTTL time.Duration `mapstructure:"rule_store_write_cache_ttl"`
	Index                  IndexConfig   `mapstructure:"index"`
	logger                 *zap.Logger
}

// IndexConfig configures indexing.
// By default all indexing is enabled.
type IndexConfig struct {
	Logs         bool   `mapstructure:"logs"`
	Tags         bool   `mapstructure:"tags"`
	ProcessTags  bool   `mapstructure:"process_tags"`
	TagBlackList string `mapstructure:"tag_blacklist"`
	TagWhiteList string `mapstructure:"tag_whitelist"`
}

type namespaceConfig struct {
	config.Configuration `mapstructure:",squash"`
	servers              string
	namespace            string
	Enabled              bool `mapstructure:"-"`
}

// NewOptions creates a new Options struct
func NewOptions(primaryNamespace string, otherNamespaces ...string) *Options {
	options := &Options{
		Primary: namespaceConfig{
			Configuration: config.Configuration{
				MaxRetryAttempts:   3,
				Keyspace:           "butler",
				ProtoVersion:       4,
				ConnectionsPerHost: 2,
				ReconnectInterval:  60 * time.Second,
			},
			servers:   "127.0.0.1",
			namespace: primaryNamespace,
			Enabled:   true,
		},
		others: make(map[string]*namespaceConfig, len(otherNamespaces)),
	}
	for _, namespace := range otherNamespaces {
		options.others[namespace] = &namespaceConfig{namespace: namespace}
	}
	return options
}

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	addFlags(flagSet, opt.Primary)
	for _, cfg := range opt.others {
		addFlags(flagSet, *cfg)
	}
	flagSet.Duration(opt.Primary.namespace+suffixRuleStoreWriteCacheTTL,
		opt.RuleStoreWriteCacheTTL,
		"The duration to wait before rewriting an existing service or operation name")
	flagSet.String(
		opt.Primary.namespace+suffixIndexTagsBlacklist,
		opt.Index.TagBlackList,
		"The comma-separated list of rule tags to blacklist from being indexed. All other tags will be indexed. Mutually exclusive with the whitelist option.")
	flagSet.String(
		opt.Primary.namespace+suffixIndexTagsWhitelist,
		opt.Index.TagWhiteList,
		"The comma-separated list of rule tags to whitelist for being indexed. All other tags will not be indexed. Mutually exclusive with the blacklist option.")
	flagSet.Bool(
		opt.Primary.namespace+suffixIndexLogs,
		!opt.Index.Logs,
		"Controls log field indexing. Set to false to disable.")
	flagSet.Bool(
		opt.Primary.namespace+suffixIndexTags,
		!opt.Index.Tags,
		"Controls tag indexing. Set to false to disable.")
	flagSet.Bool(
		opt.Primary.namespace+suffixIndexProcessTags,
		!opt.Index.ProcessTags,
		"Controls process tag indexing. Set to false to disable.")
}

func addFlags(flagSet *flag.FlagSet, nsConfig namespaceConfig) {

	if nsConfig.namespace != primaryStorageConfig {
		flagSet.Bool(
			nsConfig.namespace+suffixEnabled,
			false,
			"Enable extra storage")
	}
	flagSet.Int(
		nsConfig.namespace+suffixConnPerHost,
		nsConfig.ConnectionsPerHost,
		"The number of scylla connections from a single backend instance")
	flagSet.Int(
		nsConfig.namespace+suffixMaxRetryAttempts,
		nsConfig.MaxRetryAttempts,
		"The number of attempts when reading from scylla")
	flagSet.Duration(
		nsConfig.namespace+suffixTimeout,
		nsConfig.Timeout,
		"Timeout used for queries. A Timeout of zero means no timeout")
	flagSet.Duration(
		nsConfig.namespace+suffixConnectTimeout,
		nsConfig.ConnectTimeout,
		"Timeout used for connections to scylla Servers")
	flagSet.Duration(
		nsConfig.namespace+suffixReconnectInterval,
		nsConfig.ReconnectInterval,
		"Reconnect interval to retry connecting to downed hosts")
	flagSet.String(
		nsConfig.namespace+suffixServers,
		nsConfig.servers,
		"The comma-separated list of scylla servers")
	flagSet.Int(
		nsConfig.namespace+suffixPort,
		nsConfig.Port,
		"The port for scylla")
	flagSet.String(
		nsConfig.namespace+suffixKeyspace,
		nsConfig.Keyspace,
		"The scylla keyspace for Butler data")
	flagSet.String(
		nsConfig.namespace+suffixDC,
		nsConfig.LocalDC,
		"The name of the scylla local data center for DC Aware host selection")
	flagSet.String(
		nsConfig.namespace+suffixConsistency,
		nsConfig.Consistency,
		"The scylla consistency level, e.g. ANY, ONE, TWO, THREE, QUORUM, ALL, LOCAL_QUORUM, EACH_QUORUM, LOCAL_ONE (default LOCAL_ONE)")
	flagSet.Bool(
		nsConfig.namespace+suffixDisableCompression,
		false,
		"Disables the use of the default Snappy Compression while connecting to the scylla Cluster if set to true. This is useful for connecting to scylla Clusters(like Azure Cosmos Db with scylla API) that do not support SnappyCompression")
	flagSet.Int(
		nsConfig.namespace+suffixProtoVer,
		nsConfig.ProtoVersion,
		"The scylla protocol version")
	flagSet.Duration(
		nsConfig.namespace+suffixSocketKeepAlive,
		nsConfig.SocketKeepAlive,
		"scylla's keepalive period to use, enabled if > 0")
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	opt.Primary.initFromViper(v)
	for _, cfg := range opt.others {
		cfg.initFromViper(v)
	}
	opt.RuleStoreWriteCacheTTL = v.GetDuration(opt.Primary.namespace + suffixRuleStoreWriteCacheTTL)
	opt.Index.Tags = v.GetBool(opt.Primary.namespace + suffixIndexTags)
	opt.Index.Logs = v.GetBool(opt.Primary.namespace + suffixIndexLogs)
	opt.Index.ProcessTags = v.GetBool(opt.Primary.namespace + suffixIndexProcessTags)
}

func (cfg *namespaceConfig) initFromViper(v *viper.Viper) {
	if cfg.namespace != primaryStorageConfig {
		cfg.Enabled = v.GetBool(cfg.namespace + suffixEnabled)
	}
	cfg.ConnectionsPerHost = v.GetInt(cfg.namespace + suffixConnPerHost)
	cfg.MaxRetryAttempts = v.GetInt(cfg.namespace + suffixMaxRetryAttempts)
	cfg.Timeout = v.GetDuration(cfg.namespace + suffixTimeout)
	cfg.ConnectTimeout = v.GetDuration(cfg.namespace + suffixConnectTimeout)
	cfg.ReconnectInterval = v.GetDuration(cfg.namespace + suffixReconnectInterval)
	cfg.Port = v.GetInt(cfg.namespace + suffixPort)
	cfg.Keyspace = v.GetString(cfg.namespace + suffixKeyspace)
	cfg.LocalDC = v.GetString(cfg.namespace + suffixDC)
	cfg.Consistency = v.GetString(cfg.namespace + suffixConsistency)
	cfg.ProtoVersion = v.GetInt(cfg.namespace + suffixProtoVer)
	cfg.SocketKeepAlive = v.GetDuration(cfg.namespace + suffixSocketKeepAlive)
	cfg.DisableCompression = v.GetBool(cfg.namespace + suffixDisableCompression)
}

func (opt *Options) GetPrimary() *config.Configuration {
	opt.Primary.Hosts = strings.Split(opt.Primary.servers, ",")
	return &opt.Primary.Configuration
}
