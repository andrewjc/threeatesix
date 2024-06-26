package intel8086

import (
	"fmt"
	"log"
)

func INSTR_MOV(core *CpuCore) {
	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0xA0:
		{
			// mov al, moffs8*
			offset, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr)
			if err != nil {
				core.logInstruction(fmt.Sprintf("Error reading memory: %s", err))
				return
			}
			core.currentByteAddr++

			segment := core.getSegmentOverride() // Get overridden segment register

			segOff := uint16(offset)
			byteValue, err := core.memoryAccessController.ReadMemoryValue8(core.SegmentAddressToLinearAddress(segment, segOff))
			if err != nil {
				core.logInstruction(fmt.Sprintf("Error reading memory: %s", err))
				return
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV al, byte ptr cs:%#02x", core.GetCurrentlyExecutingInstructionAddress(), segOff))
			core.registers.AL = byteValue
		}
	case 0xA1:
		{

			// mov ax, moffs16*
			offset, err := core.memoryAccessController.ReadMemoryValue16(core.currentByteAddr)

			if offset == 0x82da {
				print("test")
			}

			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			byteValue, err := core.memoryAccessController.ReadMemoryValue16(uint32(offset))
			if err != nil {
				goto eof
			}
			core.logInstruction(fmt.Sprintf("[%#04x] MOV ax, byte ptr cs:%#02x", core.GetCurrentlyExecutingInstructionAddress(), offset))

			core.registers.AX = byteValue
		}
	case 0xA2:
		{
			// mov moffs8*, al
			offset, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr)
			if err != nil {
				goto eof
			}

			core.currentByteAddr++

			segOff := uint16(offset)

			err = core.memoryAccessController.WriteMemoryAddr8(uint32(segOff), core.registers.AL)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV byte ptr cs:%#02x, al", core.GetCurrentlyExecutingInstructionAddress(), segOff))
		}
	case 0xA3:
		{
			// mov moffs16*, ax
			offset, err := core.memoryAccessController.ReadMemoryValue16(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			segOff := uint16(offset)

			err = core.memoryAccessController.WriteMemoryAddr16(uint32(segOff), core.registers.AX)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV byte ptr cs:%#02x, ax", core.GetCurrentlyExecutingInstructionAddress(), segOff))

		}
	case 0xA4:
		{
			// movsb
			src := core.registers.DS
			srcOff := uint32(core.registers.SI)
			dest := core.registers.ES
			destOff := uint32(core.registers.DI)

			srcAddr := src.Base + srcOff
			destAddr := dest.Base + destOff

			srcData, err := core.memoryAccessController.ReadMemoryValue8(srcAddr)
			if err != nil {
				goto eof
			}

			err = core.memoryAccessController.WriteMemoryAddr8(destAddr, srcData)
			if err != nil {
				goto eof
			}

			core.registers.SI++
			core.registers.DI++

			core.logInstruction(fmt.Sprintf("[%#04x] MOVSB", core.GetCurrentlyExecutingInstructionAddress()))
		}
	case 0xA5:
		{
			// movsw
			src := core.registers.DS
			srcOff := uint32(core.registers.SI)
			dest := core.registers.ES
			destOff := uint32(core.registers.DI)

			srcAddr := src.Base + srcOff
			destAddr := dest.Base + destOff

			srcData, err := core.memoryAccessController.ReadMemoryValue16(srcAddr)
			if err != nil {
				goto eof
			}

			err = core.memoryAccessController.WriteMemoryAddr16(destAddr, srcData)
			if err != nil {
				goto eof
			}

			core.registers.SI += 2
			core.registers.DI += 2

			core.logInstruction(fmt.Sprintf("[%#04x] MOVSW", core.GetCurrentlyExecutingInstructionAddress()))

		}
	case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
		{
			// mov r8, imm8
			r8, r8Str := core.registers.registers8Bit[core.currentOpCodeBeingExecuted-0xB0], core.registers.index8ToString(core.currentOpCodeBeingExecuted-0xB0)
			val, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr++
			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionAddress(), r8Str, val))
			*r8 = val
		}
	case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
		{
			immVal, _, err := core.GetImmediate16()
			if err != nil {
				goto eof
			}

			rIdx := uint8(core.currentOpCodeBeingExecuted - 0xB8)

			regName, err := core.SetRegister16(rIdx, uint16(immVal))

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), regName, immVal))
		}
	case 0x85:
		{
			/* MOV r/m16, r16 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			var src *uint16
			var srcName string
			var dstName string

			src = core.registers.registers16Bit[modrm.rm]
			srcName = core.registers.index16ToString(modrm.rm)

			dstName, err = core.writeRm16(&modrm, src)

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))
		}
	case 0x88:
		{
			/* MOV r/m8, r8 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			var src *uint8
			var srcName string
			var dstName string

			src = core.registers.registers8Bit[modrm.rm]
			srcName = core.registers.index8ToString(modrm.rm)

			dstName, err = core.writeRm8(&modrm, src)

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))
		}
	case 0x89:
		{
			/* MOV r/m16, r16 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			var destName string
			src := core.registers.registers16Bit[modrm.reg]

			srcName := core.registers.index16ToString(modrm.reg)

			destName, err = core.writeRm16(&modrm, src)

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

		}
	case 0x8A:
		{
			/* 	MOV r8,r/m8 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			rm8, srcName, err := core.readRm8(&modrm)
			if err != nil {
				log.Fatal("Error reading rm8")
			}

			var destName string
			core.registers.registers8Bit[modrm.reg] = rm8

			destName = core.registers.index8ToString(modrm.reg)

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

		}
	case 0x8B:
		{
			/* mov r16, r/m16 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			rm16, srcName, err := core.readRm16(&modrm)
			if err != nil {
				log.Fatal("Error reading rm16")
			}

			var destName string
			core.registers.registers16Bit[modrm.reg] = rm16

			destName = core.registers.index16ToString(modrm.reg)

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

		}
	case 0x8C:
		{
			/* MOV r/m16,Sreg */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			// read Sreg

			src := *core.registers.registersSegmentRegisters[modrm.reg]
			srcName := core.registers.indexSegmentToString(modrm.reg)

			regBase := uint16(src.Base)
			destName, err := core.writeRm16(&modrm, &regBase)

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

		}
	case 0x8D:
		{
			/* LEA r16, m */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			dest := core.registers.registers16Bit[modrm.reg]
			destName := core.registers.index16ToString(modrm.reg)

			destName, err = core.writeRm16(&modrm, dest)

			core.logInstruction(fmt.Sprintf("[%#04x] LEA %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, "m"))
		}
	case 0x8E:
		{
			/* MOV Sreg,r/m16 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			src, srcName, err := core.readRm16(&modrm)

			dest := core.registers.registersSegmentRegisters[modrm.reg]
			dstName := core.registers.indexSegmentToString(modrm.reg)

			(*dest).Base = uint32(*src)

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s,%s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))

		}
	case 0x20:
		{
			/* MOV r32, cr0 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			dst := core.registers.registers32Bit[modrm.rm]
			dstName := core.registers.index32ToString(modrm.rm)

			var srcName string

			switch { //note: these might be wrong?
			case modrm.reg == 0:
				*dst = core.registers.CR0
				srcName = "CR0"
			case modrm.reg == 1:
				*dst = core.registers.CR1
				srcName = "CR1"
			case modrm.reg == 2:
				*dst = core.registers.CR2
				srcName = "CR2"
			case modrm.reg == 3:
				*dst = core.registers.CR3
				srcName = "CR3"
			case modrm.reg == 4:
				*dst = core.registers.CR4
				srcName = "CR4"
			default:
				log.Fatal("Unknown cr mov instruction")
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s,%s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))

		}
	case 0x22:
		{
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			src := core.registers.registers32Bit[modrm.rm]
			srcName := core.registers.index32ToString(modrm.rm)

			var dstName string

			switch { //note: these might be wrong?
			case modrm.reg == 0:
				core.registers.CR0 = *src
				core.updateSystemFlags(core.registers.CR0)
				dstName = "CR0"
			case modrm.reg == 1:
				core.registers.CR1 = *src
				dstName = "CR1"
			case modrm.reg == 2:
				core.registers.CR2 = *src
				dstName = "CR2"
			case modrm.reg == 3:
				core.registers.CR3 = *src
				dstName = "CR3"
			case modrm.reg == 4:
				core.registers.CR4 = *src
				dstName = "CR4"
			default:
				log.Fatal("Unknown cr mov instruction")
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s,%s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))

		}
	default:
		log.Fatal("Unrecognised MOV instruction!")
	}

eof:
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}
