package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/hertg/egpu-switcher/internal/service"
	"github.com/hertg/egpu-switcher/internal/xorg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "[root required][experimental] Remove GPU and restart display manager",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !isRoot {
			return fmt.Errorf("you need root privileges to remove egpu")
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("please note that this is an experimental feature, continue? (y/N): ")
		answer, _ := reader.ReadString('\n')
		if answer != "y\n" && answer != "Y\n" {
			os.Exit(0)
		}

		ctx := context.Background()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		id := uint64(viper.GetInt("egpu.id"))
		driver := viper.GetString("egpu.driver")
		gpu := pci.Find(id)
		if gpu == nil {
			return fmt.Errorf("the egpu is not connected")
		}

		init, err := service.GetInitSystem()
		if err != nil {
			return err
		}

		err = xorg.RemoveEgpuFile(x11ConfPath, verbose)
		if err != nil {
			return err
		}

		errChan := make(chan error)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("goroutine panicked: %+v", r)
					// logger.Error("goroutine panicked: %+v", r)
					// done <- true
				}
			}()

			// wait on signal
			sig := <-sigChan
			logger.Debug("got signal: %s", sig)

			// block until display manager is stopped
			for {
				stopped, err := init.IsDisplayManagerStopped(ctx)
				if err != nil {
					errChan <- err
				}
				if stopped {
					break
				}
				<-time.After(500 * time.Millisecond)
			}

			logger.Debug("display-manager has become inactive")

			if driver == "nvidia" {
				// systemctl stop nvidia-persistenced.service
				err := init.StopUnit(ctx, "nvidia-persistenced.service")
				if err != nil {
					errChan <- err
				}
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
				errChan <- err
			}
			logger.Info("pci device was removed")

			// load kernel modules again, if another gpu requires the same driver
			gpus := pci.ReadGPUs()
			for _, gpu := range gpus {
				if *gpu.PciDevice.Driver == driver {
					// another gpu still requires the now unloaded driver, so reload it here

					// modprobe ${driver}
					// if err := modprobe.Load(driver, ""); err != nil {
					// 	errChan <- err
					// }
					if driver == "nvidia" {
						// modprobe nvidia_drm
						// if err := modprobe.Load("nvidia_drm", ""); err != nil {
						// 	errChan <- err
						// }
					}
					// sleep 1s
					<-time.After(1 * time.Second)
					break
				}
			}

			logger.Info("starting display manager again")
			err = init.StartDisplayManager(ctx)
			if err != nil {
				logger.Error("unable to start display-manager: %s", err)
			}

			errChan <- nil
		}()

		// systemctl stop display-manager.service
		err = init.StopDisplayManager(ctx)
		if err != nil {
			logger.Error("unable to stop display-manager: %s", err)
		}

		return <-errChan
	},
}

func init() {
	rootCmd.AddCommand(removeCommand)
}
