package intel8086

import (
	"fmt"
	"log"
)

func INSTR_TEST(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var term2 uint32
	var result uint32
	var dataSize uint8

	// No need for a bitLength variable if not used earlier in flag calculations.
	switch core.currentOpCodeBeingExecuted {
	case 0x84:
		// TEST r/m8, r8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			log.Println("Error consuming ModR/M byte: ", err)
			return
		}
		core.currentByteAddr += bytesConsumed

		rm8, rm8Str, err := core.readRm8(&modrm)
		if err != nil {
			log.Println("Error reading r/m8: ", err)
			return
		}
		r8, r8Str := core.readR8(&modrm)
		term1 = uint32(*rm8)
		term2 = uint32(*r8)
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] TEST %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str))

	case 0x85:
		// TEST r/m16, r16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			log.Println("Error consuming ModR/M byte: ", err)
			return
		}
		core.currentByteAddr += bytesConsumed

		rm16, rm16Str, err := core.readRm16(&modrm)
		if err != nil {
			log.Println("Error reading r/m16: ", err)
			return
		}
		r16, r16Str := core.readR16(&modrm)
		term1 = uint32(*rm16)
		term2 = uint32(*r16)
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] TEST %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm16Str, r16Str))

	case 0xA8:
		// TEST al, imm8
		term1 = uint32(core.registers.AL)
		imm8, err := core.readImm8()
		if err != nil {
			log.Println("Error reading immediate value: ", err)
			return
		}
		term2 = uint32(imm8)
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] TEST AL, %#02x", core.GetCurrentlyExecutingInstructionAddress(), imm8))

	case 0xA9:
		// TEST ax, imm16
		term1 = uint32(core.registers.AX)
		imm16, err := core.readImm16()
		if err != nil {
			log.Println("Error reading immediate value: ", err)
			return
		}
		term2 = uint32(imm16)
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] TEST AX, %#04x", core.GetCurrentlyExecutingInstructionAddress(), imm16))

	case 0xF6, 0xF7:
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			log.Println("Error consuming ModR/M byte: ", err)
			return
		}
		core.currentByteAddr += bytesConsumed

		if core.currentOpCodeBeingExecuted == 0xF6 {
			// TEST r/m8, imm8
			rm8, rmStr, err := core.readRm8(&modrm)
			if err != nil {
				log.Println("Error reading r/m8: ", err)
				return
			}
			imm8, err := core.readImm8()
			if err != nil {
				log.Println("Error reading immediate value: ", err)
				return
			}
			term1 = uint32(*rm8)
			term2 = uint32(imm8)
			core.logInstruction(fmt.Sprintf("[%#04x] TEST %s, %#02x", core.GetCurrentlyExecutingInstructionAddress(), rmStr, imm8))
		} else {
			// TEST r/m16, imm16
			rm16, rmStr, err := core.readRm16(&modrm)
			if err != nil {
				log.Println("Error reading r/m16: ", err)
				return
			}
			imm16, err := core.readImm16()
			if err != nil {
				log.Println("Error reading immediate value: ", err)
				return
			}
			term1 = uint32(*rm16)
			term2 = uint32(imm16)
			core.logInstruction(fmt.Sprintf("[%#04x] TEST %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), rmStr, imm16))
		}
		if core.currentOpCodeBeingExecuted == 0xF6 {
			dataSize = 8
		} else {
			dataSize = 16
		}

	default:
		core.logInstruction("Unsupported TEST opcode: %#x\n", core.currentOpCodeBeingExecuted)
		return
	}

	// Perform the bitwise AND, but don't store the result
	result = term1 & term2

	// Set flags
	core.registers.SetFlag(CarryFlag, false) // Always cleared
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(SignFlag, (result>>(dataSize-1))&0x01 == 1)
	core.registers.SetFlag(OverFlowFlag, false) // Always cleared

	// Update IP to reflect bytes read

}

