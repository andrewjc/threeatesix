package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
    "math/bits"
)

/* CPU OPCODE IMPLEMENTATIONS */

func mapOpCodes(c *CpuCore) {

	c.opCodeMap[0xEA] = INSTR_JMP_FAR_PTR16

	c.opCodeMap[0xE9] = INSTR_JMP_NEAR_REL16

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



	c.opCodeMap[0xB0] = INSTR_MOV
	c.opCodeMap[0xB1] = INSTR_MOV
	c.opCodeMap[0xB2] = INSTR_MOV
	c.opCodeMap[0xB3] = INSTR_MOV
	c.opCodeMap[0xB4] = INSTR_MOV
	c.opCodeMap[0xB5] = INSTR_MOV
	c.opCodeMap[0xB6] = INSTR_MOV
	c.opCodeMap[0xB7] = INSTR_MOV
	c.opCodeMap[0xB8] = INSTR_MOV

	c.opCodeMap[0xBA] = INSTR_MOV
	c.opCodeMap[0xBB] = INSTR_MOV
	c.opCodeMap[0xBC] = INSTR_MOV

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

	c.opCodeMap[0xD1] = INSTR_SHL

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
		c.opCodeMap[0x50 + i] = INSTR_PUSH
	}
	c.opCodeMap[0x50] = INSTR_PUSH

	c.opCodeMap[0x6A] = INSTR_PUSH
	c.opCodeMap[0x68] = INSTR_PUSH
	c.opCodeMap[0x0E] = INSTR_PUSH
	c.opCodeMap[0x16] = INSTR_PUSH
	c.opCodeMap[0x1E] = INSTR_PUSH
	c.opCodeMap[0x06] = INSTR_PUSH

	// opcode 0F A0 and 0F A8 are push .. todo

}

type OpCodeImpl func(*CpuCore)


func INSTR_FF_OPCODES(core *CpuCore) {

	core.IncrementIP()
	modrm := consumeModRm(core)

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

	core.IncrementIP()
	modrm := consumeModRm(core)

	switch modrm.reg {
	case 1:
		INSTR_OR(core)
	case 4:
		INSTR_AND(core)
	case 7:
		INSTR_CMP(core)
	default:
		log.Println(fmt.Sprintf("INSTR_80_OPCODE UNHANDLED OPTION %d\n\n", modrm))
		doCoreDump(core)
	}
}

func INSTR_81_OPCODES(core *CpuCore) {

	core.IncrementIP()
	modrm := consumeModRm(core)

	switch modrm.reg {
	case 0:
		INSTR_ADD(core)
	case 1:
		INSTR_OR(core)
	case 4:
		INSTR_AND(core)
	case 7:
		INSTR_CMP(core)
	default:
		log.Println(fmt.Sprintf("INSTR_81_OPCODE UNHANDLED OPTION %d\n\n", modrm))
		doCoreDump(core)
	}
}

func INSTR_83_OPCODES(core *CpuCore) {

	core.IncrementIP()
	modrm := consumeModRm(core)

	switch modrm.reg {
	case 0:
		INSTR_ADD(core)
	case 1:
		INSTR_OR(core)
	case 4:
		INSTR_AND(core)
	case 5:
		INSTR_SUB(core)
	case 7:
		INSTR_CMP(core)
	default:
		log.Println(fmt.Sprintf("INSTR_83_OPCODE UNHANDLED OPTION %d\n\n", modrm))
		doCoreDump(core)
	}
}


func INSTR_PUSH(core *CpuCore) {
	core.IncrementIP()


	switch core.currentByteAtCodePointer {
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
		{
			// PUSH r16
			val, valName := core.registers.registers16Bit[core.currentByteAtCodePointer - 0x50], core.registers.index16ToString(core.currentByteAtCodePointer - 0x50)

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionPointer(), valName)

		}
	case 0x6A:
		{
			// PUSH imm8

			core.IncrementIP()

			val := core.readImm8()

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr8(uint32(core.registers.SP), val)

			log.Printf("[%#04x] push %#04x", core.GetCurrentlyExecutingInstructionPointer(), val)
		}
	case 0x68:
		{
			// PUSH imm16

			core.IncrementIP()

			val := core.readImm16()

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), val)

			log.Printf("[%#04x] push %#04x", core.GetCurrentlyExecutingInstructionPointer(), val)
		}
	case 0x0E:
		{
			// PUSH CS

			core.IncrementIP()

			val := core.registers.registers16Bit[core.registers.CS]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionPointer(), "CS")
		}
	case 0x16:
		{
			// PUSH SS
			core.IncrementIP()

			val := core.registers.registers16Bit[core.registers.SS]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionPointer(), "SS")
		}
	case 0x1E:
		{
			// PUSH DS
			core.IncrementIP()

			val := core.registers.registers16Bit[core.registers.DS]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionPointer(), "DS")
		}
	case 0x06:
		{
			// PUSH ES
			core.IncrementIP()

			val := core.registers.registers16Bit[core.registers.ES]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionPointer(), "ES")
		}
	default:
		log.Println(fmt.Printf("Unhandled PUSH instruction: %s", core.currentByteAtCodePointer))
		doCoreDump(core)
	}
}

