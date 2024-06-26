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

	initializeRegisters(cpuCore)

	initializeSegmentRegisters(cpuCore)

	cpuCore.registers.IP = 0x0000 // Instruction pointer set to the start
	cpuCore.registers.SP = 0xFFFE // Stack pointer set to the top of the stack

	cpuCore.opCodeMap = make([]OpCodeImpl, 256)
	cpuCore.opCodeMap2Byte = make([]OpCodeImpl, 256)

	mapOpCodes(cpuCore)

	return cpuCore
}

func initializeSegmentRegisters(cpuCore *CpuCore) {
	cpuCore.registers.ES = SegmentRegister{Base: 0, Limit: 0xFFFF}
	cpuCore.registers.CS = SegmentRegister{Base: 0xFFFF0, Limit: 0xFFFF}
	cpuCore.registers.SS = SegmentRegister{Base: 0, Limit: 0xFFFF}
	cpuCore.registers.DS = SegmentRegister{Base: 0, Limit: 0xFFFF}
	cpuCore.registers.FS = SegmentRegister{Base: 0, Limit: 0xFFFF}
	cpuCore.registers.GS = SegmentRegister{Base: 0, Limit: 0xFFFF}
}

func initializeRegisters(cpuCore *CpuCore) {
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
}

type CpuCore struct {
	bus                    *bus.Bus
	partId                 uint8
	memoryAccessController *memmap.MemoryAccessController
	ioPortAccessController *io.IOPortAccessController

	registers      *CpuRegisters
	opCodeMap      []OpCodeImpl
	opCodeMap2Byte []OpCodeImpl

	mode  uint8
	flags CpuExecutionFlags

	busId            uint32
	interruptChannel chan bus.BusMessage

	currentByteDecodeStart         uint32  //the start addr of the instruction being decoded (including prefixes etc)
	currentPrefixBytes             []uint8 //current prefix bytes read for the byte being decoded in the instruction
	currentByteAddr                uint32  //the current address of the byte being decoded in the current instruction
	currentOpCodeBeingExecuted     uint8   //the opcode of the instruction currently being exected
	lastExecutedInstructionPointer uint32
	is2ByteOperand                 bool
	halt                           bool
}

type CpuExecutionFlags struct {
	OperandSizeOverrideEnabled bool //treat operand size as 32bit
	AddressSizeOverrideEnabled bool //treat address size as 32bit

	MemorySegmentOverride uint32
	LockPrefixEnabled     bool
	RepPrefixEnabled      bool
}

func (device *CpuCore) GetDeviceBusId() uint32 {
	return device.busId
}

func (device *CpuCore) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *CpuCore) SetBus(bus *bus.Bus) {
	device.bus = bus
}

func (device *CpuCore) OnReceiveMessage(message bus.BusMessage) {
	switch {
	case message.Subject == common.MESSAGE_REQUEST_CPU_MODESWITCH:
		device.EnterMode(message.Data[0])
	case message.Subject == common.MESSAGE_INTERRUPT_RAISE:
		device.AcknowledgeInterrupt(message)
	case message.Subject == common.MESSAGE_INTERRUPT_EXECUTE:
		device.HandleInterrupt(message)
	}
}

func (core *CpuCore) GetPortMap() *bus.DevicePortMap {
	return nil
}

func (core *CpuCore) ReadAddr8(addr uint16) uint8 {
	//TODO implement me
	panic("implement me")
}

func (core *CpuCore) WriteAddr8(addr uint16, data uint8) {
	//TODO implement me
	panic("implement me")
}

func (core *CpuCore) SetCS(addr uint32) {
	core.registers.CS.Base = addr
}

func (core *CpuCore) SetIP(addr uint16) {
	core.registers.IP = addr
}

func (core *CpuCore) GetIP() uint16 {
	return core.registers.IP
}

