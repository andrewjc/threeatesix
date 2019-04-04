package intel8086

import (
	"fmt"
	"log"
)

/* CPU OPCODE IMPLEMENTATIONS */

func mapOpCodes(c *CpuCore) {

	c.opCodeMap[0xEA] = INSTR_JMP_FAR_PTR16

	c.opCodeMap[0xE9] = INSTR_JMP_NEAR_REL16

	c.opCodeMap[0xEB] = INSTR_JMP_SHORT_REL8

	c.opCodeMap[0xE3] = INSTR_JCXZ_SHORT_REL8

	c.opCodeMap[0x74] = INSTR_JZ_SHORT_REL8
	c.opCodeMap[0x75] = INSTR_JNZ_SHORT_REL8

	c.opCodeMap[0xFA] = INSTR_CLI
	c.opCodeMap[0xFC] = INSTR_CLD

	c.opCodeMap[0xE4] = INSTR_IN //imm to AL
	c.opCodeMap[0xE5] = INSTR_IN //DX to AL
	c.opCodeMap[0xEC] = INSTR_IN //imm to AX
	c.opCodeMap[0xED] = INSTR_IN //DX to AX

	c.opCodeMap[0xE6] = INSTR_OUT //AL to imm
	c.opCodeMap[0xE7] = INSTR_OUT //AX to imm
	c.opCodeMap[0xEE] = INSTR_OUT //AL to DX
	c.opCodeMap[0xEF] = INSTR_OUT //AX to DX

	c.opCodeMap[0xA8] = INSTR_TEST
	c.opCodeMap[0xA9] = INSTR_TEST
	c.opCodeMap[0xF6] = INSTR_TEST
	c.opCodeMap[0xF7] = INSTR_TEST
	c.opCodeMap[0x84] = INSTR_TEST
	c.opCodeMap[0x85] = INSTR_TEST

	c.opCodeMap[0xA0] = INSTR_MOV
	c.opCodeMap[0xA1] = INSTR_MOV
	c.opCodeMap[0xA2] = INSTR_MOV
	c.opCodeMap[0xA3] = INSTR_MOV
	for i := 0; i < len(c.registers.registers8Bit); i++ {
		c.opCodeMap[0xB0+i] = INSTR_MOV
	}

	for i := 0; i < len(c.registers.registers16Bit); i++ {
		c.opCodeMap[0xB8+i] = INSTR_MOV
	}

	c.opCodeMap[0x8A] = INSTR_MOV
	c.opCodeMap[0x8B] = INSTR_MOV
	c.opCodeMap[0x8C] = INSTR_MOV
	c.opCodeMap[0x8E] = INSTR_MOV

	c.opCodeMap[0x3A] = INSTR_CMP
	c.opCodeMap[0x3B] = INSTR_CMP
	c.opCodeMap[0x3C] = INSTR_CMP
	c.opCodeMap[0x3D] = INSTR_CMP
	//c.opCodeMap[0x80] = INSTR_CMP // handled by 80 opcode switch
	//c.opCodeMap[0x81] = INSTR_CMP // handled by 81 opcode switch
	//c.opCodeMap[0x83] = INSTR_CMP // handled by 83 opcode switch
	c.opCodeMap[0x38] = INSTR_CMP
	c.opCodeMap[0x39] = INSTR_CMP
	c.opCodeMap[0xA6] = INSTR_CMP
	c.opCodeMap[0xA7] = INSTR_CMP

	c.opCodeMap[0x86] = INSTR_XCHG
	c.opCodeMap[0x87] = INSTR_XCHG
	for i := 0; i < len(c.registers.registers16Bit); i++ {
		c.opCodeMap[0x90+i] = INSTR_XCHG
	}

	//c.opCodeMap[0x90] = INSTR_NOP // we don't define an NOP because NOP = xchg ax, ax

	c.opCodeMap[0xC3] = INSTR_RET_NEAR

	c.opCodeMap[0x28] = INSTR_SUB
	c.opCodeMap[0x29] = INSTR_SUB
	c.opCodeMap[0x2A] = INSTR_SUB
	c.opCodeMap[0x2B] = INSTR_SUB
	c.opCodeMap[0x2C] = INSTR_SUB
	c.opCodeMap[0x2D] = INSTR_SUB
	c.opCodeMap[0x80] = INSTR_SUB
	c.opCodeMap[0x81] = INSTR_SUB
	c.opCodeMap[0x83] = INSTR_SUB

	c.opCodeMap[0x04] = INSTR_ADD
	c.opCodeMap[0x05] = INSTR_ADD
	c.opCodeMap[0x00] = INSTR_ADD
	c.opCodeMap[0x01] = INSTR_ADD
	c.opCodeMap[0x02] = INSTR_ADD
	c.opCodeMap[0x03] = INSTR_ADD

	c.opCodeMap[0x24] = INSTR_AND
	c.opCodeMap[0x25] = INSTR_AND
	c.opCodeMap[0x20] = INSTR_AND
	c.opCodeMap[0x21] = INSTR_AND
	c.opCodeMap[0x22] = INSTR_AND
	c.opCodeMap[0x23] = INSTR_AND

	c.opCodeMap[0x14] = INSTR_ADC
	c.opCodeMap[0x15] = INSTR_ADC
	c.opCodeMap[0x10] = INSTR_ADC
	c.opCodeMap[0x11] = INSTR_ADC
	c.opCodeMap[0x12] = INSTR_ADC
	c.opCodeMap[0x13] = INSTR_ADC

	c.opCodeMap[0xD0] = INSTR_SHIFT
	c.opCodeMap[0xD1] = INSTR_SHIFT
	c.opCodeMap[0xD2] = INSTR_SHIFT
	c.opCodeMap[0xD3] = INSTR_SHIFT
	c.opCodeMap[0xC0] = INSTR_SHIFT
	c.opCodeMap[0xC1] = INSTR_SHIFT

	c.opCodeMap[0x30] = INSTR_XOR
	c.opCodeMap[0x31] = INSTR_XOR
	c.opCodeMap[0x32] = INSTR_XOR
	c.opCodeMap[0x33] = INSTR_XOR
	c.opCodeMap[0x34] = INSTR_XOR
	c.opCodeMap[0x35] = INSTR_XOR

	// opcodes that handle multiple instructions (handled by modrm byte)
	c.opCodeMap[0xFF] = INSTR_FF_OPCODES
	c.opCodeMap[0x80] = INSTR_80_OPCODES
	c.opCodeMap[0x81] = INSTR_81_OPCODES
	c.opCodeMap[0x83] = INSTR_83_OPCODES

	c.opCodeMap[0x0c] = INSTR_OR
	c.opCodeMap[0x0d] = INSTR_OR
	//	c.opCodeMap[0x80] = INSTR_OR // handled by 80 opcode switch
	//	c.opCodeMap[0x81] = INSTR_OR // handled by 81 opcode switch
	//	c.opCodeMap[0x83] = INSTR_OR // handled by 83 opcode switch
	c.opCodeMap[0x08] = INSTR_OR
	c.opCodeMap[0x09] = INSTR_OR
	c.opCodeMap[0x0A] = INSTR_OR
	c.opCodeMap[0x0B] = INSTR_OR

	for i := 0; i < len(c.registers.registers16Bit); i++ {
		c.opCodeMap[0x50+i] = INSTR_PUSH
	}
	c.opCodeMap[0x50] = INSTR_PUSH

	c.opCodeMap[0x6A] = INSTR_PUSH
	c.opCodeMap[0x68] = INSTR_PUSH
	c.opCodeMap[0x0E] = INSTR_PUSH
	c.opCodeMap[0x16] = INSTR_PUSH
	c.opCodeMap[0x1E] = INSTR_PUSH
	c.opCodeMap[0x06] = INSTR_PUSH


	// 2 byte opcodes
	c.opCodeMap2Byte[0x01] = INSTR_SMSW
}

