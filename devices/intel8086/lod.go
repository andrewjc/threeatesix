package intel8086

import (
	"fmt"
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
				m8, err := core.memoryAccessController.ReadMemoryAddr8(core.SegmentAddressToLinearAddress(core.registers.DS, uint16(core.registers.SI)))
				if err != nil {
					goto eof
				}
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
				//log.Printf("Reading from %#04x", core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.SI))
				m8, err := core.memoryAccessController.ReadMemoryAddr16(core.SegmentAddressToLinearAddress(core.registers.DS, uint16(core.registers.SI)))
				if err != nil {
					goto eof
				}
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

	core.logInstruction(fmt.Sprintf("[%#04x] %s %s %s", core.GetCurrentlyExecutingInstructionAddress(), prefixStr, operStr, extras))

eof:
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}

func INSTR_STOSD(core *CpuCore) {
	core.currentByteAddr++

	// Check if there is a repetition prefix
	if core.flags.RepPrefixEnabled {
		// Execute the operation for the number of times specified in the CX register
		for core.registers.CX > 0 {
			// Perform the STOSD operation
			core.memoryAccessController.WriteMemoryAddr32(core.SegmentAddressToLinearAddress32(core.registers.ES, core.registers.EDI), core.registers.EAX)

			// Update the DI register depending on the direction flag
			if core.registers.GetFlag(DirectionFlag) {
				core.registers.EDI -= 4
			} else {
				core.registers.EDI += 4
			}

			core.registers.CX--
		}

		// Reset the repetition prefix
		core.flags.RepPrefixEnabled = false

	} else {
		// No repetition prefix, just perform the STOSD operation once
		core.memoryAccessController.WriteMemoryAddr32(core.SegmentAddressToLinearAddress32(core.registers.ES, core.registers.EDI), core.registers.EAX)

		// Update the DI register depending on the direction flag
		if core.registers.GetFlag(DirectionFlag) {
			core.registers.EDI -= 4
		} else {
			core.registers.EDI += 4
		}
	}

	// Log the instruction
	core.logInstruction(fmt.Sprintf("[%#04x] STOSD", core.GetCurrentlyExecutingInstructionAddress()))

	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}
