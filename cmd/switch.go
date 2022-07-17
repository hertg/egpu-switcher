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

const path = "/etc/X11/xorg.conf.d/99-egpu-switcher.conf"

var switchCommand = &cobra.Command{
	Use:   "switch [auto|internal|external]",
	Short: "todo",
	RunE: func(cmd *cobra.Command, args []string) error {
		id := viper.GetInt("egpu.id")
		if id == 0 {
			return fmt.Errorf("egpu-switcher has not been configured, we don't know what gpu to look out for...")
		}

		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("unable to create missing directories for %s", path)
		}

		gpu := pci.Find(uint64(id))
		green := color.New(color.FgHiGreen).SprintFunc()
		red := color.New(color.FgHiRed).SprintFunc()

		if gpu != nil {
			logger.Info("the eGPU is %s", green("connected"))
			return switchEgpu(gpu)
		} else {
			logger.Info("the eGPU is %s", red("disconnected"))
			return switchInternal()
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCommand)
}

func switchEgpu(gpu *pci.GPU) error {
	driver := viper.GetString("egpu.driver")
	conf := xorg.GenerateConf("Device0", driver, gpu.XorgPCIString())
	_, err := os.Stat(path)
	if err != nil {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("unable to create file %s", path)
		}
		_, err = f.Write([]byte(conf))
		if err != nil {
			return fmt.Errorf("unable to write config to file %s", path)
		}
		return nil
	}
	logger.Info("the file %s already exists", path)
	return nil
}

func switchInternal() error {
	f, _ := os.Stat(path)
	if f != nil {
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("unable to remove file %s", path)
		}
	}
	logger.Info("the file %s is already absent", path)
	return nil
}
