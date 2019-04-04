package intel8086

import (
	"fmt"
	"log"
)

func INSTR_MOV(core *CpuCore) {

	switch core.currentByteAtCodePointer {
	case 0xA0:
		{
			// mov al, moffs8*
			offset := uint8(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

			segOff := uint16(offset)

			byteValue := core.memoryAccessController.ReadAddr8(uint32(segOff))

			log.Print(fmt.Sprintf("[%#04x] MOV al, byte ptr cs:%#02x", core.GetCurrentlyExecutingInstructionPointer(), segOff))

			core.registers.AL = byteValue

			core.registers.IP = uint16(core.GetIP() + 1)
		}
	case 0xA1:
		{
			// mov ax, moffs16*
			offset := uint16(core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1))

			byteValue := core.memoryAccessController.ReadAddr16(uint32(offset))

			log.Print(fmt.Sprintf("[%#04x] MOV ax, byte ptr cs:%#02x", core.GetCurrentlyExecutingInstructionPointer(), offset))

			core.registers.AX = byteValue

			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xA2:
		{
			// mov moffs8*, al
			offset := uint8(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

			segOff := uint16(offset)

			core.memoryAccessController.WriteAddr8(uint32(segOff), core.registers.AL)

			log.Print(fmt.Sprintf("[%#04x] MOV byte ptr cs:%#02x, al", core.GetCurrentlyExecutingInstructionPointer(), segOff))

			core.registers.IP = uint16(core.GetIP() + 1)
		}
	case 0xA3:
		{
			// mov moffs16*, ax
			offset := uint16(core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1))

			segOff := uint16(offset)

			core.memoryAccessController.WriteAddr16(uint32(segOff), core.registers.AX)

			log.Print(fmt.Sprintf("[%#04x] MOV byte ptr cs:%#02x, ax", core.GetCurrentlyExecutingInstructionPointer(), segOff))

			core.registers.IP = uint16(core.GetIP() + 1)
		}
	case 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
		{
			// mov r8, imm8
			r8, r8Str := core.registers.registers8Bit[core.currentByteAtCodePointer-0xB0], core.registers.index8ToString(core.currentByteAtCodePointer-0xB0)
			val := core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1)
			log.Print(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionPointer(), r8Str, val))
			*r8 = val
			core.registers.IP = uint16(core.GetIP() + 2)
		}
	case 0xB8, 0xB9, 0xBA, 0xBB, 0xBC, 0xBD, 0xBE, 0xBF:
		{
			// mov r16, imm16
			r16, r16Str := core.registers.registers16Bit[core.currentByteAtCodePointer-0xB8], core.registers.index8ToString(core.currentByteAtCodePointer-0xB8)
			val := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
			log.Print(fmt.Sprintf("[%#04x] MOV %s, %#02x", core.GetCurrentlyExecutingInstructionPointer(), r16Str, val))
			*r16 = val
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


