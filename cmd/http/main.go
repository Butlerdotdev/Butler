package main

import (
	"fmt"
	"github.com/butdotdev/butler/cmd/flags"
	"github.com/butdotdev/butler/cmd/http/app"
	"github.com/butdotdev/butler/ports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

const serviceName = "butler-http"

func main()  {
	svc := flags.NewService(ports.WebAdminHTTP)

	v := viper.New()
	command := &cobra.Command {
		Use:   "butler-http",
		Short: "butler http is the main http server for butler",
		Long:  `Butler http is the server that runs and serves the butler frontend.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := svc.Start(v); err != nil {
				return err
			}
			logger := svc.Logger

			w := app.New(&app.WebParams{
				ServiceName: serviceName,
				Logger:      logger,
			})
			webOpts := new(app.WebOptions).InitFromViper(v)
			if err := w.Start(webOpts); err != nil {
				logger.Fatal("Failed to start web server", zap.Error(err))
			}

		return nil

		},
	}
	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}