func (core *CpuCore) GetCS() uint32 {
	return core.registers.CS.Base
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
	core.registers.CS.Base = 0xF000 // This sets the base of the code segment to 0xF0000 when multiplied by 16.
	core.registers.IP = 0xFFF0      // Instruction pointer set to 0xFFF0.
	core.registers.CR0 = 0          // Set to real mode
	core.registers.FLAGS = 0x0002   // Set default flags
	core.bus.SendMessage(bus.BusMessage{Subject: common.MESSAGE_GLOBAL_LOCK_BIOS_MEM_REGION, Data: []byte{}})

	core.shadowBios()
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
	core.logDebug(fmt.Sprintf("%s entered %s", processorString, modeString))
}

// Gets the current code segment + IP addr in memory
func (core *CpuCore) GetCurrentCodePointer() uint32 {
	addr := core.SegmentAddressToLinearAddress(core.registers.CS, core.registers.IP)
	return addr
}

func (core *CpuCore) GetCurrentDataPointer() uint32 {
	addr := core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.IP)
	return addr
}

func (core *CpuCore) GetCurrentSegmentWidth() uint8 {
	segment := core.registers.registersSegmentRegisters[common.SEGMENT_CS]

	if core.flags.MemorySegmentOverride > 0 {
		segment = core.registers.registersSegmentRegisters[core.flags.MemorySegmentOverride]
	}

	if segment.Limit == 0xFFFF {
		return 16
	} else {
		return 32
	}
}

func (core *CpuCore) SegmentAddressToLinearAddress(segment SegmentRegister, offset uint16) uint32 {
	addr := core.SegmentAddressToLinearAddress_NoMask(segment, offset)
	linearAddr := addr & 0xFFFFF // Mask to 20 bits to simulate real-mode address wrapping.

	return linearAddr
}

func (core *CpuCore) SegmentAddressToLinearAddress_NoMask(segment SegmentRegister, offset uint16) uint32 {

	if core.flags.MemorySegmentOverride > 0 {
		// default segment override
		switch core.flags.MemorySegmentOverride {
		case common.SEGMENT_CS:
			return uint32(core.registers.CS.Base)<<4 + uint32(offset)
		case common.SEGMENT_SS:
			return uint32(core.registers.SS.Base)<<4 + uint32(offset)
		case common.SEGMENT_DS:
			return uint32(core.registers.DS.Base)<<4 + uint32(offset)
		case common.SEGMENT_ES:
			return uint32(core.registers.ES.Base)<<4 + uint32(offset)
		case common.SEGMENT_FS:
			return uint32(core.registers.FS.Base)<<4 + uint32(offset)
		case common.SEGMENT_GS:
			return uint32(core.registers.GS.Base)<<4 + uint32(offset)
		default:
			panic("Unhandled segment register override")
		}
	}

	addr := uint32(segment.Base)<<4 + uint32(offset)

	return addr
}

func (core *CpuCore) SegmentAddressToLinearAddress32(segment SegmentRegister, offset uint32) uint32 {

	// Determine the base address of the segment
	if core.flags.MemorySegmentOverride > 0 {
		switch core.flags.MemorySegmentOverride {
		case common.SEGMENT_CS:
			return uint32(core.registers.CS.Base)<<4 + uint32(offset)
		case common.SEGMENT_SS:
			return uint32(core.registers.SS.Base)<<4 + uint32(offset)
		case common.SEGMENT_DS:
			return uint32(core.registers.DS.Base)<<4 + uint32(offset)
		case common.SEGMENT_ES:
			return uint32(core.registers.ES.Base)<<4 + uint32(offset)
		case common.SEGMENT_FS:
			return uint32(core.registers.FS.Base)<<4 + uint32(offset)
		case common.SEGMENT_GS:
			return uint32(core.registers.GS.Base)<<4 + uint32(offset)
		default:
			panic("Unhandled segment register override")
		}
	}

	addr := uint32(segment.Base)<<4 + uint32(offset)

	return addr
}

// Returns the address in memory of the instruction currently executing.
// This is different from GetCurrentCodePointer in that the currently executing
// instruction can update the CS and IP registers.
func (core *CpuCore) GetCurrentlyExecutingInstructionAddress() uint32 {
	return core.currentByteDecodeStart
}

