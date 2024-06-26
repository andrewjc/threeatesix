package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func INSTR_CALL_M16(core *CpuCore) {
	modrm, _, err := core.consumeModRm()
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error reading ModR/M byte: %s", err))
		return
	}

	addr, addrName, err := core.readRm16(&modrm)
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error reading address: %s", err))
		return
	}

	stackPush16(core, uint16(core.GetIP()+2))
	core.registers.IP = uint16(*addr)
	core.logInstruction(fmt.Sprintf("[%#04x] CALL %s (%#04x)", core.GetCurrentlyExecutingInstructionAddress(), addrName, uint16(*addr)))
}

func INSTR_CALL_RM16(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		core.logInstruction("Error reading modrm: %s", err.Error())
		doCoreDump(core)
		panic(0)
	}

	addr, addrName, err := core.readRm16(&modrm)
	if err != nil {
		core.logInstruction("Error reading address: %s", err.Error())
		doCoreDump(core)
		panic(0)
	}

	stackPush16(core, uint16(core.GetIP()+2))
	core.registers.IP = uint16(*addr)
	core.logInstruction(fmt.Sprintf("[%#04x] CALL %s (%#04x)", core.GetCurrentlyExecutingInstructionAddress(), addrName, uint16(*addr)))
}

func INSTR_DEC_COUNT_JMP_SHORT(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)
	destAddr += uint16(offset) // Correctly calculate destination address

	term1 := uint16(core.registers.CX)
	result := term1 - 1 // Decrement CX
	core.registers.CX = result

	// Update flags
	dataSize := uint16(16) // Assuming CX is a 16-bit register
	sign1 := int16(term1 >> (dataSize - 1))
	signr := int16((result >> (dataSize - 1)) & 0x01)

	core.registers.SetFlag(CarryFlag, result > term1) // Set if borrow occurs

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr == 1)

	core.registers.SetFlag(OverFlowFlag, (sign1 == 0) && (signr == 1)) // Set if result exceeds maximum negative value

	silenceLogging := false
	if core.lastExecutedInstructionPointer == core.GetCurrentlyExecutingInstructionAddress() {
		silenceLogging = true
	}

	if !silenceLogging {
		core.logInstruction(fmt.Sprintf("[%#04x] LOOP %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
	if result != 0 {
		core.registers.IP = uint16(destAddr)
		if !silenceLogging {
			core.logInstruction(fmt.Sprintf("[%#04x]   |-> jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
		}
	} else {
		core.registers.IP += 2
		if !silenceLogging {
			core.logInstruction(fmt.Sprintf("[%#04x]   |-> not jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
		}
	}

}

func INSTR_DEC_COUNT_JMP_SHORT_Z(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)
	destAddr += uint16(offset) // Correctly calculate destination address

	term1 := uint16(core.registers.CX)
	result := term1 - 1 // Decrement CX
	core.registers.CX = result

	// Update flags
	dataSize := uint16(16) // Assuming CX is a 16-bit register
	sign1 := int16(term1 >> (dataSize - 1))
	signr := int16((result >> (dataSize - 1)) & 0x01)

	core.registers.SetFlag(CarryFlag, result > term1) // Set if borrow occurs

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr == 1)

	core.registers.SetFlag(OverFlowFlag, (sign1 == 0) && (signr == 1)) // Set if result exceeds maximum negative value

	additional := true
	extraStr := ""

	if core.currentOpCodeBeingExecuted == 0xE0 {
		extraStr = "ZF=0"
		additional = !core.registers.GetFlag(ZeroFlag) // Check for not zero flag
	}
	if core.currentOpCodeBeingExecuted == 0xE1 {
		extraStr = "ZF=1"
		additional = core.registers.GetFlag(ZeroFlag) // Check for zero flag
	}

	core.logInstruction(fmt.Sprintf("[%#04x] LOOP %#04x (SHORT REL8) %s", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr), extraStr))
	if result != 0 && additional {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x]   |-> jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
	} else {
		core.registers.IP += 2
		core.logInstruction(fmt.Sprintf("[%#04x]   |-> not jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
	}

}

func INSTR_JMP_FAR_PTR16(core *CpuCore) {

	destAddr, err := core.memoryAccessController.ReadMemoryValue16(uint32(core.GetCurrentCodePointer()) + 1)
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error reading memory address: %s", err))
		return
	}
	segment, err := core.memoryAccessController.ReadMemoryValue16(uint32(core.GetCurrentCodePointer()) + 3)
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error reading memory address: %s", err))
		return
	}

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x:%#04x (FAR_PTR16)", core.GetCurrentlyExecutingInstructionAddress(), segment, destAddr))
	if err == nil && segment < 0xFFFF {
		segmentBase := uint32(segment)
		core.writeSegmentRegister(&core.registers.CS, segmentBase)
	}

	core.registers.IP = uint16(destAddr)
}

func INSTR_JMP_FAR_M16(core *CpuCore) {

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	if err != nil {
		core.logInstruction("Error consuming ModR/M byte: %v\n", err)
		return // Exit early on error
	}
	core.currentByteAddr--

	addr, addrName, err := core.readRm16(&modrm)
	if err != nil {
		log.Fatalf("Error reading address: %s", err.Error())
		doCoreDump(core)
		panic(0)
	}

	newCSAddr, _ := core.getEffectiveAddress16(&modrm)
	newCS, err := core.memoryAccessController.ReadMemoryValue16(uint32(newCSAddr + 2))
	if err != nil {
		core.logInstruction(fmt.Sprintf("Error reading new CS: %s", err.Error()))
		return
	}

	core.registers.IP = uint16(*addr)
	core.writeSegmentRegister(&core.registers.CS, uint32(newCS))

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %s (JMP_FAR_M16) (dst=%#04x:%#04x)",
		core.GetCurrentlyExecutingInstructionAddress(), addrName, newCS, uint16(*addr)))
}

func INSTR_JMP_NEAR_REL16(core *CpuCore) {

	offset, err := core.memoryAccessController.ReadMemoryValue16(uint32(core.GetCurrentCodePointer()) + 1)
	if err != nil {
		return
	}

	var destAddr = core.registers.IP + 3
	destAddr2 := int32(destAddr) + int32(int16(offset))

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr2)))
	core.registers.IP = uint16(destAddr2)
}