type OpCodeImpl func(*CpuCore)

func INSTR_SMSW(core *CpuCore) {
	core.IncrementIP()
	modrm := core.consumeModRm()

	value := uint16(core.registers.CR0)

	core.writeRm16(&modrm, &value)
	log.Printf("[%#04x] smsw %s", core.GetCurrentlyExecutingInstructionPointer(), "r/m16")
	core.registers.IP = uint16(core.GetIP() + 1)
}

func INSTR_FF_OPCODES(core *CpuCore) {

	tmp := core.registers.IP
	core.IncrementIP()
	modrm := core.consumeModRm()
	core.registers.IP = tmp

	switch {
	case modrm.rm == 0:
		{
			// inc rm32
		}
	case modrm.rm == 1:
		{
			// dec rm32
		}
	case modrm.rm == 2:
		{
			// call rm32
		}
	case modrm.rm == 3:
		{
			// call m16
		}
	case modrm.rm == 4:
		{
			// jmp rm32
			INSTR_JMP_FAR_M16(core, &modrm)
		}
	case modrm.rm == 5:
		{
			// jmp m16
			INSTR_JMP_FAR_M16(core, &modrm)
		}
	case modrm.rm == 6:
		{
			// push rm32
			INSTR_PUSH(core)
		}
	}
}

