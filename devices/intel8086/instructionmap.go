package intel8086

import "fmt"

/* CPU OPCODE IMPLEMENTATIONS */

type OpCodeImpl func(*CpuCore)

func mapOpCodes(c *CpuCore) {

	c.opCodeMap[0xEA] = INSTR_JMP_FAR_PTR16

	c.opCodeMap[0xE9] = INSTR_JMP_NEAR_REL16
	c.opCodeMap[0xE8] = INSTR_CALL_NEAR_REL16

	c.opCodeMap[0xEB] = INSTR_JMP_SHORT_REL8

	c.opCodeMap[0xE0] = INSTR_DEC_COUNT_JMP_SHORT_ECX
	c.opCodeMap[0xE1] = INSTR_DEC_COUNT_JMP_SHORT_ECX
	c.opCodeMap[0xE2] = INSTR_DEC_COUNT_JMP_SHORT_ECX
	c.opCodeMap[0xE3] = INSTR_JCXZ_SHORT_REL8

	c.opCodeMap[0x70] = INSTR_JO_SHORT_REL8
	c.opCodeMap[0x71] = INSTR_JNO_SHORT_REL8
	c.opCodeMap[0x74] = INSTR_JZ_SHORT_REL8
	c.opCodeMap[0x75] = INSTR_JNZ_SHORT_REL8
	c.opCodeMap[0x76] = INSTR_JBE_SHORT_REL8
	c.opCodeMap[0x78] = INSTR_JS_SHORT_REL8
	c.opCodeMap[0x79] = INSTR_JNS_SHORT_REL8

	c.opCodeMap[0xFA] = INSTR_CLI
	c.opCodeMap[0xFC] = INSTR_CLD
	c.opCodeMap[0xFE] = INSTR_INC_SHORT_REL8
	c.opCodeMap[0xF4] = INSTR_HLT
	c.opCodeMap[0xF9] = INSTR_STC

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
	c.opCodeMap[0x8D] = INSTR_MOV

	c.opCodeMap[0x8E] = INSTR_MOV

	c.opCodeMap[0x88] = INSTR_MOV
	c.opCodeMap[0x89] = INSTR_MOV

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

	c.opCodeMap[0x1a] = INSTR_SBB
	c.opCodeMap[0x1b] = INSTR_SBB
	c.opCodeMap[0x1c] = INSTR_SBB
	c.opCodeMap[0x1d] = INSTR_SBB
	c.opCodeMap[0x80] = INSTR_SBB
	c.opCodeMap[0x81] = INSTR_SBB
	c.opCodeMap[0x83] = INSTR_SBB
	c.opCodeMap[0x18] = INSTR_SBB
	c.opCodeMap[0x19] = INSTR_SBB

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

	// Confirm these, they seem to be the same???
	c.opCodeMap[0x80] = INSTR_80_OPCODES
	c.opCodeMap[0x82] = INSTR_80_OPCODES

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
		c.opCodeMap[0x40+i] = INSTR_INC
	}
	c.opCodeMap[0x40] = INSTR_INC

	for i := 0; i < len(c.registers.registers16Bit); i++ {
		c.opCodeMap[0x50+i] = INSTR_PUSH
	}
	c.opCodeMap[0x50] = INSTR_PUSH

	c.opCodeMap[0x58] = INSTR_POP
	c.opCodeMap[0x59] = INSTR_POP
	c.opCodeMap[0x5A] = INSTR_POP
	c.opCodeMap[0x5B] = INSTR_POP
	c.opCodeMap[0x5C] = INSTR_POP
	c.opCodeMap[0x5D] = INSTR_POP
	c.opCodeMap[0x5E] = INSTR_POP
	c.opCodeMap[0x5F] = INSTR_POP
	c.opCodeMap[0x61] = INSTR_POP
	c.opCodeMap[0x8F] = INSTR_POP

	c.opCodeMap[0x60] = INSTR_PUSH
	c.opCodeMap[0x6A] = INSTR_PUSH
	c.opCodeMap[0x68] = INSTR_PUSH
	c.opCodeMap[0x0E] = INSTR_PUSH
	c.opCodeMap[0x16] = INSTR_PUSH
	c.opCodeMap[0x1E] = INSTR_PUSH
	c.opCodeMap[0x06] = INSTR_PUSH

	c.opCodeMap[0xAB] = INSTR_STOSD
	c.opCodeMap[0xAC] = INSTR_LODS
	c.opCodeMap[0xAD] = INSTR_LODS

	// 2 byte opcodes
	c.opCodeMap2Byte[0x01] = INSTR_SMSW
	c.opCodeMap2Byte[0x20] = INSTR_MOV
	c.opCodeMap2Byte[0x22] = INSTR_MOV
	c.opCodeMap2Byte[0x85] = INSTR_TEST
	c.opCodeMap2Byte[0x09] = INSTR_WBINVD
}

func INSTR_HLT(core *CpuCore) {
	core.halt = true

	core.logInstruction(fmt.Sprintf("[%#04x] HLT", core.GetCurrentlyExecutingInstructionAddress()))
}
