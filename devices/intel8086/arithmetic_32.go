package intel8086

import (
	"fmt"
)

func INSTR_ADD_RM32(core *CpuCore, modrm ModRm, immediate uint32) {

	if core.Is32BitOperand() {

		src := &immediate
		srcName := fmt.Sprintf("0x%08x", immediate)

		dest, destName, err := core.readRm32(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest += *src
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>31)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // Assume no overflow for ADD

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "ADD", destName, srcName))
	} else {
		src := immediate
		srcName := fmt.Sprintf("0x%08x", immediate)
		dest, destName, err := core.readRm16(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest += uint16(src)
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>15)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // Assume no overflow for ADD

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "ADD", destName, srcName))
	}
	return
}

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

	return

eof:
}

func INSTR_OR_RM32(core *CpuCore, modrm ModRm, immediate uint32) {
	if core.Is32BitOperand() {
		src := immediate
		srcName := fmt.Sprintf("0x%08x", immediate)

		dest, destName, err := core.readRm32(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest |= src
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>31)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // OR never causes overflow
		core.registers.SetFlag(CarryFlag, false)    // OR always clears carry flag

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "OR", destName, srcName))
	} else {
		src := uint16(immediate)
		srcName := fmt.Sprintf("0x%04x", src)
		dest, destName, err := core.readRm16(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest |= src
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>15)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // OR never causes overflow
		core.registers.SetFlag(CarryFlag, false)    // OR always clears carry flag

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "OR", destName, srcName))
	}
	return
}

func INSTR_AND_RM32(core *CpuCore, modrm ModRm, immediate uint32) {
	if core.Is32BitOperand() {
		src := immediate
		srcName := fmt.Sprintf("0x%08x", immediate)

		dest, destName, err := core.readRm32(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest &= src
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>31)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // AND never causes overflow
		core.registers.SetFlag(CarryFlag, false)    // AND always clears carry flag

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "AND", destName, srcName))
	} else {
		src := uint16(immediate)
		srcName := fmt.Sprintf("0x%04x", src)
		dest, destName, err := core.readRm16(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest &= src
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>15)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // AND never causes overflow
		core.registers.SetFlag(CarryFlag, false)    // AND always clears carry flag

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "AND", destName, srcName))
	}
	return
}

func INSTR_XOR_RM32(core *CpuCore, modrm ModRm, immediate uint32) {
	if core.Is32BitOperand() {
		src := immediate
		srcName := fmt.Sprintf("0x%08x", immediate)

		dest, destName, err := core.readRm32(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest ^= src
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>31)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // XOR never causes overflow
		core.registers.SetFlag(CarryFlag, false)    // XOR always clears carry flag

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "XOR", destName, srcName))
	} else {
		src := uint16(immediate)
		srcName := fmt.Sprintf("0x%04x", src)
		dest, destName, err := core.readRm16(&modrm)
		if err != nil {
			core.dumpAndExit()
		}
		*dest ^= src
		core.registers.SetFlag(ZeroFlag, *dest == 0)
		core.registers.SetFlag(SignFlag, (*dest>>15)&0x01 == 1)
		core.registers.SetFlag(OverFlowFlag, false) // XOR never causes overflow
		core.registers.SetFlag(CarryFlag, false)    // XOR always clears carry flag

		core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), "XOR", destName, srcName))
	}
	return
}
