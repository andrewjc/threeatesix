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

	segment := core.getSegmentOverride() // Get overridden segment register

	for (core.registers.CX > 0 && core.flags.RepPrefixEnabled) || !core.flags.RepPrefixEnabled {
		switch core.currentOpCodeBeingExecuted {
		case 0xAC:
			m8, err := core.memoryAccessController.ReadMemoryValue8(core.SegmentAddressToLinearAddress(segment, uint16(core.registers.SI)))
			if err != nil {
				core.logInstruction(fmt.Sprintf("Error reading memory: %s", err))
				return
			}
			core.registers.AL = m8
			if core.registers.GetFlag(DirectionFlag) {
				core.registers.SI -= 1
			} else {
				core.registers.SI += 1
			}
		case 0xAD:
			m16, err := core.memoryAccessController.ReadMemoryValue16(core.SegmentAddressToLinearAddress(segment, uint16(core.registers.SI)))
			if err != nil {
				core.logInstruction(fmt.Sprintf("Error reading memory: %s", err))
				return
			}
			core.registers.AX = m16
			if core.registers.GetFlag(DirectionFlag) {
				core.registers.SI -= 2
			} else {
				core.registers.SI += 2
			}
		}
		if !core.flags.RepPrefixEnabled {
			break
		} else {
			core.registers.CX--
		}
	}

	core.logInstruction(fmt.Sprintf("[%#04x] %s %s %s", core.GetCurrentlyExecutingInstructionAddress(), prefixStr, operStr, extras))

}

func INSTR_STOSB(core *CpuCore) {
	core.currentByteAddr++

	// Check if there is a repetition prefix
	if core.flags.RepPrefixEnabled {
		// Execute the operation for the number of times specified in the CX register
		for core.registers.CX > 0 {
			// Perform the STOSB operation
			core.memoryAccessController.WriteMemoryAddr8(core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI), core.registers.AL)

			// Update the DI register depending on the direction flag
			if core.registers.GetFlag(DirectionFlag) {
				core.registers.DI--
			} else {
				core.registers.DI++
			}

			core.registers.CX--
		}

		// Reset the repetition prefix
		core.flags.RepPrefixEnabled = false

	} else {
		// No repetition prefix, just perform the STOSB operation once
		core.memoryAccessController.WriteMemoryAddr8(core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI), core.registers.AL)

		// Update the DI register depending on the direction flag
		if core.registers.GetFlag(DirectionFlag) {
			core.registers.DI--
		} else {
			core.registers.DI++
		}
	}

	// Log the instruction
	core.logInstruction(fmt.Sprintf("[%#04x] STOSB", core.GetCurrentlyExecutingInstructionAddress()))

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

}
