package monitor

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type HardwareMonitor struct {
	busId              uint32
	logCpuInstructions bool
}

func NewHardwareMonitor() *HardwareMonitor {
	device := &HardwareMonitor{}
	device.logCpuInstructions = false

	return device
}

func (device *HardwareMonitor) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *HardwareMonitor) OnReceiveMessage(message bus.BusMessage) {
	if device.logCpuInstructions && message.Subject == common.MESSAGE_GLOBAL_CPU_INSTRUCTION_LOG {
		log.Printf("[%#04x] %s", device.busId, message.Data)
	}
}
