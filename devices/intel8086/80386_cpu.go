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

	register_locations_8 := core.registers.registers8Bit
	register_locations_16 := core.registers.registers16Bit
	register_locations_32 := core.registers.registers32Bit

	core.currentByteAddr = core.GetCurrentCodePointer()
	tmp := core.currentByteAddr
	if core.currentByteAddr == core.lastExecutedInstructionPointer {
		log.Fatalf("CPU appears to be in a loop! Did you forget to increment the IP register?")
		doCoreDump(core)
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

func (device *CpuCore) getEffectiveAddress32(modRm *ModRm) (*uint32, string) {
	if modRm.mod == 3 {
		return device.registers.registers32Bit[modRm.rm], device.registers.index32ToString(modRm.rm)
	}

	addressMode := modRm.getAddressMode32(device)
	return &addressMode, fmt.Sprintf("dword_F%#04x", addressMode)
}

func (core *CpuCore) readImm8() (uint8, error) {
	retVal, err := core.memoryAccessController.ReadMemoryAddr8(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr++
	return retVal, nil
}

func (core *CpuCore) readImm16() (uint16, error) {
	retVal, err := core.memoryAccessController.ReadMemoryAddr16(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr += 2
	return retVal, nil
}

func (core *CpuCore) readImm32() (uint32, error) {
	retVal, err := core.memoryAccessController.ReadMemoryAddr32(uint32(core.currentByteAddr))
	if err != nil {
		return 0, err
	}
	core.currentByteAddr += 4
	return retVal, nil
}

func (core *CpuCore) readRm8(modrm *ModRm) (*uint8, string, error) {
	if modrm.mod == 3 {
		dest := core.registers.registers8Bit[modrm.rm]
		destName := core.registers.index8ToString(modrm.rm)
		return dest, destName, nil

	} else {
		addressMode := modrm.getAddressMode16(core)
		destValue, err := core.memoryAccessController.ReadMemoryAddr8(uint32(addressMode))
		destName := fmt.Sprintf("dword_F%#04x", addressMode)
		return &destValue, destName, err
	}
}

func (core *CpuCore) readRm16(modrm *ModRm) (*uint16, string, error) {
	if modrm.mod == 3 {
		dest := uint16(*core.registers.registers16Bit[modrm.rm])
		destName := core.registers.index16ToString(modrm.rm)
		return &dest, destName, nil

	} else {
		addressMode := modrm.getAddressMode16(core)
		destValue, err := core.memoryAccessController.ReadMemoryAddr16(uint32(addressMode))
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
		err := core.memoryAccessController.WriteMemoryAddr8(uint32(addressMode), *value)
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
		err := core.memoryAccessController.WriteMemoryAddr16(uint32(addressMode), *value)
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

func (cpuCore *CpuCore) logInstruction(logMessage string) {
	// encode logMessage to utf-8 bytes
	logMessageBytes := []byte(logMessage)

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
	core.registers.CS.Base, _ = core.memoryAccessController.ReadMemoryAddr32(uint32(vector * 4))
	core.registers.IP, _ = core.memoryAccessController.ReadMemoryAddr16(uint32(vector*4 + 2))

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
		log.Printf("CPU: Error sending EOI message: %v", err)
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

func INSTR_HLT(core *CpuCore) {
	core.halt = true
	core.logInstruction(fmt.Sprintf("[%#04x] HLT", core.GetCurrentlyExecutingInstructionAddress()))
}
