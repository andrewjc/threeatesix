package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"github.com/andrewjc/threeatesix/devices/io"
	"github.com/andrewjc/threeatesix/devices/memmap"
	"log"
)

func New80386CPU() *CpuCore {

	cpuCore := &CpuCore{}
	cpuCore.partId = common.MODULE_PRIMARY_PROCESSOR

	cpuCore.registers = &CpuRegisters{}

	// index of 8 bit registers
	cpuCore.registers.registers8Bit = []*uint8{
		&cpuCore.registers.AL,
		&cpuCore.registers.CL,
		&cpuCore.registers.DL,
		&cpuCore.registers.BL,
		&cpuCore.registers.AH,
		&cpuCore.registers.CH,
		&cpuCore.registers.DH,
		&cpuCore.registers.BH,
	}

	// index of 16 bit registers
	cpuCore.registers.registers16Bit = []*uint16{
		&cpuCore.registers.AX,
		&cpuCore.registers.CX,
		&cpuCore.registers.DX,
		&cpuCore.registers.BX,
		&cpuCore.registers.SP,
		&cpuCore.registers.BP,
		&cpuCore.registers.SI,
		&cpuCore.registers.DI,
	}

	cpuCore.registers.registers32Bit = []*uint32{
		&cpuCore.registers.EAX,
		&cpuCore.registers.ECX,
		&cpuCore.registers.EDX,
		&cpuCore.registers.EBX,
		&cpuCore.registers.ESP,
		&cpuCore.registers.EBP,
		&cpuCore.registers.ESI,
		&cpuCore.registers.EDI,
	}

	cpuCore.registers.registersSegmentRegisters = []*SegmentRegister{
		&cpuCore.registers.ES,
		&cpuCore.registers.CS,
		&cpuCore.registers.SS,
		&cpuCore.registers.DS,
		&cpuCore.registers.FS,
		&cpuCore.registers.GS,
	}

	cpuCore.opCodeMap = make([]OpCodeImpl, 256)
	cpuCore.opCodeMap2Byte = make([]OpCodeImpl, 256)

	mapOpCodes(cpuCore)

	return cpuCore
}

type CpuCore struct {
	partId                 uint8
	bus                    *bus.Bus
	memoryAccessController *memmap.MemoryAccessController
	ioPortAccessController *io.IOPortAccessController

	registers      *CpuRegisters
	opCodeMap      []OpCodeImpl
	opCodeMap2Byte []OpCodeImpl

	mode  uint8
	flags CpuExecutionFlags

	busId uint32

	currentByteDecodeStart         uint32  //the start addr of the instruction being decoded (including prefixes etc)
	currentPrefixBytes             []uint8 //current prefix bytes read for the byte being decoded in the instruction
	currentByteAddr                uint32  //the current address of the byte being decoded in the current instruction
	currentOpCodeBeingExecuted     uint8   //the opcode of the instruction currently being exected
	lastExecutedInstructionPointer uint32
}

type CpuExecutionFlags struct {
	OperandSizeOverrideEnabled bool //treat operand size as 32bit
	AddressSizeOverrideEnabled bool //treat address size as 32bit

	MemorySegmentOverride uint32
	LockPrefixEnabled     bool
	RepPrefixEnabled      bool
}

func (device *CpuCore) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *CpuCore) OnReceiveMessage(message bus.BusMessage) {
	switch {
	case message.Subject == common.MESSAGE_REQUEST_CPU_MODESWITCH:
		device.EnterMode(message.Data[0])
	}
}

func (core *CpuCore) SetCS(addr uint16) {
	core.registers.CS.base = addr
}

func (core *CpuCore) SetIP(addr uint16) {
	core.registers.IP = addr
}

func (core *CpuCore) GetIP() uint16 {
	return core.registers.IP
}

func (core *CpuCore) GetCS() uint16 {
	return core.registers.CS.base
}

func (core *CpuCore) IncrementIP() {
	core.registers.IP++
}

