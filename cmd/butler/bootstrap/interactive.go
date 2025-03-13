package bootstrap

import (
	"butler/internal/logger"
	"butler/internal/tui"

	"github.com/spf13/cobra"
)

// NewInteractiveCmd starts bootstrap in interactive mode
func NewInteractiveCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interactive",
		Short: "Run bootstrap in interactive TUI mode",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.GetLogger()
			log.Info("Starting interactive TUI mode")
			tui.StartTUI(rootCmd, log)
		},
	}
	return cmd
}