func INSTR_CALL_NEAR_REL16(core *CpuCore) {

	// Get the current IP and read the offset
	currentIP := core.GetIP()
	offset, err := core.readImm16()

	// Calculate the address after the instruction and the destination address
	nextIP := currentIP + 3                           // Size of this CALL instruction
	destAddr := uint16(int32(nextIP) + int32(offset)) // nextIP is the base from which the offset is applied

	// Push the return address (next instruction address) onto the stack
	err = stackPush16(core, nextIP)
	if err != nil {
		core.logInstruction("Error pushing to stack: %s", err.Error())
		return
	}

	core.logInstruction(fmt.Sprintf("[%#04x] CALL %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	core.registers.IP = destAddr
}

func INSTR_JS_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	if core.registers.GetFlag(SignFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JS %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP = uint16(core.GetIP() + 2)
		core.logInstruction(fmt.Sprintf("[%#04x] JS %#04x (SHORT REL8) (Skipped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
}

func INSTR_JNS_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	// Jump if NOT Sign (SignFlag is not set)
	if !core.registers.GetFlag(SignFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JNS %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP += 2
		core.logInstruction(fmt.Sprintf("[%#04x] JNS %#04x (SHORT REL8) (Not Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
}

func INSTR_JZ_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	if core.registers.GetFlag(ZeroFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JZ %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP = uint16(core.GetIP() + 2)
		core.logInstruction(fmt.Sprintf("[%#04x] JZ %#04x (SHORT REL8) (Skipped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
}

func INSTR_JNZ_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	if !core.registers.GetFlag(ZeroFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JNZ %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP = uint16(core.GetIP() + 2)
		core.logInstruction(fmt.Sprintf("[%#04x] JNZ %#04x (SHORT REL8) (Skipped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
}

func INSTR_JBE_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	if core.registers.GetFlag(CarryFlag) || core.registers.GetFlag(ZeroFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JBE %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP += 2
		core.logInstruction(fmt.Sprintf("[%#04x] JBE %#04x (SHORT REL8) (Skipped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}

}

func INSTR_JCXZ_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	if core.registers.CX == 0 {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JCXZ %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP += 2
		core.logInstruction(fmt.Sprintf("[%#04x] JCXZ %#04x (SHORT REL8) (Skipped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}

}

func INSTR_JMP_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	core.registers.IP = uint16(destAddr)

}

func INSTR_JO_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	if core.registers.GetFlag(OverFlowFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JO %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP += 2
		core.logInstruction(fmt.Sprintf("[%#04x] JO %#04x (SHORT REL8) (Not Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
}

func INSTR_JNO_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	if !core.registers.GetFlag(OverFlowFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JNO %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP += 2
		core.logInstruction(fmt.Sprintf("[%#04x] JNO %#04x (SHORT REL8) (Not Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
}

func INSTR_JLE_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if core.registers.GetFlag(ZeroFlag) || core.registers.GetFlag(SignFlag) != core.registers.GetFlag(OverFlowFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JLE %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}

func INSTR_JG_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if !core.registers.GetFlag(ZeroFlag) && core.registers.GetFlag(SignFlag) == core.registers.GetFlag(OverFlowFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JG %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}

func INSTR_JB_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if core.registers.GetFlag(CarryFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JB %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}

func INSTR_JNB_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if !core.registers.GetFlag(CarryFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JNB %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}

func INSTR_JL_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if core.registers.GetFlag(SignFlag) != core.registers.GetFlag(OverFlowFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JL %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}

func INSTR_JGE_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if core.registers.GetFlag(SignFlag) == core.registers.GetFlag(OverFlowFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JGE %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}

func INSTR_JPE_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if core.registers.GetFlag(ParityFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JPE %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}

func INSTR_JPO_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryValue8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}
	destAddr := uint16(int16(core.registers.IP) + 2 + int16(offset))
	if !core.registers.GetFlag(ParityFlag) {
		core.registers.IP = destAddr
		core.logInstruction(fmt.Sprintf("[%#04x] JPO %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	} else {
		core.registers.IP += 2
	}
}
