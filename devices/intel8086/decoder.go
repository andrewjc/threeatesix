package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func (core *CpuCore) resetFlags() {
	core.flags = CpuExecutionFlags{}
	core.is2ByteOperand = false
	core.currentPrefixBytes = []byte{}
}

func (core *CpuCore) handlePrefixes() {
	for {
		prefixByte, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr)
		if err != nil || !isPrefixByte(&prefixByte) {
			break
		}
		core.handlePrefix(prefixByte)
		core.currentByteAddr++
		core.currentPrefixBytes = append(core.currentPrefixBytes, prefixByte)
	}
}

func (core *CpuCore) handlePrefix(prefixByte byte) {
	switch prefixByte {
	case 0x2e:
		core.flags.MemorySegmentOverride = common.SEGMENT_CS
	case 0x36:
		core.flags.MemorySegmentOverride = common.SEGMENT_SS
	case 0x3e:
		core.flags.MemorySegmentOverride = common.SEGMENT_DS
	case 0x26:
		core.flags.MemorySegmentOverride = common.SEGMENT_ES
	case 0x64:
		core.flags.MemorySegmentOverride = common.SEGMENT_FS
	case 0x65:
		core.flags.MemorySegmentOverride = common.SEGMENT_GS
	case 0xf0:
		core.flags.LockPrefixEnabled = true
	case 0xf2, 0xf3:
		core.flags.RepPrefixEnabled = true
	case 0x66:
		core.flags.OperandSizeOverrideEnabled = true
	case 0x67:
		core.flags.AddressSizeOverrideEnabled = true
	}
}

func (core *CpuCore) applyFlagsToMemoryController() {
	core.memoryAccessController.SetSegmentOverride(core.flags.MemorySegmentOverride)
	core.memoryAccessController.SetAddressSizeOverride(core.flags.AddressSizeOverrideEnabled)
	core.memoryAccessController.SetOperandSizeOverride(core.flags.OperandSizeOverrideEnabled)
	core.memoryAccessController.SetLockPrefix(core.flags.LockPrefixEnabled)
	core.memoryAccessController.SetRepPrefix(core.flags.RepPrefixEnabled)
}

func (core *CpuCore) readInstructionByte() (byte, error) {
	return core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr)
}

func (core *CpuCore) handleInstructionReadError(err error) {
	core.logInstruction("Error reading instruction byte: %s\n", err)
	doCoreDump(core)
	panic(fmt.Sprintf("Instruction read error: %v", err))
}

func (core *CpuCore) handleNullInstruction() {
	core.logInstruction("Null instruction encountered")
	doCoreDump(core)
	panic("Null instruction encountered")
}

func (core *CpuCore) handle2ByteOpcode() OpCodeImpl {
	core.currentByteAddr++
	secondByte, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr)
	if err != nil {
		core.handleInstructionReadError(err)
		return nil
	}
	core.currentByteAddr++

	core.currentOpCodeBeingExecuted = secondByte
	instructionImpl := core.opCodeMap2Byte[core.currentOpCodeBeingExecuted]

	if instructionImpl == nil {
		core.logInstruction("[%#04x] Unrecognised 2-byte opcode: 0x0F %#02x\n", core.registers.IP, secondByte)
		doCoreDump(core)
		panic(fmt.Sprintf("Unrecognized 2-byte opcode: 0x0F %#02x", secondByte))
	}

	core.currentPrefixBytes = append(core.currentPrefixBytes, 0x0F)
	core.is2ByteOperand = true

	return instructionImpl
}

func (core *CpuCore) handleUnrecognizedOpcode(instrByte byte) {
	core.logDebug(fmt.Sprintf("[%#04x] Unrecognised opcode: %#02x %v\n", core.registers.IP, instrByte, core.currentPrefixBytes))
	core.logDebug("CPU CORE ERROR!!!")
	doCoreDump(core)
	panic(fmt.Sprintf("Unrecognized opcode: %#02x", instrByte))
}

func (core *CpuCore) updateInstructionPointer() {
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}

func (core *CpuCore) decodeInstruction() uint8 {

	var instrByte uint8
	var err error
	nextInstructionAddr := core.SegmentAddressToLinearAddress32(core.registers.CS, uint32(core.registers.IP))
	core.currentByteAddr = nextInstructionAddr
	core.currentByteDecodeStart = nextInstructionAddr

	core.resetFlags()
	core.handlePrefixes()
	core.applyFlagsToMemoryController()

	instrByte, err = core.readInstructionByte()
	if err != nil {
		core.handleInstructionReadError(err)
		return 0
	}

	var instructionImpl OpCodeImpl
	switch instrByte {
	case 0x00:
		core.handleNullInstruction()
		return 0
	case 0x0F:
		instructionImpl = core.handle2ByteOpcode()
	case 0xFF:
		handleGroup5Opcode(core)
		return 0
	case 0x80:
		handleGroup80opcode(core)
		return 0
	case 0x81:
		handleGroup81opcode(core)
		return 0
	default:
		core.currentOpCodeBeingExecuted = instrByte
		instructionImpl = core.opCodeMap[core.currentOpCodeBeingExecuted]
	}

	if instructionImpl != nil {
		instructionImpl(core)
	} else {
		core.handleUnrecognizedOpcode(instrByte)
	}

	core.updateInstructionPointer()
	return 0
}

func isPrefixByte(b *uint8) bool {
	switch *b {
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
		return core.GetCurrentSegmentWidth() == 32 && !core.flags.OperandSizeOverrideEnabled
	} else {
		return false
	}
}

func handleGroup3OpCode_byte(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	if err != nil {
		core.logInstruction("Error consuming ModR/M byte: %v\n", err)
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
		core.logInstruction("INSTR_80_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm)
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
		core.logInstruction(fmt.Sprintf("Error consuming ModR/M byte: %v\n", err))
		return // Exit early on error
	}
	core.currentByteAddr--

	switch modrm.reg {
	case 0:
		// INC rm16
		INSTR_INC_RM16(core)
	case 1:
		// DEC rm16
		INSTR_DEC_RM16(core)
	case 2:
		// CALL rm16
		INSTR_CALL_RM16(core)
	case 3:
		// CALL m16:16
		INSTR_CALL_M16(core)
	case 4:
		// JMP rm16
		INSTR_JMP_FAR_M16(core)
	case 5:
		// JMP m16:16
		INSTR_JMP_FAR_M16(core)
	case 6:
		// PUSH rm16
		INSTR_PUSH_RM16(core)
	case 7:
		core.logInstruction(fmt.Sprintf("Invalid Group 5 opcode: reg = 7 is undefined\n"))
	default:
		core.logInstruction(fmt.Sprintf("INSTR_FF_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm))
	}

	// Update IP
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}

func handleGroup5Opcode_32(core *CpuCore) {

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error consuming ModR/M byte: %v\n", err))
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
		INSTR_JMP_FAR_M16(core)
	case 6:
		// push rm32
		INSTR_PUSH_32(core)
	default:
		core.logInstruction("INSTR_FF_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm)
		doCoreDump(core)
		panic("Unhandled operation")
	}
}

func handleGroup80opcode(core *CpuCore) {

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

func handleGroup81opcode(core *CpuCore) {

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
