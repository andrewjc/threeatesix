package monitor

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type HardwareMonitor struct {
	bus                *bus.Bus
	busId              uint32
	logCpuInstructions bool
	logDebugMessage    bool
	instructionLog     []string
}

func (device *HardwareMonitor) GetPortMap() *bus.DevicePortMap {
	return nil
}

func (device *HardwareMonitor) ReadAddr8(addr uint16) uint8 {
	//TODO implement me
	panic("implement me")
}

func (device *HardwareMonitor) WriteAddr8(addr uint16, data uint8) {

}

const MAX_LOG_LENGTH = 64

func NewHardwareMonitor() *HardwareMonitor {
	device := &HardwareMonitor{}
	device.logCpuInstructions = false
	device.logDebugMessage = true
	device.instructionLog = make([]string, 0)

	return device
}

func (device *HardwareMonitor) GetDeviceBusId() uint32 {
	return device.busId
}

func (device *HardwareMonitor) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *HardwareMonitor) SetBus(bus *bus.Bus) {
	device.bus = bus
}

func (device *HardwareMonitor) OnReceiveMessage(message bus.BusMessage) {
	if device.logCpuInstructions && message.Subject == common.MESSAGE_GLOBAL_CPU_INSTRUCTION_LOG {
		log.Output(4, fmt.Sprintf("[%#04x] %s", device.busId, message.Data))
	} else if device.logDebugMessage && message.Subject == common.MESSAGE_GLOBAL_DEBUG_MESSAGE_LOG {
		log.Output(4, fmt.Sprintf("[%#04x] %s", device.busId, message.Data))
	}
}

func (device *HardwareMonitor) GetInstructionLog() []string {
	return device.instructionLog
}
