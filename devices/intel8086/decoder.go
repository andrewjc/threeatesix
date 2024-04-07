package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func (core *CpuCore) decodeInstruction() uint8 {

	var instrByte uint8
	var err error

	core.flags.MemorySegmentOverride = 0
	core.flags.OperandSizeOverrideEnabled = false
	core.flags.AddressSizeOverrideEnabled = false
	core.flags.LockPrefixEnabled = false
	core.flags.RepPrefixEnabled = false
	core.is2ByteOperand = true

	core.currentPrefixBytes = []byte{}
	for isPrefixByte(core.memoryAccessController.PeekNextBytes(uint32(core.currentByteAddr), 1)[0]) {

		prefixByte := core.memoryAccessController.PeekNextBytes(uint32(core.currentByteAddr), 1)[0]
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
	}

	instrByte, err = core.memoryAccessController.ReadAddr8(uint32(core.currentByteAddr))
	if err != nil {
		log.Printf("Error reading instruction byte: %s\n", err)
		doCoreDump(core)
		panic(0)
	}

	var instructionImpl OpCodeImpl
	if core.memoryAccessController.PeekNextBytes(uint32(core.currentByteAddr), 1)[0] == 0x0F {
		// 2 byte opcode
		core.IncrementIP()
		instrByte, err = core.memoryAccessController.ReadAddr8(uint32(core.currentByteAddr + 1))
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
	if core.mode == common.REAL_MODE && core.flags.OperandSizeOverrideEnabled {
		return true
	}

	if core.mode == common.PROTECTED_MODE && !core.flags.OperandSizeOverrideEnabled {
		return true
	}

	return false
}

func INSTR_SMSW(core *CpuCore) {
	var value uint16
	var rm_str string

	core.IncrementIP()
	modrm, bytesConsumed, err := core.consumeModRm()
	if err != nil {
		goto eof
	}
	core.currentByteAddr += bytesConsumed

	value = uint16(core.registers.CR0)

	err = core.writeRm16(&modrm, &value)

	if modrm.mod == 3 {
		rm_str = core.registers.index16ToString(modrm.rm)
	} else {
		rm_str = "r/m16"
	}
eof:
	core.logInstruction(fmt.Sprintf("[%#04x] smsw %s", core.GetCurrentlyExecutingInstructionAddress(), rm_str))
	core.registers.IP = uint16(core.GetIP() + 1)
}

func INSTR_FF_OPCODES(core *CpuCore) {

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	switch modrm.reg {
	/*case modrm.rm == 0:
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
	  	}*/
	//case 3:
	//	{
	// call m16
	//		INSTR_CALLF_M16(core, &modrm)
	//	}
	case 4:
		{
			// jmp rm32
			INSTR_JMP_FAR_M16(core, &modrm)
		}
	case 5:
		{
			// jmp m16
			INSTR_JMP_FAR_M16(core, &modrm)
		}
	case 6:
		{
			// push rm32
			INSTR_PUSH(core)
		}
	default:
		log.Println(fmt.Sprintf("INSTR_FF_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm))
		doCoreDump(core)
		panic(0)
	}
eof:
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

func INSTR_83_OPCODES(core *CpuCore) {

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
	case 5:
		INSTR_SUB(core)
	case 6:
		INSTR_XOR(core)
	case 7:
		INSTR_CMP(core)
	default:
		log.Println(fmt.Sprintf("INSTR_83_OPCODE UNHANDLED OPER: (modrm: base:%d, reg:%d, mod:%d, rm: %d)\n\n", modrm.base, modrm.reg, modrm.mod, modrm.rm))
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
