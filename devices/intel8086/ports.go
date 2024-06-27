package intel8086

import (
	"fmt"
	"log"
)

func INSTR_IN(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0xE4:
		{
			// Read from port (imm) to AL
			imm, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr + 1)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			data := core.ioPortAccessController.ReadAddr8(uint16(imm))

			if imm == 0x80 {
				// do nothing
			} else {
				core.registers.AL = data
			}
			core.logInstruction(fmt.Sprintf("[%#04x] IN AL, IMM8 (Port: %#04x, data = %#08x)", core.GetCurrentlyExecutingInstructionAddress(), imm, data))
		}
	case 0xE5:
		{
			// Read from port (DX) to AL

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.registers.AL = data
			core.logInstruction(fmt.Sprintf("[%#04x] IN AL, DX (Port: %#04x, data = %#08x)", core.GetCurrentlyExecutingInstructionAddress(), dx, data))
		}
	case 0xEC:
		{
			// Read from port (imm) to AX

			imm, err := core.memoryAccessController.ReadMemoryValue16(core.currentByteAddr + 1)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			data := core.ioPortAccessController.ReadAddr16(imm)

			core.registers.AX = data
			core.logInstruction(fmt.Sprintf("[%#04x] IN AX, IMM16 (Port: %#04x, data = %#16x)", core.GetCurrentlyExecutingInstructionAddress(), imm, data))
		}
	case 0xED:
		{
			// Read from port (DX) to AX

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr16(uint16(dx))

			core.registers.AX = data
			core.logInstruction(fmt.Sprintf("[%#04x] IN AX, DX (Port: %#04x, data = %#16x)", core.GetCurrentlyExecutingInstructionAddress(), dx, data))
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}

eof:
}

func INSTR_INS(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0x6C:
		{
			// Read from port (imm) to AL

			imm, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr + 1)
			if err != nil {
				goto eof
			}
			core.currentByteAddr++

			data := core.ioPortAccessController.ReadAddr8(uint16(imm))

			core.memoryAccessController.WriteMemoryAddr8(core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI), data)
			core.registers.DI += 1
			core.logInstruction(fmt.Sprintf("[%#04x] INS AL, IMM8 (Port: %#04x, data = %#08x)", core.GetCurrentlyExecutingInstructionAddress(), imm, data))
		}
	case 0x6D:
		{
			// Read from port (DX) to AL

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.memoryAccessController.WriteMemoryAddr8(core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI), data)
			core.registers.DI += 1
			core.logInstruction(fmt.Sprintf("[%#04x] INS AL, DX (Port: %#04x, data = %#08x)", core.GetCurrentlyExecutingInstructionAddress(), dx, data))
		}
	case 0x6E:
		{
			// Read from port (imm) to AX

			imm, err := core.memoryAccessController.ReadMemoryValue16(core.currentByteAddr + 1)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += 2

			data := core.ioPortAccessController.ReadAddr16(imm)

			core.memoryAccessController.WriteMemoryAddr16(core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI), data)
			core.registers.DI += 2
			core.logInstruction(fmt.Sprintf("[%#04x] INS AX, IMM16 (Port: %#04x, data = %#16x)", core.GetCurrentlyExecutingInstructionAddress(), imm, data))
		}
	case 0x6F:
		{
			// Read from port (DX) to AX

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr16(uint16(dx))

			core.memoryAccessController.WriteMemoryAddr16(core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI), data)
			core.registers.DI += 2
			core.logInstruction(fmt.Sprintf("[%#04x] INS AX, DX (Port: %#04x, data = %#16x)", core.GetCurrentlyExecutingInstructionAddress(), dx, data))
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}

eof:
}

func INSTR_OUT(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0xE6:
		{
			// Write value in AL to port addr imm8
			imm, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr + 1)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] OUT %#08x, AL (data = %#08x)", core.GetCurrentlyExecutingInstructionAddress(), imm, core.registers.AL))
			core.ioPortAccessController.WriteAddr8(uint16(imm), core.registers.AL)

			core.currentByteAddr += 2
		}
	case 0xE7:
		{
			// Write value in AX to port addr imm8
			imm, err := core.memoryAccessController.ReadMemoryValue8(core.currentByteAddr + 1)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] OUT %#04x, AX (data = %#16x)", core.GetCurrentlyExecutingInstructionAddress(), imm, core.registers.AX))
			core.ioPortAccessController.WriteAddr16(uint16(imm), core.registers.AX)
			core.currentByteAddr++

		}
	case 0xEE:
		{
			// Use value of DX as io port addr, and write value in AL

			core.logInstruction(fmt.Sprintf("[%#04x] OUT DX, AL (Port: %#16x, data = %#08x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AL))
			core.ioPortAccessController.WriteAddr8(uint16(core.registers.DX), core.registers.AL)

		}
	case 0xEF:
		{
			// Use value of DX as io port addr, and write value in AX

			core.logInstruction(fmt.Sprintf("[%#04x] OUT DX, AX (Port: %#16x, data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AX))
			core.ioPortAccessController.WriteAddr16(uint16(core.registers.DX), core.registers.AX)

		}
	default:
		log.Fatal("Unrecognised OUT (port read) instruction!")
	}

eof:
}

func INSTR_OUTS(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0x6E:
		{
			// Write value in AL to port addr imm8
			core.logInstruction(fmt.Sprintf("[%#04x] OUTS DX, AL (Port: %#16x, data = %#08x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AL))
			core.ioPortAccessController.WriteAddr8(uint16(core.registers.DX), core.registers.AL)
			core.registers.DI += 1
		}
	case 0x6F:
		{
			// Write value in AX to port addr imm8
			core.logInstruction(fmt.Sprintf("[%#04x] OUTS DX, AX (Port: %#16x, data = %#16x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AX))
			core.ioPortAccessController.WriteAddr16(uint16(core.registers.DX), core.registers.AX)
			core.registers.DI += 2
		}
	default:
		log.Fatal("Unrecognised OUTS (port read) instruction!")
	}
}
