package cmd

import (
	"fmt"
	"os"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/xorg"
	"github.com/spf13/cobra"
)

var cleanupCommand = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove any xorg files egpu-switcher might have created",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !isRoot {
			return fmt.Errorf("you need root privileges to cleanup egpu-switcher")
		}

		// remove the file from /etc/X11/xorg.conf.d (if present)
		err := xorg.RemoveEgpuFile(x11ConfPath, verbose)
		if err != nil {
			return err
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
	rootCmd.AddCommand(cleanupCommand)
	cleanupCommand.PersistentFlags().BoolVar(&hard, "hard", false, "remove configuration files too")
}
