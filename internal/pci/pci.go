package pci

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

func Test() {
	info, err := ghw.GPU()
	cobra.CheckErr(err)
	fmt.Println(info.JSONString(true))

	for _, card := range info.GraphicsCards {
		fmt.Println(card.Address)
		fmt.Println(card.DeviceInfo.Product.Name)
		fmt.Println(card.DeviceInfo.ProgrammingInterface.Name)
	}

	/*p, err := ghw.PCI()
	cobra.CheckErr(err)
	fmt.Println(p.JSONString(true))*/

	/*topo, err := ghw.Topology()
	cobra.CheckErr(err)
	fmt.Println(topo.JSONString(true))*/

}
