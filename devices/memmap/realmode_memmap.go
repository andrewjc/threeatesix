package memmap

type RealModeAccessProvider struct {
	*MemoryAccessController
}

func (r *RealModeAccessProvider) ReadAddr8(addr uint32) uint8 {

	var byteData uint8

	inBiosSpace := (addr > uint32(len(*r.biosImage)-1) - (0xF000FFFF - addr)) && addr < 0xF000FFFF
	if r.resetVectorBaseAddr > 0 && inBiosSpace {
		return r.ReadFromBiosAddressSpace(addr)
	} else {
		byteData = (*r.backingRam)[addr]
	}

	return byteData
}

func (r *RealModeAccessProvider) ReadAddr16(addr uint32) uint16 {

	b1 := uint16(r.ReadAddr8(addr))
	b2 := uint16(r.ReadAddr8(addr + 1))
	return b2<<8 | b1

}

func (r *RealModeAccessProvider) ReadAddr32(addr uint32) uint32 {

	b1 := uint32(r.ReadAddr16(addr))
	b2 := uint32(r.ReadAddr16(addr + 1))
	return b2<<16 | b1

}

func (r *RealModeAccessProvider) PeekNextBytesImpl(addr uint32, numBytes uint32) []byte {
	buffer := make([]byte, numBytes)

	for i := uint32(0); i < numBytes; i++ {

		if r.resetVectorBaseAddr > 0 {
			buffer[i] = r.ReadFromBiosAddressSpace(addr+i)
		} else {
			buffer[i] = (*r.backingRam)[addr+i]
		}
	}

	return buffer
}

func (r *RealModeAccessProvider) ReadFromBiosAddressSpace(addr uint32) uint8 {

	ddd := 0xF000FFFF - addr

	biosImageLength := uint32(len(*r.biosImage)-1)

	offs := biosImageLength - ddd

	byteData := (*r.biosImage)[offs]

	return byteData
}