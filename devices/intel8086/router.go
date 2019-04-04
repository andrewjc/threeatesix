package intel8086

import "log"

func (core *CpuCore) routeInstruction() uint8 {

	var instrByte uint8

	core.flags.CS_OVERRIDE = 0x0
	core.flags.CS_OVERRIDE_ENABLE = false

	if core.memoryAccessController.PeekNextBytes(uint32(core.currentlyExecutingInstructionPointer), 1)[0] == 0x2E {
		// Prefix byte
		// cs segment override

		core.flags.CS_OVERRIDE = 0x0
		core.flags.CS_OVERRIDE_ENABLE = true
		core.IncrementIP()

		instrByte = core.memoryAccessController.ReadAddr8(uint32(core.currentlyExecutingInstructionPointer+1))

	} else {

		instrByte = core.memoryAccessController.ReadAddr8(uint32(core.currentlyExecutingInstructionPointer))

	}

	var instructionImpl OpCodeImpl
	if core.memoryAccessController.PeekNextBytes(uint32(core.currentlyExecutingInstructionPointer), 1)[0] == 0x0F {
		// 2 byte opcode
		core.IncrementIP()
		instrByte = core.memoryAccessController.ReadAddr8(uint32(core.currentlyExecutingInstructionPointer+1))
		core.currentByteAtCodePointer = instrByte
		instructionImpl = core.opCodeMap2Byte[core.currentByteAtCodePointer]
	} else {
		core.currentByteAtCodePointer = instrByte
		instructionImpl = core.opCodeMap[core.currentByteAtCodePointer]
	}

	if instructionImpl != nil {
		instructionImpl(core)
	} else {
		log.Printf("[%#04x] Unrecognised opcode: %#2x\n", core.currentlyExecutingInstructionPointer, instrByte)

		log.Printf("CPU CORE ERROR!!!")

		doCoreDump(core)
		panic(0)
	}

	return 0
}

