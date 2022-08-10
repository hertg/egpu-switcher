package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "configure egpu-switcher",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {

		gpus := pci.ReadGPUs()
		amount := int(len(gpus))

		/*if amount < 2 {
			logger.Error("only one GPU found... please plug in your eGPU to continue")
			os.Exit(1)
		}*/

		fmt.Println()
		fmt.Printf("Found %d possible GPU(s)...\n", amount)
		fmt.Println()
		for i, gpu := range gpus {
			fmt.Printf("%d: %s\n", i+1, gpu.DisplayName())
			fmt.Printf("%s\n", gpu.XorgPCIString())
			hasDisp, _ := gpu.HasDisplaysConnected()
			numDisp, _ := gpu.NumOfConnectedDisplays()
			fmt.Printf("%v\n", hasDisp)
			fmt.Printf("%d\n", numDisp)
			fmt.Printf("%d\n", gpu.Identifier())
		}
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Which one is your %s GPU? [%d-%d]: ", green("external"), 1, len(gpus))
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSuffix(answer, "\n")
		num, err := strconv.ParseInt(answer, 10, 8)
		if err != nil {
			return fmt.Errorf("invalid number '%s'", answer)
		}
		if num < 1 || int(num) > len(gpus) {
			return fmt.Errorf("number '%s' out of range", answer)
		}

		selected := gpus[num-1]
		driver := selected.PciDevice.Driver
		// if err != nil {
		// 	logger.Info(err.Error())
		// 	fmt.Printf("Please manually enter the driver to be used: ")
		// 	answer, _ = reader.ReadString('\n')
		// 	driver = strings.TrimSuffix(answer, "\n")
		// }

		viper.Set("egpu.driver", driver)
		viper.Set("egpu.id", selected.Identifier())
		viper.WriteConfig()

		fmt.Println()

		logger.Success("Your selection has been saved")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCommand)
}
