package config

import (
	"fmt"
	"github.com/butdotdev/butler/pkg/scylla"
	gocqlx "github.com/butdotdev/butler/pkg/scylla/gocql"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"time"
)

// Configuration holds config values for a scylla cluster
type Configuration struct {
	Hosts              []string      `validate:"nonzero" mapstructure:"hosts"`
	Keyspace           string        `validate:"nonzero" mapstructure:"keyspace"`
	LocalDC            string        `yaml:"local_dc" mapstructure:"local_dc"`
	ConnectionsPerHost int           `validate:"min=1" yaml:"connections_per_host" mapstructure:"connections_per_host"`
	Timeout            time.Duration `validate:"min=500" mapstructure:"-"`
	ConnectTimeout     time.Duration `yaml:"connect_timeout" mapstructure:"connection_timeout"`
	ReconnectInterval  time.Duration `validate:"min=500" yaml:"reconnect_interval" mapstructure:"reconnect_interval"`
	SocketKeepAlive    time.Duration `validate:"min=0" yaml:"socket_keep_alive" mapstructure:"socket_keep_alive"`
	MaxRetryAttempts   int           `validate:"min=0" yaml:"max_retry_attempt" mapstructure:"max_retry_attempts"`
	ProtoVersion       int           `yaml:"proto_version" mapstructure:"proto_version"`
	Consistency        string        `yaml:"consistency" mapstructure:"consistency"`
	DisableCompression bool          `yaml:"disable-compression" mapstructure:"disable_compression"`
	Port               int           `yaml:"port" mapstructure:"port"`
}

// ApplyDefaults copies settings from source unless its own value is non-zero.
func (c *Configuration) ApplyDefaults(source *Configuration) {
	if c.ConnectionsPerHost == 0 {
		c.ConnectionsPerHost = source.ConnectionsPerHost
	}
	if c.MaxRetryAttempts == 0 {
		c.MaxRetryAttempts = source.MaxRetryAttempts
	}
	if c.Timeout == 0 {
		c.Timeout = source.Timeout
	}
	if c.ReconnectInterval == 0 {
		c.ReconnectInterval = source.ReconnectInterval
	}
	if c.Port == 0 {
		c.Port = source.Port
	}
	if c.Keyspace == "" {
		c.Keyspace = source.Keyspace
	}
	if c.ProtoVersion == 0 {
		c.ProtoVersion = source.ProtoVersion
	}
	if c.SocketKeepAlive == 0 {
		c.SocketKeepAlive = source.SocketKeepAlive
	}
}

// SessionBuilder creates new scylla.Session
type SessionBuilder interface {
	NewSession(logger *zap.Logger) (scylla.Session, error)
}

// NewSession creates a new scylla session
func (c *Configuration) NewSession(logger *zap.Logger) (scylla.Session, error) {
	cluster, err := c.NewCluster(logger)
	if err != nil {
		return nil, err
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return gocqlx.WrapCQLSession(session), nil
}

// NewCluster defines a new scylla cluster
func (c *Configuration) NewCluster(logger *zap.Logger) (*gocql.ClusterConfig, error) {
	cluster := gocql.NewCluster(c.Hosts...)
	cluster.Keyspace = c.Keyspace
	cluster.NumConns = c.ConnectionsPerHost
	cluster.Timeout = c.Timeout
	cluster.ConnectTimeout = c.ConnectTimeout
	cluster.ReconnectInterval = c.ReconnectInterval
	cluster.SocketKeepalive = c.SocketKeepAlive
	if c.ProtoVersion > 0 {
		cluster.ProtoVersion = c.ProtoVersion
	}
	if c.MaxRetryAttempts > 1 {
		cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: c.MaxRetryAttempts - 1}
	}
	if c.Port != 0 {
		cluster.Port = c.Port
	}

	if !c.DisableCompression {
		cluster.Compressor = gocql.SnappyCompressor{}
	}

	if c.Consistency == "" {
		cluster.Consistency = gocql.LocalOne
	} else {
		cluster.Consistency = gocql.ParseConsistency(c.Consistency)
	}

	fallbackHostSelectionPolicy := gocql.RoundRobinHostPolicy()
	if c.LocalDC != "" {
		fallbackHostSelectionPolicy = gocql.DCAwareRoundRobinPolicy(c.LocalDC)
	}
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(fallbackHostSelectionPolicy, gocql.ShuffleReplicas())

	return cluster, nil
}

// String returns formatted string
func (c *Configuration) String() string {
	return fmt.Sprintf("%+v", *c)
}
