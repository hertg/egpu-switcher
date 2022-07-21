package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCommand = &cobra.Command{
	Use:   "setup",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		// todo: trigger config if no config exists (unless --no-prompt is used)
		// todo: create init system service
		return fmt.Errorf("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
