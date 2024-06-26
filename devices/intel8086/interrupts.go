package intel8086

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
)

func INSTR_INT3(core *CpuCore) {
	core.logInstruction("INT 3")
	core.registers.IP++
	interruptMessage := bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_RAISE,
		Sender:  core.busId,
		Data:    []byte{byte(3)},
	}
	err := core.bus.SendMessageSingle(common.MODULE_INTERRUPT_CONTROLLER_1, interruptMessage)
	if err != nil {
		core.logInstruction("8259A: Error sending interrupt request message: %v", err)
	}
}
