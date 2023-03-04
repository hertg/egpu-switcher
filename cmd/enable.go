package cmd

import (
	"context"
	"fmt"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var noPrompt bool

var setupCommand = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"setup"}, // backwards compatibility
	Short:   "Enable egpu-switcher to run at startup",
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
		egpuId := viper.GetInt("egpu.id")
		if egpuId == 0 {
			logger.Info("no eGPU has been configured yet")
			if noPrompt {
				logger.Warn("please run 'egpu-switcher config' to configure your eGPU")
				return fmt.Errorf("setup aborted")
			} else {
				err := configCommand.RunE(cmd, []string{})
				if err != nil {
					return err
				}
			}
		}

		if err := init.CreateService(ctx, verbose); err != nil {
			return err
		}

		logger.Info("created egpu bootup service to autorun 'egpu-switcher switch'")

		logger.Success("setup successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
	setupCommand.PersistentFlags().BoolVar(&noPrompt, "no-prompt", false, "Don't interactively prompt to configure egpu-switcher")
}
