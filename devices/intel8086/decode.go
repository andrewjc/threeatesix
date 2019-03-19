package intel8086

import (
	"fmt"
	"log"
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

	c.opCodeMap[0xA8] = INSTR_TEST_AL

	c.opCodeMap[0xB0] = INSTR_MOV
	c.opCodeMap[0xBB] = INSTR_MOV
	c.opCodeMap[0xBC] = INSTR_MOV

	c.opCodeMap[0xB4] = INSTR_MOV
	c.opCodeMap[0x8B] = INSTR_MOV
	c.opCodeMap[0x8C] = INSTR_MOV
	c.opCodeMap[0x8E] = INSTR_MOV

	c.opCodeMap[0x3C] = INSTR_CMP

	c.opCodeMap[0x87] = INSTR_XCHG

	c.opCodeMap[0x90] = INSTR_NOP
}

type OpCodeImpl func(*CpuCore)

func INSTR_NOP(core *CpuCore) {
	// Clear interrupts
	log.Printf("[%#04x] NOP", core.GetCurrentlyExecutingInstructionPointer())

	core.registers.IP = uint16(core.GetIP() + 1)
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

func INSTR_TEST_AL(core *CpuCore) {

	val := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1)

	val2 := core.registers.AL

	tmp := val & val2
	core.registers.SF = uint16(getMSB(tmp))

	if tmp == 0 {
		core.registers.ZF = 1
	} else {
		core.registers.ZF = 0
	}

	core.registers.PF = 1
	for i := uint8(0); i < 8; i++ {
		core.registers.PF ^= uint16(getBitValue(tmp, i))
	}

	core.registers.CF = 0
	core.registers.OF = 0

	core.registers.IP = uint16(uint16(core.GetIP() + 2))
	log.Printf("[%#04x] Test AL, %d", core.GetCurrentlyExecutingInstructionPointer(), val)
}

func INSTR_XCHG(core *CpuCore) {
	core.IncrementIP()
	modrm := consumeModRm(core)

	reg1 := *core.registers.registers16Bit[modrm.mod]
	reg2 := *core.registers.registers16Bit[modrm.reg]

	regName := core.registers.index16ToString(int(modrm.mod))
	regName2 := core.registers.index16ToString(int(modrm.reg))

	log.Print(fmt.Sprintf("[%#04x] XCHG %s, %s", core.GetCurrentlyExecutingInstructionPointer(), regName, regName2))

	tmp := reg2

	reg2 = reg1

	reg1 = tmp
}

func INSTR_MOV(core *CpuCore) {

	switch {
	case 0xB0 == core.currentByteAtCodePointer:
		{
			val := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1)
			core.registers.AL = val
			log.Print(fmt.Sprintf("[%#04x] MOV AL, %#02x", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xBB == core.currentByteAtCodePointer:
		{
			val := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)
			core.registers.BX = val
			log.Print(fmt.Sprintf("[%#04x] MOV BX, %#04x", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}
	case 0xBC == core.currentByteAtCodePointer:
		{
			val := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)
			core.registers.SP = val
			log.Print(fmt.Sprintf("[%#04x] MOV SP, %#04x", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 3)
		}
	case 0xB4 == core.currentByteAtCodePointer:
		{
			val := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1)
			core.registers.AH = val
			log.Print(fmt.Sprintf("[%#04x] MOV AH, %v", core.GetCurrentlyExecutingInstructionPointer(), val))
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0x8B == core.currentByteAtCodePointer:
		{
			/* mov r16, rm16 */
			core.IncrementIP()
			modrm := consumeModRm(core)

			// dest
			dest := core.registers.registers16Bit[modrm.reg]
			src := core.registers.registers16Bit[modrm.rm]

			*dest = *src

			dstName := core.registers.index16ToString(int(modrm.reg))
			srcName := core.registers.index16ToString(int(modrm.rm))

			log.Print(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionPointer(), dstName, srcName))

			core.registers.IP = uint16(core.GetIP())
		}
	case 0x8C == core.currentByteAtCodePointer:
		{
			/* MOV r/m16,Sreg */
			core.IncrementIP()
			modrm := consumeModRm(core)

			dest := core.registers.registers16Bit[modrm.rm]
			src := core.registers.registersSegmentRegisters[modrm.reg]

			*dest = *src

			srcName := core.registers.indexSegmentToString(int(modrm.reg))
			dstName := core.registers.index16ToString(int(modrm.rm))

			log.Print(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionPointer(), dstName, srcName))

			core.registers.IP = uint16(core.GetIP())
		}
	case 0x8E == core.currentByteAtCodePointer:
		{
			/* MOV Sreg,r/m16 */
			core.IncrementIP()
			modrm := consumeModRm(core)

			dest := core.registers.registersSegmentRegisters[modrm.reg]
			src := core.registers.registers16Bit[modrm.rm]

			*dest = *src

			dstName := core.registers.indexSegmentToString(int(modrm.reg))
			srcName := core.registers.index16ToString(int(modrm.rm))

			log.Print(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionPointer(), dstName, srcName))

			core.registers.IP = uint16(core.GetIP())
		}
	default:
		log.Fatal("Unrecognised MOV instruction!")
	}

}

