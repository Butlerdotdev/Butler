package flags

import (
	"fmt"
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	AdminPort      int
	Logger         *zap.Logger
	signalsChannel chan os.Signal

	Admin *AdminServer
}

func NewService(adminPort int) *Service {
	signalsChannel := make(chan os.Signal, 1)
	//hc

	signal.Notify(signalsChannel, os.Interrupt, syscall.SIGTERM)

	return &Service{
		Admin:          NewAdminServer(ports.PortToHostPort(adminPort)),
		signalsChannel: signalsChannel,
		//hc
	}
}

func (s *Service) Start(v *viper.Viper) error {
	if err := TryLoadConfigFile(v); err != nil {
		return fmt.Errorf("cannot load config file: %w", err)
	}

	sFlags := new(SharedFlags).InitFromViper(v)
	newProdConfig := zap.NewProductionConfig()
	newProdConfig.Sampling = nil
	if logger, err := sFlags.NewLogger(newProdConfig); err == nil {
		s.Logger = logger
	} else {
		return fmt.Errorf("cannot create logger: %w", err)
	}

	s.Admin.initFromViper(v, s.Logger)

	if err := s.Admin.Serve(); err != nil {
		return fmt.Errorf("cannot start the admin server: %w", err)
	}
	return nil
}
