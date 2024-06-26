package intel8086

/* CPU OPCODE IMPLEMENTATIONS */

type OpCodeImpl func(*CpuCore)

func mapOpCodes(c *CpuCore) {
	// Mapping individual opcodes to instruction handlers
	opcodeHandlers := map[byte]OpCodeImpl{
		0xEA: INSTR_JMP_FAR_PTR16,
		0xE9: INSTR_JMP_NEAR_REL16,
		0xE8: INSTR_CALL_NEAR_REL16,
		0xEB: INSTR_JMP_SHORT_REL8,
		0xE0: INSTR_DEC_COUNT_JMP_SHORT_Z,
		0xE1: INSTR_DEC_COUNT_JMP_SHORT_Z,
		0xE2: INSTR_DEC_COUNT_JMP_SHORT,
		0xE3: INSTR_JCXZ_SHORT_REL8,
		0xFA: INSTR_CLI,
		0xFC: INSTR_CLD,
		0xFE: INSTR_INC_SHORT_REL8,
		0xF4: INSTR_HLT,
		0xF5: INSTR_CMC,
		0xF8: INSTR_CLC,
		0xF9: INSTR_STC,
		0x70: INSTR_JO_SHORT_REL8,
		0x71: INSTR_JNO_SHORT_REL8,
		0x74: INSTR_JZ_SHORT_REL8,
		0x75: INSTR_JNZ_SHORT_REL8,
		0x76: INSTR_JBE_SHORT_REL8,
		0x78: INSTR_JS_SHORT_REL8,
		0x79: INSTR_JNS_SHORT_REL8,
		0x7e: INSTR_JLE_SHORT_REL8,
		0x7f: INSTR_JG_SHORT_REL8,
		0x72: INSTR_JB_SHORT_REL8,
		0x73: INSTR_JNB_SHORT_REL8,
		0x7c: INSTR_JL_SHORT_REL8,
		0x7d: INSTR_JGE_SHORT_REL8,
		0x7a: INSTR_JPE_SHORT_REL8,
		0x7b: INSTR_JPO_SHORT_REL8,
		0xC2: INSTR_RET_FAR,
		0xC3: INSTR_RET_NEAR,
		// Input/Output Instructions
		0xE4: INSTR_IN,
		0xE5: INSTR_IN,
		0xEC: INSTR_IN,
		0xED: INSTR_IN,
		0xE6: INSTR_OUT,
		0xE7: INSTR_OUT,
		0xEE: INSTR_OUT,
		0xEF: INSTR_OUT,
		// Data Movement Instructions
		0xA0: INSTR_MOV,
		0xA1: INSTR_MOV,
		0xA2: INSTR_MOV,
		0xA3: INSTR_MOV,
		0xA4: INSTR_MOV,
		0xA5: INSTR_MOV,
		0x8A: INSTR_MOV,
		0x8B: INSTR_MOV,
		0x8C: INSTR_MOV,
		0x8D: INSTR_MOV,
		0x8E: INSTR_MOV,
		0x88: INSTR_MOV,
		0x89: INSTR_MOV,
		// Comparison Instructions
		0x3A: INSTR_CMP,
		0x3B: INSTR_CMP,
		0x3C: INSTR_CMP,
		0x3D: INSTR_CMP,
		0x38: INSTR_CMP,
		0x39: INSTR_CMP,
		0xA6: INSTR_CMP,
		0xA7: INSTR_CMP,
		0x86: INSTR_XCHG,
		0x87: INSTR_XCHG,
		// Arithmetic Instructions
		0x04: INSTR_ADD,
		0x05: INSTR_ADD,
		0x00: INSTR_ADD,
		0x01: INSTR_ADD,
		0x02: INSTR_ADD,
		0x03: INSTR_ADD,
		0x24: INSTR_AND,
		0x25: INSTR_AND,
		0x20: INSTR_AND,
		0x21: INSTR_AND,
		0x22: INSTR_AND,
		0x23: INSTR_AND,
		0x2a: INSTR_SUB,
		0x2b: INSTR_SUB,
		0x2c: INSTR_SUB,
		0x2d: INSTR_SUB,
		0x28: INSTR_SUB,
		0x29: INSTR_SUB,
		0x18: INSTR_SBB,
		0x19: INSTR_SBB,
		0x1a: INSTR_SBB,
		0x1b: INSTR_SBB,
		0x14: INSTR_ADC,
		0x15: INSTR_ADC,
		0x10: INSTR_ADC,
		0x11: INSTR_ADC,
		0x12: INSTR_ADC,
		0x13: INSTR_ADC,
		0x2F: INSTR_DAS, // Decimal Adjust AL after Subtraction
		0x27: INSTR_DAA, // Decimal Adjust AL after Addition
		0x37: INSTR_AAA, // ASCII Adjust AL After Addition
		0x3F: INSTR_AAS, // ASCII Adjust AL After Subtraction
		0xD0: INSTR_SHIFT,
		0xD1: INSTR_SHIFT,
		0xD2: INSTR_SHIFT,
		0xD3: INSTR_SHIFT,
		0xC0: INSTR_SHIFT,
		0xC1: INSTR_SHIFT,
		0x30: INSTR_XOR,
		0x31: INSTR_XOR,
		0x32: INSTR_XOR,
		0x33: INSTR_XOR,
		0x34: INSTR_XOR,
		0x35: INSTR_XOR,
		0x0c: INSTR_OR,
		0x0d: INSTR_OR,
		0x08: INSTR_OR,
		0x09: INSTR_OR,
		0x0A: INSTR_OR,
		0x0B: INSTR_OR,
		// Stack Operations
		0x60: INSTR_PUSH,
		0x6A: INSTR_PUSH,

		0x6C: INSTR_INS,
		0x6D: INSTR_INS,
		0x6E: INSTR_OUTS, //ADDED
		0x6F: INSTR_OUTS, //ADDED

		0x68: INSTR_PUSH,
		0x69: INSTR_IMUL,
		0x0E: INSTR_PUSH,
		0x16: INSTR_PUSH,
		0x1E: INSTR_PUSH,
		0x06: INSTR_PUSH,
		0x61: INSTR_POP,
		0x8F: INSTR_POP,
		0x17: INSTR_POP,
		// Special Instructions
		0xAA: INSTR_STOSB, //Store string byte
		0xAB: INSTR_STOSD,
		0xAC: INSTR_LODS,
		0xAD: INSTR_LODS,

		0xA8: INSTR_TEST_IMM8_AL,  // Test immediate 8-bit with AL
		0xA9: INSTR_TEST_IMM16_AX, // Test immediate 16-bit with AX

		0x1D: INSTR_SBB, // Subtract with Borrow from AX with immediate 16-bit
		0x1C: INSTR_SBB, // Subtract with Borrow from AL with immediate 8-bit
		0x1F: INSTR_POP, // Subtract with Borrow from AX with register16

		// Group 1 opcodes, dynamically handled based on ModR/M byte
		0x80: INSTR_80_OPCODES,
		0x81: INSTR_80_OPCODES,
		0x82: INSTR_80_OPCODES,
		0x83: INSTR_80_OPCODES,
		// Test opcodes, handled based on ModR/M byte
		0x84: INSTR_TEST,              // Test 8-bit register/memory with 8-bit register
		0x85: INSTR_TEST,              // Test 16-bit register/memory with 16-bit register
		0xF6: handleGroup3OpCode_byte, // Group 3 byte operations (TEST, NOT, NEG, MUL, IMUL, DIV, IDIV)
		0xF7: handleGroup5Opcode_word, // Group 3 word operations (TEST, NOT, NEG, MUL, IMUL, DIV, IDIV)

		// Software interrupts
		//0xCD: INSTR_INT,
		0xCC: INSTR_INT3,
		//0xCE: INSTR_INT,
		//0xCF: INSTR_IRET,

	}

	// Register-based opcodes dynamically generated from register lists
	for i, _ := range c.registers.registers8Bit {
		opcodeHandlers[0xB0+byte(i)] = INSTR_MOV
	}

	for i, _ := range c.registers.registers16Bit {
		opcodeHandlers[0xB8+byte(i)] = INSTR_MOV
		opcodeHandlers[0x40+byte(i)] = INSTR_INC
		opcodeHandlers[0x48+byte(i)] = INSTR_DEC
		opcodeHandlers[0x50+byte(i)] = INSTR_PUSH
		opcodeHandlers[0x58+byte(i)] = INSTR_POP
		opcodeHandlers[0x90+byte(i)] = INSTR_XCHG
	}

	// Transfer opcodes into the CPU core map
	for k, v := range opcodeHandlers {
		c.opCodeMap[k] = v
	}

	// Two-byte opcode map
	opCodeMap2ByteHandlers := map[byte]OpCodeImpl{
		0x01: INSTR_SMSW,
		0x20: INSTR_MOV,
		0x22: INSTR_MOV,
		0x85: INSTR_TEST,
		0x09: INSTR_WBINVD,
	}

	// Transfer two-byte opcodes into the CPU core map
	for k, v := range opCodeMap2ByteHandlers {
		c.opCodeMap2Byte[k] = v
	}
}
