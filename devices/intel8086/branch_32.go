package intel8086

import (
	"fmt"
)

func INSTR_CALL_RM32(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		core.logInstruction("Error in INSTR_CALL_RM32: %s\n", err)
		doCoreDump(core)
		panic(0)
	}

	addr, addrname := core.getEffectiveAddress32(&modrm)
	stackPush32(core, uint32(core.GetIP()+2))
	core.registers.IP = uint16(addr)
	core.logInstruction(fmt.Sprintf("[%#04x] CALL %s (%#08x)", core.GetCurrentlyExecutingInstructionAddress(), addrname, addr))

}

func INSTR_JMP_FAR_M32(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		core.logInstruction("Error in INSTR_CALL_RM32: %s\n", err)
		doCoreDump(core)
		panic(0)
	}

	addr, addrName := core.getEffectiveAddress16(&modrm)
	core.registers.IP = addr
	core.logInstruction(fmt.Sprintf("[%#04x] JMP %s (JMP_FAR_M16)", core.GetCurrentlyExecutingInstructionAddress(), addrName))
}
