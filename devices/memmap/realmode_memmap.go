package memmap

import "github.com/andrewjc/threeatesix/common"

type RealModeAccessProvider struct {
	*MemoryAccessController
}

func (r *RealModeAccessProvider) ReadMemoryAddr8(addr uint32) (uint8, error) {
	// Define BIOS memory range constants
	const BIOS_START uint32 = 0xF0000
	const BIOS_END uint32 = 0xFFFFF // Adjust based on actual BIOS size

	// Check if the address is within the BIOS mapped space
	inBiosSpace := addr >= BIOS_START && addr <= BIOS_END

	if inBiosSpace {
		// Calculate index into BIOS data by subtracting the start of the BIOS address space
		biosIndex := addr - BIOS_START
		if int(biosIndex) >= len(*r.biosImage) {
			return 0, common.GeneralProtectionFault{} // or another appropriate error
		}
		return (*r.biosImage)[biosIndex], nil
	} else {
		// Ensure address is within RAM bounds
		if int(addr) >= len(*r.backingRam) || addr < 0 {
			return 0, common.GeneralProtectionFault{}
		}
		return (*r.backingRam)[addr], nil
	}
}

func (r *RealModeAccessProvider) ReadMemoryAddr16(addr uint32) (uint16, error) {
	b1, err := r.ReadMemoryAddr8(addr)
	if err != nil {
		return 0, err
	}

	b2, err2 := r.ReadMemoryAddr8(addr + 1)
	if err2 != nil {
		return 0, err2
	}

	return uint16(b2)<<8 | uint16(b1), nil

}

func (r *RealModeAccessProvider) ReadMemoryAddr32(addr uint32) (uint32, error) {

	b1, err := r.ReadMemoryAddr8(addr)
	if err != nil {
		return 0, err
	}

	b2, err2 := r.ReadMemoryAddr8(addr + 1)
	if err2 != nil {
		return 0, err2
	}

	b3, err3 := r.ReadMemoryAddr8(addr + 2)
	if err3 != nil {
		return 0, err3
	}

	b4, err4 := r.ReadMemoryAddr8(addr + 3)
	if err4 != nil {
		return 0, err4
	}

	return uint32(b4)<<24 | uint32(b3)<<16 | uint32(b2)<<8 | uint32(b1), nil
}

func (r *RealModeAccessProvider) PeekNextBytesImpl(addr uint32, numBytes uint32) []byte {
	buffer := make([]byte, numBytes)

	for i := uint32(0); i < numBytes; i++ {

		iBuff, err := r.ReadMemoryAddr8(addr)
		if err != nil {
			break
		} else {
			buffer[i] = iBuff
		}
	}

	return buffer
}

func (r *RealModeAccessProvider) ReadFromBiosAddressSpace(addr uint32) (uint8, error) {

	ddd := 0xF000FFFF - addr

	biosImageLength := uint32(len(*r.biosImage) - 1)

	offs := biosImageLength - ddd

	if int(offs) > len(*r.biosImage) || offs < 0 {
		return 0, common.GeneralProtectionFault{}
	}
	byteData := (*r.biosImage)[offs]

	return byteData, nil
}
