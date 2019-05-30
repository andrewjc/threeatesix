package intel8086

import "log"

func INSTR_IN(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0xE4:
		{
			// Read from port (imm) to AL
			imm, err := core.memoryAccessController.ReadAddr8(core.currentByteAddr + 1)
			if err != nil { goto eof }
			core.currentByteAddr++

			data := core.ioPortAccessController.ReadAddr8(uint16(imm))

			core.registers.AL = data
			log.Printf("[%#04x] IN AL, IMM8 (Port: %#04x, data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), imm, data)
		}
	case 0xE5:
		{
			// Read from port (DX) to AL

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.registers.AL = data
			log.Printf("[%#04x] IN AL, DX (Port: %#04x, data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), dx, data)
		}
	case 0xEC:
		{
			// Read from port (imm) to AX

			imm, err := core.memoryAccessController.ReadAddr16(core.currentByteAddr + 1)
			if err != nil { goto eof }
			core.currentByteAddr += 2

			data := core.ioPortAccessController.ReadAddr16(imm)

			core.registers.AX = data
			log.Printf("[%#04x] IN AX, IMM16 (Port: %#04x, data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), imm, data)
		}
	case 0xED:
		{
			// Read from port (DX) to AX

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr16(uint16(dx))

			core.registers.AX = data
			log.Printf("[%#04x] IN AX, DX (Port: %#04x, data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), dx, data)
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}

	eof:
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart) +1
}

func INSTR_OUT(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0xE6:
		{
			// Write value in AL to port addr imm8
			imm, err := core.memoryAccessController.ReadAddr8(core.currentByteAddr + 1)
			if err != nil { goto eof }
			core.currentByteAddr++

			core.ioPortAccessController.WriteAddr8(uint16(imm), core.registers.AL)

			log.Printf("[%#04x] OUT %#04x, AL (data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), imm, core.registers.AL)
		}
	case 0xE7:
		{
			// Write value in AX to port addr imm8
			imm, err := core.memoryAccessController.ReadAddr8(core.currentByteAddr + 1)
			if err != nil { goto eof }
			core.currentByteAddr++

			core.ioPortAccessController.WriteAddr16(uint16(imm), core.registers.AX)

			log.Printf("[%#04x] OUT %#04x, AX (data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), imm, core.registers.AX)
		}
	case 0xEE:
		{
			// Use value of DX as io port addr, and write value in AL

			core.ioPortAccessController.WriteAddr8(uint16(core.registers.DX), core.registers.AL)

			log.Printf("[%#04x] OUT DX, AL (Port: %#04x, data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AL)
		}
	case 0xEF:
		{
			// Use value of DX as io port addr, and write value in AX

			core.ioPortAccessController.WriteAddr16(uint16(core.registers.DX), core.registers.AX)

			log.Printf("[%#04x] OUT DX, AX (Port: %#04x, data = %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AX)
		}
	default:
		log.Fatal("Unrecognised OUT (port read) instruction!")
	}


	eof:
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart) +1

}
