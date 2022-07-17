package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "configure egpu-switcher",
	RunE: func(cmd *cobra.Command, args []string) error {

		gpus := pci.ReadGPUs()
		printGpuList(gpus)

		reader := bufio.NewReader(os.Stdin)
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Choose your preferred %s GPU [%d-%d]: ", green("external"), 1, len(gpus))
		answer, _ := reader.ReadString('\n')
		fmt.Print(answer)

		return fmt.Errorf("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(configCommand)
}

func printGpuList(gpus []*pci.GPU) {
	fmt.Println()
	for i, gpu := range gpus {
		fmt.Printf("%d: %s\n", i+1, gpu.DisplayName())
	}
	fmt.Println()
}
