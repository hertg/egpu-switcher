package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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
			if err := switchInternal(); err != nil {
				logger.Error("switch failed")
				return err
			}
			logger.Success("switch successful")
			return nil
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

	modesetting := !viper.GetBool("egpu.disableModesetting") // note: absent config defaults to 'false'
	conf := xorg.RenderConf("Device0", driver, gpu.XorgPCIString(), modesetting)
	if err := xorg.CreateEgpuFile(x11ConfPath, conf, verbose); err != nil {
		return err
	}
	if post := viper.GetString("hooks.egpu"); post != "" {
		if err := runHook(post); err != nil {
			logger.Error("egpu hook error: %s", err)
		}
	}
	return nil
}

func switchInternal() error {
	if err := xorg.RemoveEgpuFile(x11ConfPath, verbose); err != nil {
		return err
	}
	if post := viper.GetString("hooks.internal"); post != "" {
		if err := runHook(post); err != nil {
			logger.Error("internal hook error: %s", err)
		}
	}
	return nil
}

func switchAuto(gpu *pci.GPU) error {
	green := color.New(color.FgHiGreen).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	if gpu != nil {
		logger.Info("the egpu is %s", green("connected"))
		err := switchEgpu(gpu)
		return err
	}
	logger.Info("the egpu is %s", red("disconnected"))
	return switchInternal()
}

func runHook(script string) error {
	if !permissionCheck(script) {
		return fmt.Errorf("hook %s will not be executed", script)
	}
	cmd := exec.Command("/bin/sh", script)
	err := cmd.Run()
	if err == nil {
		logger.Info("hook script '%s' executed", script)
	}
	return err
}

func permissionCheck(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		logger.Error("%s", err)
		return false
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if stat.Uid != 0 {
			logger.Error("hook script '%s' must be owned by root", file)
			return false
		}
		if info.Mode() != 0700 {
			logger.Error("hook script '%s' must have permission 0700", file)
			return false
		}
		return true
	}
	return false
}
