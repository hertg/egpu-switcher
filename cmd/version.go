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

		if full {
			buildtime := buildinfo.BuildTime
			if buildtime == "" {
				buildtime = "unknown"
			}
			origin := buildinfo.Origin
			if origin == "" {
				origin = "unknown"
			}
			fmt.Printf("%s_%s_%s\n", buildinfo.Version, buildtime, origin)
			return nil
		}

		fmt.Println(buildinfo.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().BoolVar(&full, "full", false, "display all build information")
}
