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
	var destName string

	// Get the Machine Status Word (MSW), which is the lower 16 bits of CR0
	msw := uint16(core.registers.CR0 & 0xFFFF)

	modrm, bytesConsumed, err := core.consumeModRm()
	if err != nil {
		goto eof
	}
	core.currentByteAddr += bytesConsumed

	switch modrm.reg {
	case 4, 5, 7: // All variants of SMSW
		if modrm.mod == 3 { // Register operand
			if modrm.reg == 5 { // 32-bit register for SMSW r32/m16
				tmp := uint32(modrm.rm)
				var regValue uint32
				regValue, destName, _ = core.GetRegister32(&tmp)
				_, err := core.SetRegister32(modrm.rm, (regValue&0xFFFF0000)|uint32(msw))
				if err != nil {
					return
				}
			} else { // 16-bit register for SMSW r/m16
				tmp := uint16(modrm.rm)
				_, destName, _ = core.GetRegister16(&tmp)
				_, err := core.SetRegister16(modrm.rm, msw)
				if err != nil {
					return
				}
			}
		} else { // Memory operand
			var addr uint32

			// Calculate the effective address when accessing memory
			addr, destName = core.getEffectiveAddress32(&modrm) // Discard the address description here as it's not used

			// Write the 16-bit value to memory
			err := core.memoryAccessController.WriteMemoryAddr16(uint32(addr), msw)
			if err != nil {
				return
			}
		}
	default:
		core.logInstruction(fmt.Sprintf("Unrecognized SMSW opcode: 0x0F 0x01 0x%02X", modrm.reg))
		doCoreDump(core)
		panic(fmt.Sprintf("Unrecognized SMSW opcode: 0x0F 0x01 0x%02X", modrm.reg))
	}

eof:
	core.logInstruction(fmt.Sprintf("[%#04x] smsw %s", core.GetCurrentlyExecutingInstructionAddress(), destName))
}
