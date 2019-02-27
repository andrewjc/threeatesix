package common

import "log"

/*
	Memory interconnect - provides memory access between cpu and ram
*/

const (
	REAL_MODE = iota
	PROTECTED_MODE
)

type MemoryAccessController struct {
	backingRam *[]byte
	biosImage  *[]byte

	cpuRegisterController CpuRegisterController
	memoryAccessProvider  MemoryAccessProvider

	resetVectorBaseAddr uint32
	//used during initial boot. this simulates the 'hack' that motherboards do
}

type CpuRegisterController interface {
	GetIP() uint16
	GetCS() uint16
	IncrementIP()
	SetIP(addr uint16)
	SetCS(addr uint16)
}

func (mem *MemoryAccessController) SetCpuRegisterController(controller CpuRegisterController) {
	mem.cpuRegisterController = controller
}

func (mem *MemoryAccessController) GetNextInstruction() interface{} {
	return mem.memoryAccessProvider.ReadNextInstruction()
}

func (mem *MemoryAccessController) ReadAddr8(address uint16) uint8 {
	return mem.memoryAccessProvider.ReadAddr8(address)
}

func (mem *MemoryAccessController) ReadAddr16(address uint16) uint16 {
	return mem.memoryAccessProvider.ReadAddr16(address)
}

func (mem *MemoryAccessController) ReadAddr32(address uint16) uint32 {
	return mem.memoryAccessProvider.ReadAddr32(address)
}

func (mem *MemoryAccessController) WriteAddr(address uint16, value uint8) {
	(*mem.backingRam)[address] = value
}

func (mem *MemoryAccessController) LockBootVector() {
	mem.resetVectorBaseAddr = 0xFFFF0000
}

type MemoryAccessProvider interface {

	// Read the next instruction at CS:IP
	ReadNextInstruction() interface{}

	ReadAddr8(u uint16) uint8
	ReadAddr16(u uint16) uint16
	ReadAddr32(u uint16) uint32

}

type RealModeAccessProvider struct {
	*MemoryAccessController
}

func (r *RealModeAccessProvider) ReadAddr8(addr uint16) uint8 {

	var byteData uint8
	if r.resetVectorBaseAddr > 0 {
		byteData = (*r.biosImage)[addr]
	} else {
		byteData = (*r.backingRam)[addr]
	}

	return byteData
}

func (r *RealModeAccessProvider) ReadAddr16(addr uint16) uint16 {

	b1 := uint16(r.ReadAddr8(addr))
	b2 := uint16(r.ReadAddr8(addr + 1))
	return b2<<8 | b1

}

func (r *RealModeAccessProvider) ReadAddr32(addr uint16) uint32 {

	b1 := uint32(r.ReadAddr16(addr))
	b2 := uint32(r.ReadAddr16(addr + 1))
	return b2<<16 | b1

}

func (r *RealModeAccessProvider) ReadNextInstruction() interface{} {

	ip := (r.cpuRegisterController).GetIP()
	cs := r.cpuRegisterController.GetCS()

	addr := uint32(uint32(cs)<<4 + uint32(ip))

	// bios address space
	var byteData uint8
	if r.resetVectorBaseAddr > 0 {
		addr = uint32(ip)
		byteData = (*r.biosImage)[addr]

		log.Printf("[BIOS MAP] Reading instruction at %#4x: next byte: %#2x\n", addr, byteData)

	} else {
		byteData = (*r.backingRam)[addr]
		log.Printf("[RAM MAP] Reading instruction at %#4x: next byte: %#2x\n", addr, byteData)
	}

	return byteData
}

func (mem *MemoryAccessController) EnterMode(mode uint8) {
	switch {
	case mode == REAL_MODE:
		mem.memoryAccessProvider = &RealModeAccessProvider{mem}
	}
}

func (mem *MemoryAccessController) SetIP(addr uint16) {
	mem.cpuRegisterController.SetIP(addr)
}

func (mem *MemoryAccessController) SetCS(addr uint16) {
	mem.cpuRegisterController.SetCS(addr)
}

func (mem *MemoryAccessController) PeekNextBytes(numBytes int) []byte {

	buffer := make([]byte, numBytes)

	for i := 0; i < numBytes; i++ {

		ip := (mem.cpuRegisterController).GetIP() + uint16(i)
		cs := mem.cpuRegisterController.GetCS()
		addr := uint32(uint32(cs)<<4 + uint32(ip))

		if mem.resetVectorBaseAddr > 0 {
			addr = uint32(ip)
			buffer[i] = (*mem.biosImage)[addr]
		} else {
			buffer[i] = (*mem.backingRam)[addr]
		}
	}

	return buffer
}

func CreateMemoryController(ram *[]byte, bios *[]byte) *MemoryAccessController {
	return &MemoryAccessController{ram, bios, nil, nil, 0}
}