func INSTR_OR(core *CpuCore) {
	core.IncrementIP()

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x0c:
		{
			// OR AL,imm8
			term2 = uint32(core.readImm8())
			term1 = uint32(core.registers.AL)
			result = uint32(term1) | uint32(term2)
			core.registers.AL = uint8(result)

			log.Printf("[%#04x] or al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x0d:
		{
			// OR AX,imm16
			term2 = uint32(core.readImm16())
			term1 = uint32(core.registers.AX)
			result = uint32(term1) | uint32(term2)
			core.registers.AX = uint16(result)

			log.Printf("[%#04x] or ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x80:
		{
			// OR r/m8,imm8
			modrm := consumeModRm(core)
			rm, rmStr := core.readRm8(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm8())
			tmp := uint8( uint32(term1) | uint32(term2) )

			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x81:
		{
			// OR r/m16,imm16
			modrm := consumeModRm(core)
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm16())
			tmp := uint16( uint32(term1) | uint32(term2) )

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x83:
		{
			// OR r/m16,imm8
			modrm := consumeModRm(core)
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm8())
			tmp := uint16( uint32(term1) | uint32(term2) )

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x08:
		{
			// OR r/m8,r8
			modrm := consumeModRm(core)
			rm, rmStr := core.readRm8(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR8(&modrm)
			term2 = uint32(*rm2)
			tmp := uint8( uint32(term1) | uint32(term2) )

			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x09:
		{
			// OR r/m16,r16
			modrm := consumeModRm(core)
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR16(&modrm)
			term2 = uint32(*rm2)
			tmp := uint16( uint32(term1) | uint32(term2) )

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x0A:
		{
			// OR r8,r/m8
			modrm := consumeModRm(core)
			rm, rmStr := core.readR8(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readRm8(&modrm)
			term2 = uint32(*rm2)
			tmp := uint8( uint32(term1) | uint32(term2) )

			core.writeR8(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x0B:
		{
			// OR r16,r/m16
			modrm := consumeModRm(core)
			rm, rmStr := core.readR16(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readRm16(&modrm)
			term2 = uint32(*rm2)
			tmp := uint16( uint32(term1) | uint32(term2) )

			core.writeR16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags

	core.registers.OF = 0
	core.registers.CF = 0
	core.registers.SF = uint16( common.Bool2Uint8((result >> bitLength) != 0 ) )
	core.registers.ZF = uint16( common.Bool2Uint8( result == 0 ) )
	core.registers.PF = uint16( common.Bool2Uint8( bits.OnesCount32(result) % 2 == 0 ) )
}

func INSTR_NOP(core *CpuCore) {
	// Clear interrupts
	log.Printf("[%#04x] NOP", core.GetCurrentlyExecutingInstructionPointer())

	core.registers.IP = uint16(core.GetIP() + 1)
}


func INSTR_RET_NEAR(core *CpuCore) {

	log.Printf("[%#04x] retn", core.GetCurrentCodePointer())

	stackPntrAddr := core.registers.SP

	core.registers.IP = uint16(stackPntrAddr)

	core.registers.SP += 2
}


func INSTR_CLI(core *CpuCore) {
	// Clear interrupts
	log.Printf("[%#04x] TODO: Write CLI (Clear interrupts implementation!", core.GetCurrentCodePointer())

	core.registers.IP = uint16(uint16(core.GetIP() + 1))
}

func INSTR_CLD(core *CpuCore) {
	// Clear direction flag
	core.registers.DF = 0
	log.Printf("[%#04x] CLD", core.GetCurrentCodePointer())
	core.registers.IP = uint16(uint16(core.GetIP() + 1))
}

func INSTR_TEST(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0xA8: {
		//  TEST al, imm8
		core.IncrementIP()
		term1 = uint32(core.registers.AL)
		term2 = uint32(core.readImm8())

		log.Printf("[%#04x] test al, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0xA9: {
		// TEST ax, imm16
		core.IncrementIP()
		term1 = uint32(core.registers.AX)
		term2 = uint32(core.readImm16())

		log.Printf("[%#04x] test ax, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0xF6: {
		// TEST r/m8, imm8
		core.IncrementIP()

		modrm := consumeModRm(core)
		rm, rmStr := core.readRm8(&modrm)

		term1 = uint32(*rm)

		term2 = uint32(core.readImm8())

		log.Printf("[%#04x] test %s, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
	}
	case 0xF7: {
		// TEST r/m16, imm16
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm, rmStr := core.readRm16(&modrm)

		term1 = uint32(*rm)

		term2 = uint32(core.readImm16())

		log.Printf("[%#04x] test %s, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
	}
	case 0x84: {
		// TEST r/m8, r8
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm, rmStr := core.readRm8(&modrm)

		term1 = uint32(*rm)

		rm2, rm2Str := core.readR8(&modrm)

		term2 = uint32(*rm2)

		log.Printf("[%#04x] test %s, %s", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
	}
	case 0x85: {
		// TEST r/m16, r16
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm, rmStr := core.readRm16(&modrm)

		term1 = uint32(*rm)

		rm2, rm2Str := core.readR16(&modrm)

		term2 = uint32(*rm2)

		log.Printf("[%#04x] test %s, %s", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
	}
	}

	bitLength = uint32(bits.Len32(result))


	tmp := term1 & term2
	core.registers.SF = uint16((tmp >> bitLength) & 1)

	if tmp == 0 {
		core.registers.ZF = 1
	} else {
		core.registers.ZF = 0
	}

	core.registers.PF = 1
	for i := uint32(0); i < 8; i++ {
		core.registers.PF ^= uint16((tmp >> bitLength) & i)
	}

	core.registers.CF = 0
	core.registers.OF = 0

}

func INSTR_CMP(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0xA6: {
		//  CMPS m8, m8
		core.IncrementIP()
		term1 = uint32(core.memoryAccessController.ReadAddr8(core.BuildAddress(core.registers.DS, core.registers.SI)))
		term2 = uint32(core.memoryAccessController.ReadAddr8(core.BuildAddress(core.registers.DS, core.registers.DI)))
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp m8, m8", core.GetCurrentlyExecutingInstructionPointer())
	}
	case 0xA7: {
		// CMPS m16, m16
		core.IncrementIP()
		term1 = uint32(core.memoryAccessController.ReadAddr16(core.BuildAddress(core.registers.DS, core.registers.SI)))
		term2 = uint32(core.memoryAccessController.ReadAddr16(core.BuildAddress(core.registers.DS, core.registers.DI)))
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp m16, m16", core.GetCurrentlyExecutingInstructionPointer())
	}
	case 0x3C: {
		// CMP AL, imm8
		core.IncrementIP()
		term1 = uint32(core.registers.AL)
		term2 = uint32(core.readImm8())
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp AL, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x3D: {
		//	CMP AX, imm16
		core.IncrementIP()
		term1 = uint32(core.registers.AX)
		term2 = uint32(core.readImm16())
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp AX, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x80: {
		// CMP r/m8, imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm8, rm8Str := core.readRm8(&modrm)
		term1 = uint32(*rm8)
		term2 = uint32(core.readImm8())
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp %s, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), rm8Str, term2)
	}
	case 0x81: {
		// CMP r/m16, imm16
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm8, rm8Str := core.readRm16(&modrm)
		term1 = uint32(*rm8)
		term2 = uint32(core.readImm16())
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp %s, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), rm8Str, term2)
	}
	case 0x83: {
		// CMP r/m16,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm8, rm8Str := core.readRm16(&modrm)
		term1 = uint32(*rm8)
		term2 = uint32(core.readImm8())
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp %s, [%#04x]", core.GetCurrentlyExecutingInstructionPointer(), rm8Str, term2)
	}
	case 0x38: {
		// CMP r/m8,r8
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm8, rm8Str := core.readRm8(&modrm)
		term1 = uint32(*rm8)

		r8, r8Str := core.readR8(&modrm)
		term2 = uint32(*r8)
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionPointer(), rm8Str, r8Str)
	}
	case 0x39: {
		// CMP r/m16,r16
		core.IncrementIP()
		modrm := consumeModRm(core)
		rm8, rm8Str := core.readRm16(&modrm)
		term1 = uint32(*rm8)

		r8, r8Str := core.readR16(&modrm)
		term2 = uint32(*r8)
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionPointer(), rm8Str, r8Str)
	}
	case 0x3A: {
		// CMP r8,r/m8
		core.IncrementIP()
		modrm := consumeModRm(core)
		r8, r8Str := core.readR8(&modrm)
		term1 = uint32(*r8)

		rm8, rm8Str := core.readRm8(&modrm)
		term2 = uint32(*rm8)
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionPointer(), r8Str, rm8Str)
	}
	case 0x3B: {
		// CMP r16,r/m16
		core.IncrementIP()
		modrm := consumeModRm(core)
		r8, r8Str := core.readR16(&modrm)
		term1 = uint32(*r8)

		rm8, rm8Str := core.readRm16(&modrm)
		term2 = uint32(*rm8)
		result = uint32(term1) - uint32(term2)

		log.Printf("[%#04x] cmp %s, %s", core.GetCurrentlyExecutingInstructionPointer(), r8Str, rm8Str)
	}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.CF = uint16(common.Bool2Uint8( result >> (bitLength) != 0 ))

	core.registers.ZF = uint16(common.Bool2Uint8(result == 0))

	core.registers.SF = uint16(common.Bool2Uint8(signr != 0))

	core.registers.OF = uint16(common.Bool2Uint8((sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0)))

}

func INSTR_ADD(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x04: {
		// 	add AL,imm8
		core.IncrementIP()
		term2 = uint32(core.readImm8())
		term1 = uint32(core.registers.AL)
		result = uint32(term1) + uint32(term2)
		core.registers.AL = uint8(term1)

		log.Printf("[%#04x] add al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x05: {
		// 		add AX,imm16
		core.IncrementIP()
		term2 = uint32(core.readImm16())
		term1 = uint32(core.registers.AX)
		result = uint32(term1) + uint32(term2)
		core.registers.AX = uint16(term1)

		log.Printf("[%#04x] add ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x80: {
		// add r/m8,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) + uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x81: {
		// add r/m16,imm16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm16())
		result = uint32(term1) + uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x83: {
		// add r/m16,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) + uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x00: {
		// add r/m8,r8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x01: {
		// add r/m16,r16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x02: {
		// add r8,r/m8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x03: {
		// add r16,r/m16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2)
		tmp := uint16(result)
		core.writeR16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.CF = uint16(common.Bool2Uint8( result >> (bitLength) != 0 ))

	core.registers.ZF = uint16(common.Bool2Uint8(result == 0))

	core.registers.SF = uint16(common.Bool2Uint8(signr != 0))

	core.registers.OF = uint16(common.Bool2Uint8((sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0)))

}

func INSTR_AND(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x24: {
		// 	and AL,imm8
		core.IncrementIP()
		term2 = uint32(core.readImm8())
		term1 = uint32(core.registers.AL)
		result = uint32(term1) & uint32(term2)
		core.registers.AL = uint8(term1)

		log.Printf("[%#04x] add al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x25: {
		// 	and AX,imm16
		core.IncrementIP()
		term2 = uint32(core.readImm16())
		term1 = uint32(core.registers.AX)
		result = uint32(term1) & uint32(term2)
		core.registers.AX = uint16(term1)

		log.Printf("[%#04x] add ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x80: {
		// and r/m8,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) & uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x81: {
		// and r/m16,imm16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm16())
		result = uint32(term1) & uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x83: {
		// and r/m16,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) & uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x20: {
		// and r/m8,r8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) & uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x21: {
		// and r/m16,r16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) & uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x22: {
		// add r8,r/m8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) & uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x23: {
		// add r16,r/m16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) & uint32(term2)
		tmp := uint16(result)
		core.writeR16(&modrm, &tmp )

		log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.CF = uint16(common.Bool2Uint8( result >> (bitLength) != 0 ))

	core.registers.ZF = uint16(common.Bool2Uint8(result == 0))

	core.registers.SF = uint16(common.Bool2Uint8(signr != 0))

	core.registers.OF = uint16(common.Bool2Uint8((sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0)))

}

func INSTR_ADC(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x14: {
		// 	adc AL,imm8
		core.IncrementIP()
		term2 = uint32(core.readImm8())
		term1 = uint32(core.registers.AL)
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		core.registers.AL = uint8(term1)

		log.Printf("[%#04x] adc al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x15: {
		// adc AX,imm16
		core.IncrementIP()
		term2 = uint32(core.readImm16())
		term1 = uint32(core.registers.AX)
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		core.registers.AX = uint16(term1)

		log.Printf("[%#04x] adc ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x80: {
		// adc r/m8,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x81: {
		// adc r/m16,imm16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm16())
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x83: {
		// adc r/m16,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x10: {
		// adc r/m8,r8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x11: {
		// adc r/m16,r16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x12: {
		// adc r8,r/m8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x13: {
		// adc r16,r/m16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) + uint32(term2) + uint32(core.registers.CF)
		tmp := uint16(result)
		core.writeR16(&modrm, &tmp )

		log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.CF = uint16(common.Bool2Uint8( result >> (bitLength) != 0 ))

	core.registers.ZF = uint16(common.Bool2Uint8(result == 0))

	core.registers.SF = uint16(common.Bool2Uint8(signr != 0))

	core.registers.OF = uint16(common.Bool2Uint8((sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0)))

}

func INSTR_SUB(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x2c: {
		// 	SUB AL,imm8
		core.IncrementIP()
		term2 = uint32(core.readImm8())
		term1 = uint32(core.registers.AL)
		result = uint32(term1) - uint32(term2)
		core.registers.AL = uint8(term1)

		log.Printf("[%#04x] sub al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x2d: {
		// 		SUB AX,imm16
		core.IncrementIP()
		term2 = uint32(core.readImm16())
		term1 = uint32(core.registers.AX)
		result = uint32(term1) - uint32(term2)
		core.registers.AX = uint16(term1)

		log.Printf("[%#04x] sub ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
	}
	case 0x80: {
		// SUB r/m8,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) - uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] sub %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x81: {
		// SUB r/m16,imm16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm16())
		result = uint32(term1) - uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] sub %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x83: {
		// SUB r/m16,imm8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(core.readImm8())
		result = uint32(term1) - uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] sub %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
	}
	case 0x28: {
		// SUB r/m8,r8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) - uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x29: {
		// SUB r/m16,r16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readRm16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readR16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) - uint32(term2)
		tmp := uint16(result)
		core.writeRm16(&modrm, &tmp )

		log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x2A: {
		// SUB r8,r/m8
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR8(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm8(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) - uint32(term2)
		tmp := uint8(result)
		core.writeRm8(&modrm, &tmp )

		log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	case 0x2B: {
		// SUB r16,r/m16
		core.IncrementIP()
		modrm := consumeModRm(core)
		t1, t1Name := core.readR16(&modrm)
		term1 = uint32(*t1)
		t2, t2Name := core.readRm16(&modrm)
		term2 = uint32(*t2)
		result = uint32(term1) - uint32(term2)
		tmp := uint16(result)
		core.writeR16(&modrm, &tmp )

		log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
	}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.CF = uint16(common.Bool2Uint8( result >> (bitLength) != 0 ))

	core.registers.ZF = uint16(common.Bool2Uint8(result == 0))

	core.registers.SF = uint16(common.Bool2Uint8(signr != 0))

	core.registers.OF = uint16(common.Bool2Uint8((sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0)))

}

func INSTR_SHL(core *CpuCore) {

	core.IncrementIP()
	modrm := consumeModRm(core)

	value := *core.registers.registers16Bit[modrm.rm]

	regName := core.registers.index16ToString(modrm.rm)
	count := 1

	tmpCount := count & 0x1f
	tempDest := value

	for tmpCount != 0 {
		if modrm.reg == 4 { // SAL OR SHL
			core.registers.CF = (value >> 8) & 1
			value = value << 1
		} else {
			// TODO: if SAR then signed divide, rounding towards negative infinity, otherwise SHR - unsigned divide
			core.registers.CF = (value >> 0) & 1
			value = value / 2
		}
		tmpCount = tmpCount -1
	}

	if count & 0x1f == 1 {
		if modrm.reg == 4 { // SAL OR SHL
			core.registers.OF = ((value >> 8) & 1) ^ core.registers.CF
		} else if modrm.reg == 7 {
			core.registers.OF = 0
		} else {
			core.registers.OF =  (tempDest >> 8) & 1
		}
	}


	log.Printf("[%#04x] shl %s, %d", core.GetCurrentlyExecutingInstructionPointer(), regName, count)

}

func INSTR_XCHG(core *CpuCore) {
	core.IncrementIP()

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x90,0x91,0x92,0x93,0x94,0x95,0x96,0x97,0x98: {
		// xchg ax, 16
		term1 = uint32(core.registers.AX)
		r16,r16Str := core.registers.registers16Bit[core.currentByteAtCodePointer - 0x90], core.registers.index16ToString(core.currentByteAtCodePointer - 0x90)
		term2 = uint32(*r16)
		log.Printf("[%#04x] xchg AX, %s", core.GetCurrentlyExecutingInstructionPointer(), r16Str)
	}
	case 0x86: {
		// xchg r/m8, r8
		term1 = uint32(core.registers.AX)
		r16,r16Str := core.registers.registers16Bit[core.currentByteAtCodePointer - 0x90], core.registers.index16ToString(core.currentByteAtCodePointer - 0x90)
		term2 = uint32(*r16)
		log.Printf("[%#04x] xchg AX, %s", core.GetCurrentlyExecutingInstructionPointer(), r16Str)
	}
	default:
		log.Println("Unrecognised xchg instruction!")
		doCoreDump(core)
	}


	core.IncrementIP()

	modrm := consumeModRm(core)

	reg1 := *core.registers.registers16Bit[modrm.mod]
	reg2 := *core.registers.registers16Bit[modrm.reg]

	regName := core.registers.index16ToString(modrm.mod)
	regName2 := core.registers.index16ToString(modrm.reg)

	log.Print(fmt.Sprintf("[%#04x] XCHG %s, %s", core.GetCurrentlyExecutingInstructionPointer(), regName, regName2))

	tmp := reg2

	reg2 = reg1

	reg1 = tmp
}

func INSTR_MOV(core *CpuCore) {

	switch core.currentByteAtCodePointer {
	case 0xB0:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[0] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(0), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB1:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[1] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(1), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB2:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[2] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(2), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB3:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[3] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(3), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB4:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[4] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s,  %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(4), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB5:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[5] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s,  %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(5), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB6:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[6] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s,  %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(6), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB7:
		{
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers8Bit[7] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s,  %#02x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index8ToString(7), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}


	case 0xB8:
		{
			val := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
			*core.registers.registers16Bit[0] = val
			log.Print(fmt.Sprintf("[%#04x] MOV %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), core.registers.index16ToString(0), val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}

	case 0xBA:
		{
			val := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
			core.registers.DX = val
			log.Print(fmt.Sprintf("[%#04x] MOV DX, %#04x", core.GetCurrentlyExecutingInstructionPointer(),  val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}
	case 0xBB:
		{
			val := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
			core.registers.BX = val
			log.Print(fmt.Sprintf("[%#04x] MOV BX, %#04x", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}
	case 0xBC:
		{
			val := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
			core.registers.SP = val
			log.Print(fmt.Sprintf("[%#04x] MOV SP, %#04x", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}
	case 0xBD:
		{
			val := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
			core.registers.BP = val
			log.Print(fmt.Sprintf("[%#04x] MOV BP, %#04x", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}
	case 0xBE:
		{
			val := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
			core.registers.SI = val
			log.Print(fmt.Sprintf("[%#04x] MOV SI, %#04x", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}
	case 0x8A:
		{
			/* 	MOV r8,r/m8 */
			core.IncrementIP()
			modrm := consumeModRm(core)

			var src *uint8
			var srcName string
			dest := core.registers.registers8Bit[modrm.reg]

			dstName := core.registers.index8ToString(modrm.reg)

			if modrm.mod == 3 {
				src = core.registers.registers8Bit[modrm.rm]
				srcName = core.registers.index8ToString(modrm.rm)
				*dest = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				data := core.memoryAccessController.ReadAddr8(uint32(addressMode))
				src = &data
				srcName = "r/m8"
				*dest = *src
			}

			log.Print(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionPointer(), dstName, srcName))

			core.registers.IP = uint16(core.GetIP())
		}
	case 0x8B:
		{
			/* mov r16, r/m16 */
			core.IncrementIP()
			modrm := consumeModRm(core)

			// dest
			dest := core.registers.registers16Bit[modrm.reg]
			dstName := core.registers.index16ToString(modrm.reg)
			var src *uint16
			var srcName string
			if modrm.mod == 3 {
				src = core.registers.registers16Bit[modrm.rm]
				srcName = core.registers.index16ToString(modrm.rm)
				*dest = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				data := core.memoryAccessController.ReadAddr16(uint32(addressMode))
				src = &data
				*dest = *src
				srcName = "rm/16"
			}

			log.Print(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionPointer(), dstName, srcName))

			core.registers.IP = uint16(core.GetIP())
		}
	case 0x8C:
		{
			/* MOV r/m16,Sreg */
			core.IncrementIP()
			modrm := consumeModRm(core)

			src := core.registers.registersSegmentRegisters[modrm.reg]
			srcName := core.registers.indexSegmentToString(modrm.reg)

			var dest *uint16
			var destName string
			if modrm.mod == 3 {
				dest = core.registers.registers16Bit[modrm.rm]
				destName = core.registers.index16ToString(modrm.rm)
				*dest = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				core.memoryAccessController.WriteAddr16(uint32(addressMode), *src)
				srcName = "rm/16"
			}

			log.Print(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionPointer(), destName, srcName))

			core.registers.IP = uint16(core.GetIP())
		}
	case 0x8E:
		{
			/* MOV Sreg,r/m16 */
			core.IncrementIP()
			modrm := consumeModRm(core)

			dest := core.registers.registersSegmentRegisters[modrm.reg]
			dstName := core.registers.indexSegmentToString(modrm.reg)

			var src *uint16
			var srcName string
			if modrm.mod == 3 {
				src = core.registers.registers16Bit[modrm.rm]
				srcName = core.registers.index16ToString(modrm.rm)
				*dest = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				data := core.memoryAccessController.ReadAddr16(uint32(addressMode))
				src = &data
				*dest = *src
				srcName = "rm/16"
			}

			log.Print(fmt.Sprintf("[%#04x] MOV %s,%s", core.GetCurrentlyExecutingInstructionPointer(), dstName, srcName))

			core.registers.IP = uint16(core.GetIP())
		}


	default:
		log.Fatal("Unrecognised MOV instruction!")
	}

}

func INSTR_CMPOLD(core *CpuCore) {
	/*
		cmp dst, src	ZF	CF
		dst = src	1	0
		dst < src	0	1
		dst > src	0	0

	*/
	switch core.currentByteAtCodePointer {
	case 0x3C:
		{
			src := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			dst := core.registers.AL

			if src == dst {
				core.registers.ZF = 1
				core.registers.CF = 0
			} else if dst < src {
				core.registers.ZF = 0
				core.registers.CF = 1
			} else if dst > src {
				core.registers.ZF = 0
				core.registers.CF = 0
			}

			log.Print(fmt.Sprintf("[%#04x] CMP AL, %v", core.GetCurrentlyExecutingInstructionPointer(), src))
		}

	default:
		log.Fatal("Unrecognised CMP instruction!")
	}

	core.registers.IP = uint16(core.GetIP() + 2)
}

func INSTR_IN(core *CpuCore) {
	// Read from port

	switch core.currentByteAtCodePointer {
	case 0xE4:
		{
			// Read from port (imm) to AL
			imm := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)

			data := core.ioPortAccessController.ReadAddr8(uint16(imm))

			core.registers.AL = data
			log.Printf("[%#04x] Port IN addr: imm addr %04X to AL (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), imm, data)
		}
	case 0xE5:
		{
			// Read from port (DX) to AL

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.registers.AL = data
			log.Printf("[%#04x] Port IN addr: DX VAL %04X to AL (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), dx, data)
		}
	case 0xEC:
		{
			// Read from port (imm) to AX

			imm := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)

			data := core.ioPortAccessController.ReadAddr16(imm)

			core.registers.AX = data
			log.Printf("[%#04x] Port IN addr: imm addr %04X to AX (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), imm, data)
		}
	case 0xED:
		{
			// Read from port (DX) to AX

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr16(uint16(dx))

			core.registers.AX = data
			log.Printf("[%#04x] Port IN addr: DX VAL %04X to AX (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), dx, data)
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}

	core.registers.IP = uint16(core.GetIP() + 2)
}

func INSTR_OUT(core *CpuCore) {
	// Read from port

	switch core.currentByteAtCodePointer {
	case 0xE6:
		{
			// Write value in AL to port addr imm8
			imm := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)

			core.ioPortAccessController.WriteAddr8(uint16(imm), core.registers.AL)

			log.Printf("[%#04x] out %04X, al", core.GetCurrentlyExecutingInstructionPointer(), imm)
		}
	case 0xE7:
		{
			// Write value in AX to port addr imm8
			imm := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)

			core.ioPortAccessController.WriteAddr16(uint16(imm), core.registers.AX)

			log.Printf("[%#04x] out %04X, ax", core.GetCurrentlyExecutingInstructionPointer(), imm)
		}
	case 0xEE:
		{
			// Use value of DX as io port addr, and write value in AL

			core.ioPortAccessController.WriteAddr8(uint16(core.registers.DX), core.registers.AL)

			log.Printf("[%#04x] Port out addr: DX addr to io port imm addr %04X (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), core.registers.DX, core.registers.AL)
		}
	case 0xEF:
		{
			// Use value of DX as io port addr, and write value in AX

			core.ioPortAccessController.WriteAddr16(uint16(core.registers.DX), core.registers.AX)

			log.Printf("[%#04x] Port out addr: DX addr to io port imm addr %04X (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), core.registers.DX, core.registers.AX)
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}

	core.registers.IP = uint16(core.GetIP() + 2)
}

func INSTR_JMP_FAR_PTR16(core *CpuCore) {
	destAddr := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
	segment := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 3)

	log.Printf("[%#04x] JMP %#04x:%#04x (FAR_PTR16)", core.GetCurrentlyExecutingInstructionPointer(), segment, destAddr)
	core.registers.CS = segment
	core.registers.IP = destAddr
}

func INSTR_JMP_FAR_M16(core *CpuCore, modrm *ModRm) {
	if modrm.mod == 3 {
		addr := core.registers.registers16Bit[modrm.rm]
		core.registers.IP = *addr
		log.Printf("[%#04x] JMP %#04x (JMP_FAR_M16)", core.GetCurrentlyExecutingInstructionPointer(), uint16(*addr))
	} else {
		addr := modrm.getAddressMode16(core)
		core.registers.IP = addr
		log.Printf("[%#04x] JMP %#04x (JMP_FAR_M16)", core.GetCurrentlyExecutingInstructionPointer(), uint16(addr))
	}

}


func INSTR_JMP_NEAR_REL16(core *CpuCore) {

	offset := int16(core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1))

	var destAddr = int16(core.registers.IP + 3)

	destAddr = destAddr + int16(offset)

	log.Printf("[%#04x] JMP %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionPointer(), uint16(destAddr))
	core.registers.IP = uint16(destAddr)
}

func INSTR_JZ_SHORT_REL8(core *CpuCore) {

	offset := int16(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	var destAddr = int16(core.registers.IP + 2)

	destAddr = destAddr + int16(offset)

	if core.registers.ZF == 0 {
		log.Printf("[%#04x] JZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionPointer(), uint16(destAddr))
		core.registers.IP = uint16(destAddr)
	} else {
		log.Printf("[%#04x] JZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionPointer(), uint16(core.GetIP()+1))
		core.registers.IP = uint16(uint16(core.GetIP() + 2))
	}

}

func INSTR_JNZ_SHORT_REL8(core *CpuCore) {

	offset := int16(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	var destAddr = int16(core.registers.IP + 2)

	destAddr = destAddr + (offset)

	if core.registers.ZF != 0 {
		log.Printf("[%#04x] JNZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionPointer(), uint16(destAddr))
		core.registers.IP = uint16(destAddr)
	} else {
		log.Printf("[%#04x] JNZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionPointer(), uint16(core.GetIP()+2))
		core.registers.IP = uint16(core.GetIP() + 2)
	}

}

func INSTR_JCXZ_SHORT_REL8(core *CpuCore) {

	offset := int16(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	var destAddr = int16(core.registers.IP + 2)

	destAddr = destAddr + int16(offset)

	if core.registers.CX == 0 {
		log.Printf("[%#04x] JCXZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionPointer(), uint16(destAddr))
		core.registers.IP = uint16(destAddr)
	} else {
		log.Printf("[%#04x] JCXZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionPointer(), uint16(core.GetIP()+2))
		core.registers.IP = uint16(core.GetIP() + 2)
	}

}
