package memmap

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
)

/*
	Memory interconnect - provides memory access between intel8086 and ram
*/

type MemoryAccessController struct {
	backingRam *[]byte
	biosImage  *[]byte

	memoryAccessProvider MemoryAccessProvider

	resetVectorBaseAddr uint32

	bus                 *bus.Bus
	busId               uint32
}


type MemoryAccessProvider interface {
	ReadAddr8(u uint16) uint8
	ReadAddr16(u uint16) uint16
	ReadAddr32(u uint16) uint32
	PeekNextBytesImpl(numBytes int) []byte
}

func (mem *MemoryAccessController) SetDeviceBusId(id uint32) {
	mem.busId = id
}

func (mem *MemoryAccessController) OnReceiveMessage(message bus.BusMessage) {
	switch {
	case message.Subject == common.MESSAGE_GLOBAL_LOCK_BIOS_MEM_REGION:
		mem.LockBootVector()
	case message.Subject == common.MESSAGE_GLOBAL_UNLOCK_BIOS_MEM_REGION:
		mem.UnlockBootVector()
	case message.Subject == common.MESSAGE_GLOBAL_CPU_MODESWITCH:
		mem.HandleMemoryMapSwitch(message.Data[0])
	}
}


func CreateMemoryController(ram *[]byte, bios *[]byte) *MemoryAccessController {
	return &MemoryAccessController{ram, bios, nil, 0, nil, 0}
}

func (mem *MemoryAccessController) HandleMemoryMapSwitch(modeSwitch byte) {
	switch {
	case modeSwitch == common.REAL_MODE:
		mem.memoryAccessProvider = &RealModeAccessProvider{mem}
	}
}

func (controller *MemoryAccessController) SetBus(bus *bus.Bus) {
	controller.bus = bus
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
func (mem *MemoryAccessController) UnlockBootVector() {
	mem.resetVectorBaseAddr = 0x0
}


func (mem *MemoryAccessController) PeekNextBytes(numBytes int) []byte {
	return mem.memoryAccessProvider.PeekNextBytesImpl(numBytes)
}
