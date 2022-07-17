package pci

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/pcidb"
)

// todo: eliminate external dependency to ghw by using sysfs directly

var db *pcidb.PCIDB
var once sync.Once

type IdName struct {
	id   string
	name string
}

type GPU struct {
	hexAddress    string
	id            uint64
	vendorId      uint16
	vendorName    string
	deviceId      uint16
	deviceName    string
	classId       uint16
	className     string
	subclassId    uint16
	subclassName  string
	subvendorId   uint16
	subvendorName string
	subdeviceId   uint16
	subdeviceName string
	revision      uint8
}

func (g *GPU) XorgPCIString() string {
	domain, _ := strconv.ParseInt(g.hexAddress[:4], 16, 16)
	bus, _ := strconv.ParseInt(g.hexAddress[5:7], 16, 8)
	device, _ := strconv.ParseInt(g.hexAddress[8:10], 16, 8)
	function, _ := strconv.ParseInt(g.hexAddress[11:], 16, 8)
	return fmt.Sprintf("PCI:%d@%d:%d:%d", bus, domain, device, function)
}

func (g *GPU) Identifier() uint64 {
	return uint64(g.vendorId)<<48 | uint64(g.deviceId)<<32 | uint64(g.subvendorId)<<16 | uint64(g.subdeviceId)
}

func (g *GPU) DisplayName() string {
	bold := color.New(color.Bold).SprintFunc()
	return fmt.Sprintf("\t%s (rev %02x)\n\t%s (%s)", bold(g.deviceName), g.revision, g.vendorName, g.subvendorName)
}

func (g *GPU) LspciDisplayName() string {
	return fmt.Sprintf("%s %s: %s (%s) %s (rev %02x)", g.hexAddress, g.subclassName, g.vendorName, g.subvendorName, g.deviceName, g.revision)
}

func (g *GPU) GuessDriver() (string, error) {
	if contains(g.vendorName, "nvidia") {
		return "nvidia", nil
	} else if contains(g.vendorName, "amd") {
		return "amdgpu", nil
	} else if contains(g.vendorName, "intel") {
		return "intel", nil // ?
	}
	return "", fmt.Errorf("unable to guess driver for vendor %s", g.vendorName)
}

func contains(str string, substr string) bool {
	return strings.Contains(strings.ToLower(str), substr)
}

func getPCIDB() *pcidb.PCIDB {
	once.Do(func() {
		var err error
		db, err = pcidb.New()
		if err != nil {
			panic("unable to get pcidb")
		}
	})
	return db
}

func ReadGPUs() []*GPU {
	pci, err := ghw.PCI()
	if err != nil {
		panic("unable to get pci informaticon")
	}
	pdb := getPCIDB()
	var gpus []*GPU
	for _, d := range pci.Devices {
		if d.Class.ID == "03" && (d.Subclass.ID == "00" || d.Subclass.ID == "02") {
			vendorId, _ := strconv.ParseUint(d.Vendor.ID, 16, 16)
			deviceId, _ := strconv.ParseUint(d.Product.ID, 16, 16)
			classId, _ := strconv.ParseUint(d.Class.ID, 16, 16)
			subclassId, _ := strconv.ParseUint(d.Subclass.ID, 16, 16)
			subvendorId, _ := strconv.ParseUint(d.Subsystem.VendorID, 16, 16)
			subvendorName := ""
			subdeviceId, _ := strconv.ParseUint(d.Subsystem.ID, 16, 16)
			revision, _ := strconv.ParseUint(strings.TrimPrefix(d.Revision, "0x"), 16, 8)
			if subvendorId != 0 {
				if vinfo, ok := pdb.Vendors[d.Subsystem.VendorID]; ok {
					subvendorName = vinfo.Name
				}
			}
			dev := &GPU{
				hexAddress:    d.Address,
				id:            vendorId<<48 | deviceId<<32 | subvendorId<<16 | subdeviceId,
				vendorId:      uint16(vendorId),
				vendorName:    d.Vendor.Name,
				deviceId:      uint16(deviceId),
				deviceName:    d.Product.Name,
				classId:       uint16(classId),
				className:     d.Class.Name,
				subclassId:    uint16(subclassId),
				subclassName:  d.Subclass.Name,
				subvendorId:   uint16(subvendorId),
				subvendorName: subvendorName,
				subdeviceId:   uint16(subdeviceId),
				subdeviceName: d.Subsystem.Name,
				revision:      uint8(revision),
			}
			gpus = append(gpus, dev)
		}
	}
	return gpus
}

func IsPresent(id uint64) bool {
	gpus := ReadGPUs()
	for _, gpu := range gpus {
		if gpu.Identifier() == id {
			return true
		}
	}
	return false
}
