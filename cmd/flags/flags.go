package flags

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevel = "log-level"
	configFile = "config-file"
)

func AddConfigFileFlag(flagSet *flag.FlagSet) {
	flagSet.String(configFile, "",  "Configuration file in JSON, TOML, YAML, HCL, or Java properties formats (default none). See spf13/viper for precedence.")
}

func TryLoadConfigFile(v *viper.Viper) error {
	if file := v.GetString(configFile); file != "" {
		v.SetConfigFile(file)
		err := v.ReadInConfig()
		if err != nil {
			return fmt.Errorf("cannot load config file %s: %w", file, err)
		}
	}
	return nil
}

type SharedFlags struct {
	Logging logging
}

type logging struct {
	Level string
}

func AddLoggingFlag(flagSet *flag.FlagSet) {
	flagSet.String(logLevel, "info", "Minimal allowed log level")
}

func (flags *SharedFlags) InitFromViper(v *viper.Viper) *SharedFlags {
	flags.Logging.Level = v.GetString(logLevel)
	return flags
}

func (flags *SharedFlags) NewLogger(conf zap.Config, options ...zap.Option) (*zap.Logger, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(flags.Logging.Level))
	if err != nil {
		return nil, err
	}
	conf.Level = zap.NewAtomicLevelAt(level)
	return conf.Build(options...)
}