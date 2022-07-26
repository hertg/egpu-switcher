package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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

		prefix, _, found := strings.Cut(gpu.Address(), ".")
		if !found {
			return fmt.Errorf("unable to get device id from pci address %s", gpu.Address())
		}
		pattern := fmt.Sprintf("/sys/bus/pci/devices/%s.[0-9]*/remove", prefix)
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}

		systemd, err := dbus.NewSystemConnection()
		if err != nil {
			return fmt.Errorf("unable to connect to dbus")
		}
		defer systemd.Close()

		go func() {
			sig := <-sigChan
			_ = sig
			//fmt.Println(sig)

			// todo: wait until display-manager is no longer 'active'
			dm, err := os.Readlink("/etc/systemd/system/display-manager.service")
			dm = filepath.Base(dm)
			//dm = strings.TrimSuffix(dm, ".service")
			fmt.Printf("%+v\n", dm)
			// todo: Requires systemd v230 or higher. (https://github.com/systemd/systemd/releases/tag/v230)

			for {
				status, err := systemd.ListUnitsByNamesContext(context.Background(), []string{dm})
				if err != nil {
					panic("unable to get display manager status information")
				}
				if len(status) != 1 {
					panic("expected to find only one display manager service")
				}
				if status[0].ActiveState != "active" {
					break
				}
				<-time.After(500 * time.Millisecond)
			}

			fmt.Println("display manager is no longer active")

			done <- true
			return

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
					// todo: modprobe -r $mod
				}
			}

			for _, path := range matches {
				f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0220)
				if err != nil {
					panic(err)
				}
				// todo: enable
				/*_, err = f.Write([]byte{1})
				if err != nil {
					return err
				}*/
				fmt.Println(f)
			}

			// systemctl start display-manager.service
			_, err = systemd.StartUnit("display-manager", "replace", nil)
			if err != nil {
				logger.Error("unable to start display-manager: %s", err)
			}

			done <- true
		}()

		// systemctl stop display-manager.service
		fmt.Println("todo: stop display-manager")
		/*_, err = systemd.StopUnit("display-manager", "replace", nil)
		if err != nil {
			logger.Error("unable to stop display-manager: %s", err)
		}*/

		<-done
		fmt.Println("exiting")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCommand)
}
