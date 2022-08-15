package intel8086

import (
	"fmt"
	"log"
)

func stackPush8(core *CpuCore, val uint8) error {
	core.registers.SP = core.registers.SP - 1
	err := core.memoryAccessController.WriteAddr8(uint32(core.registers.SP), val)
	return err
}

func stackPush16(core *CpuCore, val uint16) error {
	core.registers.SP = core.registers.SP - 2
	err := core.memoryAccessController.WriteAddr16(uint32(core.registers.SP), val)
	return err
}

func stackPop8(core *CpuCore) (uint8, error) {
	val, err := core.memoryAccessController.ReadAddr8(uint32(core.registers.SP))
	if err != nil {
		goto eof
	}
	core.registers.SP = core.registers.SP + 1
	return val, nil
eof:
	return 0, err
}

func stackPop16(core *CpuCore) (uint16, error) {
	addr := core.registers.SS.base + core.registers.SP

	val, err := core.memoryAccessController.ReadAddr16(uint32(addr))
	if err != nil {
		goto eof
	}
	core.registers.SP = core.registers.SP + 2
	return val, nil
eof:
	return 0, err
}

func INSTR_RET_NEAR(core *CpuCore) {

	current_ip := core.GetCurrentCodePointer()

	stackPntrAddr, err := stackPop16(core)
	if err != nil {
		goto eof
	}

	if stackPntrAddr == 0 {
		stackPntrAddr = core.registers.SP
	}

	core.registers.IP = uint16(stackPntrAddr)

	core.logInstruction(fmt.Sprintf("[%#04x] retn (current code pointer: %#08x / new ip: %#08x)", current_ip, current_ip, core.registers.IP))

eof:
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}

func INSTR_PUSH(core *CpuCore) {
	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
		{
			// PUSH r16
			val, valName := core.registers.registers16Bit[core.currentOpCodeBeingExecuted-0x50], core.registers.index16ToString(core.currentOpCodeBeingExecuted-0x50)

			err := stackPush16(core, *val)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), valName))

		}
	case 0x6A:
		{
			// PUSH imm8

			val, err := core.readImm8()
			if err != nil {
				goto eof
			}

			err = stackPush8(core, val)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] push %#04x", core.GetCurrentlyExecutingInstructionAddress(), val))
		}
	case 0x68:
		{
			// PUSH imm16

			val, err := core.readImm16()
			if err != nil {
				goto eof
			}

			err = stackPush16(core, val)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] push %#04x", core.GetCurrentlyExecutingInstructionAddress(), val))
		}
	case 0x0E:
		{
			// PUSH CS

			val := core.registers.CS.base

			err := stackPush16(core, val)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "CS"))
		}
	case 0x16:
		{
			// PUSH SS

			val := core.registers.SS.base

			err := stackPush16(core, val)

			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "SS"))
		}
	case 0x1E:
		{
			// PUSH DS

			val := core.registers.DS.base

			err := stackPush16(core, val)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "DS"))
		}
	case 0x06:
		{
			// PUSH ES

			val := core.registers.ES.base

			err := stackPush16(core, val)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] push %s", core.GetCurrentlyExecutingInstructionAddress(), "ES"))
		}
	case 0x60:
		{
			// Push all general purpose registers
			tempSp := core.registers.SP

			err := stackPush16(core, core.registers.AX)
			err = stackPush16(core, core.registers.CX)
			err = stackPush16(core, core.registers.DX)
			err = stackPush16(core, core.registers.BX)
			err = stackPush16(core, tempSp)
			err = stackPush16(core, core.registers.BP)
			err = stackPush16(core, core.registers.SI)
			err = stackPush16(core, core.registers.DI)

			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] pusha", core.GetCurrentlyExecutingInstructionAddress()))
		}

	default:
		log.Println(fmt.Printf("Unhandled PUSH instruction:  %#04x", core.currentOpCodeBeingExecuted))
		doCoreDump(core)
	}

eof:
	core.registers.IP += uint16(core.currentByteAddr - core.currentByteDecodeStart)
}
