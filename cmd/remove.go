package cmd

// TODO: I was not able to get this to work properly.
// So this whole feature is disabled for now, because
// I can't support something that doesn't work on my machine.

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"os/exec"
// 	"os/signal"
// 	"strings"
// 	"syscall"
// 	"time"

// 	"github.com/hertg/egpu-switcher/internal/logger"
// 	"github.com/hertg/egpu-switcher/internal/pci"
// 	"github.com/hertg/egpu-switcher/internal/service"
// 	"github.com/hertg/egpu-switcher/internal/xorg"
// 	"github.com/pmorjan/kmod"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// 	"github.com/ulikunitz/xz"
// 	"golang.org/x/sys/unix"
// )

// var removeCommand = &cobra.Command{
// 	Use:    "remove",
// 	Short:  "Remove GPU and restart display manager [experimental]",
// 	Hidden: true, // TODO
// 	RunE: func(cmd *cobra.Command, args []string) error {

// 		if !isRoot {
// 			return fmt.Errorf("you need root privileges to remove egpu")
// 		}

// 		reader := bufio.NewReader(os.Stdin)
// 		fmt.Print("please note that this is an experimental feature, continue? (y/N): ")
// 		answer, _ := reader.ReadString('\n')
// 		if answer != "y\n" && answer != "Y\n" {
// 			os.Exit(0)
// 		}

// 		ctx := context.Background()

// 		sigChan := make(chan os.Signal, 1)
// 		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

// 		id := uint64(viper.GetInt("egpu.id"))
// 		driver := viper.GetString("egpu.driver")
// 		gpu := pci.Find(id)
// 		if gpu == nil {
// 			return fmt.Errorf("the egpu is not connected")
// 		}

// 		init, err := service.GetInitSystem()
// 		if err != nil {
// 			return err
// 		}

// 		err = xorg.RemoveEgpuFile(x11ConfPath, verbose)
// 		if err != nil {
// 			return err
// 		}

// 		errChan := make(chan error)
// 		k, err := kmod.New(kmod.SetInitFunc(modInit))
// 		if err != nil {
// 			return err
// 		}

// 		go func() {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					errChan <- fmt.Errorf("goroutine panicked: %+v", r)
// 				}
// 			}()

// 			// wait on signal
// 			sig := <-sigChan
// 			logger.Debug("got signal: %s", sig)

// 			// block until display manager is stopped
// 			for {
// 				stopped, err := init.IsDisplayManagerStopped(ctx)
// 				if err != nil {
// 					errChan <- err
// 				}
// 				if stopped {
// 					break
// 				}
// 				<-time.After(500 * time.Millisecond)
// 			}

// 			logger.Debug("display-manager has become inactive")

// 			// check if same module still required after removal
// 			// this needs to be done before unloading, otherwise
// 			// we can't get the information from sysfs which pci
// 			// devices also use the same module
// 			gpus := pci.ReadGPUs()
// 			driverUsed := 0
// 			for _, gpu := range gpus {
// 				if *gpu.PciDevice.Driver == driver {
// 					driverUsed += 1
// 					if driverUsed > 1 {
// 						break
// 					}
// 				}
// 			}
// 			loadModAgain := driverUsed > 1

// 			if driver == "nvidia" {
// 				err := init.StopUnit(ctx, "nvidia-persistenced.service")
// 				if err != nil {
// 					errChan <- err
// 				}
// 				modules := []string{"nvidia_uvm", "nvidia_drm", "nvidia_modeset", "nvidia"}
// 				for _, mod := range modules {
// 					if err := unloadMod(k, mod); err != nil {
// 						panic(err)
// 					}
// 				}
// 			}

// 			logger.Info("removing pci device...")
// 			err = gpu.PciDevice.Remove()
// 			if err != nil {
// 				logger.Error("unable to remove pci device: %s", err)
// 				errChan <- err
// 			}
// 			logger.Info("pci device was removed")

// 			// load kernel modules again, if another gpu requires the same driver
// 			if loadModAgain {
// 				if err := loadMod(k, driver); err != nil {
// 					errChan <- err
// 				}
// 				if driver == "nvidia" {
// 					if err := loadMod(k, "nvidia_drm"); err != nil {
// 						errChan <- err
// 					}
// 				}
// 				<-time.After(1 * time.Second)
// 			}

// 			logger.Info("starting display manager again")
// 			err = init.StartDisplayManager(ctx)
// 			if err != nil {
// 				logger.Error("unable to start display-manager: %s", err)
// 			}

// 			errChan <- nil
// 		}()

// 		err = init.StopDisplayManager(ctx)
// 		if err != nil {
// 			logger.Error("unable to stop display-manager: %s", err)
// 		}

// 		select {
// 		case err := <-errChan:
// 			if err != nil {
// 				logger.Error("got error: %s", err)
// 			} else {
// 				logger.Success("removal was finished")
// 			}
// 			return err
// 		case <-time.After(10 * time.Second):
// 			return fmt.Errorf("exiting due to timeout")
// 		}
// 	},
// }

// func modInit(path, params string, flags int) error {
// 	logger.Debug("try to init module at %s", path)
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	var b []byte
// 	if strings.HasSuffix(f.Name(), ".xz") {
// 		reader, err := xz.NewReader(f)
// 		if err != nil {
// 			return err
// 		}
// 		// NOTE: It seems to hang on the next line
// 		b, err = ioutil.ReadAll(reader)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		b, err = ioutil.ReadAll(f)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return unix.InitModule(b, params)
// }

// func loadMod(k *kmod.Kmod, name string) error {
// 	logger.Debug("attempting to load module '%s'...", name)

// 	// NOTE: Couldn't get that to work, it would block on ioutil.ReadAll(xzreader)
// 	// if err := k.Load(name, "", 0); err != nil {

// 	// TODO: At least this seems to work better, but it requries modprobe
// 	cmd := exec.Command("modprobe", name)
// 	if err := cmd.Run(); err != nil {
// 		logger.Error("loading module '%s' failed: %s", name, err)
// 		return err
// 	}
// 	logger.Debug("loading module '%s' successful", name)
// 	return nil
// }

// func unloadMod(k *kmod.Kmod, name string) error {
// 	logger.Debug("attempting to unload module '%s'...", name)

// 	// NOTE: This actually worked, but when to Load() counterpart
// 	// doesn't there's not much of a point using it
// 	// if err := k.Unload(name); err != nil {

// 	// TODO: At least this seems to work better, but it requires modprobe
// 	cmd := exec.Command("modprobe", "-r", name)
// 	if err := cmd.Run(); err != nil {
// 		logger.Error("unloading module '%s' failed: %s", name, err)
// 		return err
// 	}
// 	logger.Debug("unloading module '%s' successful", name)
// 	return nil
// }

// func init() {
// 	rootCmd.AddCommand(removeCommand)
// }
