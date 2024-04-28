package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func INSTR_CALLF_M16(core *CpuCore, modrm *ModRm) {
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
}

func INSTR_DEC_COUNT_JMP_SHORT_ECX(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadMemoryAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)

	destAddr = destAddr + uint16(offset)

	term1 := uint16(core.registers.CX)
	term1 = term1 - 1
	core.registers.CX = uint16(term1)

	additional := true
	extraStr := ""

	if core.currentOpCodeBeingExecuted == 0xE0 {
		extraStr = "ZF=0"
		additional = !core.registers.GetFlag(ZeroFlag)
	}
	if core.currentOpCodeBeingExecuted == 0xE1 {
		extraStr = "ZF=1"
		additional = core.registers.GetFlag(ZeroFlag)
	}

	core.logInstruction(fmt.Sprintf("[%#04x] LOOP %#04x (SHORT REL8) %s", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr), extraStr))
	if term1 != 0 && additional {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x]   |-> jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
	} else {
		core.registers.IP = uint16(uint16(core.GetIP() + 2))
		core.logInstruction(fmt.Sprintf("[%#04x]   |-> not jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX))
	}

}

func INSTR_JMP_FAR_PTR16(core *CpuCore) {
	destAddr, err := core.memoryAccessController.ReadMemoryAddr16(uint32(core.GetCurrentCodePointer()) + 1)
	segment, err := core.memoryAccessController.ReadMemoryAddr16(uint32(core.GetCurrentCodePointer()) + 3)

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x:%#04x (FAR_PTR16)", core.GetCurrentlyExecutingInstructionAddress(), segment, destAddr))
	if err == nil {
		core.writeSegmentRegister(&core.registers.CS, uint16(segment))
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

	var destAddr = core.registers.IP + 2
	destAddr2 := int32(destAddr) + int32(int16(offset))

	core.logInstruction(fmt.Sprintf("[%#04x] JMP %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr2)))
	core.registers.IP = uint16(destAddr2)
}

func INSTR_CALL_NEAR_REL16(core *CpuCore) {

	offset, err := core.memoryAccessController.ReadMemoryAddr16(uint32(core.GetCurrentCodePointer()) + 1)

	if err != nil {
		return
	}

	var destAddr = core.registers.IP + 3

	var destAddrTest = int32(destAddr)

	destAddr2 := destAddrTest + int32(offset)

	err = stackPush16(core, uint16(uint16(core.GetIP()+3)))
	if err != nil {
		log.Printf("Error pushing to stack: %s", err.Error())
		return
	}

	core.logInstruction(fmt.Sprintf("[%#04x] CALL %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	core.registers.IP = uint16(destAddr2)
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

	if !core.registers.GetFlag(SignFlag) {
		core.registers.IP = uint16(destAddr)
		core.logInstruction(fmt.Sprintf("[%#04x] JS %#04x (SHORT REL8) (Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	} else {
		core.registers.IP = uint16(core.GetIP() + 2)
		core.logInstruction(fmt.Sprintf("[%#04x] JS %#04x (SHORT REL8) (Skipped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
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
		core.registers.IP = uint16(core.GetIP() + 2)
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
		core.registers.IP = uint16(core.GetIP() + 2)
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
		core.registers.IP = uint16(core.GetIP() + 2)
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
		core.registers.IP = uint16(core.GetIP() + 2)
		core.logInstruction(fmt.Sprintf("[%#04x] JNO %#04x (SHORT REL8) (Not Jumped)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr)))
	}
}
