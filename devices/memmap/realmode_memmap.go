package memmap

import (
	"encoding/binary"
	"errors"
	"github.com/andrewjc/threeatesix/common"
)

type RealModeAccessProvider struct {
	*MemoryAccessController
}

func (r *RealModeAccessProvider) ReadMemoryAddr8(addr uint32) (*uint8, error) {
	// Define BIOS memory range constants
	const BIOS_START uint32 = 0xF0000
	const BIOS_END uint32 = 0xFFFFF // Adjust based on actual BIOS size

	// Check if the address is within the BIOS mapped space
	inBiosSpace := addr >= BIOS_START && addr <= BIOS_END

	if inBiosSpace {
		// Calculate index into BIOS data by subtracting the start of the BIOS address space
		biosIndex := addr - BIOS_START

		if int(biosIndex) >= len(*r.biosImage) {
			return nil, common.GeneralProtectionFault{} // or another appropriate error
		}
		return &(*r.biosImage)[biosIndex], nil
	} else {
		// Ensure address is within RAM bounds
		if int(addr) >= len(*r.backingRam) || addr < 0 {
			return nil, common.GeneralProtectionFault{}
		}
		return &(*r.backingRam)[addr], nil
	}
}

func (r *RealModeAccessProvider) ReadMemoryAddr16(addr uint32) (*uint16, error) {
	offset := uint32(addr & 0xFFFF)

	base := uint32((addr - uint32(offset)) >> 4)

	if offset&0x1 == 0 { // Check if the address is aligned to a 16-bit boundary
		addr = uint32(base<<4 + (offset & 0xffff))
		bArray := r.ReadSequential(addr, 2)
		if len(bArray) < 2 {
			return nil, errors.New("not enough data read")
		}
		bRes := make([]byte, 2)
		bRes[0] = *bArray[0]
		bRes[1] = *bArray[1]
		retVal := binary.LittleEndian.Uint16(bRes)
		return &retVal, nil
	} else { // Address is not aligned
		addr = uint32(base<<4 + (offset & 0xffff))
		b1, err1 := r.ReadMemoryAddr8(addr)
		if err1 != nil {
			return nil, err1
		}
		addr = uint32(base<<4 + ((offset + 1) & 0xffff))
		b2, err2 := r.ReadMemoryAddr8(addr)
		if err2 != nil {
			return nil, err2
		}
		retVal := uint16(*b1) | uint16(*b2)<<8
		return &retVal, nil
	}
}

func (r *RealModeAccessProvider) ReadMemoryAddr16A(addr uint32) (*uint16, error) {
	offset := uint16(addr & 0xFFFF)

	base := uint16((addr - uint32(offset)) >> 4)

	if offset&0x1 == 0 { // Check if the address is aligned to a 16-bit boundary
		addr = uint32(base<<4 + (offset & 0xffff))
		bArray := r.ReadSequential(addr, 2)
		if len(bArray) < 2 {
			return nil, errors.New("not enough data read")
		}
		bRes := make([]byte, 2)
		bRes[0] = *bArray[0]
		bRes[1] = *bArray[1]
		retVal := binary.LittleEndian.Uint16(bRes)
		return &retVal, nil
	} else { // Address is not aligned
		addr = uint32(base<<4 + (offset & 0xffff))
		b1, err1 := r.ReadMemoryAddr8(addr)
		if err1 != nil {
			return nil, err1
		}
		addr = uint32((base<<4 + ((offset + 1) & 0xffff)) << 8)
		b2, err2 := r.ReadMemoryAddr8(addr)
		if err2 != nil {
			return nil, err2
		}
		retVal := uint16(*b1) | uint16(*b2)<<8
		return &retVal, nil
	}
}

func (r *RealModeAccessProvider) ReadMemoryAddr32(addr uint32) (*uint32, error) {
	bArray := r.ReadSequential(addr, 4)
	if len(bArray) < 4 {
		return nil, errors.New("not enough data read")
	}
	bRes := make([]byte, 4)
	bRes[0] = *bArray[0]
	bRes[1] = *bArray[1]
	bRes[2] = *bArray[2]
	bRes[3] = *bArray[3]
	retVal := binary.LittleEndian.Uint32(bRes)
	return &retVal, nil
}

func (r *RealModeAccessProvider) ReadSequential(addr uint32, numBytes uint32) []*uint8 {
	buffer := make([]*uint8, numBytes)

	for i := uint32(0); i < numBytes; i++ {

		iBuff, err := r.ReadMemoryAddr8(addr + i)
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

func (r *RealModeAccessProvider) WriteMemoryAddr8(addr uint32, value uint8) error {
	offset := addr & 0xF
	base := (addr - offset) >> 4

	address := base<<4 + (offset & 0xffff)

	if int(address) > len(*r.backingRam) || address < 0 {
		return common.GeneralProtectionFault{}
	}

	(*r.backingRam)[address] = value & 0xff

	return nil
}
