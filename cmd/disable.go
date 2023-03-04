package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/service"
	"github.com/hertg/egpu-switcher/internal/xorg"
	"github.com/spf13/cobra"
)

var disableCommand = &cobra.Command{
	Use:     "disable",
	Aliases: []string{"cleanup"}, // backwards compatibility
	Short:   "Disable egpu-switcher from running at startup",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !isRoot {
			return fmt.Errorf("you need root privileges to cleanup egpu-switcher")
		}

		ctx := context.Background()

		// remove the file from /etc/X11/xorg.conf.d (if present)
		err := xorg.RemoveEgpuFile(x11ConfPath, verbose)
		if err != nil {
			return err
		}

		// remove egpu service
		init, err := service.GetInitSystem()
		if err != nil {
			return err
		}
		if err := init.TeardownService(ctx, verbose); err != nil {
			return fmt.Errorf("unable to tear down service: %s", err)
		}

		if hard {
			// remove /etc/egpu-switcher
			err = os.RemoveAll(configPath)
			if err != nil {
				return err
			}
		}

		logger.Success("cleanup successful")
		return nil
	},
}

var hard bool

func init() {
	rootCmd.AddCommand(disableCommand)
	disableCommand.PersistentFlags().BoolVar(&hard, "hard", false, "remove configuration files too")
}
