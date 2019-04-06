package intel8086

import "log"

func INSTR_IN(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0xE4:
		{
			// Read from port (imm) to AL
			imm := core.memoryAccessController.ReadAddr8(core.currentByteAddr + 1)
			core.currentByteAddr++

			data := core.ioPortAccessController.ReadAddr8(uint16(imm))

			core.registers.AL = data
			log.Printf("[%#04x] Port IN addr: imm addr %04X to AL (data = %04X)", core.GetCurrentlyExecutingInstructionAddress(), imm, data)
		}
	case 0xE5:
		{
			// Read from port (DX) to AL

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr8(uint16(dx))

			core.registers.AL = data
			log.Printf("[%#04x] Port IN addr: DX VAL %04X to AL (data = %04X)", core.GetCurrentlyExecutingInstructionAddress(), dx, data)
		}
	case 0xEC:
		{
			// Read from port (imm) to AX

			imm := core.memoryAccessController.ReadAddr16(core.currentByteAddr + 1)
			core.currentByteAddr += 2

			data := core.ioPortAccessController.ReadAddr16(imm)

			core.registers.AX = data
			log.Printf("[%#04x] Port IN addr: imm addr %04X to AX (data = %04X)", core.GetCurrentlyExecutingInstructionAddress(), imm, data)
		}
	case 0xED:
		{
			// Read from port (DX) to AX

			dx := core.registers.DX

			data := core.ioPortAccessController.ReadAddr16(uint16(dx))

			core.registers.AX = data
			log.Printf("[%#04x] Port IN addr: DX VAL %04X to AX (data = %04X)", core.GetCurrentlyExecutingInstructionAddress(), dx, data)
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}

	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart) +1

}

func INSTR_OUT(core *CpuCore) {
	// Read from port

	switch core.currentOpCodeBeingExecuted {
	case 0xE6:
		{
			// Write value in AL to port addr imm8
			imm := core.memoryAccessController.ReadAddr8(core.currentByteAddr + 1)
			core.currentByteAddr++


			core.ioPortAccessController.WriteAddr8(uint16(imm), core.registers.AL)

			log.Printf("[%#04x] out %04X, al", core.GetCurrentlyExecutingInstructionAddress(), imm)
		}
	case 0xE7:
		{
			// Write value in AX to port addr imm8
			imm := core.memoryAccessController.ReadAddr8(core.currentByteAddr + 1)
			core.currentByteAddr++

			core.ioPortAccessController.WriteAddr16(uint16(imm), core.registers.AX)

			log.Printf("[%#04x] out %04X, ax", core.GetCurrentlyExecutingInstructionAddress(), imm)
		}
	case 0xEE:
		{
			// Use value of DX as io port addr, and write value in AL

			core.ioPortAccessController.WriteAddr8(uint16(core.registers.DX), core.registers.AL)

			log.Printf("[%#04x] Port out addr: DX addr to io port imm addr %04X (data = %04X)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AL)
		}
	case 0xEF:
		{
			// Use value of DX as io port addr, and write value in AX

			core.ioPortAccessController.WriteAddr16(uint16(core.registers.DX), core.registers.AX)

			log.Printf("[%#04x] Port out addr: DX addr to io port imm addr %04X (data = %04X)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.DX, core.registers.AX)
		}
	default:
		log.Fatal("Unrecognised IN (port read) instruction!")
	}


	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart) +1

}
