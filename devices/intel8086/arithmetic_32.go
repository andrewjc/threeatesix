package intel8086

import (
	"fmt"
)

func INSTR_INC_RM32(core *CpuCore) {
	var src *uint32
	var srcName string

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	src, srcName, err = core.readRm32(&modrm)
	if err != nil {
		goto eof
	}

	*src++

	err = core.writeRm32(&modrm, src)
	if err != nil {
		return
	}

	core.registers.SetFlag(ZeroFlag, *src == 0)
	core.registers.SetFlag(SignFlag, (*src>>31)&0x01 == 1)
	core.registers.SetFlag(OverFlowFlag, false) // Assume no overflow for INC

	core.logInstruction(fmt.Sprintf("[%#04x] %s %s", core.GetCurrentlyExecutingInstructionAddress(), "INC", srcName))
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
	return

eof:
}

func INSTR_DEC_RM32(core *CpuCore) {
	var src *uint32
	var srcName string

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	src, srcName, err = core.readRm32(&modrm)
	if err != nil {
		goto eof
	}

	*src--

	err = core.writeRm32(&modrm, src)
	if err != nil {
		return
	}

	core.logInstruction(fmt.Sprintf("[%#04x] %s %s", core.GetCurrentlyExecutingInstructionAddress(), "INC", srcName))
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
	return

eof:
}
