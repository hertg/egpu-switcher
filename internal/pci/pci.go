package pci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/hertg/gopci/pkg/pci"
)

// todo: eliminate external dependency to ghw by using sysfs directly

type IdName struct {
	id   string
	name string
}

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

func (g *GPU) Outputs() {
	// glob := filepath.Join("drm", "card*", "card*-*", "status")
}

func (g *GPU) DisplayName() string {
	bold := color.New(color.Bold).SprintFunc()
	deviceName := g.PciDevice.Product.Label
	if g.PciDevice.Subdevice != nil {
		deviceName = g.PciDevice.Subdevice.Label
	}
	return fmt.Sprintf(
		"\t%s (rev %02x)\n\t%s (%s)",
		bold(deviceName),
		g.PciDevice.Config.Revision(),
		g.PciDevice.Vendor.Label,
		g.PciDevice.Subvendor.Label,
	)
}

func (g *GPU) HasDisplaysConnected() (bool, error) {
	pattern := fmt.Sprintf("%s/drm/card[0-9]*/card[0-9]*-*/status", g.PciDevice.SysfsPath())
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false, err
	}
	for _, path := range matches {
		contents, err := os.ReadFile(path)
		if err != nil {
			return false, err
		}
		if strings.TrimSuffix(string(contents), "\n") == "connected" {
			return true, nil
		}
	}
	return false, nil
}

func (g *GPU) NumberOfDisplays() (int, error) {
	pattern := fmt.Sprintf("%s/drm/card[0-9]*/card[0-9]*-*/status", g.PciDevice.SysfsPath())
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0, err
	}
	return len(matches), nil
}

//func (g *GPU) LspciDisplayName() string {
//	return fmt.Sprintf("%s %s: %s (%s) %s (rev %02x)", g.hexAddress, g.subclassName, g.vendorName, g.subvendorName, g.deviceName, g.revision)
//}

func contains(str string, substr string) bool {
	return strings.Contains(strings.ToLower(str), substr)
}

func ReadGPUs() []*GPU {
	gpuFilter := func(device *pci.Device) bool {
		return device.Class.Class() == 0x03
	}
	devices, err := pci.Scan(gpuFilter)
	if err != nil {
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
