package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/monitor"
	"log"
	"strings"
)

func doCoreDump(core *CpuCore) {

	log.Printf("Instruction log:")
	for _, instruction := range core.bus.FindSingleDevice(common.MODULE_DEBUG_MONITOR).(*monitor.HardwareMonitor).GetInstructionLog() {
		log.Printf("%s", instruction)
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
	log.Printf("Next 10 bytes at instruction pointer: " + stb.String())

	peekBytes = core.memoryAccessController.PeekNextBytes(core.currentByteDecodeStart-10, 20)
	stb = strings.Builder{}
	for _, b := range peekBytes {
		stb.WriteString(fmt.Sprintf("%#2x ", b))
	}
	log.Printf("Previous 10 bytes at instruction pointer: " + stb.String())

	log.Printf("CS: %#2x, IP: %#2x", core.registers.CS, core.registers.IP)

	log.Printf("8 Bit registers:")
	for x, y := range core.registers.registers8Bit {
		log.Printf("%v %#2x (pntr: %#2x)", core.registers.index8ToString(uint8(x)), *y, y)
	}
	log.Printf("16 Bit registers:")
	for x, y := range core.registers.registers16Bit {
		log.Printf("%v %#2x (pntr: %#2x)", core.registers.index16ToString(uint8(x)), *y, y)
	}
	log.Printf("Segment registers:")
	for x, y := range core.registers.registersSegmentRegisters {
		log.Printf("%v %#2x (pntr: %#2x)", core.registers.indexSegmentToString(uint8(x)), *y, y)
	}

	log.Printf("Flags:")
	log.Printf("Z: %t", core.registers.GetFlag(ZeroFlag))
	log.Printf("D: %t", core.registers.GetFlag(DirectionFlag))
	log.Printf("C: %t", core.registers.GetFlag(CarryFlag))
	log.Printf("O: %t", core.registers.GetFlag(OverFlowFlag))

	log.Printf("Control flags:")
	log.Printf("CR0[pe] = %b", core.registers.CR0>>0&1)
	log.Printf("CR0[mp] = %b", core.registers.CR0>>1&1)
	log.Printf("CR0[em] = %b", core.registers.CR0>>2&1)
	log.Printf("CR0[ts] = %b", core.registers.CR0>>3&1)
	log.Printf("CR0[et] = %b", core.registers.CR0>>4&1)
	log.Printf("CR0[ne] = %b", core.registers.CR0>>5&1)

	log.Print("Other details:")
	log.Printf("Is in protected mode: %t", core.mode == common.PROTECTED_MODE)
	log.Printf("Is in real mode: %t", core.mode == common.REAL_MODE)
	log.Printf("Is decoding 2 byte instruction: %t", core.is2ByteOperand)
}
