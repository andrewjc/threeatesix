package monitor

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type HardwareMonitor struct {
	busId              uint32
	logCpuInstructions bool
	instructionLog     []string
}

const MAX_LOG_LENGTH = 64

func NewHardwareMonitor() *HardwareMonitor {
	device := &HardwareMonitor{}
	device.logCpuInstructions = false
	device.instructionLog = make([]string, 0)

	return device
}

func (device *HardwareMonitor) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *HardwareMonitor) OnReceiveMessage(message bus.BusMessage) {
	if device.logCpuInstructions && message.Subject == common.MESSAGE_GLOBAL_CPU_INSTRUCTION_LOG {
		log.Printf("[%#04x] %s", device.busId, message.Data)
	} else {
		device.instructionLog = append(device.instructionLog, fmt.Sprintf("%s", message.Data))
		if len(device.instructionLog) > MAX_LOG_LENGTH {
			device.instructionLog = device.instructionLog[1:]
		}
	}
}

func (device *HardwareMonitor) GetInstructionLog() []string {
	return device.instructionLog
}
