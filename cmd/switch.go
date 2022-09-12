package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	Use:   "switch [auto|internal|egpu]",
	Short: "Check if eGPU is present and configure X.org accordingly",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !isRoot {
			return fmt.Errorf("you need root privileges to switch gpu")
		}

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
			logger.Error("it seems that there is no configuration present, we don't know what gpu to look out for...")

			oldPath := "/etc/egpu-switcher/egpu-switcher.conf"
			_, err := os.Stat(oldPath)
			if err == nil {
				logger.Info("there seems to be an old configuration present at %s", oldPath)
				logger.Info("this config format is no longer valid and can not be automatically migrated")
			}

			logger.Info("please connect your eGPU and run 'egpu-switcher config' to resolve this issue")

			return fmt.Errorf("no configuration found")
		}

		if arg == "internal" {
			logger.Success("switch successful")
			return switchInternal()
		}

		gpu := pci.Find(uint64(id))
	Outer:
		switch arg {
		case "egpu", "external":
			if gpu == nil {
				return fmt.Errorf("the eGPU is not connected, unable to switch")
			}
			err = switchEgpu(gpu)
			break
		case "auto":
			logger.Info("looking for eGPU...")

			if gpu != nil {
				err = switchAuto(gpu)
				break Outer
			}

			maxRetries := viper.GetInt("detection.retries")
			interval := viper.GetInt("detection.interval")
			attempt := 0
			for {
				if attempt > maxRetries {
					logger.Info("giving up after %d retries", maxRetries)
					err = switchAuto(gpu)
					break Outer
				}
				<-time.After(time.Duration(interval) * time.Millisecond)
				gpu = pci.Find(uint64(id))
				if gpu != nil {
					err = switchAuto(gpu)
					break Outer
				}
				attempt += 1
			}
		default:
			return fmt.Errorf("unknown value %s", arg)
		}

		if err == nil {
			logger.Success("switch completed")
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(switchCommand)
	switchCommand.PersistentFlags().BoolVar(&override, "override", false, "switch to the eGPU even if there are no displays attached") // todo: usage
}

func switchEgpu(gpu *pci.GPU) error {
	driver := viper.GetString("egpu.driver")

	if driver != "nvidia" {
		outputs, err := gpu.Outputs()
		if err != nil {
			return err
		}
		connectedOutputs, err := gpu.ConnectedOutputs()
		if err != nil {
			return err
		}
		if outputs > 0 && connectedOutputs == 0 {
			logger.Warn("No eGPU attached display detected with open source drivers. (Of %d eGPU outputs detected) Internal mode and setting DRI_PRIME variable are recommended for this configuration.\n", outputs)
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

func switchAuto(gpu *pci.GPU) error {
	green := color.New(color.FgHiGreen).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	if gpu != nil {
		logger.Info("the egpu is %s", green("connected"))
		return switchEgpu(gpu)
	}
	logger.Info("the egpu is %s", red("disconnected"))
	return switchInternal()
}