func INSTR_XCHG(core *CpuCore) {
	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98:
		{
			// xchg ax, 16
			term1 := core.registers.AX
			r16, r16Str := core.registers.registers16Bit[core.currentOpCodeBeingExecuted-0x90], core.registers.index16ToString(core.currentOpCodeBeingExecuted-0x90)
			term2 := *r16
			core.registers.AX = term2
			*r16 = term1
			core.logInstruction(fmt.Sprintf("[%#04x] xchg AX, %s", core.GetCurrentlyExecutingInstructionAddress(), r16Str))
			goto eof
		}
	case 0x86:
		{
			// XCHG r/m8, r8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm8, rm8Str, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			r8, r8Str := core.readR8(&modrm)

			tmp := *rm8

			*rm8 = *r8

			*r8 = tmp

			core.logInstruction(fmt.Sprintf("[%#04x] xchg %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str))
			goto eof
		}
	case 0x87:
		{
			// XCHG r/m16, r16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm8, rm8Str, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			r8, r8Str := core.readR16(&modrm)

			tmp := *rm8

			*rm8 = *r8

			*r8 = tmp

			core.logInstruction(fmt.Sprintf("[%#04x] xchg %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str))
			goto eof
		}
	default:
		log.Println("Unrecognised xchg instruction!")
		doCoreDump(core)
	}

eof:
}
func INSTR_CMP(core *CpuCore) {
	var term1, term2, result uint32
	var dataSize uint8

	switch core.currentOpCodeBeingExecuted {
	case 0xA6: // CMPS m8, m8
		address1 := core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.SI)
		address2 := core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.DI)
		tmp1, err := core.memoryAccessController.ReadMemoryValue8(address1)
		if err != nil {
			return
		}
		tmp2, err := core.memoryAccessController.ReadMemoryValue8(address2)
		if err != nil {
			return
		}
		term1 = uint32(tmp1)
		term2 = uint32(tmp2)
		result = term1 - term2
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] CMP (m8) %d with %d", core.GetCurrentlyExecutingInstructionAddress(), tmp1, tmp2))

	case 0xA7: // CMPS m16, m16
		address1 := core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.SI)
		address2 := core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.DI)
		tmp1, err := core.memoryAccessController.ReadMemoryValue16(address1)
		if err != nil {
			return
		}
		tmp2, err := core.memoryAccessController.ReadMemoryValue16(address2)
		if err != nil {
			return
		}
		term1 = uint32(tmp1)
		term2 = uint32(tmp2)
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP (m16) %d with %d", core.GetCurrentlyExecutingInstructionAddress(), tmp1, tmp2))

	case 0x3A: // CMP r8, r/m8
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		r8, r8Str := core.readR8(&modrm)
		rm8, rm8Str, err := core.readRm8(&modrm)
		if err != nil {
			return
		}
		term1 = uint32(*r8)
		term2 = uint32(*rm8)
		result = term1 - term2
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %s", core.GetCurrentlyExecutingInstructionAddress(), r8Str, rm8Str))

	case 0x3B: // CMP r16, r/m16
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		r16, r16Str := core.readR16(&modrm)
		rm16, rm16Str, err := core.readRm16(&modrm)
		if err != nil {
			return
		}
		term1 = uint32(*r16)
		term2 = uint32(*rm16)
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %s", core.GetCurrentlyExecutingInstructionAddress(), r16Str, rm16Str))

	case 0x3C: // CMP AL, imm8
		core.currentByteAddr++
		imm8, err := core.readImm8()
		if err != nil {
			return
		}
		term1 = uint32(core.registers.AL)
		term2 = uint32(imm8)
		result = term1 - term2
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] CMP AL, %#02x", core.GetCurrentlyExecutingInstructionAddress(), imm8))

	case 0x3D: // CMP AX, imm16
		core.currentByteAddr++
		imm16, err := core.readImm16()
		if err != nil {
			return
		}
		term1 = uint32(core.registers.AX)
		term2 = uint32(imm16)
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP AX, %#04x", core.GetCurrentlyExecutingInstructionAddress(), imm16))

	case 0x38: // CMP r/m8, r8
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		rm8, rm8Str, err := core.readRm8(&modrm)
		if err != nil {
			return
		}
		r8, r8Str := core.readR8(&modrm)
		term1 = uint32(*rm8)
		term2 = uint32(*r8)
		result = term1 - term2
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str))

	case 0x39: // CMP r/m16, r16
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		rm16, rm16Str, err := core.readRm16(&modrm)
		if err != nil {
			return
		}
		r16, r16Str := core.readR16(&modrm)
		term1 = uint32(*rm16)
		term2 = uint32(*r16)
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm16Str, r16Str))

	case 0x80: // CMP r/m8, imm8
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		rm8, rm8Str, err := core.readRm8(&modrm)
		if err != nil {
			return
		}
		imm8, err := core.readImm8()
		if err != nil {
			return
		}
		term1 = uint32(*rm8)
		term2 = uint32(imm8)
		result = term1 - term2
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %#02x", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, imm8))

	case 0x81: // CMP r/m16, imm16
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		rm16, rm16Str, err := core.readRm16(&modrm)
		if err != nil {
			return
		}
		imm16, err := core.readImm16()
		if err != nil {
			return
		}
		term1 = uint32(*rm16)
		term2 = uint32(imm16)
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), rm16Str, imm16))

	case 0x83: // CMP r/m16, imm8
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		rm16, rm16Str, err := core.readRm16(&modrm)
		if err != nil {
			return
		}
		imm8, err := core.readImm8()
		if err != nil {
			return
		}
		term1 = uint32(*rm16)
		term2 = uint32(int8(imm8)) // Sign-extend the 8-bit immediate to 32-bit if necessary
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), rm16Str, imm8))
	case 0xBC:
		// CMP AL, imm8
		core.currentByteAddr++
		imm8, err := core.readImm8()
		if err != nil {
			return
		}
		term1 = uint32(core.registers.AL)
		term2 = uint32(imm8)
		result = term1 - term2
		dataSize = 8
		core.logInstruction(fmt.Sprintf("[%#04x] CMP AL, %#02x", core.GetCurrentlyExecutingInstructionAddress(), imm8))
	case 0xBD:
		// CMP AX, imm16
		core.currentByteAddr++
		imm16, err := core.readImm16()
		if err != nil {
			return
		}
		term1 = uint32(core.registers.AX)
		term2 = uint32(imm16)
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP AX, %#04x", core.GetCurrentlyExecutingInstructionAddress(), imm16))
	case 0x75: //via handleGroup80opcode
		// CMP AL, imm8
		core.currentByteAddr++
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed

		rm16, rm16Str, err := core.readRm8(&modrm)
		if err != nil {
			return
		}
		imm8, err := core.readImm8()
		if err != nil {
			return
		}
		term1 = uint32(*rm16)
		term2 = uint32(int8(imm8)) // Sign-extend the 8-bit immediate to 32-bit if necessary
		result = term1 - term2
		dataSize = 16
		core.logInstruction(fmt.Sprintf("[%#04x] CMP %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), rm16Str, imm8))
	default:
		fmt.Printf("Unsupported opcode %#x\n", core.currentOpCodeBeingExecuted)
		doCoreDump(core)
		panic(0)
	}

	// Update flags
	core.registers.SetFlag(CarryFlag, term1 < term2)
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(SignFlag, (result>>(dataSize-1))&0x01 == 1)
	core.registers.SetFlag(OverFlowFlag, ((term1^term2)&(term1^result)&(1<<(dataSize-1))) != 0)
}