func INSTR_CMP(core *CpuCore) {
	/*
		cmp dst, src	ZF	CF
		dst = src	1	0
		dst < src	0	1
		dst > src	0	0

	*/
	switch {
	case 0x3C == core.currentByteAtCodePointer:
		{
			src := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1)
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

	switch {
	case 0xE4 == core.currentByteAtCodePointer:
		{
			// Read from port (imm) to AL
			imm := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1)

			data := core.ioPortAccessController.ReadAddr8(uint16(imm))

			core.registers.AL = data
			log.Printf("[%#04x] Port IN addr: imm addr %04X to AL (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), imm, data)
		}
	case 0xE5 == core.currentByteAtCodePointer:
		{
			// Read from port (DX) to AL

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.registers.AL = data
			log.Printf("[%#04x] Port IN addr: DX VAL %04X to AL (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), dx, data)
		}
	case 0xEC == core.currentByteAtCodePointer:
		{
			// Read from port (imm) to AX

			imm := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)

			data := core.ioPortAccessController.ReadAddr16(imm)

			core.registers.AX = data
			log.Printf("[%#04x] Port IN addr: imm addr %04X to AX (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), imm, data)
		}
	case 0xED == core.currentByteAtCodePointer:
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

	switch {
	case 0xE6 == core.currentByteAtCodePointer:
		{
			// Write value in AL to port addr imm8
			imm := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1)

			core.ioPortAccessController.WriteAddr8(uint16(imm), core.registers.AL)

			log.Printf("[%#04x] Port out addr: AL to io port imm addr %04X (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), imm, core.registers.AL)
		}
	case 0xE7 == core.currentByteAtCodePointer:
		{
			// Write value in AX to port addr imm8
			imm := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1)

			core.ioPortAccessController.WriteAddr16(uint16(imm), core.registers.AX)

			log.Printf("[%#04x] Port out addr: AX to io port imm addr %04X (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), imm, core.registers.AX)
		}
	case 0xEE == core.currentByteAtCodePointer:
		{
			// Use value of DX as io port addr, and write value in AL

			core.ioPortAccessController.WriteAddr8(uint16(core.registers.DX), core.registers.AL)

			log.Printf("[%#04x] Port out addr: DX addr to io port imm addr %04X (data = %04X)", core.GetCurrentlyExecutingInstructionPointer(), core.registers.DX, core.registers.AL)
		}
	case 0xEF == core.currentByteAtCodePointer:
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
	destAddr := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1)
	segment := core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 3)

	log.Printf("[%#04x] JMP %#04x:%#04x (FAR_PTR16)", core.GetCurrentlyExecutingInstructionPointer(), segment, destAddr)
	core.registers.CS = segment
	core.registers.IP = destAddr
}

func INSTR_JMP_NEAR_REL16(core *CpuCore) {

	offset := int16(core.memoryAccessController.ReadAddr16(core.GetCurrentCodePointer() + 1))

	var destAddr = int16(core.registers.IP + 3)

	destAddr = destAddr + int16(offset)

	log.Printf("[%#04x] JMP %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionPointer(), uint16(destAddr))
	core.registers.IP = uint16(destAddr)
}

func INSTR_JZ_SHORT_REL8(core *CpuCore) {

	offset := int16(core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1))

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

	offset := int16(core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1))

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

	offset := int16(core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer() + 1))

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
