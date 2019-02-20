package common

/*
	Memory interconnect - provides memory access between cpu and ram
 */


type CpuMemInterconnect struct {
	backingRam *[]byte
	cpuRegisterController CpuRegisterController
}

type CpuRegisterController interface {
	IncrementEIP()
	GetEIP() uint32
}

func (mem *CpuMemInterconnect) SetCpuRegisterController(controller CpuRegisterController) {
	mem.cpuRegisterController = controller
}

func (mem *CpuMemInterconnect) ReadHalfWord() uint8 {
	loc := (*mem.backingRam)[(mem.cpuRegisterController).GetEIP()]

	(mem.cpuRegisterController).IncrementEIP()

	return loc
}

func (mem *CpuMemInterconnect) ReadNextWord() uint16 {
	b1 := uint16(mem.ReadHalfWord())
	b2 := uint16(mem.ReadHalfWord())
	return b2<<8 | b1
}

func (mem *CpuMemInterconnect) ReadNextDWord() uint32 {
	b1 := uint32(mem.ReadNextWord())
	b2 := uint32(mem.ReadNextWord())
	return b2<<16 | b1
}

func (mem *CpuMemInterconnect) ReadAddr(address uint16) uint8 {
	return (*mem.backingRam)[address]
}

func (mem *CpuMemInterconnect) WriteAddr(address uint16, value uint8) {
	(*mem.backingRam)[address] = value
}

func CreateCpuMemoryInterconnect(ram *[]byte) *CpuMemInterconnect {
	return &CpuMemInterconnect{ram, nil}
}