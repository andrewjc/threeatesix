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
			offset, err := core.memoryAccessController.ReadAddr8(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr++

			segOff := uint16(offset)

			byteValue, err := core.memoryAccessController.ReadAddr8(uint32(segOff))
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV al, byte ptr cs:%#02x", core.GetCurrentlyExecutingInstructionAddress(), segOff))

			core.registers.AL = byteValue
		}
	case 0xA1:
		{
			// mov ax, moffs16*
			offset, err := core.memoryAccessController.ReadAddr16(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			byteValue, err := core.memoryAccessController.ReadAddr16(uint32(offset))
			if err != nil {
				goto eof
			}
			core.logInstruction(fmt.Sprintf("[%#04x] MOV ax, byte ptr cs:%#02x", core.GetCurrentlyExecutingInstructionAddress(), offset))

			core.registers.AX = byteValue
		}
	case 0xA2:
		{
			// mov moffs8*, al
			offset, err := core.memoryAccessController.ReadAddr8(core.currentByteAddr)
			if err != nil {
				goto eof
			}

			core.currentByteAddr++

			segOff := uint16(offset)

			err = core.memoryAccessController.WriteAddr8(uint32(segOff), core.registers.AL)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV byte ptr cs:%#02x, al", core.GetCurrentlyExecutingInstructionAddress(), segOff))
		}
	case 0xA3:
		{
			// mov moffs16*, ax
			offset, err := core.memoryAccessController.ReadAddr16(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			segOff := uint16(offset)

			err = core.memoryAccessController.WriteAddr16(uint32(segOff), core.registers.AX)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV byte ptr cs:%#02x, ax", core.GetCurrentlyExecutingInstructionAddress(), segOff))

		}
	case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
		{
			// mov r8, imm8
			r8, r8Str := core.registers.registers8Bit[core.currentOpCodeBeingExecuted-0xB0], core.registers.index8ToString(core.currentOpCodeBeingExecuted-0xB0)
			val, err := core.memoryAccessController.ReadAddr8(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr++
			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionAddress(), r8Str, val))
			*r8 = val
		}
	case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
		{

			// mov r16, imm16
			r16, r16Str := core.registers.registers16Bit[core.currentOpCodeBeingExecuted-0xB8], core.registers.index16ToString(core.currentOpCodeBeingExecuted-0xB8)
			val, err := core.memoryAccessController.ReadAddr16(core.currentByteAddr)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2
			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionAddress(), r16Str, val))
			*r16 = val
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
			dest := core.registers.registers16Bit[modrm.reg]

			dstName := core.registers.index16ToString(modrm.reg)

			if modrm.mod == 3 {
				src = core.registers.registers16Bit[modrm.rm]
				srcName = core.registers.index16ToString(modrm.rm)
				*dest = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				var data, err = core.memoryAccessController.ReadAddr16(uint32(addressMode))
				if err != nil {
					goto eof
				}
				src = &data
				*dest = *src
				srcName = "rm/16"
			}

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

			var dest *uint8
			var destName string
			src := core.registers.registers8Bit[modrm.reg]

			srcName := core.registers.index8ToString(modrm.reg)

			if modrm.mod == 3 {
				dest = core.registers.registers8Bit[modrm.rm]
				destName = core.registers.index8ToString(modrm.rm)
				*dest = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				err := core.memoryAccessController.WriteAddr8(uint32(addressMode), *src)
				if err != nil {
					goto eof
				}
				destName = "r/m8"
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

		}
	case 0x89:
		{
			/* MOV r/m16, r16 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			var dest *uint16
			var destName string
			src := core.registers.registers16Bit[modrm.reg]

			srcName := core.registers.index16ToString(modrm.reg)

			if modrm.mod == 3 {
				dest = core.registers.registers16Bit[modrm.rm]
				destName = core.registers.index16ToString(modrm.rm)
				*dest = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				err := core.memoryAccessController.WriteAddr16(uint32(addressMode), *src)
				if err != nil {
					goto eof
				}
				destName = "r/m16"
			}

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
				var data, err = core.memoryAccessController.ReadAddr8(uint32(addressMode))
				if err != nil {
					goto eof
				}
				src = &data
				srcName = "r/m8"
				*dest = *src
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))

		}
	case 0x8B:
		{
			/* mov r16, r/m16 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

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
				var data, err = core.memoryAccessController.ReadAddr16(uint32(addressMode))
				if err != nil {
					goto eof
				}
				src = &data
				*dest = *src
				srcName = "rm/16"
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s, %s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))

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

			var dest *uint16
			var destName string
			if modrm.mod == 3 {
				dest = core.registers.registers16Bit[modrm.rm]
				destName = core.registers.index16ToString(modrm.rm)
				*dest = (*src).base
			} else {
				addressMode := modrm.getAddressMode16(core)
				err = core.memoryAccessController.WriteAddr16(uint32(addressMode), (*src).base)
				if err != nil {
					goto eof
				}
				srcName = "rm/16"
			}

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

			addressMode := modrm.getAddressMode16(core)
			*dest = uint16(addressMode)

			core.logInstruction(fmt.Sprintf("[%#04x] LEA %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), destName, uint16(addressMode)))
		}
	case 0x8E:
		{
			/* MOV Sreg,r/m16 */
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed

			dest := core.registers.registersSegmentRegisters[modrm.reg]
			dstName := core.registers.indexSegmentToString(modrm.reg)

			var src *uint16
			var srcName string
			if modrm.mod == 3 {
				src = core.registers.registers16Bit[modrm.rm]
				srcName = core.registers.index16ToString(modrm.rm)
				(*dest).base = *src
			} else {
				addressMode := modrm.getAddressMode16(core)
				var data, err = core.memoryAccessController.ReadAddr16(uint32(addressMode))
				if err != nil {
					goto eof
				}
				src = &data
				(*dest).base = *src
				srcName = "rm/16"
			}

			core.logInstruction(fmt.Sprintf("[%#04x] MOV %s,%s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))

		}
	case 0x20:
		{
			/* MOV r32, cr0 */
			core.currentByteAddr++

			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr--
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
			core.currentByteAddr++
			/* MOV cr0, r32 */

			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			core.currentByteAddr--
			core.currentByteAddr += bytesConsumed

			src := core.registers.registers32Bit[modrm.rm]
			srcName := core.registers.index32ToString(modrm.rm)

			var dstName string

			switch { //note: these might be wrong?
			case modrm.reg == 0:
				core.registers.CR0 = *src
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
