package ps2

import (
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type Ps2Controller struct {
	bus                *bus.Bus
	busId              uint32
	statusRegister     uint8
	bufferedOutputData []uint8
	bufferedInputData  []uint8

	internalRam          uint8 // used for storing config bytes
	pendingOperation     uint8
	port1_enabled        bool
	port2_enabled        bool
	dataPortWriteEnabled bool
	dataPortReadEnabled  bool
	systemFlag           bool
	isCommand            bool
	keyboardLocked       bool
	auxiliaryBufferFull  bool
	timeout              bool
	parityError          bool
}

func (controller *Ps2Controller) SetDeviceBusId(id uint32) {
	controller.busId = id
}

func (controller *Ps2Controller) OnReceiveMessage(message bus.BusMessage) {

}

func CreatePS2Controller() *Ps2Controller {
	return &Ps2Controller{}
}

func (controller *Ps2Controller) GetBus() *bus.Bus {
	return controller.bus
}

func (controller *Ps2Controller) SetBus(bus *bus.Bus) {
	controller.bus = bus
}

func (controller *Ps2Controller) ReadStatusRegister() uint8 {
	controller.updateStatusRegister()
	return controller.statusRegister
}

func (controller *Ps2Controller) resetStatusRegister(value uint8) error {
	log.Printf("PS2 controller write status register: [%#04x]", value)
	controller.statusRegister = value

	// Bit 2: System Flag
	controller.systemFlag = false
	if (controller.statusRegister & 0x04) != 0 {
		controller.systemFlag = true
	}

	// Bit 3: Command/Data
	controller.isCommand = false
	if (controller.statusRegister & 0x08) != 0 {
		controller.isCommand = true
	}

	// Bit 4: Keyboard Lock
	controller.keyboardLocked = false
	if (controller.statusRegister & 0x10) != 0 {
		controller.keyboardLocked = true
	}

	// Bit 5: Auxiliary Output Buffer Full
	controller.auxiliaryBufferFull = false
	if (controller.statusRegister & 0x20) != 0 {
		controller.auxiliaryBufferFull = true
	}

	// Bit 6: Time-out
	controller.timeout = false
	if (controller.statusRegister & 0x40) != 0 {
		controller.timeout = true
	}

	// Bit 7: Parity error
	controller.parityError = false
	if (controller.statusRegister & 0x80) != 0 {
		controller.parityError = true
	}

	return nil
}

func (controller *Ps2Controller) updateStatusRegister() {
	controller.statusRegister = 0

	// Bit 0: Output Buffer Status (OBF)
	if len(controller.bufferedOutputData) > 0 {
		controller.statusRegister |= 0x01 // Set OBF
	}

	// Bit 1: Input Buffer Status (IBF)
	if len(controller.bufferedInputData) == 0 {
		controller.statusRegister |= 0x02 // Set IBF
	}

	// Bit 2: System Flag
	// Assuming that system flag is represented by a bool variable 'systemFlag'
	if controller.systemFlag {
		controller.statusRegister |= 0x04 // Set System Flag
	}

	// Bit 3: Command/Data
	// Assuming a bool variable 'isCommand' indicates whether last write was a command
	if controller.isCommand {
		controller.statusRegister |= 0x08 // Set Command/Data
	}

	// Bit 4: Keyboard Lock
	// Assuming a bool variable 'keyboardLocked' indicates the state of keyboard lock
	if controller.keyboardLocked {
		controller.statusRegister |= 0x10 // Set Keyboard Lock
	}

	// Bit 5: Auxiliary Output Buffer Full
	// Assuming a bool variable 'auxiliaryBufferFull' indicates the state of the auxiliary buffer
	if controller.auxiliaryBufferFull {
		controller.statusRegister |= 0x20 // Set Auxiliary Output Buffer Full
	}

	// Bit 6: Time-out
	// Assuming a bool variable 'timeout' indicates the state of time-out
	if controller.timeout {
		controller.statusRegister |= 0x40 // Set Time-out
	}

	// Bit 7: Parity error
	// Assuming a bool variable 'parityError' indicates the state of parity error
	if controller.parityError {
		controller.statusRegister |= 0x80 // Set Parity error
	}
}

func (controller *Ps2Controller) ReadDataPort() uint8 {
	if len(controller.bufferedOutputData) > 0 {
		response := controller.bufferedOutputData[0]
		controller.bufferedOutputData = nil
		controller.DisableDataPortReadyForRead()
		return response
	}
	return 0x00
}

func (controller *Ps2Controller) WriteCommandRegister(value uint8) {
	if controller.isCommand {
		log.Printf("PS2 controller write command: [%#04x]", value)
		if value == 0xAA {
			// Test PS2 Controller
			controller.SendBufferedResponse([]uint8{0x55}) //OK Message
			return
		}

		if value == 0xAD {
			// Disable first ps2 port
			controller.port1_enabled = false
			controller.DisableDataPortReadyForRead()
			controller.DisableDataPortReadyForWrite()
			return
		}

		if value == 0xAE {
			// Enable first ps2 port
			controller.port1_enabled = true
			controller.EnableDataPortReadyForRead()
			controller.EnableDataPortReadyForWrite()
			return
		}

		if value == 0xA7 {
			// Disable second ps2 port
			controller.port2_enabled = false
			return
		}

		if value == 0xA8 {
			// Enable second ps2 port
			controller.port2_enabled = true
			return
		}

		if value == 0x60 {
			// Write next byte to byte 0 of internal data
			controller.pendingOperation = 0x60
			controller.EnableDataPortReadyForWrite()
			return
		}

		if controller.pendingOperation == 0x60 {
			controller.resetStatusRegister(value)
			return
		}

		log.Printf("Unknown PS2 controller write command: [%#04x]", value)
	} else {
		// The value is intended for the device plugged into PS2 PORT 0 or 1
		log.Printf("PS2 Port Output Data: [%#04x]", value)
	}

}

func (controller *Ps2Controller) SendBufferedResponse(response []uint8) {
	controller.bufferedOutputData = response
	controller.EnableDataPortReadyForRead()
}

func (controller *Ps2Controller) EnableDataPortReadyForWrite() {
	controller.dataPortWriteEnabled = true
	controller.updateStatusRegister()
}

func (controller *Ps2Controller) DisableDataPortReadyForWrite() {
	controller.dataPortWriteEnabled = false
	controller.updateStatusRegister()
}
func (controller *Ps2Controller) EnableDataPortReadyForRead() {
	controller.dataPortReadEnabled = true
	controller.updateStatusRegister()
}

func (controller *Ps2Controller) DisableDataPortReadyForRead() {
	controller.dataPortReadEnabled = false
	controller.updateStatusRegister()
}
