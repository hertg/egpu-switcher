package cmd

import (
	"fmt"

	"github.com/hertg/egpu-switcher/internal/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s (%s)", buildinfo.Version, buildinfo.BuildTime)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
