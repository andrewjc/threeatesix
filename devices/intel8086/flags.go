package intel8086

import (
	"fmt"
)

const (
	CarryFlag            = 0x0001
	ParityFlag           = 0x0004
	AdjustFlag           = 0x0010
	ZeroFlag             = 0x0040
	SignFlag             = 0x0080
	TrapFlag             = 0x0100
	InterruptFlag        = 0x0200
	DirectionFlag        = 0x0400
	OverFlowFlag         = 0x0800
	IoPrivilegeLevelFlag = 0x3000
	NestedTaskFlag       = 0x4000
)

func (core *CpuRegisters) GetFlag(mask uint16) bool {
	return core.GetFlagInt(mask) == uint16(mask)
}

func (core *CpuRegisters) GetFlagInt(mask uint16) uint16 {
	if mask == 0x0002 {
		return 1
	} //Reserved, always 1 in EFLAGS
	if mask == 0x8000 {
		return 0
	} // Reserved, always 1 on 8086 and 186, always 0 on later models

	return core.FLAGS & mask
}

func (core *CpuRegisters) SetFlag(mask uint16, status bool) {
	if status {
		core.FLAGS = core.FLAGS | mask
	} else {
		core.FLAGS &= ^mask
	}
}

func INSTR_CLI(core *CpuCore) {
	// Clear interrupts

	core.logInstruction(fmt.Sprintf("[%#04x] CLI", core.GetCurrentCodePointer()))
	core.registers.SetFlag(InterruptFlag, false)
	core.currentByteAddr++
}

func INSTR_CLD(core *CpuCore) {
	// Clear direction flag
	core.currentByteAddr++
	core.logInstruction(fmt.Sprintf("[%#04x] CLD", core.GetCurrentCodePointer()))
	core.registers.SetFlag(DirectionFlag, false)
}

func INSTR_STI(core *CpuCore) {
	// Set interrupts
	core.currentByteAddr++
	core.logInstruction(fmt.Sprintf("[%#04x] STI", core.GetCurrentCodePointer()))
	core.registers.SetFlag(InterruptFlag, true)

	// Important: Handle delayed interrupt enabling
	core.interruptEnableDelay = 1 // Enable interrupts after the next instruction

}

func INSTR_CMC(core *CpuCore) {
	// Complement carry flag
	carryFlag := core.registers.GetFlag(CarryFlag)
	core.logInstruction(fmt.Sprintf("[%#04x] CMC", core.GetCurrentCodePointer()))
	core.registers.SetFlag(CarryFlag, !carryFlag)

	core.currentByteAddr++
}

func INSTR_CLC(core *CpuCore) {
	core.logInstruction(fmt.Sprintf("[%#04x] CLC", core.GetCurrentCodePointer()))
	core.registers.SetFlag(CarryFlag, false)
	core.currentByteAddr++
}

func INSTR_SMSW(core *CpuCore) {
	var value uint16
	var destName string

	core.currentByteAddr++
	modrm, bytesConsumed, err := core.consumeModRm()
	if err != nil {
		goto eof
	}
	core.currentByteAddr += bytesConsumed

	value = uint16(core.registers.CR0)

	// mask out the reserved bits
	value = value & 0xFFFF

	destName, err = core.writeRm16(&modrm, &value)
	if err != nil {
		core.logInstruction("Error in INSTR_SMSW: %s\n", err)
		doCoreDump(core)
		panic(0)
		return
	}

eof:
	core.logInstruction(fmt.Sprintf("[%#04x] smsw %s", core.GetCurrentlyExecutingInstructionAddress(), destName))
}
