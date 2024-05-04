package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func INSTR_CALL_M16(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	if modrm.mod == 3 {
		reg, reg_str := core.registers.registers16Bit[modrm.rm], core.registers.index16ToString(modrm.rm)

		stackPush16(core, uint16(uint16(core.GetIP()+2)))

		core.registers.IP = uint16(*reg)
		core.logInstruction(fmt.Sprintf("[%#04x] CALLF %s (%#04x)", core.GetCurrentlyExecutingInstructionAddress(), reg_str, uint16(*reg)))
	} else {
		addr := modrm.getAddressMode16(core)
		stackPush16(core, uint16(uint16(core.GetIP()+2)))

		core.registers.IP = uint16(addr)
		core.logInstruction(fmt.Sprintf("[%#04x] CALLF %#04x", core.GetCurrentlyExecutingInstructionAddress(), uint16(addr)))
	}

eof:
}

func INSTR_CALL_RM16(core *CpuCore) {
	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	core.currentByteAddr--
	if err != nil {
		goto eof
	}

	if modrm.mod == 3 {
		reg, reg_str := core.registers.registers16Bit[modrm.rm], core.registers.index16ToString(modrm.rm)

		stackPush16(core, uint16(core.GetIP()+2))

		core.registers.IP = uint16(*reg)
		core.logInstruction(fmt.Sprintf("[%#04x] CALL %s (%#04x)", core.GetCurrentlyExecutingInstructionAddress(), reg_str, uint16(*reg)))
	} else {
		addr := modrm.getAddressMode16(core)
		stackPush16(core, uint16(core.GetIP()+2))

		core.registers.IP = uint16(addr)
		core.logInstruction(fmt.Sprintf("[%#04x] CALL %#04x", core.GetCurrentlyExecutingInstructionAddress(), uint16(addr)))
	}

eof:
}

func INSTR_DEC_COUNT_JMP_SHORT_ECX(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)
	destAddr += uint16(offset) // Correctly calculate destination address

	term1 := uint16(core.registers.CX)
	term1-- // Decrement CX
	core.registers.CX = term1

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
	if term1 != 0 && additional {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x]   |-> jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
	} else {
		core.registers.IP += 2
		core.logInstruction(fmt.Sprintf("[%#04x]   |-> not jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
	}

}

func INSTR_JMP_FAR_PTR16(core *CpuCore) {
	core.currentByteAddr++
	destAddr, err := core.memoryAccessController.ReadMemoryAddr16(uint32(core.GetCurrentCodePointer()) + 1)
	if err != nil {
		fmt.Println("Error reading memory address")
		return
	}
	segment, err := core.memoryAccessController.ReadMemoryAddr16(uint32(core.GetCurrentCodePointer()) + 3)
	if err != nil {
		fmt.Println("Error reading memory address")
		return
	}

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x:%#04x (FAR_PTR16)", core.GetCurrentlyExecutingInstructionAddress(), segment, destAddr))
	if err == nil {
		segmentBase := uint32(segment)
		core.writeSegmentRegister(&core.registers.CS, segmentBase)
	}

	core.registers.IP = uint16(destAddr)
}

func INSTR_JMP_FAR_M16(core *CpuCore, modrm *ModRm) {
	if modrm.mod == 3 {
		reg, reg_str := core.registers.registers16Bit[modrm.rm], core.registers.index16ToString(modrm.rm)
		core.registers.IP = uint16(*reg)
		core.logInstruction(fmt.Sprintf("[%#04x] JMP %s (%#04x) (JMP_FAR_M16)", core.GetCurrentlyExecutingInstructionAddress(), reg_str, uint16(*reg)))
	} else {
		addr := modrm.getAddressMode16(core)
		core.registers.IP = uint16(addr)
		core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x (JMP_FAR_M16)", core.GetCurrentlyExecutingInstructionAddress(), uint16(addr)))
	}
}

func INSTR_JMP_NEAR_REL16(core *CpuCore) {

	offset, err := core.memoryAccessController.ReadMemoryAddr16(uint32(core.GetCurrentCodePointer()) + 1)
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
	offsetAddress := uint32(currentIP) + 1 // Convert currentIP to uint32 before adding

	// Read the offset directly as a signed 16-bit integer
	offsetBytes, err := core.memoryAccessController.ReadMemoryAddr16(offsetAddress)
	if err != nil {
		return
	}

	// Convert the offset from unsigned to signed
	offset := int16(offsetBytes)

	// Calculate the address after the instruction and the destination address
	nextIP := currentIP + 3                           // Size of this CALL instruction
	destAddr := uint16(int32(nextIP) + int32(offset)) // nextIP is the base from which the offset is applied

	// Push the return address (next instruction address) onto the stack
	err = stackPush16(core, nextIP)
	if err != nil {
		log.Printf("Error pushing to stack: %s", err.Error())
		return
	}

	core.logInstruction(fmt.Sprintf("[%#04x] CALL %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionAddress(), destAddr))
	core.registers.IP = destAddr
}

func INSTR_JS_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

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

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

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

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

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

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

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

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

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

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

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

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	tmp := int16(core.registers.IP + 2)
	destAddr := uint16(tmp + int16(offset))

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	core.registers.IP = uint16(destAddr)

}

func INSTR_JO_SHORT_REL8(core *CpuCore) {
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))
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
