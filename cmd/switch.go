package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/hertg/egpu-switcher/internal/xorg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const x11ConfPath = "/etc/X11/xorg.conf.d/99-egpu-switcher.conf"

var override bool

var switchCommand = &cobra.Command{
	Use:   "switch [auto|internal|external]",
	Short: "todo",
	RunE: func(cmd *cobra.Command, args []string) error {

		// if no argument is provided, default to 'auto'
		if len(args) == 0 {
			args = []string{"auto"}
		}

		arg := args[0]

		// create /etc/X11/xorg.conf.d/ if directory doesn't exist
		dir := filepath.Dir(x11ConfPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("unable to create missing directories for %s", x11ConfPath)
		}

		id := viper.GetInt("egpu.id")
		if id == 0 {
			return fmt.Errorf("egpu-switcher has not been configured, we don't know what gpu to look out for...")
		}
		gpu := pci.Find(uint64(id))

		switch arg {
		case "internal":
			return switchInternal()
		case "external":
			if gpu == nil {
				return fmt.Errorf("the eGPU is not connected, unable to switch")
			}
			return switchEgpu(gpu)
		case "auto":
			green := color.New(color.FgHiGreen).SprintFunc()
			red := color.New(color.FgHiRed).SprintFunc()
			if gpu != nil {
				logger.Info("the eGPU is %s", green("connected"))
				return switchEgpu(gpu)
			} else {
				logger.Info("the eGPU is %s", red("disconnected"))
				return switchInternal()
			}
		default:
			return fmt.Errorf("unknown value %s", arg)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCommand)
	switchCommand.PersistentFlags().BoolVar(&override, "override", false, "switch to the eGPU even if there are no displays attached") // todo: usage
}

func switchEgpu(gpu *pci.GPU) error {
	driver := viper.GetString("egpu.driver")

	if driver != "nvidia" {
		displayCount, err := gpu.NumberOfDisplays()
		if err != nil {
			return err
		}
		isConnected, err := gpu.HasDisplaysConnected()
		if err != nil {
			return err
		}
		if displayCount > 0 && !isConnected {
			logger.Warn("No eGPU attached display detected with open source drivers. (Of %d eGPU outputs detected) Internal mode and setting DRI_PRIME variable are recommended for this configuration.\n", displayCount)
			if !override {
				return fmt.Errorf("Not setting eGPU mode. Run the command with the '--override' flag to force loading eGPU mode")
			}
			logger.Debug("-> Overridden: setting eGPU mode")
		}
	}

	conf := xorg.RenderConf("Device0", driver, gpu.XorgPCIString())
	return xorg.CreateEgpuFile(x11ConfPath, conf, verbose)
}

func switchInternal() error {
	return xorg.RemoveEgpuFile(x11ConfPath, verbose)
}
