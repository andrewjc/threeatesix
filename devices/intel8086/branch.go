package intel8086

import (
	"github.com/andrewjc/threeatesix/common"
	"log"
)

func INSTR_RET_NEAR(core *CpuCore) {

	current_ip := core.GetCurrentCodePointer()

	stackPntrAddr := core.registers.SP

	core.registers.IP = uint16(stackPntrAddr)

	core.registers.SP += 2

	log.Printf("[%#04x] retn (current code pointer: %#08x / new ip: %#08x)", current_ip, current_ip, core.registers.IP)

}

func INSTR_DEC_COUNT_JMP_SHORT_ECX(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

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

	log.Printf("[%#04x] LOOP %#04x (SHORT REL8) %s", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr), extraStr)
	if term1 != 0 && additional {
		core.registers.IP = uint16(destAddr)
		log.Printf("[%#04x]   |-> jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX)
	} else {
		core.registers.IP = uint16(uint16(core.GetIP() + 2))
		log.Printf("[%#04x]   |-> not jumped (CX %#04x)", core.GetCurrentlyExecutingInstructionAddress(), core.registers.CX)
	}

}

func INSTR_JMP_FAR_PTR16(core *CpuCore) {
	destAddr, err := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1)
	segment, err := core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 3)

	log.Printf("[%#04x] JMP %#04x:%#04x (FAR_PTR16)", core.GetCurrentlyExecutingInstructionAddress(), segment, destAddr)
	if err == nil {
		core.writeSegmentRegister(&core.registers.CS, segment)
	}

	core.registers.IP = destAddr
}

func INSTR_JMP_FAR_M16(core *CpuCore, modrm *ModRm) {
	if modrm.mod == 3 {
		addr := core.registers.registers16Bit[modrm.rm]
		core.registers.IP = *addr
		log.Printf("[%#04x] JMP %#04x (JMP_FAR_M16)", core.GetCurrentlyExecutingInstructionAddress(), uint16(*addr))
	} else {
		addr := modrm.getAddressMode16(core)
		core.registers.IP = addr
		log.Printf("[%#04x] JMP %#04x (JMP_FAR_M16)", core.GetCurrentlyExecutingInstructionAddress(), uint16(addr))
	}

}

func INSTR_JMP_NEAR_REL16(core *CpuCore) {

	offset, err := common.Int16Err(core.memoryAccessController.ReadAddr16(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	var destAddr = int16(core.registers.IP + 3)

	destAddr = destAddr + int16(offset)

	log.Printf("[%#04x] JMP %#04x (NEAR_REL16)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr))
	core.registers.IP = uint16(destAddr)
}

func INSTR_JZ_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)

	destAddr = destAddr + uint16(offset)

	log.Printf("[%#04x] JZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr))
	if core.registers.GetFlag(ZeroFlag) {
		core.registers.IP = uint16(destAddr)
		log.Printf("[%#04x]   |-> jumped", core.GetCurrentlyExecutingInstructionAddress())
	} else {
		core.registers.IP = uint16(uint16(core.GetIP() + 2))
		log.Printf("[%#04x]   |-> no jump", core.GetCurrentlyExecutingInstructionAddress())
	}
}

func INSTR_JNZ_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)

	destAddr = destAddr + uint16(offset)

	log.Printf("[%#04x] JNZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr))
	if !core.registers.GetFlag(ZeroFlag) {
		core.registers.IP = uint16(destAddr)
		log.Printf("[%#04x]   |-> jumped", core.GetCurrentlyExecutingInstructionAddress())
	} else {
		core.registers.IP = uint16(core.GetIP() + 2)
		log.Printf("[%#04x]   |-> no jump", core.GetCurrentlyExecutingInstructionAddress())
	}

}

func INSTR_JCXZ_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)

	destAddr = destAddr + uint16(offset)

	log.Printf("[%#04x] JCXZ %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr))
	if core.registers.CX == 0 {
		core.registers.IP = uint16(destAddr)
		log.Printf("[%#04x]   |-> jumped", core.GetCurrentlyExecutingInstructionAddress())
	} else {
		core.registers.IP = uint16(core.GetIP() + 2)
		log.Printf("[%#04x]   |-> no jump", core.GetCurrentlyExecutingInstructionAddress())
	}

}

func INSTR_JMP_SHORT_REL8(core *CpuCore) {

	offset, err := common.Int8Err(core.memoryAccessController.ReadAddr8(uint32(core.GetCurrentCodePointer()) + 1))

	if err != nil {
		return
	}

	var destAddr = uint16(core.registers.IP + 2)

	destAddr = destAddr + uint16(offset)

	log.Printf("[%#04x] JMP %#04x (SHORT REL8)", core.GetCurrentlyExecutingInstructionAddress(), uint16(destAddr))
	core.registers.IP = uint16(destAddr)

}
