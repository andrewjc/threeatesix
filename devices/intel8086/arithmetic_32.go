package intel8086

import "fmt"

func INSTR_INC_RM32(core *CpuCore) {
	var addr *uint32
	var addrDesc string

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	switch modrm.mod {
	case 3:
		panic("INSTR_INC_RM32: mod 3 not implemented")
		//core.registers.SetReg32(modrm.rm, core.registers.GetReg32(modrm.rm)+1)
	default:
		addr, addrDesc = core.getEffectiveAddress32(&modrm)
		if *addr == 0 {
			goto eof
		}
		val, err := core.memoryAccessController.ReadMemoryAddr32(*addr)
		if err != nil {
			goto eof
		}
		val++
		err = core.memoryAccessController.WriteMemoryAddr32(*addr, val)
		if err != nil {
			goto eof
		}
	}

	core.logInstruction(fmt.Sprintf("[%#04x] %s %s", core.GetCurrentlyExecutingInstructionAddress(), "INC", addrDesc))
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
	return

eof:
}

func INSTR_DEC_RM32(core *CpuCore) {
	var addr *uint32
	var addrDesc string

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	switch modrm.mod {
	case 3:
		panic("INSTR_DEC_RM32: mod 3 not implemented")
		//core.registers.SetReg32(modrm.rm, core.registers.GetReg32(modrm.rm)-1)
	default:
		addr, addrDesc = core.getEffectiveAddress32(&modrm)
		if *addr == 0 {
			goto eof
		}
		val, err := core.memoryAccessController.ReadMemoryAddr32(*addr)
		if err != nil {
			goto eof
		}
		val--
		err = core.memoryAccessController.WriteMemoryAddr32(*addr, val)
		if err != nil {
			goto eof
		}
	}

	core.logInstruction(fmt.Sprintf("[%#04x] %s %s", core.GetCurrentlyExecutingInstructionAddress(), "DEC", addrDesc))
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
	return

eof:
}
