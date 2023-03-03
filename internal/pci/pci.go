package pci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/hertg/egpu-switcher/internal/logger"
	"github.com/hertg/gopci/pkg/pci"
)

type GPU struct {
	PciDevice *pci.Device
}

func (g *GPU) XorgPCIString() string {
	return fmt.Sprintf(
		"PCI:%d@%d:%d:%d",
		g.PciDevice.Address.Bus(),
		g.PciDevice.Address.Domain(),
		g.PciDevice.Address.Device(),
		g.PciDevice.Address.Function(),
	)
}

func (g *GPU) Identifier() uint64 {
	subven := uint64(0)
	if g.PciDevice.Subvendor != nil {
		subven = uint64(g.PciDevice.Subvendor.ID)
	}
	subdev := uint64(0)
	if g.PciDevice.Subdevice != nil {
		subdev = uint64(g.PciDevice.Subdevice.ID)
	}
	return uint64(g.PciDevice.Vendor.ID)<<48 | uint64(g.PciDevice.Product.ID)<<32 | subven<<16 | subdev
}

func (g *GPU) DisplayName() string {
	bold := color.New(color.Bold).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	str := fmt.Sprintf("\t%s %s", bold(g.PciDevice.Vendor.Label), g.PciDevice.Product.Label)
	if sv := g.PciDevice.Subvendor; sv != nil && !strings.HasPrefix(sv.Label, "Subvendor") {
		if sd := g.PciDevice.Subdevice; sd != nil && !strings.HasPrefix(sd.Label, "Subdevice") {
			str = fmt.Sprintf("\t%s %s", bold(sv.Label), sd.Label)
		}
	}
	if d := g.PciDevice.Driver; d != nil {
		str = fmt.Sprintf("%s (%s)", str, cyan(*d))
	} else {
		str = fmt.Sprintf("%s (%s)", str, red("unknown"))
	}
	return str
}

func (g *GPU) ConnectedOutputs() (uint, error) {
	pattern := fmt.Sprintf("%s/drm/card[0-9]*/card[0-9]*-*/status", g.PciDevice.SysfsPath())
	matches, err := filepath.Glob(pattern)
	num := uint(0)
	if err != nil {
		return num, err
	}
	for _, path := range matches {
		contents, err := os.ReadFile(path)
		if err != nil {
			return num, err
		}
		if strings.HasPrefix(string(contents), "connected") {
			num += 1
		}
	}
	return num, nil
}

func (g *GPU) Outputs() (int, error) {
	pattern := fmt.Sprintf("%s/drm/card[0-9]*/card[0-9]*-*/status", g.PciDevice.SysfsPath())
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, err
	}
	return len(matches), nil
}

func contains(str string, substr string) bool {
	return strings.Contains(strings.ToLower(str), substr)
}

func ReadGPUs() []*GPU {
	gpuFilter := func(device *pci.Device) bool {
		return device.Class.Class() == 0x03
		// return true
	}
	devices, err := pci.Scan(gpuFilter)
	if err != nil {
		logger.Error("unable to read pci information from sysfs: %s", err)
		panic("unable to read pci information from sysfs")
	}
	var gpus []*GPU
	for _, device := range devices {
		gpus = append(gpus, &GPU{PciDevice: device})
	}
	return gpus
}

func Find(id uint64) *GPU {
	gpus := ReadGPUs()
	for _, gpu := range gpus {
		if gpu.Identifier() == id {
			return gpu
		}
	}
	return nil
}
