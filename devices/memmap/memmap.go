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

	bus   *bus.Bus
	busId uint32

	segmentOverride     uint32
	addressSizeOverride bool
	operandSizeOverride bool
	setLockPrefix       bool
	setRepPrefix        bool
}

type MemoryAccessProvider interface {
	ReadMemoryAddr8(u uint32) (uint8, error)
	ReadMemoryAddr16(u uint32) (uint16, error)
	ReadMemoryAddr32(u uint32) (uint32, error)
	PeekNextBytesImpl(addr uint32, numBytes uint32) []byte
}

func NewMemoryController(ram *[]byte, bios *[]byte) *MemoryAccessController {
	return &MemoryAccessController{ram, bios, nil, 0, nil, 0, 0, false, false, false, false}
}

func (mem *MemoryAccessController) GetDeviceBusId() uint32 {
	return mem.busId
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

func (mem *MemoryAccessController) GetPortMap() *bus.DevicePortMap {
	return nil
}

func (controller *MemoryAccessController) ReadAddr8(addr uint16) uint8 {
	//TODO implement me
	panic("implement me")
}

func (mem *MemoryAccessController) WriteAddr8(addr uint16, data uint8) {
	//TODO implement me
	panic("implement me")
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

func (mem *MemoryAccessController) ReadMemoryAddr8(address uint32) (uint8, error) {

	return mem.memoryAccessProvider.ReadMemoryAddr8(address)
}

func (mem *MemoryAccessController) ReadMemoryAddr16(address uint32) (uint16, error) {
	return mem.memoryAccessProvider.ReadMemoryAddr16(address)
}

func (mem *MemoryAccessController) ReadMemoryAddr32(address uint32) (uint32, error) {
	return mem.memoryAccessProvider.ReadMemoryAddr32(address)
}

func (mem *MemoryAccessController) WriteMemoryAddr8(address uint32, value uint8) error {
	if int(address) > len(*mem.backingRam) || address < 0 {
		return common.GeneralProtectionFault{}
	}

	(*mem.backingRam)[address] = value

	return nil
}

func (mem *MemoryAccessController) WriteMemoryAddr16(address uint32, value uint16) error {
	for i := uint32(0); i < 2; i++ {
		err := mem.WriteMemoryAddr8(address+i, uint8(value>>uint32(i*8)&0xFF))
		if err != nil {
			return err
		}
	}

	return nil
}

func (mem *MemoryAccessController) WriteMemoryAddr32(address uint32, value uint32) error {
	for i := uint32(0); i < 4; i++ {
		err := mem.WriteMemoryAddr8(address+i, uint8(value>>uint32(i*8)&0xFF))
		if err != nil {
			return err
		}
	}

	return nil
}

func (mem *MemoryAccessController) LockBootVector() {
	mem.resetVectorBaseAddr = 0xFFFF0000
}
func (mem *MemoryAccessController) UnlockBootVector() {
	mem.resetVectorBaseAddr = 0x0
}

func (mem *MemoryAccessController) PeekNextBytes(addr uint32, numBytes uint32) []byte {
	return mem.memoryAccessProvider.PeekNextBytesImpl(addr, numBytes)
}

func (mem *MemoryAccessController) SetSegmentOverride(override uint32) {
	mem.segmentOverride = override
}

func (mem *MemoryAccessController) SetAddressSizeOverride(enabled bool) {
	mem.addressSizeOverride = enabled
}

func (mem *MemoryAccessController) SetOperandSizeOverride(enabled bool) {
	mem.operandSizeOverride = enabled
}

func (mem *MemoryAccessController) SetLockPrefix(enabled bool) {
	mem.setLockPrefix = enabled
}

func (mem *MemoryAccessController) SetRepPrefix(enabled bool) {
	mem.setRepPrefix = enabled
}
