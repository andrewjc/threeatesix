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

			byteValue, err := core.memoryAccessController.ReadMemoryValue8(core.SegmentAddressToLinearAddress(segment, uint16(offset)))
			if err != nil {
				core.logInstruction(fmt.Sprintf("Error reading memory: %s", err))
				return
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV al, byte ptr %s:[%#02x]", core.GetCurrentlyExecutingInstructionAddress(), core.registers.indexSegmentToString(uint8(segment.Base)), offset))
			core.registers.AL = byteValue
		}
	case 0xA1:
		{

			// mov ax, moffs16*
			offset, err := core.memoryAccessController.ReadMemoryValue16(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			segment := core.getSegmentOverride()
			byteValue, err := core.memoryAccessController.ReadMemoryValue16(core.SegmentAddressToLinearAddress(segment, offset))
			if err != nil {
				goto eof
			}
			core.logInstruction(fmt.Sprintf("[%#04x] MOV ax, word ptr %s:[%#04x]", core.GetCurrentlyExecutingInstructionAddress(), core.registers.indexSegmentToString(uint8(segment.Base)), offset))

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

			segment := core.getSegmentOverride()
			err = core.memoryAccessController.WriteMemoryAddr8(core.SegmentAddressToLinearAddress(segment, uint16(offset)), core.registers.AL)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV byte ptr %s:[%#02x], al", core.GetCurrentlyExecutingInstructionAddress(), core.registers.indexSegmentToString(uint8(segment.Base)), offset))
		}
	case 0xA3:
		{
			// mov moffs16*, ax
			offset, err := core.memoryAccessController.ReadMemoryValue16(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			segment := core.getSegmentOverride()
			err = core.memoryAccessController.WriteMemoryAddr16(core.SegmentAddressToLinearAddress(segment, offset), core.registers.AX)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV word ptr %s:[%#04x], ax", core.GetCurrentlyExecutingInstructionAddress(), core.registers.indexSegmentToString(uint8(segment.Base)), offset))
		}
	case 0xA4:
		{
			// movsb
			src := core.getSegmentOverride()
			dest := core.registers.ES
			srcAddr := core.SegmentAddressToLinearAddress(src, core.registers.SI)
			destAddr := core.SegmentAddressToLinearAddress(dest, core.registers.DI)

			srcData, err := core.memoryAccessController.ReadMemoryValue8(srcAddr)
			if err != nil {
				goto eof
			}

			err = core.memoryAccessController.WriteMemoryAddr8(destAddr, srcData)
			if err != nil {
				goto eof
			}

			if !core.registers.GetFlag(DirectionFlag) {
				core.registers.SI++
				core.registers.DI++
			} else {
				core.registers.SI--
				core.registers.DI--
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOVSB", core.GetCurrentlyExecutingInstructionAddress()))
		}
	case 0xA5:
		{
			// movsw
			src := core.getSegmentOverride()
			dest := core.registers.ES
			srcAddr := core.SegmentAddressToLinearAddress(src, core.registers.SI)
			destAddr := core.SegmentAddressToLinearAddress(dest, core.registers.DI)

			srcData, err := core.memoryAccessController.ReadMemoryValue16(srcAddr)
			if err != nil {
				goto eof
			}

			err = core.memoryAccessController.WriteMemoryAddr16(destAddr, srcData)
			if err != nil {
				goto eof
			}

			if !core.registers.GetFlag(DirectionFlag) {
				core.registers.SI += 2
				core.registers.DI += 2
			} else {
				core.registers.SI -= 2
				core.registers.DI -= 2
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOVSW", core.GetCurrentlyExecutingInstructionAddress()))

		}
	case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
		{
			// mov r8, imm8
			regIndex := core.currentOpCodeBeingExecuted - 0xB0
			r8 := core.registers.registers8Bit[regIndex]
			r8Str := core.registers.index8ToString(regIndex)
			val, err := core.readImm8()
			if err != nil {
				goto eof
			}
			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionAddress(), r8Str, val))
			*r8 = val
		}
	case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
		{
			// mov r16, imm16
			regIndex := uint8(core.currentOpCodeBeingExecuted - 0xB8)
			immVal, _, err := core.GetImmediate16()
			if err != nil {
				goto eof
			}

			regName, err := core.SetRegister16(regIndex, uint16(immVal))
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), regName, immVal))
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

			src := core.registers.registers16Bit[modrm.reg]
			srcName := core.registers.index16ToString(modrm.reg)

			destName, err := core.writeRm16(&modrm, src)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))
		}
	case 0x8A:
		{
			/* MOV r8,r/m8 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			rm8, srcName, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}

			dest := core.registers.registers8Bit[modrm.reg]
			destName := core.registers.index8ToString(modrm.reg)

			*dest = *rm8

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
				goto eof
			}

			dest := core.registers.registers16Bit[modrm.reg]
			destName := core.registers.index16ToString(modrm.reg)

			*dest = *rm16

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

			src := core.registers.registersSegmentRegisters[modrm.reg]
			srcName := core.registers.indexSegmentToString(modrm.reg)

			regBase := uint16(src.Base >> 4)
			destName, err := core.writeRm16(&modrm, &regBase)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))
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
			if err != nil {
				goto eof
			}

			dest := core.registers.registersSegmentRegisters[modrm.reg]
			dstName := core.registers.indexSegmentToString(modrm.reg)

			dest.Base = uint32(*src) << 4

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))
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
}