func (core *CpuCore) Step() {

	core.currentByteAddr = core.GetCurrentCodePointer()
	tmp := core.currentByteAddr
	if core.currentByteAddr == core.lastExecutedInstructionPointer {
		//log.Fatalf("CPU appears to be in a loop! Did you forget to increment the IP register?")
		//doCoreDump(core)
		//core.logInstruction("looping...")
	}

	core.currentByteDecodeStart = core.currentByteAddr

	status := core.decodeInstruction()

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
func (core *CpuCore) getEffectiveAddress8(modrm *ModRm) (uint32, string) {
	var addr uint32
	var addrDesc string

	// Compute the effective address based on the mod and rm fields
	switch modrm.mod {
	case 0:
		if modrm.rm == 6 { // Special case for 16-bit displacement when rm = 6
			addr = uint32(modrm.disp16)
			addrDesc = fmt.Sprintf("Disp16 [0x%X]", modrm.disp16)
		} else {
			addr = *core.registers.registers32Bit[modrm.rm] // Using 32-bit register for simplification in handling
			addrDesc = fmt.Sprintf("[%s]", core.registers.index8ToString(modrm.rm))
		}
	case 1: // 8-bit displacement added to the register
		addr = *core.registers.registers32Bit[modrm.rm] + uint32(int32(int8(modrm.disp8)))
		addrDesc = fmt.Sprintf("[%s + 0x%X]", core.registers.index8ToString(modrm.rm), int8(modrm.disp8))
	case 2: // 16-bit displacement added to the register
		addr = *core.registers.registers32Bit[modrm.rm] + uint32(modrm.disp16)
		addrDesc = fmt.Sprintf("[%s + 0x%X]", core.registers.index8ToString(modrm.rm), modrm.disp16)
	case 3: // Register direct mode, no memory addressing involved
		addr = 0 // No effective memory address
		addrDesc = fmt.Sprintf("Register %s directly", core.registers.index8ToString(modrm.rm))
	}

	return addr, addrDesc
}

func (core *CpuCore) getEffectiveAddress16(modrm *ModRm) (uint16, string) {
	var addr uint16
	var addrDesc string

	// Check for SIB byte usage
	if modrm.mod != 3 && modrm.rm == 4 { // SIB byte is present
		baseAddr := *core.registers.registers16Bit[modrm.base]
		indexValue := *core.registers.registers16Bit[modrm.index] * (1 << modrm.scale)

		// Compute the effective address based on mod value
		switch modrm.mod {
		case 0:
			if modrm.base == 5 { // Special case, base is absent, only displacement
				addr = modrm.disp16
				addrDesc = fmt.Sprintf("Disp32 [0x%X]", modrm.disp16) // Bug: Should be Disp16, not Disp32
			} else {
				addr = baseAddr + indexValue
				addrDesc = fmt.Sprintf("[%s + %s*%d]", core.registers.index16ToString(modrm.base), core.registers.index16ToString(modrm.index), 1<<modrm.scale)
			}
		case 1:
			addr = baseAddr + indexValue + uint16(int16(int8(modrm.disp8)))
			addrDesc = fmt.Sprintf("[%s + %s*%d + 0x%X]", core.registers.index16ToString(modrm.base), core.registers.index16ToString(modrm.index), 1<<modrm.scale, int8(modrm.disp8))
		case 2:
			addr = baseAddr + indexValue + modrm.disp16
			addrDesc = fmt.Sprintf("[%s + %s*%d + 0x%X]", core.registers.index16ToString(modrm.base), core.registers.index16ToString(modrm.index), 1<<modrm.scale, modrm.disp32) // Bug: Should use modrm.disp16
		}
	} else {
		// No SIB byte, direct register addressing or displacement
		switch modrm.mod {
		case 0:
			if modrm.rm == 6 { // Displacement-only addressing
				addr = modrm.disp16
				addrDesc = fmt.Sprintf("Disp16 [0x%X]", modrm.disp16)
			} else {
				addr = *core.registers.registers16Bit[modrm.rm]
				addrDesc = fmt.Sprintf("[%s]", core.registers.index16ToString(modrm.rm))
			}
		case 1:
			addr = *core.registers.registers16Bit[modrm.rm] + uint16(int16(int8(modrm.disp8)))
			addrDesc = fmt.Sprintf("[%s + 0x%X]", core.registers.index16ToString(modrm.rm), int8(modrm.disp8))
		case 2:
			addr = *core.registers.registers16Bit[modrm.rm] + modrm.disp16
			addrDesc = fmt.Sprintf("[%s + 0x%X]", core.registers.index16ToString(modrm.rm), modrm.disp16)
		}
	}

	return addr, addrDesc
}

func (core *CpuCore) getEffectiveAddress32(modrm *ModRm) (uint32, string) {
	var addr uint32
	var addrDesc string

	// Check for SIB byte usage
	if modrm.mod != 3 && modrm.rm == 4 { // SIB byte is present
		baseAddr := *core.registers.registers32Bit[modrm.base]
		indexValue := *core.registers.registers32Bit[modrm.index] * (1 << modrm.scale)

		// Compute the effective address based on mod value
		switch modrm.mod {
		case 0:
			if modrm.base == 5 { // Special case, base is absent, only displacement
				addr = modrm.disp32
				addrDesc = fmt.Sprintf("Disp32 [0x%X]", modrm.disp32)
			} else {
				addr = baseAddr + indexValue
				addrDesc = fmt.Sprintf("[%s + %s*%d]", core.registers.index32ToString(modrm.base), core.registers.index32ToString(modrm.index), 1<<modrm.scale)
			}
		case 1:
			addr = baseAddr + indexValue + uint32(int32(int8(modrm.disp8)))
			addrDesc = fmt.Sprintf("[%s + %s*%d + 0x%X]", core.registers.index32ToString(modrm.base), core.registers.index32ToString(modrm.index), 1<<modrm.scale, int8(modrm.disp8))
		case 2:
			addr = baseAddr + indexValue + modrm.disp32
			addrDesc = fmt.Sprintf("[%s + %s*%d + 0x%X]", core.registers.index32ToString(modrm.base), core.registers.index32ToString(modrm.index), 1<<modrm.scale, modrm.disp32)
		}
	} else {
		// No SIB byte, direct register addressing or displacement
		switch modrm.mod {
		case 0:
			if modrm.rm == 5 { // Displacement-only addressing
				addr = modrm.disp32
				addrDesc = fmt.Sprintf("Disp32 [0x%X]", modrm.disp32)
			} else {
				addr = *core.registers.registers32Bit[modrm.rm]
				addrDesc = fmt.Sprintf("[%s]", core.registers.index32ToString(modrm.rm))
			}
		case 1:
			addr = *core.registers.registers32Bit[modrm.rm] + uint32(int32(int8(modrm.disp8)))
			addrDesc = fmt.Sprintf("[%s + 0x%X]", core.registers.index32ToString(modrm.rm), int8(modrm.disp8))
		case 2:
			addr = *core.registers.registers32Bit[modrm.rm] + modrm.disp32
			addrDesc = fmt.Sprintf("[%s + 0x%X]", core.registers.index32ToString(modrm.rm), modrm.disp32)
		}
	}

	return addr, addrDesc
}

func (core *CpuCore) readImm8() (uint8, error) {
	retVal, err := core.memoryAccessController.ReadMemoryValue8(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr++
	return retVal, nil
}

func (core *CpuCore) readImm16() (uint16, error) {
	retVal, err := core.memoryAccessController.ReadMemoryValue16(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr += 2
	return retVal, nil
}

func (core *CpuCore) readImm32() (uint32, error) {
	retVal, err := core.memoryAccessController.ReadMemoryValue32(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr += 4
	return retVal, nil
}

func (core *CpuCore) readRm8(modrm *ModRm) (*uint8, string, error) {
	if modrm.mod == 3 {
		// Directly accessing the register when Mod = 3
		dest := core.registers.registers8Bit[modrm.rm]
		destName := core.registers.index8ToString(modrm.rm)
		return dest, destName, nil
	} else {
		// Calculating the effective address when accessing memory
		addr, addrDesc := core.getEffectiveAddress8(modrm)

		// Reading the value from memory using the calculated address
		destValue, err := core.memoryAccessController.ReadMemoryValue8(addr)
		if err != nil {
			return nil, "", err // Properly return nil and error if the read fails
		}

		// Updating the address description to be more informative
		destName := fmt.Sprintf("Memory at %s", addrDesc)
		return &destValue, destName, nil
	}
}

func (core *CpuCore) readRm16(modrm *ModRm) (*uint16, string, error) {
	if modrm.mod == 3 {
		// Direct register access when Mod is 3
		dest := core.registers.registers16Bit[modrm.rm] // Direct pointer to the register
		destName := core.registers.index16ToString(modrm.rm)
		return dest, destName, nil
	} else {
		// Calculate the effective address when accessing memory
		addr, addrDesc := core.getEffectiveAddress32(modrm)

		// Read the 16-bit value from memory
		destValue, err := core.memoryAccessController.ReadMemoryValue16(addr)
		if err != nil {
			return nil, "", err // Handle error by returning nil and the error
		}

		// Update the destination name to be more descriptive
		return &destValue, addrDesc, nil
	}
}

func (core *CpuCore) readRm32(modrm *ModRm) (*uint32, string, error) {
	if modrm.mod == 3 {
		// Direct register access when Mod is 3
		dest := core.registers.registers32Bit[modrm.rm] // Direct pointer to the register
		destName := core.registers.index32ToString(modrm.rm)
		return dest, destName, nil
	} else {
		// Calculate the effective address when accessing memory
		addr, addrDesc := core.getEffectiveAddress32(modrm)

		// Read the 16-bit value from memory
		destValue, err := core.memoryAccessController.ReadMemoryValue32(addr)
		if err != nil {
			return nil, "", err // Handle error by returning nil and the error
		}

		// Update the destination name to be more descriptive
		destName := fmt.Sprintf("Memory at %s", addrDesc)
		return &destValue, destName, nil
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

func (core *CpuCore) writeRm8(modrm *ModRm, value *uint8) (string, error) {
	var destName string
	var addr uint32
	if modrm.mod == 3 {
		// Direct register access when Mod is 3
		*core.registers.registers8Bit[modrm.rm] = *value
	} else {
		// Calculate the effective address when accessing memory
		addr, destName = core.getEffectiveAddress32(modrm) // Discard the address description here as it's not used

		// Write the 8-bit value to memory
		err := core.memoryAccessController.WriteMemoryAddr8(uint32(addr), *value)
		if err != nil {
			return destName, err // Properly return the error
		}
	}

	return destName, nil // Return nil explicitly when no error occurs
}

func (core *CpuCore) writeRm16(modrm *ModRm, value *uint16) (string, error) {
	var destName string
	var addr uint32
	if modrm.mod == 3 {
		// Direct register access when Mod is 3
		*core.registers.registers16Bit[modrm.rm] = *value
		destName = core.registers.index16ToString(modrm.rm)
	} else {
		// Calculate the effective address when accessing memory
		addr, destName = core.getEffectiveAddress32(modrm) // Discard the address description here as it's not used

		// Write the 16-bit value to memory
		err := core.memoryAccessController.WriteMemoryAddr16(uint32(addr), *value)
		if err != nil {
			return destName, err // Properly return the error
		}
	}

	return destName, nil // Return nil explicitly when no error occurs
}

func (core *CpuCore) writeRm32(modrm *ModRm, value *uint32) error {
	if modrm.mod == 3 {
		// Direct register access when Mod is 3
		*core.registers.registers32Bit[modrm.rm] = *value
	} else {
		// Calculate the effective address when accessing memory
		addr, _ := core.getEffectiveAddress32(modrm) // Discard the address description here as it's not used

		// Write the 16-bit value to memory
		err := core.memoryAccessController.WriteMemoryAddr32(uint32(addr), *value)
		if err != nil {
			return err // Properly return the error
		}
	}

	return nil // Return nil explicitly when no error occurs
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

func (core *CpuCore) writeSegmentRegister(register *SegmentRegister, value uint32) {
	if core.flags.MemorySegmentOverride > 0 {
		// default segment override
		switch core.flags.MemorySegmentOverride {
		case common.SEGMENT_CS:
			core.registers.CS.Base = value
		case common.SEGMENT_SS:
			core.registers.SS.Base = value
		case common.SEGMENT_DS:
			core.registers.DS.Base = value
		case common.SEGMENT_ES:
			core.registers.ES.Base = value
		case common.SEGMENT_FS:
			core.registers.FS.Base = value
		case common.SEGMENT_GS:
			core.registers.GS.Base = value
		default:
			panic("Unhandled segment register override")
		}
	}

	register.Base = value

}

func (cpuCore *CpuCore) logInstruction(logMessage string, a ...any) {
	// encode logMessage to utf-8 bytes

	if len(a) > 0 {
		logMessage = fmt.Sprintf(logMessage, a)
	}
	logMessageString := fmt.Sprintf("[op=%#04x]"+logMessage, cpuCore.currentOpCodeBeingExecuted)

	// encode logMessage to utf-8 bytes
	logMessageBytes := []byte(logMessageString)

	cpuCore.bus.SendMessageToAll(common.MODULE_DEBUG_MONITOR, bus.BusMessage{Subject: common.MESSAGE_GLOBAL_CPU_INSTRUCTION_LOG, Data: logMessageBytes})
}

func (cpuCore *CpuCore) logDebug(logMessage string) {
	// encode logMessage to utf-8 bytes
	logMessageBytes := []byte(logMessage)

	cpuCore.bus.SendMessageToAll(common.MODULE_DEBUG_MONITOR, bus.BusMessage{Subject: common.MESSAGE_GLOBAL_DEBUG_MESSAGE_LOG, Data: logMessageBytes})
}

func (core *CpuCore) HandleInterrupt(message bus.BusMessage) {

	vector := message.Data[0]

	// if vector comes from the second interrupt controller, subtract 8
	if message.Sender == common.MODULE_INTERRUPT_CONTROLLER_2 {
		vector -= 8
	}

	// Disable interrupts
	core.registers.SetFlag(InterruptFlag, false)

	// Push the current flags and CS:IP onto the stack
	err := stackPush16(core, core.registers.FLAGS)
	if err != nil {
		return
	}
	err = stackPush32(core, core.registers.CS.Base)
	if err != nil {
		return
	}
	err = stackPush16(core, core.registers.IP)
	if err != nil {
		return
	}

	// Set the necessary flags
	core.registers.SetFlag(TrapFlag, false)

	// Set the CS:IP to the interrupt vector

	vectorAddr := uint16(vector) << 2

	// Read IP and CS from the interrupt vector table

	vectorAddr2 := uint32(0x000020)<<4 + uint32(vectorAddr)

	core.registers.IP, _ = core.memoryAccessController.ReadMemoryValue16(uint32(vectorAddr2))
	newBase, _ := core.memoryAccessController.ReadMemoryValue16(uint32(vectorAddr2 + 3))

	core.registers.CS.Base = uint32(newBase)

	// Re-enable interrupts
	core.registers.SetFlag(InterruptFlag, true)

	// Send a message to the debug monitor
	core.logDebug(fmt.Sprintf("Interrupt %d handled", vector))

	// Send EOI message to the 8259A
	//core.sendEOI(vector)
	// Send EOI command to the 8259A
	eoiMessage := bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_COMPLETE,
		Sender:  message.Sender,
		Data:    []byte{vector},
	}
	err = core.bus.SendMessageToDeviceById(message.Sender, eoiMessage)
	if err != nil {
		core.logDebug(fmt.Sprintf("CPU: Error sending EOI message: %v", err))
	}
}

func (core *CpuCore) AcknowledgeInterrupt(message bus.BusMessage) {
	// send message to the interrupt controller that raised the interrupt
	err := core.bus.SendMessageToDeviceById(message.Sender, bus.BusMessage{Subject: common.MESSAGE_INTERRUPT_ACKNOWLEDGE, Data: message.Data})
	if err != nil {
		log.Fatalf("Failed to acknowledge interrupt: %s", err)
		return
	}
}

func (core *CpuCore) SetMemoryAccessController(controller *memmap.MemoryAccessController) {
	core.memoryAccessController = controller
}

func (core *CpuCore) shadowBios() {

	core.memoryAccessController.InitShadowBios()
}

func (core *CpuCore) GetRegister16(register *uint16) (uint32, string, uint8) {

	var registerIndex uint8
	for i, reg := range core.registers.registers16Bit {
		if reg == register {
			registerIndex = uint8(i)
			break
		}
	}

	if core.flags.OperandSizeOverrideEnabled {

		// find the index of the pointer address in registers16Bit
		return *core.registers.registers32Bit[registerIndex], core.registers.index32ToString(uint8(registerIndex)), 32

	} else {
		return uint32(*core.registers.registers16Bit[registerIndex]), core.registers.index16ToString(uint8(registerIndex)), 16
	}
}

func (core *CpuCore) SetRegister16(registerIndex uint8, value uint16) (string, error) {

	core.registers.registers16Bit[registerIndex] = &value
	return core.registers.index16ToString(registerIndex), nil
}

func (core *CpuCore) GetImmediate16() (uint32, uint8, error) {
	if core.flags.OperandSizeOverrideEnabled {
		imm32, err := core.readImm32()
		return imm32, 32, err
	} else {
		imm16, err := core.readImm16()
		return uint32(imm16), 16, err
	}
}

func (core *CpuCore) updateSystemFlags(cr0 uint32) {
	core.registers.CR0 = cr0

	// Toggle between protected and real mode based on the PE bit
	if cr0&0x1 == 0x1 {
		core.EnterMode(common.PROTECTED_MODE)
	} else {
		core.EnterMode(common.REAL_MODE)
	}

	// Handle Paging
	if cr0&0x80000000 != 0 {
		core.EnablePaging()
	} else {
		core.DisablePaging()
	}

	// Cache control
	if cr0&0x40000000 != 0 {
		core.DisableCache()
	} else {
		core.EnableCache()
	}

	// Write Protection
	if cr0&0x10000 != 0 {
		core.EnableWriteProtection()
	} else {
		core.DisableWriteProtection()
	}

	// Numeric Error
	if cr0&0x20 != 0 {
		core.EnableNumericError()
	} else {
		core.DisableNumericError()
	}
}

func (core *CpuCore) EnablePaging() {
	log.Printf("Paging enabled")
}

func (core *CpuCore) DisablePaging() {
	log.Printf("Paging disabled")
}

func (core *CpuCore) EnableCache() {
	log.Printf("Cache enabled")
}

func (core *CpuCore) DisableCache() {
	log.Printf("Cache disabled")
}

func (core *CpuCore) EnableWriteProtection() {
	log.Printf("Write protection enabled")
}

func (core *CpuCore) DisableWriteProtection() {
	log.Printf("Write protection disabled")
}

func (core *CpuCore) EnableNumericError() {
	log.Printf("Numeric error enabled")
}

func (core *CpuCore) DisableNumericError() {
	log.Printf("Numeric error disabled")
}

func (device *CpuCore) getSegmentOverride() SegmentRegister {
	if device.flags.MemorySegmentOverride > 0 {
		switch device.flags.MemorySegmentOverride {
		case common.SEGMENT_CS:
			return device.registers.CS
		case common.SEGMENT_SS:
			return device.registers.SS
		case common.SEGMENT_DS:
			return device.registers.DS
		case common.SEGMENT_ES:
			return device.registers.ES
		case common.SEGMENT_FS:
			return device.registers.FS
		case common.SEGMENT_GS:
			return device.registers.GS
		}
	}

	return device.registers.CS
}

func INSTR_HLT(core *CpuCore) {
	core.halt = true
	core.logInstruction(fmt.Sprintf("[%#04x] HLT", core.GetCurrentlyExecutingInstructionAddress()))
}
