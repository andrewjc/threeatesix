package intel8086

import (
	"fmt"
)

func INSTR_CALL_RM32(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	if modrm.mod == 3 {
		reg, reg_str := core.registers.registers32Bit[modrm.rm], core.registers.index32ToString(modrm.rm)

		stackPush32(core, uint32(core.GetIP()+2))

		core.registers.IP = uint16(*reg)
		core.logInstruction(fmt.Sprintf("[%#04x] CALL %s (%#08x)", core.GetCurrentlyExecutingInstructionAddress(), reg_str, *reg))
	} else {
		addr := modrm.getAddressMode32(core)
		stackPush32(core, uint32(core.GetIP()+2))

		core.registers.IP = uint16(addr)
		core.logInstruction(fmt.Sprintf("[%#04x] CALL %#08x", core.GetCurrentlyExecutingInstructionAddress(), addr))
	}

eof:
}

func INSTR_JMP_FAR_M32(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	if modrm.mod == 3 {
		reg, reg_str := core.registers.registers32Bit[modrm.rm], core.registers.index32ToString(modrm.rm)
		core.registers.IP = uint16(*reg)
		core.logInstruction(fmt.Sprintf("[%#04x] JMP %s (%#04x) (JMP_FAR_M32)", core.GetCurrentlyExecutingInstructionAddress(), reg_str, uint32(*reg)))
	} else {
		addr := modrm.getAddressMode32(core)
		core.registers.IP = uint16(addr)
		core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x (JMP_FAR_M32)", core.GetCurrentlyExecutingInstructionAddress(), uint32(addr)))
	}

eof:
}
