package intel8086

import (
	"log"
	"math/bits"
)

func INSTR_TEST(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32
	core.currentByteAddr++
	switch core.currentOpCodeBeingExecuted {
	case 0xA8:
		{
			//  TEST al, imm8
			term1 = uint32(core.registers.AL)
			term2 = uint32(core.readImm8())

			log.Printf("[%#04x] test al, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0xA9:
		{
			// TEST ax, imm16
			term1 = uint32(core.registers.AX)
			term2 = uint32(core.readImm16())

			log.Printf("[%#04x] test ax, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0xF6:
		{
			// TEST r/m8, imm8
			modrm, bytesConsumed := core.consumeModRm()
			rm, rmStr := core.readRm8(&modrm)
			core.currentByteAddr += bytesConsumed

			term1 = uint32(*rm)

			term2 = uint32(core.readImm8())

			log.Printf("[%#04x] test %s, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), rmStr, term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0xF7:
		{
			// TEST r/m16, imm16
			modrm, bytesConsumed := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			core.currentByteAddr += bytesConsumed

			term1 = uint32(*rm)

			term2 = uint32(core.readImm16())

			log.Printf("[%#04x] test %s, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), rmStr, term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x84:
		{
			// TEST r/m8, r8
			modrm, bytesConsumed := core.consumeModRm()
			rm, rmStr := core.readRm8(&modrm)
			core.currentByteAddr += bytesConsumed

			term1 = uint32(*rm)

			rm2, rm2Str := core.readR8(&modrm)

			term2 = uint32(*rm2)

			log.Printf("[%#04x] test %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rmStr, rm2Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x85:
		{
			// TEST r/m16, r16
			modrm, bytesConsumed := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			core.currentByteAddr += bytesConsumed

			term1 = uint32(*rm)

			rm2, rm2Str := core.readR16(&modrm)

			term2 = uint32(*rm2)

			log.Printf("[%#04x] test %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rmStr, rm2Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	}

	bitLength = uint32(bits.Len32(result))

	result = term1 & term2

	core.registers.SetFlag(OverFlowFlag,  false)

	core.registers.SetFlag(CarryFlag, false)

	core.registers.SetFlag(SignFlag, (result >> bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(ParityFlag, bits.OnesCount32(result)%2 == 0)

}


func INSTR_XCHG(core *CpuCore) {
	core.IncrementIP()

	switch core.currentOpCodeBeingExecuted {
	case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98:
		{
			// xchg ax, 16
			term1 := core.registers.AX
			r16, r16Str := core.registers.registers16Bit[core.currentOpCodeBeingExecuted-0x90], core.registers.index16ToString(core.currentOpCodeBeingExecuted-0x90)
			term2 := *r16
			core.registers.AX = term2
			*r16 = term1

			log.Printf("[%#04x] xchg AX, %s", core.GetCurrentlyExecutingInstructionAddress(), r16Str)
		}
	case 0x86:
		{
			// XCHG r/m8, r8
			modrm, bytesConsumed := core.consumeModRm()
			rm8, rm8Str := core.readRm8(&modrm)
			core.currentByteAddr += bytesConsumed

			r8, r8Str := core.readR8(&modrm)

			tmp := *rm8

			*rm8 = *r8

			*r8 = tmp

			log.Printf("[%#04x] xchg %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x87:
		{
			// XCHG r/m16, r16
			modrm, bytesConsumed := core.consumeModRm()
			rm8, rm8Str := core.readRm16(&modrm)
			core.currentByteAddr += bytesConsumed

			r8, r8Str := core.readR16(&modrm)

			tmp := *rm8

			*rm8 = *r8

			*r8 = tmp

			log.Printf("[%#04x] xchg %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	default:
		log.Println("Unrecognised xchg instruction!")
		doCoreDump(core)
	}

}

func INSTR_CMP(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentOpCodeBeingExecuted {
	case 0xA6:
		{
			//  CMPS m8, m8
			term1 = uint32(core.memoryAccessController.ReadAddr8(core.BuildAddress(core.registers.DS, core.registers.SI)))
			term2 = uint32(core.memoryAccessController.ReadAddr8(core.BuildAddress(core.registers.DS, core.registers.DI)))
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp m8, m8", core.GetCurrentlyExecutingInstructionAddress())
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0xA7:
		{
			// CMPS m16, m16
			term1 = uint32(core.memoryAccessController.ReadAddr16(core.BuildAddress(core.registers.DS, core.registers.SI)))
			term2 = uint32(core.memoryAccessController.ReadAddr16(core.BuildAddress(core.registers.DS, core.registers.DI)))
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp m16, m16", core.GetCurrentlyExecutingInstructionAddress())
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x3C:
		{
			// CMP AL, imm8
			term1 = uint32(core.registers.AL)
			term2 = uint32(core.readImm8())
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp AL, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x3D:
		{
			//	CMP AX, imm16
			term1 = uint32(core.registers.AX)
			term2 = uint32(core.readImm16())
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp AX, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x80:
		{
			// CMP r/m8, imm8
			modrm, bytesConsumed := core.consumeModRm()
			rm8, rm8Str := core.readRm8(&modrm)
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm8)
			term2 = uint32(core.readImm8())
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp %s, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x81:
		{
			// CMP r/m16, imm16
			modrm, bytesConsumed := core.consumeModRm()
			rm8, rm8Str := core.readRm16(&modrm)
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm8)
			term2 = uint32(core.readImm16())
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp %s, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x83:
		{
			// CMP r/m16,imm8
			modrm, bytesConsumed := core.consumeModRm()
			rm8, rm8Str := core.readRm16(&modrm)
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm8)
			term2 = uint32(core.readImm8())
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp %s, [%#04x]", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, term2)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x38:
		{
			// CMP r/m8,r8
			modrm, bytesConsumed := core.consumeModRm()
			rm8, rm8Str := core.readRm8(&modrm)
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm8)
			r8, r8Str := core.readR8(&modrm)
			term2 = uint32(*r8)
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x39:
		{
			// CMP r/m16,r16
			modrm, bytesConsumed := core.consumeModRm()
			rm8, rm8Str := core.readRm16(&modrm)
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm8)
			r8, r8Str := core.readR16(&modrm)
			term2 = uint32(*r8)
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rm8Str, r8Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x3A:
		{
			// CMP r8,r/m8
			modrm, bytesConsumed := core.consumeModRm()
			r8, r8Str := core.readR8(&modrm)
			term1 = uint32(*r8)

			rm8, rm8Str := core.readRm8(&modrm)
			core.currentByteAddr += bytesConsumed
			term2 = uint32(*rm8)
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionAddress(), r8Str, rm8Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	case 0x3B:
		{
			// CMP r16,r/m16
			modrm, bytesConsumed := core.consumeModRm()
			r8, r8Str := core.readR16(&modrm)
			term1 = uint32(*r8)

			rm8, rm8Str := core.readRm16(&modrm)
			core.currentByteAddr += bytesConsumed
			term2 = uint32(*rm8)
			result = uint32(term1) - uint32(term2)

			log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionAddress(), r8Str, rm8Str)
			core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.SetFlag(CarryFlag, result>>(bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr != 0)

	core.registers.SetFlag(OverFlowFlag,  (sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0))

}
