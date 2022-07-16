package pci

import (
	"fmt"
	"strconv"

	"github.com/jaypipes/ghw"
)

// todo: eliminate external dependency to ghw by using sysfs directly

type IdName struct {
	id   string
	name string
}

type GPU struct {
	hexAddress           string
	driver               *string
	vendor               IdName
	product              IdName
	revision             string
	subsystem            IdName
	class                IdName
	subclass             IdName
	programmingInterface IdName
}

func (g *GPU) XorgPCIString() string {
	domain, _ := strconv.ParseInt(g.hexAddress[:4], 16, 16)
	bus, _ := strconv.ParseInt(g.hexAddress[5:7], 16, 8)
	device, _ := strconv.ParseInt(g.hexAddress[8:10], 16, 8)
	function, _ := strconv.ParseInt(g.hexAddress[11:], 16, 8)
	return fmt.Sprintf("PCI:%d@%d:%d:%d", bus, domain, device, function)
}

func (g *GPU) Identifier() string {
	return fmt.Sprintf("%s:%s:%s:%s", g.vendor.id, g.subsystem.id, g.product.id, g.revision)
}

func (g *GPU) DisplayName() string {
	return fmt.Sprintf("%s: %s", g.vendor.name, g.product.name)
}

func ReadGPUs() []*GPU {
	info, err := ghw.PCI()
	if err != nil {
		panic("unable to get pci informaticon")
	}
	var gpus []*GPU
	for _, device := range info.Devices {
		j, _ := device.MarshalJSON()
		fmt.Println(string(j))
		//if device.Class.ID == "03" && (device.Subclass.ID == "00" || device.Subclass.ID == "02") {
		dev := &GPU{
			hexAddress: device.Address,
			driver:     &device.Driver,
			vendor: IdName{
				id:   device.Vendor.ID,
				name: device.Vendor.Name,
			},
			product: IdName{
				id:   device.Product.ID,
				name: device.Product.Name,
			},
			revision: device.Revision,
			subsystem: IdName{
				id:   device.Subsystem.ID,
				name: device.Subsystem.Name,
			},
			class: IdName{
				id:   device.Class.ID,
				name: device.Class.Name,
			},
			subclass: IdName{
				id:   device.Subclass.ID,
				name: device.Subclass.Name,
			},
			programmingInterface: IdName{
				id:   device.Product.ID,
				name: device.Product.Name,
			},
		}
		gpus = append(gpus, dev)
		//}
	}
	return gpus
}
