package memmap

import "github.com/andrewjc/threeatesix/common"

type RealModeAccessProvider struct {
    *MemoryAccessController
}

func (r *RealModeAccessProvider) ReadAddr8(addr uint32) (uint8, error) {

    var byteData uint8
    inBiosSpace := (addr > uint32(len(*r.biosImage)-1)-(0xFFFFFFF-addr)) && addr < 0xFFFFFFFF
    if r.resetVectorBaseAddr > 0 && inBiosSpace {
        return r.ReadFromBiosAddressSpace(addr)
    } else {
        if int(addr) > len(*r.backingRam) || addr < 0 {
            return 0, common.GeneralProtectionFault{}
        }
        byteData = (*r.backingRam)[addr]
    }

    return byteData, nil
}

func (r *RealModeAccessProvider) ReadAddr16(addr uint32) (uint16, error) {
    b1, err := r.ReadAddr8(addr)
    if err != nil {
        return 0, err
    }

    b2, err2 := r.ReadAddr8(addr + 1)
    if err2 != nil {
        return 0, err2
    }

    return uint16(b2)<<8 | uint16(b1), nil

}

func (r *RealModeAccessProvider) ReadAddr32(addr uint32) (uint32, error) {

    b1, err := r.ReadAddr8(addr)
    if err != nil {
        return 0, err
    }

    b2, err2 := r.ReadAddr8(addr + 1)
    if err2 != nil {
        return 0, err2
    }

    b3, err3 := r.ReadAddr8(addr + 2)
    if err3 != nil {
        return 0, err3
    }

    b4, err4 := r.ReadAddr8(addr + 3)
    if err4 != nil {
        return 0, err4
    }

    return uint32(b4)<<24 | uint32(b3)<<16 | uint32(b2)<<8 | uint32(b1), nil
}

func (r *RealModeAccessProvider) PeekNextBytesImpl(addr uint32, numBytes uint32) []byte {
    buffer := make([]byte, numBytes)

    for i := uint32(0); i < numBytes; i++ {

        if r.resetVectorBaseAddr > 0 {
            buffer[i], _ = r.ReadFromBiosAddressSpace(addr + i)
        } else {
            buffer[i] = (*r.backingRam)[addr+i]
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
