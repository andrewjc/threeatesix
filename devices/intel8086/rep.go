package intel8086

import "log"

func INSTR_REP(core *CpuCore) {
	core.currentByteAddr++

	secondaryOpcode := core.readImm8()
/*
	direction :=1
	if core.registers.GetFlag(DirectionFlag) {
		direction = -1
	}*/

	switch core.currentOpCodeBeingExecuted {
	case 0xF3:
		{
			switch secondaryOpcode {
			case 0x6C:
				// rep ins m8, dx

				for x:= uint16(0);x<core.registers.CX;x++ {
					value := core.ioPortAccessController.ReadAddr8(core.registers.DX)
					addr := core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI)
					core.memoryAccessController.WriteAddr8(addr, value)
					core.registers.DI++
				}
				log.Printf("[%#04x] rep ins m8, dx", core.GetCurrentlyExecutingInstructionAddress())

			case 0x6D:
				// rep ins m16, dx
				for x:= uint16(0);x<core.registers.CX;x++ {
					value := core.ioPortAccessController.ReadAddr16(core.registers.DX)
					addr := core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI)
					core.memoryAccessController.WriteAddr16(addr, value)
				}
				log.Printf("[%#04x] rep ins m16, dx", core.GetCurrentlyExecutingInstructionAddress())

			case 0xA4:
				// rep movs m8, m8

				for x:= uint16(0);x<core.registers.CX;x++ {
					value := core.memoryAccessController.ReadAddr8(core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.DX))
					addr := core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI)
					core.memoryAccessController.WriteAddr8(addr, value)
				}
				log.Printf("[%#04x] rep movs m8, m8", core.GetCurrentlyExecutingInstructionAddress())

			case 0xA5:
				// rep movs m16, m16

				for x:= uint16(0);x<core.registers.CX;x++ {
					value := core.memoryAccessController.ReadAddr16(core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.DX))
					addr := core.SegmentAddressToLinearAddress(core.registers.ES, core.registers.DI)
					core.memoryAccessController.WriteAddr16(addr, value)
				}
				log.Printf("[%#04x] rep movs m8, m8", core.GetCurrentlyExecutingInstructionAddress())

			case 0x6E:
				// rep out dx, r/m8

				for x:= uint16(0);x<core.registers.CX;x++ {
					value := core.memoryAccessController.ReadAddr8(core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.DX))
					core.ioPortAccessController.WriteAddr8(core.registers.DX, value)
				}
				log.Printf("[%#04x] rep out dx, r/m8", core.GetCurrentlyExecutingInstructionAddress())

			case 0x6F:
				// rep out dx, r/m16

				for x:= uint16(0);x<core.registers.CX;x++ {
					value := core.memoryAccessController.ReadAddr16(core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.DX))
					core.ioPortAccessController.WriteAddr16(core.registers.DX, value)
				}
				log.Printf("[%#04x] rep out dx, r/m16", core.GetCurrentlyExecutingInstructionAddress())

			case 0xAC:
				// rep lods al
				for x:= uint16(0);x<core.registers.CX;x++ {
					value := core.memoryAccessController.ReadAddr8(core.SegmentAddressToLinearAddress(core.registers.DS, core.registers.DX))
					core.registers.AL = value
				}
				log.Printf("[%#04x] rep out dx, r/m8", core.GetCurrentlyExecutingInstructionAddress())

			case 0xAD:
				// rep lods ax
			case 0xAA:
				// rep stos m8
			case 0xAB:
				// rep stos m16
			case 0xA6:
				// repe cmps m8, m8
			case 0xA7:
				// repe cmps m16, m16
			case 0xAE:
				// repe scas m8
			case 0xAF:
				// repe scas m16
			}
		}
	case 0xF2:
		{
			switch secondaryOpcode {
			case 0xA6:
				// repe cmps m8, m8
			case 0xA7:
				// repe cmps m16, m16
			case 0xAE:
				// repe scas m8
			case 0xAF:
				// repe scas m16
			}
		}
	}

	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}


