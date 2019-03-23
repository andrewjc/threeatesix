package memmap

import (
	"github.com/andrewjc/threeatesix/common"
	cpu2 "github.com/andrewjc/threeatesix/devices/cpu"
)

type RealModeAccessProvider struct {
	*MemoryAccessController
}

func (r *RealModeAccessProvider) ReadAddr8(addr uint32) uint8 {

	var byteData uint8
	if r.resetVectorBaseAddr > 0 {
		byteData = (*r.biosImage)[addr]
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

func (r *RealModeAccessProvider) PeekNextBytesImpl(numBytes int) []byte {
	buffer := make([]byte, numBytes)

	cpu := r.bus.FindSingleDevice(common.MODULE_PRIMARY_PROCESSOR).(cpu2.CpuController)
	cs := cpu.GetCS()
	ip := cpu.GetIP()

	for i := 0; i < numBytes; i++ {

		ip := ip + uint16(i)
		addr := uint32(uint16(cs)<<4 + uint16(ip))

		if r.resetVectorBaseAddr > 0 {
			addr = uint32(ip)
			buffer[i] = (*r.biosImage)[addr]
		} else {
			buffer[i] = (*r.backingRam)[addr]
		}
	}

	return buffer
}