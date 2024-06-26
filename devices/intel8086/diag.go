package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/monitor"
	"log"
	"strings"
)

func doCoreDump(core *CpuCore) {

	core.logInstruction("Instruction log:")
	for _, instruction := range core.bus.FindSingleDevice(common.MODULE_DEBUG_MONITOR).(*monitor.HardwareMonitor).GetInstructionLog() {
		core.logInstruction("%s", instruction)
	}

	log.Println("Dumping core: " + core.FriendlyPartName())

	if core.mode == common.REAL_MODE {
		log.Println("Cpu core in real mode")
	}

	// Gather next few bytes for debugging...
	peekBytes := core.memoryAccessController.PeekNextBytes(core.currentByteDecodeStart, 10)
	stb := strings.Builder{}
	for _, b := range peekBytes {
		stb.WriteString(fmt.Sprintf("%#2x ", b))
	}
	core.logInstruction("Next 10 bytes at instruction pointer: " + stb.String())

	peekBytes = core.memoryAccessController.PeekNextBytes(core.currentByteDecodeStart-10, 20)
	stb = strings.Builder{}
	for _, b := range peekBytes {
		stb.WriteString(fmt.Sprintf("%#2x ", b))
	}
	core.logInstruction("Previous 10 bytes at instruction pointer: " + stb.String())

	core.logInstruction("CS: %#2x, IP: %#2x", core.registers.CS, core.registers.IP)

	core.logInstruction("8 Bit registers:")
	for x, y := range core.registers.registers8Bit {
		core.logInstruction("%v %#2x (pntr: %#2x)", core.registers.index8ToString(uint8(x)), *y, y)
	}
	core.logInstruction("16 Bit registers:")
	for x, y := range core.registers.registers16Bit {
		core.logInstruction("%v %#2x (pntr: %#2x)", core.registers.index16ToString(uint8(x)), *y, y)
	}
	core.logInstruction("Segment registers:")
	for x, y := range core.registers.registersSegmentRegisters {
		core.logInstruction("%v %#2x (pntr: %#2x)", core.registers.indexSegmentToString(uint8(x)), *y, y)
	}

	core.logInstruction("Flags:")
	core.logInstruction("Z: %t", core.registers.GetFlag(ZeroFlag))
	core.logInstruction("D: %t", core.registers.GetFlag(DirectionFlag))
	core.logInstruction("C: %t", core.registers.GetFlag(CarryFlag))
	core.logInstruction("O: %t", core.registers.GetFlag(OverFlowFlag))

	core.logInstruction("Control flags:")
	core.logInstruction("CR0[pe] = %b", core.registers.CR0>>0&1)
	core.logInstruction("CR0[mp] = %b", core.registers.CR0>>1&1)
	core.logInstruction("CR0[em] = %b", core.registers.CR0>>2&1)
	core.logInstruction("CR0[ts] = %b", core.registers.CR0>>3&1)
	core.logInstruction("CR0[et] = %b", core.registers.CR0>>4&1)
	core.logInstruction("CR0[ne] = %b", core.registers.CR0>>5&1)

	log.Print("Other details:")
	core.logInstruction("Is in protected mode: %t", core.mode == common.PROTECTED_MODE)
	core.logInstruction("Is in real mode: %t", core.mode == common.REAL_MODE)
	core.logInstruction("Is decoding 2 byte instruction: %t", core.is2ByteOperand)
}
