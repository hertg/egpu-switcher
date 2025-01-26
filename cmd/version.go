package cmd

import (
	"fmt"

	"github.com/hertg/egpu-switcher/internal/buildinfo"
	"github.com/spf13/cobra"
)

var full bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(buildinfo.VersionString(full))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().BoolVar(&full, "full", false, "display all build information")
}