func (core *CpuCore) Init(bus *bus.Bus) {
	core.bus = bus

	// obtain a pointer to the memory controller on the bus
	// this is a bit of a hack but avoids a linear lookup for every
	// instruction access
	dev1 := core.bus.FindSingleDevice(common.MODULE_MEMORY_ACCESS_CONTROLLER).(*memmap.MemoryAccessController)
	core.memoryAccessController = dev1

	dev2 := core.bus.FindSingleDevice(common.MODULE_IO_PORT_ACCESS_CONTROLLER).(*io.IOPortAccessController)
	core.ioPortAccessController = dev2

	core.EnterMode(common.REAL_MODE)

	core.Reset()
}

func (core *CpuCore) Reset() {
	core.registers.CS.base = 0xF000
	core.registers.IP = 0xFFF0
	core.bus.SendMessage(bus.BusMessage{Subject: common.MESSAGE_GLOBAL_LOCK_BIOS_MEM_REGION, Data: []byte{}})
}

func (core *CpuCore) EnterMode(mode uint8) {
	core.mode = mode

	core.bus.SendMessage(bus.BusMessage{Subject: common.MESSAGE_GLOBAL_CPU_MODESWITCH, Data: []byte{mode}})

	processorString := core.FriendlyPartName()
	modeString := ""
	if core.mode == common.REAL_MODE {
		modeString = "REAL MODE"
	} else if core.mode == common.PROTECTED_MODE {
		modeString = "PROTECTED MODE"
	}
	log.Printf("%s entered %s\r\n", processorString, modeString)
}

// Gets the current code segment + IP addr in memory
func (core *CpuCore) GetCurrentCodePointer() uint32 {
	addr := core.SegmentAddressToLinearAddress(core.registers.CS, core.registers.IP)
	return addr
}

func (core *CpuCore) SegmentAddressToLinearAddress(segment SegmentRegister, offset uint16) uint32 {

	if core.flags.MemorySegmentOverride > 0 {
		// default segment override
		switch core.flags.MemorySegmentOverride {
		case common.SEGMENT_CS:
			return uint32(core.registers.CS.base)<<16 + uint32(offset)
		case common.SEGMENT_SS:
			return uint32(core.registers.SS.base)<<16 + uint32(offset)
		case common.SEGMENT_DS:
			return uint32(core.registers.DS.base)<<16 + uint32(offset)
		case common.SEGMENT_ES:
			return uint32(core.registers.ES.base)<<16 + uint32(offset)
		case common.SEGMENT_FS:
			return uint32(core.registers.FS.base)<<16 + uint32(offset)
		case common.SEGMENT_GS:
			return uint32(core.registers.GS.base)<<16 + uint32(offset)
		default:
			panic("Unhandled segment register override")
		}
	}

	addr := uint32(segment.base)<<16 + uint32(offset)

	return addr
}

// Returns the address in memory of the instruction currently executing.
// This is different from GetCurrentCodePointer in that the currently executing
// instruction can update the CS and IP registers.
func (core *CpuCore) GetCurrentlyExecutingInstructionAddress() uint32 {
	return core.currentByteDecodeStart
}

func (core *CpuCore) Step() {

	register_locations_8 := core.registers.registers8Bit
	register_locations_16 := core.registers.registers16Bit
	register_locations_32 := core.registers.registers32Bit

	core.currentByteAddr = core.GetCurrentCodePointer()
	tmp := core.currentByteAddr
	if core.currentByteAddr == core.lastExecutedInstructionPointer {
		log.Fatalf("CPU appears to be in a loop! Did you forget to increment the IP register?")
	}

	core.currentByteDecodeStart = core.currentByteAddr

	status := core.decodeInstruction()

	// Check that no register pointers were overwritten
	// checks for bugs
	for i := 0; i < 8; i++ {
		if register_locations_8[i] != core.registers.registers8Bit[i] {
			panic(fmt.Sprintf("8-bit register %s was overwritten", core.registers.index8ToString(uint8(i))))
		}
		if register_locations_16[i] != core.registers.registers16Bit[i] {
			panic(fmt.Sprintf("16-bit register %s was overwritten", core.registers.index16ToString(uint8(i))))
		}
		if register_locations_32[i] != core.registers.registers32Bit[i] {
			panic(fmt.Sprintf("32-bit register %s was overwritten", core.registers.index32ToString(uint8(i))))
		}
	}

	if status != 0 {
		panic(0)
	}

	core.lastExecutedInstructionPointer = tmp

}

