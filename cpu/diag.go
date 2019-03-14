package cpu

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
	"strings"
)

func doCoreDump(core *CpuCore) {
	if core.mode == common.REAL_MODE {
		log.Println("Cpu core in real mode")
	}

	// Gather next few bytes for debugging...
	peekBytes := core.memoryAccessController.PeekNextBytes(10)
	stb := strings.Builder{}
	for _, b := range peekBytes {
		stb.WriteString(fmt.Sprintf("%#2x ", b))
	}
	log.Printf("Next 10 bytes at instruction pointer: " + stb.String())

	log.Printf("CS: %#2x, IP: %#2x", core.registers.CS, core.registers.IP)

	log.Printf("8 Bit registers:")
	for x, y := range core.registers.registers8Bit {
		log.Printf("%v %#2x (pntr: %#2x)", core.registers.index8ToString(x), *y, y)
	}
	log.Printf("16 Bit registers:")
	for x, y := range core.registers.registers16Bit {
		log.Printf("%v %#2x (pntr: %#2x)", core.registers.index16ToString(x), *y, y)
	}
	log.Printf("Segment registers:")
	for x, y := range core.registers.registersSegmentRegisters {
		log.Printf("%v %#2x (pntr: %#2x)", core.registers.indexSegmentToString(x), *y, y)
	}

}
