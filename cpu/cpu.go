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
	registers              *CpuRegisters
	opCodeMap              []OpCodeImpl
	mode                   uint8
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
	instrByte := core.memoryAccessController.GetNextInstruction() //read 8 bit value

	instructionImpl := core.opCodeMap[instrByte.(uint8)]
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
}

func (core *CpuCore) Init(memController *common.MemoryAccessController) {
	core.memoryAccessController = memController

	core.EnterMode(common.REAL_MODE)

	core.memoryAccessController.SetCpuRegisterController(core)

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

	c.opCodeMap[0xEA] = JMP_FAR_PTR16
	c.opCodeMap[0xE9] = JMP_NEAR_REL16

}

type OpCodeImpl func(*CpuCore)

func JMP_FAR_PTR16(core *CpuCore) {
	destAddr := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)
	segment := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 3)

	log.Printf("INSTR DECODE: JMP FAR PTR %016X:%016X", segment, destAddr)

	core.memoryAccessController.SetCS(segment)
	core.memoryAccessController.SetIP(destAddr)
}

func JMP_NEAR_REL16(core *CpuCore) {

	tmp := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)

	log.Printf(fmt.Sprintf("Test: %b %08X   lastbit: %b", tmp, tmp, tmp>>15))

	var offset uint16
	offset = tmp
	var destAddr uint16 = core.registers.IP
	if (tmp>>15)&1 == 1 {
		offset = tmp << 1 >> 1
		log.Printf(fmt.Sprintf("Test: %b", offset))
		destAddr = destAddr + uint16(offset)

	} else {
		destAddr = destAddr - uint16(offset)
	}

	log.Printf("INSTR JMP NEAR REL 16 %08X", offset)

	core.memoryAccessController.SetIP(destAddr)
}
