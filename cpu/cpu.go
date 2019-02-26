package cpu

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
	"strings"
)

func New80386CPU() CpuCore {

	cpuCore := CpuCore{}
	cpuCore.registers = &CpuRegisters{}

	cpuCore.opCodeMap = make([]OpCodeImpl, 256)

	mapOpCodes(&cpuCore)

	return cpuCore
}

type CpuCore struct {
	memoryAccessController *common.MemoryAccessController
	ioPortAccessController *common.IOPortAccessController
	registers              *CpuRegisters
	opCodeMap              []OpCodeImpl
	mode                   uint8

	lastExecutedInstructionPointer uint16
	currentByteAtCodePointer       byte
}

func (core *CpuCore) SetCS(addr uint16) {
	core.registers.CS = addr
}

func (core *CpuCore) SetIP(addr uint16) {
	core.registers.IP = addr
}

func (core *CpuCore) GetIP() uint16 {
	return core.registers.IP
}

func (core *CpuCore) GetCS() uint16 {
	return core.registers.CS
}

func (core *CpuCore) IncrementIP() {
	core.registers.IP++
}

func (core *CpuCore) Step() {
	curCodePointer := core.GetCurrentCodePointer()
	if curCodePointer == core.lastExecutedInstructionPointer {
		log.Fatalf("CPU appears to be in a loop! Did you forget to increment the IP register?")
	}

	instrByte := core.memoryAccessController.GetNextInstruction() //read 8 bit value
	core.currentByteAtCodePointer = instrByte.(uint8)

	instructionImpl := core.opCodeMap[core.currentByteAtCodePointer]
	if instructionImpl != nil {
		instructionImpl(core)
	} else {
		log.Printf("CPU CORE ERROR!!!")
		core.dumpCoreInfo()

		// Gather next few bytes for debugging...
		peekBytes := core.memoryAccessController.PeekNextBytes(10)
		stb := strings.Builder{}
		for _, b := range peekBytes {
			stb.WriteString(fmt.Sprintf("%#2x ", b))
		}
		log.Printf("Next 10 bytes at instruction pointer: " + stb.String())

		log.Fatalf("CPU core failure. Unrecognised opcode: %#2x\n", instrByte)

	}

	fmt.Printf("CPU Stepped...\n")
	core.lastExecutedInstructionPointer = curCodePointer

}

func (core *CpuCore) Init(memController *common.MemoryAccessController, ioPortController *common.IOPortAccessController) {
	core.memoryAccessController = memController

	core.EnterMode(common.REAL_MODE)

	core.memoryAccessController.SetCpuRegisterController(core)
	core.ioPortAccessController = ioPortController

	core.Reset()
}

func (core *CpuCore) Reset() {
	core.registers.CS = 0xF000
	core.registers.IP = 0xFFF0
	core.memoryAccessController.LockBootVector()
}

type CpuRegisters struct {
	// 16bit registers (real mode)
	CS uint16 // code segment
	DS uint16 // data segment
	IP uint16 // instruction pointer

	// accumulator registers
	// used for I/O port access, arithmetic, interrupt calls
	AH  uint8
	AL  uint8
	AX  uint8

	// base registers
	// used as a base pointer for memory access
	BX  uint8
	BH  uint8
	BL  uint8

	// counter registers
	// used as loop counter and for shifts
	CX  uint8
	CH  uint8
	CL  uint8

	// data registers
	// used for I/O port access, arithmetic, interrupt calls
	DX  uint8
	DH  uint8
	DL  uint8



	// Flags
	DF uint16 // direction flag
	CF uint16
	OF uint16
	ZF uint16
	SF uint16
	AF uint16
	PF uint16

	// 32bit registers (protected mode)
	EIP uint32
	EAX uint32
	ECX uint32
	EDX uint32
	EBX uint32
	ESP uint32
	EBP uint32
	ESI uint32
	EDI uint32
}

func (core *CpuCore) IncrementEIP() {
	core.registers.EIP++
}

func (core *CpuCore) GetEIP() uint32 {
	return core.registers.EIP
}

func (core *CpuCore) EnterMode(mode uint8) {
	core.mode = mode
	core.memoryAccessController.EnterMode(mode)
}

func (core *CpuCore) dumpCoreInfo() {
	if core.mode == common.REAL_MODE {
		log.Println("Cpu core in real mode")
		log.Printf("Registers: IP: %016X, CS: %016X, DS: %016X", core.registers.IP, core.registers.CS, core.registers.DS)
	}
}

// Gets the current code segment + IP addr in memory
func (core *CpuCore) GetCurrentCodePointer() uint16 {
	return core.registers.CS<<4 + core.registers.IP
}

/* CPU OPCODE IMPLEMENTATIONS */

func mapOpCodes(c *CpuCore) {

	c.opCodeMap[0xEA] = INSTR_JMP_FAR_PTR16
	c.opCodeMap[0xE9] = INSTR_JMP_NEAR_REL16

	c.opCodeMap[0xFA] = INSTR_CLI
	c.opCodeMap[0xFC] = INSTR_CLD

	c.opCodeMap[0xE4] = INSTR_IN //imm to AL
	c.opCodeMap[0xE5] = INSTR_IN //DX to AL
	c.opCodeMap[0xEC] = INSTR_IN //imm to AX
	c.opCodeMap[0xED] = INSTR_IN //DX to AX

}

type OpCodeImpl func(*CpuCore)

func INSTR_CLI(core *CpuCore) {
	// Clear interrupts
	log.Printf("TODO: Write CLI (Clear interrupts implementation!")

	core.memoryAccessController.SetIP(uint16(core.GetIP()+1))
}

func INSTR_CLD(core *CpuCore) {
	// Clear direction flag
	core.registers.DF = 0

	core.memoryAccessController.SetIP(uint16(core.GetIP()+1))
}

func INSTR_IN(core *CpuCore) {
	// Read from port

	switch {
	case 0xE4 == core.currentByteAtCodePointer:
		{
			// Read from port (imm) to AL
			imm := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)

			data := core.ioPortAccessController.ReadAddr8(imm)

			core.registers.AL = data
			log.Printf("Port IN addr: imm addr %04X to AL (data = %04X)", imm, data)
		}
	case 0xE5 == core.currentByteAtCodePointer:
		{
			// Read from port (DX) to AL

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.registers.AL = data
			log.Printf("Port IN addr: DX VAL %04X to AL (data = %04X)", dx, data)
		}
	case 0xEC == core.currentByteAtCodePointer:
		{
			// Read from port (imm) to AX

			imm := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)

			data := core.ioPortAccessController.ReadAddr8(imm)

			core.registers.AX = data
			log.Printf("Port IN addr: imm addr %04X to AX (data = %04X)", imm, data)
		}
	case 0xED == core.currentByteAtCodePointer:
		{
			// Read from port (DX) to AX

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.registers.AX = data
			log.Printf("Port IN addr: DX VAL %04X to AX (data = %04X)", dx, data)
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}

	core.memoryAccessController.SetIP(uint16(core.GetIP()+2))
}


func INSTR_JMP_FAR_PTR16(core *CpuCore) {
	destAddr := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)
	segment := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 3)

	core.memoryAccessController.SetCS(segment)
	core.memoryAccessController.SetIP(destAddr)
}

func INSTR_JMP_NEAR_REL16(core *CpuCore) {

	offset := int16(core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1))

	var destAddr = int16(core.registers.IP)

	destAddr = destAddr + int16(offset)

	core.memoryAccessController.SetIP(uint16(destAddr)+3)
}
