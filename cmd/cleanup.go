package cmd

import (
	"os"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/xorg"
	"github.com/spf13/cobra"
)

var cleanupCommand = &cobra.Command{
	Use: "cleanup",
	RunE: func(cmd *cobra.Command, args []string) error {

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
	cleanupCommand.PersistentFlags().BoolVar(&hard, "hard", false, "also remove configuration files")
}
