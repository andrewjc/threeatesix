package intel8086

import (
	"fmt"
	"log"
)

func INSTR_LODS(core *CpuCore) {
	core.currentByteAddr++

	var prefixStr = ""
	var operStr = ""
	var extras = ""
	if core.flags.RepPrefixEnabled {
		prefixStr = "REP"
	}

	if core.registers.CX > 0 && core.flags.RepPrefixEnabled {
		extras = fmt.Sprintf("(%d repetitions)", core.registers.CX)
	}

	for (core.registers.CX > 0 && core.flags.RepPrefixEnabled) || !core.flags.RepPrefixEnabled {

		switch core.currentOpCodeBeingExecuted {
		case 0xAC:
			{
				operStr = "LODSB"
				m8 := uint8(core.memoryAccessController.ReadAddr8(core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.SI)))
				core.registers.AL = m8
				if core.registers.GetFlag(DirectionFlag) {
					core.registers.SI -= 1
				} else {
					core.registers.SI += 1
				}
			}
		case 0xAD:
			{
				operStr = "LODSW"
				m8 := uint16(core.memoryAccessController.ReadAddr16(core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.SI)))
				core.registers.AX = m8
				if core.registers.GetFlag(DirectionFlag) {
					core.registers.SI -= 2
				} else {
					core.registers.SI += 2
				}
			}
		}

		if !core.flags.RepPrefixEnabled {
			break
		} else {
			core.registers.CX--
		}
	}

	log.Print(fmt.Sprintf("[%#04x] %s %s %s", core.GetCurrentlyExecutingInstructionAddress(), prefixStr, operStr, extras))

	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}