func INSTR_80_OPCODES(core *CpuCore) {

	tmp := core.registers.IP
	core.IncrementIP()
	modrm := core.consumeModRm()
	core.registers.IP = tmp

	switch modrm.reg {
	case 1:
		INSTR_OR(core)
	case 4:
		INSTR_AND(core)
	case 6:
		INSTR_XOR(core)
	case 7:
		INSTR_CMP(core)
	default:
		log.Println(fmt.Sprintf("INSTR_80_OPCODE UNHANDLED OPTION %d\n\n", modrm))
		doCoreDump(core)
	}
}

func INSTR_81_OPCODES(core *CpuCore) {

	tmp := core.registers.IP
	core.IncrementIP()
	modrm := core.consumeModRm()
	core.registers.IP = tmp

	switch modrm.reg {
	case 0:
		INSTR_ADD(core)
	case 1:
		INSTR_OR(core)
	case 4:
		INSTR_AND(core)
	case 6:
		INSTR_XOR(core)
	case 7:
		INSTR_CMP(core)
	default:
		log.Println(fmt.Sprintf("INSTR_81_OPCODE UNHANDLED OPTION %d\n\n", modrm))
		doCoreDump(core)
	}
}

func INSTR_83_OPCODES(core *CpuCore) {

	tmp := core.registers.IP
	core.IncrementIP()
	modrm := core.consumeModRm()
	core.registers.IP = tmp

	switch modrm.reg {
	case 0:
		INSTR_ADD(core)
	case 1:
		INSTR_OR(core)
	case 4:
		INSTR_AND(core)
	case 5:
		INSTR_SUB(core)
	case 6:
		INSTR_XOR(core)
	case 7:
		INSTR_CMP(core)
	default:
		log.Println(fmt.Sprintf("INSTR_83_OPCODE UNHANDLED OPTION %d\n\n", modrm))
		doCoreDump(core)
	}
}

func INSTR_NOP(core *CpuCore) {
	// Clear interrupts
	log.Printf("[%#04x] NOP", core.GetCurrentlyExecutingInstructionPointer())

	core.registers.IP = uint16(core.GetIP() + 1)
}


func INSTR_CLI(core *CpuCore) {
	// Clear interrupts
	log.Printf("[%#04x] CLI", core.GetCurrentCodePointer())
	core.registers.SetFlag(InterruptFlag, false)
	core.registers.IP = uint16(uint16(core.GetIP() + 1))
}

func INSTR_CLD(core *CpuCore) {
	// Clear direction flag
	log.Printf("[%#04x] CLD", core.GetCurrentCodePointer())
	core.registers.SetFlag(DirectionFlag, false)
	core.registers.IP = uint16(uint16(core.GetIP() + 1))
}





