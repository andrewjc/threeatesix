package cpu

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func New80386CPU() CpuCore {

	cpuCore := CpuCore{}
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

	cpuCore.registers.registersSegmentRegisters = []*uint16{
		&cpuCore.registers.ES,
		&cpuCore.registers.CS,
		&cpuCore.registers.SS,
		&cpuCore.registers.DS,
	}

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

func (core *CpuCore) EnterMode(mode uint8) {
	core.mode = mode
	core.memoryAccessController.EnterMode(mode)
}

// Gets the current code segment + IP addr in memory
func (core *CpuCore) GetCurrentCodePointer() uint16 {
	return core.registers.CS<<4 + core.registers.IP
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

		log.Printf("CPU core failure. Unrecognised opcode: %#2x\n", instrByte)
		doCoreDump(core)

		//log.Fatal("Execution Terminated.")
	}

	fmt.Printf("CPU Stepped...\n")
	core.lastExecutedInstructionPointer = curCodePointer

}


func getMSB(value uint8) uint8 {
	return (value >> 8) & 1
}

func getBitValue(value uint8, place uint8) uint8 {
	return (value >> place) & 1
}
