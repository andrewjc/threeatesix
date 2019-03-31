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

	c.opCodeMap[0xD0] = INSTR_SHIFT
	c.opCodeMap[0xD1] = INSTR_SHIFT
	c.opCodeMap[0xD2] = INSTR_SHIFT
	c.opCodeMap[0xD3] = INSTR_SHIFT
	c.opCodeMap[0xC0] = INSTR_SHIFT
	c.opCodeMap[0xC1] = INSTR_SHIFT

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

	// opcode 0F A0 and 0F A8 are push .. todo

}

type OpCodeImpl func(*CpuCore)

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



func INSTR_XCHG(core *CpuCore) {
	core.IncrementIP()

	switch core.currentByteAtCodePointer {
	case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98:
		{
			// xchg ax, 16
			term1 := core.registers.AX
			r16, r16Str := core.registers.registers16Bit[core.currentByteAtCodePointer-0x90], core.registers.index16ToString(core.currentByteAtCodePointer-0x90)
			term2 := *r16
			core.registers.AX = term2
			*r16 = term1

			log.Printf("[%#04x] xchg AX, %s", core.GetCurrentlyExecutingInstructionPointer(), r16Str)
		}
	case 0x86:
		{
			// XCHG r/m8, r8
			modrm := core.consumeModRm()
			rm8, rm8Str := core.readRm8(&modrm)

			r8, r8Str := core.readR8(&modrm)

			tmp := *rm8

			*rm8 = *r8

			*r8 = tmp

			log.Printf("[%#04x] xchg %s, %s", core.GetCurrentlyExecutingInstructionPointer(), rm8Str, r8Str)
		}
	case 0x87:
		{
			// XCHG r/m16, r16
			modrm := core.consumeModRm()
			rm8, rm8Str := core.readRm16(&modrm)

			r8, r8Str := core.readR16(&modrm)

			tmp := *rm8

			*rm8 = *r8

			*r8 = tmp

			log.Printf("[%#04x] xchg %s, %s", core.GetCurrentlyExecutingInstructionPointer(), rm8Str, r8Str)
		}
	default:
		log.Println("Unrecognised xchg instruction!")
		doCoreDump(core)
	}

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
			log.Print(fmt.Sprintf("[%#04x] MOV DX, %#04x", core.GetCurrentlyExecutingInstructionPointer(), val))
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
			modrm := core.consumeModRm()

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
			modrm := core.consumeModRm()

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
			modrm := core.consumeModRm()

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
			modrm := core.consumeModRm()

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



