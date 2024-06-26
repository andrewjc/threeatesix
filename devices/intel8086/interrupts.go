package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
	"math/bits"
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

func (core *CpuCore) checkInterrupts() {
	if core.registers.GetFlag(InterruptFlag) && core.hasPendingInterrupt() {
		core.handleInterrupt()
	}
}

func (core *CpuCore) handleInterrupt() {
	if !core.registers.GetFlag(InterruptFlag) {
		return // Interrupts are disabled
	}

	intNum, hasInterrupt := core.getHighestPriorityInterrupt()
	if !hasInterrupt {
		return
	}

	// Send an interrupt message to the interrupt controller
	interruptMessage := bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_RAISE,
		Sender:  core.busId,
		Data:    []byte{intNum},
	}

	err := core.bus.SendMessageSingle(common.MODULE_INTERRUPT_CONTROLLER_1, interruptMessage)
	if err != nil {
		core.logDebug(fmt.Sprintf("CPU: Error sending interrupt request message: %v", err))
	}

	// Send a debug message
	core.logDebug(fmt.Sprintf("CPU: Interrupt %d raised", intNum))

}

func (core *CpuCore) getHighestPriorityInterrupt() (uint8, bool) {
	// Check master PIC
	masterIRQ := core.interruptControllerMaster.IrqRequest & ^core.interruptControllerMaster.IrqMask
	if masterIRQ != 0 {
		irqNum := uint8(bits.TrailingZeros8(masterIRQ))
		return core.interruptControllerMaster.InterruptVectorBase + irqNum, true
	}

	// Check slave PIC
	if core.interruptControllerMaster.IrqMask&(1<<2) == 0 {
		slaveIRQ := core.interruptControllerSlave.IrqRequest & ^core.interruptControllerSlave.IrqMask
		if slaveIRQ != 0 {
			irqNum := uint8(bits.TrailingZeros8(slaveIRQ))
			return core.interruptControllerSlave.InterruptVectorBase + irqNum, true
		}
	}

	return 0, false
}

func (core *CpuCore) hasPendingInterrupt() bool {

	masterIRQ := core.interruptControllerMaster.IrqRequest & ^core.interruptControllerMaster.IrqMask
	if masterIRQ != 0 {
		return true
	}

	// Check slave PIC only if its interrupt line on master is unmasked
	if core.interruptControllerMaster.IrqMask&(1<<2) == 0 {
		slaveIRQ := core.interruptControllerSlave.IrqRequest & ^core.interruptControllerSlave.IrqMask
		if slaveIRQ != 0 {
			return true
		}
	}

	return false
}

func (core *CpuCore) HandleInterruptBusMessage(message bus.BusMessage) {

	vector := message.Data[0]

	// if vector comes from the second interrupt controller, subtract 8
	if message.Sender == common.MODULE_INTERRUPT_CONTROLLER_2 {
		vector -= 8
	}

	// Disable interrupts
	core.registers.SetFlag(InterruptFlag, false)

	// Push the current flags and CS:IP onto the stack
	err := stackPush16(core, core.registers.FLAGS)
	if err != nil {
		return
	}
	err = stackPush32(core, core.registers.CS.Base)
	if err != nil {
		return
	}
	err = stackPush16(core, core.registers.IP)
	if err != nil {
		return
	}

	// Set the necessary flags
	core.registers.SetFlag(TrapFlag, false)

	// Set the CS:IP to the interrupt vector

	vectorAddr := uint16(vector) << 2

	// Read IP and CS from the interrupt vector table

	vectorAddr2 := uint32(0x000020)<<4 + uint32(vectorAddr)

	core.registers.IP, _ = core.memoryAccessController.ReadMemoryValue16(uint32(vectorAddr2))
	newBase, _ := core.memoryAccessController.ReadMemoryValue16(uint32(vectorAddr2 + 3))

	core.registers.CS.Base = uint32(newBase)

	// Re-enable interrupts
	core.registers.SetFlag(InterruptFlag, true)

	// Send a message to the debug monitor
	core.logDebug(fmt.Sprintf("Interrupt %d handled", vector))

	// Send EOI message to the 8259A
	//core.sendEOI(vector)
	// Send EOI command to the 8259A
	eoiMessage := bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_COMPLETE,
		Sender:  message.Sender,
		Data:    []byte{vector},
	}
	err = core.bus.SendMessageToDeviceById(message.Sender, eoiMessage)
	if err != nil {
		core.logDebug(fmt.Sprintf("CPU: Error sending EOI message: %v", err))
	}
}

func (core *CpuCore) AcknowledgeInterrupt(message bus.BusMessage) {
	// send message to the interrupt controller that raised the interrupt
	err := core.bus.SendMessageToDeviceById(message.Sender, bus.BusMessage{Subject: common.MESSAGE_INTERRUPT_ACKNOWLEDGE, Data: message.Data})
	if err != nil {
		log.Fatalf("Failed to acknowledge interrupt: %s", err)
		return
	}
}
