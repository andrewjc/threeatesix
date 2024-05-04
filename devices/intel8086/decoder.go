package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func (core *CpuCore) decodeInstruction() uint8 {

	var instrByte uint8
	var err error
	nextInstructionAddr := core.SegmentAddressToLinearAddress32(core.registers.CS, uint32(core.registers.IP))

	core.flags.MemorySegmentOverride = 0
	core.flags.OperandSizeOverrideEnabled = false
	core.flags.AddressSizeOverrideEnabled = false
	core.flags.LockPrefixEnabled = false
	core.flags.RepPrefixEnabled = false
	core.is2ByteOperand = false

	core.currentPrefixBytes = []byte{}
	if isPrefixByte(core.memoryAccessController.PeekNextBytes(nextInstructionAddr, 1)[0]) {

		prefixByte := core.memoryAccessController.PeekNextBytes(nextInstructionAddr, 1)[0]
		core.currentPrefixBytes = append(core.currentPrefixBytes, prefixByte)
		switch prefixByte {
		case 0x2e:
			// cs segment override
			core.flags.MemorySegmentOverride = common.SEGMENT_CS
		case 0x36:
			// ss segment override
			core.flags.MemorySegmentOverride = common.SEGMENT_SS
		case 0x3e:
			// ds segment override
			core.flags.MemorySegmentOverride = common.SEGMENT_DS
		case 0x26:
			// es segment override
			core.flags.MemorySegmentOverride = common.SEGMENT_ES
		case 0x64:
			// fs segment override
			core.flags.MemorySegmentOverride = common.SEGMENT_FS
		case 0x65:
			// gs segment override
			core.flags.MemorySegmentOverride = common.SEGMENT_GS
		case 0xf0:
			// lock prefix
			core.flags.LockPrefixEnabled = true
		case 0xf2:
			// repne/repnz prefix
		case 0xf3:
			// rep or repe/repz prefix
			core.flags.RepPrefixEnabled = true
		case 0x66:
			// operand size override
			core.flags.OperandSizeOverrideEnabled = true
		case 0x67:
			// address size override
			core.flags.AddressSizeOverrideEnabled = true
		}

		core.currentByteAddr++
		nextInstructionAddr++
	}

	core.memoryAccessController.SetSegmentOverride(core.flags.MemorySegmentOverride)
	core.memoryAccessController.SetAddressSizeOverride(core.flags.AddressSizeOverrideEnabled)
	core.memoryAccessController.SetOperandSizeOverride(core.flags.OperandSizeOverrideEnabled)
	core.memoryAccessController.SetLockPrefix(core.flags.LockPrefixEnabled)
	core.memoryAccessController.SetRepPrefix(core.flags.RepPrefixEnabled)

	instrByte, err = core.memoryAccessController.ReadMemoryAddr8(nextInstructionAddr)

	if err != nil {
		log.Printf("Error reading instruction byte: %s\n", err)
		doCoreDump(core)
		panic(0)
	}

	var instructionImpl OpCodeImpl
	if instrByte == 0x0F {
		// 2 byte opcode
		core.currentByteAddr++
		instrByte, err = core.memoryAccessController.ReadMemoryAddr8(uint32(core.currentByteAddr))
		if err != nil {
			log.Printf("Error reading instruction byte: %s\n", err)
			doCoreDump(core)
			panic(0)
		}

		core.currentOpCodeBeingExecuted = instrByte
		instructionImpl = core.opCodeMap2Byte[core.currentOpCodeBeingExecuted]

		if instructionImpl == nil {
			log.Printf("[%#04x] Unrecognised 2-bit opcode: %#2x %#2x\n", core.registers.IP, core.currentPrefixBytes, instrByte)
			doCoreDump(core)
			panic(0)
		}

		core.currentPrefixBytes = append(core.currentPrefixBytes, 0x0F)
		core.is2ByteOperand = true
	} else if instrByte == 0xFF {
		// 2 byte opcode dictated by modrm
		handleGroup5Opcode(core)
		return 0
	} else {
		core.currentOpCodeBeingExecuted = instrByte
		instructionImpl = core.opCodeMap[core.currentOpCodeBeingExecuted]
	}

	if instructionImpl != nil {
		instructionImpl(core)
	} else {
		log.Printf("[%#04x] Unrecognised opcode: %#2x %#2x\n", core.registers.IP, core.currentPrefixBytes, instrByte)

		log.Printf("CPU CORE ERROR!!!")

		doCoreDump(core)
		panic(0)
	}

	return 0
}

func isPrefixByte(b byte) bool {
	switch b {
	case 0x2e:
		// cs segment override
		return true
	case 0x36:
		// ss segment override
		return true
	case 0x3e:
		// ds segment override
		return true
	case 0x26:
		// es segment override
		return true
	case 0x64:
		// fs segment override
		return true
	case 0x65:
		// gs segment override
		return true
	case 0xf0:
		// lock prefix
		return true
	case 0xf2:
		// repne/repnz prefix
		return true
	case 0xf3:
		// rep or repe/repz prefix
		return true
	case 0x66:
		// operand size override
		return true
	case 0x67:
		// address size override
		return true
	}
	return false
}

