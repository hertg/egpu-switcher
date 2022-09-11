package cmd

import (
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
	"github.com/pmorjan/kmod"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCommand = &cobra.Command{
	Use:   "remove",
	Short: "[root required][experimental] Remove GPU and restart display manager",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !isRoot {
			return fmt.Errorf("you need root privileges to remove egpu")
		}

		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("please note that this is an experimental feature, continue? (y/N): ")
		// answer, _ := reader.ReadString('\n')
		// if answer != "y\n" && answer != "Y\n" {
		// 	os.Exit(0)
		// }

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
		k, err := kmod.New()
		if err != nil {
			return err
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("goroutine panicked: %+v", r)
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

			// check if same module still required after removal
			// this needs to be done before unloading, otherwise
			// we can't get the information from sysfs which pci
			// devices also use the same module
			gpus := pci.ReadGPUs()
			driverUsed := 0
			for _, gpu := range gpus {
				if *gpu.PciDevice.Driver == driver {
					driverUsed += 1
					if driverUsed > 1 {
						break
					}
				}
			}
			loadModAgain := driverUsed > 1

			// unload module
			if driver == "nvidia" {
				// systemctl stop nvidia-persistenced.service
				err := init.StopUnit(ctx, "nvidia-persistenced.service")
				if err != nil {
					errChan <- err
				}
				modules := []string{"nvidia_uvm", "nvidia_drm", "nvidia_modeset", "nvidia"}
				for _, mod := range modules {
					if err := unloadMod(k, mod); err != nil {
						panic(err)
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
			if loadModAgain {
				if err := loadMod(k, driver); err != nil {
					errChan <- err
				}
				if driver == "nvidia" {
					if err := loadMod(k, "nvidia_drm"); err != nil {
						errChan <- err
					}
				}
				<-time.After(1 * time.Second)
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

		select {
		case err := <-errChan:
			if err != nil {
				logger.Error("got error: %s", err)
			} else {
				logger.Success("removal was finished")
			}
			return err
		case <-time.After(10 * time.Second):
			return fmt.Errorf("exiting due to timeout")
		}
	},
}

func loadMod(k *kmod.Kmod, name string) error {
	logger.Debug("attempting to load module '%s'...", name)
	if err := k.Load(name, "", 0); err != nil {
		logger.Error("loading module '%s' failed: %s", name, err)
		return err
	}
	logger.Debug("loading module '%s' successful", name)
	return nil
}

func unloadMod(k *kmod.Kmod, name string) error {
	logger.Debug("attempting to unload module '%s'...", name)
	if err := k.Unload(name); err != nil {
		logger.Error("unloading module '%s' failed: %s", name, err)
		return err
	}
	logger.Debug("unloading module '%s' successful", name)
	return nil
}

func init() {
	rootCmd.AddCommand(removeCommand)
}
