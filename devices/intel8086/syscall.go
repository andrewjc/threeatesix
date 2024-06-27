package intel8086

import "fmt"

func INSTR_SYSCALL(core *CpuCore) {
	// System call
	core.logInstruction(fmt.Sprintf("[%#04x] SYSCALL", core.GetCurrentCodePointer()))
	core.currentByteAddr++
}
