package common

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
	GetEIP() uint32
	IncrementEIP()

	GetIP() uint16
	GetCS() uint16
	IncrementIP()
	SetIP(addr uint16)
}

func (mem *MemoryAccessController) SetCpuRegisterController(controller CpuRegisterController) {
	mem.cpuRegisterController = controller
}

func (mem *MemoryAccessController) ReadHalfWord() interface{} {
	return mem.memoryAccessProvider.ReadByteCode()
}

func (mem *MemoryAccessController) ReadNextWord() interface{} {
	return mem.memoryAccessProvider.ReadNextWordCode()
}

func (mem *MemoryAccessController) ReadNextDWord() interface{} {
	return mem.memoryAccessProvider.ReadNextDWordCode()
}

func (mem *MemoryAccessController) ReadAddr(address uint16) uint8 {
	return (*mem.backingRam)[address]
}

func (mem *MemoryAccessController) WriteAddr(address uint16, value uint8) {
	(*mem.backingRam)[address] = value
}

func (mem *MemoryAccessController) LockBootVector() {
	mem.resetVectorBaseAddr = 0xFFFF0000
}

type MemoryAccessProvider interface {

	// Access functions for reading instructions from code segment
	ReadByteCode() interface{}
	ReadNextWordCode() interface{}
	ReadNextDWordCode() interface{}
}

type RealModeAccessProvider struct {
	*MemoryAccessController
}

func (r *RealModeAccessProvider) ReadByteCode() interface{} {

	ip := (r.cpuRegisterController).GetIP()
	cs := r.cpuRegisterController.GetCS()

	addr := uint32(uint32(cs)<<4 + uint32(ip))

	// simulate memory mapped address space

	// bios address space
	var byteData byte
	if r.resetVectorBaseAddr > 0 {
		addr = uint32(ip)
		byteData = (*r.biosImage)[addr]
	} else {
		byteData = (*r.backingRam)[addr]
	}

	(r.cpuRegisterController).IncrementIP()

	return byteData
}

func (r *RealModeAccessProvider) ReadNextWordCode() interface{} {
	b1 := uint16(r.ReadHalfWord().(uint8))
	b2 := uint16(r.ReadHalfWord().(uint8))
	return b2<<8 | b1
}

func (r *RealModeAccessProvider) ReadNextDWordCode() interface{} {
	b1 := uint32(r.ReadNextWord().(uint16))
	b2 := uint32(r.ReadNextWord().(uint16))
	return b2<<16 | b1
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

func CreateCpuMemoryInterconnect(ram *[]byte, bios *[]byte) *MemoryAccessController {
	return &MemoryAccessController{ram, bios, nil, nil, 0}
}
