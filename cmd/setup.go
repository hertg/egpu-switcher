package cmd

import (
	"fmt"

	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/spf13/cobra"
)

var setupCommand = &cobra.Command{
	Use:   "setup",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		pci.Test()
		return fmt.Errorf("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
