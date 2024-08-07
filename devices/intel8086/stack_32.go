package intel8086

import (
	"fmt"
	"log"
)

func stackPush32(core *CpuCore, val uint32) error {
	core.registers.SP -= 4
	stackAddr := core.SegmentAddressToLinearAddress32(core.registers.SS, uint32(core.registers.SP))
	return core.memoryAccessController.WriteMemoryAddr32(stackAddr, val)
}

func stackPop32(core *CpuCore) (uint32, error) {
	stackAddr := core.SegmentAddressToLinearAddress32(core.registers.SS, uint32(core.registers.SP))
	val, err := core.memoryAccessController.ReadMemoryValue32(stackAddr)
	if err != nil {
		return 0, err
	}
	core.registers.SP += 4
	return val, nil
}

func INSTR_PUSH_32(core *CpuCore) {
	opcode := core.currentOpCodeBeingExecuted
	switch opcode {
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57: // PUSH r16
		index := opcode - 0x50
		val := core.registers.registers16Bit[index]
		if err := stackPush16(core, *val); err != nil {
			core.logInstruction("Error pushing register: %#04x\n", opcode)
			return
		}
		core.logInstruction(fmt.Sprintf("PUSH %s", core.registers.index16ToString(index)))
		core.registers.IP += 1 // Increment IP by 1 to simulate the reading of the opcode
	case 0x06: // PUSH ES
		segmentSelector := uint16(core.registers.ES.Base >> 4)
		if err := stackPush16(core, segmentSelector); err != nil {
			core.logInstruction("Error pushing ES: %s\n", err)
			return
		}
		core.logInstruction(fmt.Sprintf("PUSH ES"))
		core.registers.IP += 1 // Increment IP by 1 to simulate the reading of the opcode

	case 0x60:
		// Save the original SP to push later
		originalSP := core.registers.SP

		// The order in which registers are pushed to the stack: DI, SI, BP, original SP, BX, DX, CX, AX
		err := stackPush16(core, core.registers.DI)
		if err != nil {
			log.Println("Error pushing DI:", err)
			return
		}
		err = stackPush16(core, core.registers.SI)
		if err != nil {
			log.Println("Error pushing SI:", err)
			return
		}
		err = stackPush16(core, core.registers.BP)
		if err != nil {
			log.Println("Error pushing BP:", err)
			return
		}
		err = stackPush16(core, originalSP) // Push the original SP value
		if err != nil {
			log.Println("Error pushing original SP:", err)
			return
		}
		err = stackPush16(core, core.registers.BX)
		if err != nil {
			log.Println("Error pushing BX:", err)
			return
		}
		err = stackPush16(core, core.registers.DX)
		if err != nil {
			log.Println("Error pushing DX:", err)
			return
		}
		err = stackPush16(core, core.registers.CX)
		if err != nil {
			log.Println("Error pushing CX:", err)
			return
		}
		err = stackPush16(core, core.registers.AX)
		if err != nil {
			log.Println("Error pushing AX:", err)
			return
		}

		core.logInstruction(fmt.Sprintf("[%#04x] PUSHA", core.GetCurrentlyExecutingInstructionAddress()))
		core.registers.IP += 1 // Increment IP by 1 to simulate the reading of the opcode

	case 0x6A: // PUSH imm8
		imm8, err := core.readImm8() // Assumes implementation to read the next byte as signed 8-bit
		if err != nil {
			core.logInstruction("Error reading immediate value: %s\n", err)
			return
		}
		signExtended := uint16(int8(imm8))
		if err := stackPush16(core, signExtended); err != nil {
			core.logInstruction("Error pushing sign-extended immediate value: %s\n", err)
			return
		}
		core.logInstruction(fmt.Sprintf("PUSH %d (sign-extended)", signExtended))
		core.registers.IP += 2 // Increment IP by 2 (1 for opcode + 1 for immediate value)

	case 0x68: // PUSH imm16
		imm16, err := core.readImm16() // Assumes implementation to read the next two bytes as 16-bit
		if err != nil {
			core.logInstruction("Error reading immediate 16-bit value: %s\n", err)
			return
		}
		if err := stackPush16(core, imm16); err != nil {
			core.logInstruction("Error pushing immediate value: %s\n", err)
			return
		}
		core.logInstruction(fmt.Sprintf("PUSH %#04x", imm16))
		core.registers.IP += 3 // Increment IP by 3 (1 for opcode + 2 for immediate value)

	default:
		core.logInstruction("Unhandled PUSH opcode: %#04x\n", opcode)
	}
}
func INSTR_POP_32(core *CpuCore) {
	var instructionSize uint16

	switch core.currentOpCodeBeingExecuted {
	case 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F:
		regIndex := core.currentOpCodeBeingExecuted - 0x58
		val, valName := core.registers.registers16Bit[regIndex], core.registers.index16ToString(regIndex)
		pval, err := stackPop16(core)
		if err != nil {
			log.Println("Error popping from stack:", err)
			return
		}
		*val = pval
		core.logInstruction(fmt.Sprintf("[%#04x] POP %s", core.GetCurrentlyExecutingInstructionAddress(), valName))
		instructionSize = 1 // POP r16 instructions are 1 byte long

	case 0x07, 0x17, 0x1F, 0x0E: // POP ES, SS, DS, CS respectively
		var reg *SegmentRegister
		var regName string
		switch core.currentOpCodeBeingExecuted {
		case 0x07:
			reg = &core.registers.ES
			regName = "ES"
		case 0x17:
			reg = &core.registers.SS
			regName = "SS"
		case 0x1F:
			reg = &core.registers.DS
			regName = "DS"
		case 0x0E:
			reg = &core.registers.CS
			regName = "CS"
		}
		val, err := stackPop16(core)
		if err != nil {
			core.logInstruction("Error popping %s from stack: %s\n", regName, err)
			return
		}
		reg.Base = uint32(val) << 4
		core.logInstruction(fmt.Sprintf("[%#04x] POP %s", core.GetCurrentlyExecutingInstructionAddress(), regName))
		instructionSize = 1 // POP segment register instructions are also 1 byte long

	default:
		core.logInstruction("Unhandled POP instruction: %#04x\n", core.currentOpCodeBeingExecuted)
		return
	}

	// Increment the instruction pointer by the size of the instruction
	core.registers.IP += instructionSize
}
