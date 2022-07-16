package pci

/*type PciAddreess struct {
	address [4]byte
}

func (s *PciAddreess) Domain() [2]byte {
	var domain [2]byte
	copy(domain[:], s.address[:2])
	return domain
}

func (s *PciAddreess) Bus() byte {
	return s.address[2]
}

func (s *PciAddreess) Device() byte {
	return s.address[3] >> 3
}

func (s *PciAddreess) Function() byte {
	return s.address[3] & 0b00000111
}

func (s *PciAddreess) HexString() string {
	return fmt.Sprintf("%4x:%2x:%2x.%x", s.Domain(), s.Bus(), s.Device(), s.Function())
}
*/