func (core *CpuCore) Is32BitOperand() bool {
	if core.mode == common.PROTECTED_MODE {
		if core.GetCurrentSegmentWidth() == 32 || core.flags.OperandSizeOverrideEnabled {
			return true
		} else {
			return false
		}
	} else {
		if core.flags.OperandSizeOverrideEnabled {
			return true
		} else {
			return false
		}
	}
}

func handleGroup3OpCode_byte(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	if err != nil {
		log.Printf("Error consuming ModR/M byte: %v\n", err)
		return // Exit early on error
	}
	core.currentByteAddr--

	switch modrm.reg {
	/*
		reg = 0: TEST (Tests bits by performing a bitwise AND)
		reg = 2: NOT (Bitwise NOT)
		reg = 3: NEG (Two's complement negation)
		reg = 4: MUL (Unsigned multiply)
		reg = 5: IMUL (Signed multiply)
		reg = 6: DIV (Unsigned divide)
		reg = 7: IDIV (Signed divide)
	*/
	case 0:
		INSTR_TEST(core)
	case 2:
		INSTR_NOT(core)
	case 3:
		INSTR_NEG(core)
	case 4:
		INSTR_MUL(core)
	case 5:
		INSTR_IMUL(core)
	case 6:
		INSTR_DIV(core)
	case 7:
		INSTR_IDIV(core)

	default:
		log.Printf("INSTR_80_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm)
		doCoreDump(core)
		panic(0)
	}
}

func handleGroup5Opcode_word(core *CpuCore) {
	handleGroup5Opcode(core) //TODO! Implement this
}

func handleGroup5Opcode(core *CpuCore) {

	if core.Is32BitOperand() {
		handleGroup5Opcode_32(core)
		return
	}

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	if err != nil {
		log.Printf("Error consuming ModR/M byte: %v\n", err)
		return // Exit early on error
	}
	core.currentByteAddr--

	switch modrm.reg {
	case 0:
		if modrm.rm == 0 {
			// inc rm16
			INSTR_INC_RM16(core)
		} else {
			log.Println("Unexpected ModRM setup for INC instruction")
		}
	case 1:
		if modrm.rm == 1 {
			// dec rm16
			INSTR_DEC_RM16(core)
		} else {
			log.Println("Unexpected ModRM setup for DEC RM16 instruction")
		}
	case 2:
		if modrm.rm == 2 {
			// call rm32
			INSTR_CALL_RM16(core)
		} else {
			log.Println("Unexpected ModRM setup for CALL RM16 instruction")
		}
	case 3:
		if modrm.rm == 3 {
			// call m16:16
			INSTR_CALL_M16(core)
		} else {
			log.Println("Unexpected ModRM setup for CALL M16 instruction")
		}
	case 4:
		// jmp rm32
		INSTR_JMP_FAR_M16(core, &modrm)
	case 5:
		// jmp m16:16
		INSTR_JMP_FAR_M16(core, &modrm)
	case 6:
		// push rm32
		INSTR_PUSH_RM16(core)
	default:
		log.Printf("INSTR_FF_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm)
		doCoreDump(core)
		panic("Unhandled operation")
	}
}

func handleGroup5Opcode_32(core *CpuCore) {

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	if err != nil {
		log.Printf("Error consuming ModR/M byte: %v\n", err)
		return // Exit early on error
	}
	core.currentByteAddr--

	switch modrm.reg {
	case 0:
		if modrm.rm == 0 {
			// inc rm32
			INSTR_INC_RM32(core)
		} else {
			log.Println("Unexpected ModRM setup for INC instruction")
		}
	case 1:
		if modrm.rm == 1 {
			// dec rm32
			INSTR_DEC_RM32(core)
		} else {
			log.Println("Unexpected ModRM setup for DEC instruction")
		}
	case 2:
		if modrm.rm == 2 {
			// call rm32
			INSTR_CALL_RM32(core)
		} else {
			log.Println("Unexpected ModRM setup for CALL RM32 instruction")
		}
	case 3:
		if modrm.rm == 3 {
			// call m16:16
			INSTR_CALL_M16(core)
		} else {
			log.Println("Unexpected ModRM setup for CALL M16 instruction")
		}
	case 4:
		// jmp rm32
		INSTR_JMP_FAR_M32(core)
	case 5:
		// jmp m16:16
		INSTR_JMP_FAR_M16(core, &modrm)
	case 6:
		// push rm32
		INSTR_PUSH_32(core)
	default:
		log.Printf("INSTR_FF_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm)
		doCoreDump(core)
		panic("Unhandled operation")
	}
}

func INSTR_80_OPCODES(core *CpuCore) {

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	switch modrm.reg {
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
		log.Println(fmt.Sprintf("INSTR_80_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm))
		doCoreDump(core)
		panic(0)
	}

eof:
}

func INSTR_81_OPCODES(core *CpuCore) {

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

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
		log.Println(fmt.Sprintf("INSTR_81_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm))
		doCoreDump(core)
		panic(0)
	}
eof:
}

func (core *CpuCore) GetAddressSize() int {
	if core.mode == common.REAL_MODE {
		return 2
	} else {
		if core.flags.OperandSizeOverrideEnabled {
			return 2
		} else {
			return 4
		}
	}
}