// INSTR_TEST_IMM8_AL tests the AL register with an immediate 8-bit value.
func INSTR_TEST_IMM8_AL(core *CpuCore) {
	// Read the next 8-bit immediate value from the code segment
	core.currentByteAddr++
	imm8, err := core.readImm8()
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error reading immediate value for TEST: %s", err))
		return
	}

	// Perform the bitwise AND operation between AL and imm8
	result := uint8(core.registers.AL) & imm8

	// Log the instruction for debugging purposes
	core.logInstruction(fmt.Sprintf("[%#04x] TEST AL, %#02x", core.GetCurrentlyExecutingInstructionAddress(), imm8))

	// Update flags based on the result
	core.registers.SetFlag(CarryFlag, false) // Always cleared
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(SignFlag, (result>>7)&0x01 == 1)
	core.registers.SetFlag(OverFlowFlag, false) // Always cleared

}

// INSTR_TEST_IMM16_AX tests the AX register with an immediate 16-bit value.
func INSTR_TEST_IMM16_AX(core *CpuCore) {
	core.currentByteAddr++
	// Read the next 16-bit immediate value from the code segment
	imm16, err := core.readImm16()
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error reading immediate value for TEST: %s", err))
		return
	}

	// Perform the bitwise AND operation between AX and imm16
	result := uint16(core.registers.AX) & imm16

	// Log the instruction for debugging purposes
	core.logInstruction(fmt.Sprintf("[%#04x] TEST AX, %#04x", core.GetCurrentlyExecutingInstructionAddress(), imm16))

	// Update flags based on the result
	core.registers.SetFlag(CarryFlag, false) // Always cleared
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(SignFlag, (result>>15)&0x01 == 1)
	core.registers.SetFlag(OverFlowFlag, false) // Always cleared

}
