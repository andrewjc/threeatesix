package intel8086

import (
	"fmt"
	"log"
)

func INSTR_WBINVD(core *CpuCore) {
	core.currentByteAddr++

	operStr := "WBINVD"

	log.Print(fmt.Sprintf("[%#04x] %s", core.GetCurrentlyExecutingInstructionAddress(), operStr))

	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}
