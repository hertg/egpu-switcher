package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/hertg/egpu-switcher/internal/xorg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

var removeCommand = &cobra.Command{
	Use: "remove",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !isRoot {
			return fmt.Errorf("you need root privileges to remove egpu")
		}

		ctx := context.Background()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		done := make(chan bool, 1)

		id := uint64(viper.GetInt("egpu.id"))
		driver := viper.GetString("egpu.driver")
		gpu := pci.Find(id)
		if gpu == nil {
			return fmt.Errorf("the egpu is not connected")
		}

		systemd, err := dbus.NewSystemdConnectionContext(ctx)
		if err != nil {
			return fmt.Errorf("unable to connect to dbus")
		}
		defer systemd.Close()

		err = xorg.RemoveEgpuFile(x11ConfPath, verbose)
		if err != nil {
			return err
		}

		dmServiceName, err := os.Readlink("/etc/systemd/system/display-manager.service")
		dmServiceName = filepath.Base(dmServiceName)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("goroutine panicked: %+v", r)
					done <- true
				}
			}()

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
				ch := make(chan string)
				_, err := systemd.StopUnitContext(ctx, "nvidia-persistenced", "replace", ch)
				if err != nil {
					logger.Error("unable to stop nvidia-persistenced: %s", err.Error())
				}
				// logger.Debug("got result from stopping service: %s", <-ch)
				modules := []string{"nvidia_uvm", "nvidia_drm", "nvidia_modeset", "nvidia"}
				for _, mod := range modules {
					logger.Info("unload kernel module: %s", mod)
					err = unix.DeleteModule(mod, 0)
					if err != nil {
						logger.Error("unable to unload '%s' kernel module: %s", mod, err)
					}
				}
			}

			logger.Info("removing pci device...")
			err = gpu.PciDevice.Remove()
			if err != nil {
				logger.Error("unable to remove pci device: %s", err)
				panic(err)
			}

			// load kernel modules again, if a gpu requiring the driver is still connected
			gpus := pci.ReadGPUs()
			for _, gpu := range gpus {
				if *gpu.PciDevice.Driver == driver {
					// another gpu still requires the now unloaded driver, so reload it here
					// TODO modprobe ${driver}
					if driver == "nvidia" {
						// TODO: modprobe nvidia_drm
					}
					// TODO: sleep 1s ?
					break
				}
			}

			logger.Info("starting %s again", dmServiceName)
			_, err = systemd.StartUnitContext(ctx, dmServiceName, "replace", nil)
			if err != nil {
				logger.Error("unable to start display-manager: %s", err)
			}

			done <- true
		}()

		// systemctl stop display-manager.service
		_, err = systemd.StopUnitContext(ctx, dmServiceName, "replace", nil)
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
