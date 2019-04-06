package intel8086

import (
	"fmt"
	"log"
)

func INSTR_PUSH(core *CpuCore) {
	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
		{
			// PUSH r16
			val, valName := core.registers.registers16Bit[core.currentOpCodeBeingExecuted-0x50], core.registers.index16ToString(core.currentOpCodeBeingExecuted-0x50)

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), valName)

		}
	case 0x6A:
		{
			// PUSH imm8


			val := core.readImm8()

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr8(uint32(core.registers.SP), val)

			log.Printf("[%#04x] push %#04x", core.GetCurrentlyExecutingInstructionAddress(), val)
		}
	case 0x68:
		{
			// PUSH imm16


			val := core.readImm16()

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), val)

			log.Printf("[%#04x] push %#04x", core.GetCurrentlyExecutingInstructionAddress(), val)
		}
	case 0x0E:
		{
			// PUSH CS


			val := core.registers.registers16Bit[core.registers.CS]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "CS")
		}
	case 0x16:
		{
			// PUSH SS

			val := core.registers.registers16Bit[core.registers.SS]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "SS")
		}
	case 0x1E:
		{
			// PUSH DS

			val := core.registers.registers16Bit[core.registers.DS]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "DS")
		}
	case 0x06:
		{
			// PUSH ES

			val := core.registers.registers16Bit[core.registers.ES]

			core.registers.SP = core.registers.SP - 2

			core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), *val)

			log.Printf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "ES")
		}
	default:
		log.Println(fmt.Printf("Unhandled PUSH instruction:  %#04x", core.currentOpCodeBeingExecuted))
		doCoreDump(core)
	}

	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}

