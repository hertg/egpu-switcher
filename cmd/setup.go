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
		gpus := pci.ReadGPUs()
		for _, gpu := range gpus {
			fmt.Println(gpu.DisplayName())
			fmt.Println(gpu.XorgPCIString())
			fmt.Println(gpu.Identifier())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
