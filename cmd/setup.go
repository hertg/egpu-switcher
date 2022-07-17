package cmd

import (
	"fmt"
	"time"

	"github.com/hertg/egpu-switcher/internal/pci"
	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

var setupCommand = &cobra.Command{
	Use:   "setup",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		gpus := pci.ReadGPUs()
		for _, gpu := range gpus {
			fmt.Printf("%+v\n", gpu)
			fmt.Println(gpu.DisplayName())
			fmt.Println(gpu.XorgPCIString())
			fmt.Println("------")
		}

		a1 := time.Now()
		ghw.PCI()
		a2 := time.Now()
		am := a2.UnixMicro() - a1.UnixMicro()
		fmt.Printf("running ghw.PCI() takes %dµs\n", am)

		b1 := time.Now()
		fmt.Println(pci.Find(1153611719250962689))
		b2 := time.Now()
		fmt.Println(pci.Find(1153611719250962690))
		b3 := time.Now()

		bm1 := b2.UnixMicro() - b1.UnixMicro()
		bm2 := b3.UnixMicro() - b2.UnixMicro()

		fmt.Printf("searching for a present device took %dµs\n", bm1)
		fmt.Printf("searching for an absent device took %dµs\n", bm2)

		c1 := time.Now()
		for i := 0; i < 100; i++ {
			_ = pci.Find(1153611719250962689)
		}
		c2 := time.Now()
		cm := c2.UnixMicro() - c1.UnixMicro()
		fmt.Printf("searching for a present device 100 times took %dµs\n", cm)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
