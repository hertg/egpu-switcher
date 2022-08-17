package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

var removeCommand = &cobra.Command{
	Use: "remove",
	RunE: func(cmd *cobra.Command, args []string) error {

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		done := make(chan bool, 1)

		id := uint64(viper.GetInt("egpu.id"))
		driver := viper.GetString("egpu.driver")
		gpu := pci.Find(id)
		if gpu == nil {
			return fmt.Errorf("the egpu is not connected")
		}

		// prefix, _, found := strings.Cut(gpu.Address(), ".")
		// if !found {
		// 	return fmt.Errorf("unable to get device id from pci address %s", gpu.Address())
		// }
		// pattern := fmt.Sprintf("/sys/bus/pci/devices/%s.[0-9]*/remove", prefix)
		// matches, err := filepath.Glob(pattern)
		// if err != nil {
		// 	return err
		// }

		systemd, err := dbus.NewSystemConnection()
		if err != nil {
			return fmt.Errorf("unable to connect to dbus")
		}
		defer systemd.Close()

		dmServiceName, err := os.Readlink("/etc/systemd/system/display-manager.service")
		dmServiceName = filepath.Base(dmServiceName)

		go func() {
			sig := <-sigChan
			logger.Debug("got signal: %s", sig)
			dmStatusChange, _ := systemd.SubscribeUnitsCustom(500*time.Millisecond, 0, func(u1, u2 *dbus.UnitStatus) bool { return *u1 != *u2 }, func(s string) bool { return s != dmServiceName })
			dmInactive := make(chan bool)
			go func() {
				for {
					select {
					case status := <-dmStatusChange:
						for _, v := range status {
							if v.ActiveState != "active" {
								dmInactive <- true
								return
							}
						}

					}
				}
			}()

			<-dmInactive // block until display manager is inactive
			logger.Debug("display-manager '%s' has become inactive", dmServiceName)

			if driver == "nvidia" {
				// systemctl stop nvidia-persistenced.service
				_, err := systemd.StopUnit("nvidia-persistenced", "replace", nil)
				if err != nil {
					logger.Error("unable to stop nvidia-persistenced: %s", err.Error())
				}
				modules := []string{"nvidia_uvm", "nvidia_drm", "nvidia_modeset", "nvidia"}
				for _, mod := range modules {
					err = unix.DeleteModule(mod, 0)
					if err != nil {
						logger.Error("unable to unload '%s' kernel module: %s", mod, err)
					}
				}
			}

			err = gpu.PciDevice.Remove()
			if err != nil {
				panic(err)
			}
			// for _, path := range matches {
			// 	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0220)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// 	_, err = f.Write([]byte{1})
			// 	if err != nil {
			// 		logger.Error("unable to remove %s: %s", path, err)
			// 		return
			// 	}
			// }

			// todo: load kernel modules again, if a gpu requiring the driver is still connected
			// if [ $(lspci -k | grep -c ${vga_driver}) -gt 0 ]; then
			// 	modprobe ${vga_driver}
			// 	if [ ${vga_driver} = "nvidia" ]; then
			// 		modprobe nvidia_drm
			// 	fi
			// 	sleep 1
			// fi

			_, err = systemd.StartUnit(dmServiceName, "replace", nil)
			if err != nil {
				logger.Error("unable to start display-manager: %s", err)
			}

			done <- true
		}()

		// systemctl stop display-manager.service
		_, err = systemd.StopUnit(dmServiceName, "replace", nil)
		if err != nil {
			logger.Error("unable to stop display-manager: %s", err)
		}

		<-done
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCommand)
}
