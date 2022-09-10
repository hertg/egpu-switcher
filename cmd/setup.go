package cmd

import (
	"context"
	"fmt"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/service"
	"github.com/spf13/cobra"
)

var setupCommand = &cobra.Command{
	Use:   "setup",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isRoot {
			return fmt.Errorf("you need root privileges to setup egpu-switcher")
		}

		ctx := context.Background()

		init, err := service.GetInitSystem()
		if err != nil {
			return err
		}

		// todo: trigger config if no config exists (unless --no-prompt is used)

		if err := init.CreateService(ctx); err != nil {
			return err
		}

		logger.Success("setup successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