func (core *CpuCore) FriendlyPartName() string {
	if core.partId == common.MODULE_PRIMARY_PROCESSOR {
		return "PRIMARY PROCESSOR"
	}

	if core.partId == common.MODULE_MATH_CO_PROCESSOR {
		return "MATH CO PROCESSOR"
	}

	return "Unknown"
}

func (core *CpuCore) readImm8() (uint8, error) {
	retVal, err := core.memoryAccessController.ReadAddr8(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr++
	return retVal, nil
}

func (core *CpuCore) readImm16() (uint16, error) {
	retVal, err := core.memoryAccessController.ReadAddr16(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr += 2
	return retVal, nil
}

func (core *CpuCore) readRm8(modrm *ModRm) (*uint8, string, error) {
	if modrm.mod == 3 {
		dest := core.registers.registers8Bit[modrm.rm]
		destName := core.registers.index8ToString(modrm.rm)
		return dest, destName, nil

	} else {
		addressMode := modrm.getAddressMode16(core)
		destValue, err := core.memoryAccessController.ReadAddr8(uint32(addressMode))
		destName := fmt.Sprintf("dword_F%#04x", addressMode)
		return &destValue, destName, err
	}
}

func (core *CpuCore) readRm16(modrm *ModRm) (*uint16, string, error) {
	if modrm.mod == 3 {
		dest := core.registers.registers16Bit[modrm.rm]
		destName := core.registers.index16ToString(modrm.rm)
		return dest, destName, nil

	} else {
		addressMode := modrm.getAddressMode16(core)
		destValue, err := core.memoryAccessController.ReadAddr16(uint32(addressMode))
		destName := fmt.Sprintf("dword_F%#04x", addressMode)
		return &destValue, destName, err
	}
}

func (core *CpuCore) readR8(modrm *ModRm) (*uint8, string) {
	dest := core.registers.registers8Bit[modrm.reg]
	dstName := core.registers.index8ToString(modrm.reg)
	return dest, dstName
}

func (core *CpuCore) readR16(modrm *ModRm) (*uint16, string) {
	dest := core.registers.registers16Bit[modrm.reg]
	dstName := core.registers.index16ToString(modrm.reg)
	return dest, dstName

}

func (core *CpuCore) writeRm8(modrm *ModRm, value *uint8) error {
	if modrm.mod == 3 {
		*core.registers.registers8Bit[modrm.rm] = *value
	} else {
		addressMode := modrm.getAddressMode16(core)
		err := core.memoryAccessController.WriteAddr8(uint32(addressMode), *value)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (core *CpuCore) writeRm16(modrm *ModRm, value *uint16) error {
	if modrm.mod == 3 {
		*core.registers.registers16Bit[modrm.rm] = *value
	} else {
		addressMode := modrm.getAddressMode16(core)
		err := core.memoryAccessController.WriteAddr16(uint32(addressMode), *value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (core *CpuCore) writeR8(modrm *ModRm, value *uint8) {
	*core.registers.registers8Bit[modrm.reg] = *value
}

func (core *CpuCore) writeR16(modrm *ModRm, value *uint16) {
	*core.registers.registers16Bit[modrm.reg] = *value
}

func (core *CpuCore) SetFlag(mask uint16, status bool) {
	core.registers.SetFlag(mask, status)
}

func (core *CpuCore) GetFlag(mask uint16) bool {
	return core.registers.GetFlag(mask)
}

func (core *CpuCore) GetFlagInt(mask uint16) uint16 {
	return core.registers.GetFlagInt(mask)
}

func (core *CpuCore) GetRegisters() *CpuRegisters {
	return core.registers
}

func (core *CpuCore) writeSegmentRegister(register *SegmentRegister, value uint16) {
	if core.flags.MemorySegmentOverride > 0 {
		// default segment override
		switch core.flags.MemorySegmentOverride {
		case common.SEGMENT_CS:
			core.registers.CS.base = value
		case common.SEGMENT_SS:
			core.registers.SS.base = value
		case common.SEGMENT_DS:
			core.registers.DS.base = value
		case common.SEGMENT_ES:
			core.registers.ES.base = value
		case common.SEGMENT_FS:
			core.registers.FS.base = value
		case common.SEGMENT_GS:
			core.registers.GS.base = value
		default:
			panic("Unhandled segment register override")
		}
	}

	register.base = value

}
