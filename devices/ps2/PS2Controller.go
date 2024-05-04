package ps2

import (
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type Ps2Controller struct {
	bus                  *bus.Bus
	busId                uint32
	statusRegister       uint8
	outputBuffer         uint8
	inputBuffer          uint8
	configurationByte    uint8
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
	systemControlPort    uint8

	endpoint *Ps2Device
}

type Ps2Device interface {
	Connect(controller *Ps2Controller)
	Disconnect()
	SendData(data uint8)
	ReceiveData(data uint8)
}

func (controller *Ps2Controller) GetPortMap() *bus.DevicePortMap {
	return &bus.DevicePortMap{
		ReadPorts:  []uint16{0x60, 0x61, 0x64},
		WritePorts: []uint16{0x60, 0x64},
	}
}

func (controller *Ps2Controller) ReadAddr8(addr uint16) uint8 {
	switch addr {
	case 0x60:
		return controller.ReadDataPort()
	case 0x61:
		return controller.ReadSystemControlPort()
	case 0x64:
		return controller.ReadStatusRegister()
	default:
		log.Printf("Invalid PS/2 controller read address: [%#04x]", addr)
		return 0x00
	}
}

func (controller *Ps2Controller) WriteAddr8(addr uint16, data uint8) {
	switch addr {
	case 0x60:
		controller.WriteDataPort(data)
	case 0x61:
		controller.WriteSystemControlPort(data)
	case 0x64:
		controller.WriteControlPort(data)
	default:
		log.Printf("Invalid PS/2 controller write address: [%#04x]", addr)
	}
}

func (controller *Ps2Controller) GetDeviceBusId() uint32 {
	return controller.busId
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

func (controller *Ps2Controller) WriteDataPort(value uint8) {
	controller.isCommand = false
	controller.BufferOutputData(value)
	controller.DisableDataPortReadyForWrite()

}

func (controller *Ps2Controller) WriteControlPort(value uint8) {
	controller.isCommand = true
	controller.WriteCommandRegister(value)
}

func (controller *Ps2Controller) ReadStatusRegister() uint8 {
	return controller.statusRegister
}

func (controller *Ps2Controller) updateStatusRegister() {
	controller.statusRegister = 0

	// Bit 0: Output Buffer Status (OBF)
	if controller.dataPortReadEnabled {
		controller.statusRegister |= 0x01 // Set OBF
	}

	// Bit 1: Input Buffer Status (IBF)
	if !controller.dataPortWriteEnabled {
		controller.statusRegister |= 0x02 // Set IBF
	}

	// Bit 2: System Flag
	if controller.systemFlag {
		controller.statusRegister |= 0x04 // Set System Flag
	}

	// Bit 3: Command/Data
	if controller.isCommand {
		controller.statusRegister |= 0x08 // Set Command/Data
	}

	// Bit 4: Keyboard Lock
	if controller.keyboardLocked {
		controller.statusRegister |= 0x10 // Set Keyboard Lock
	}

	// Bit 5: Auxiliary Output Buffer Full
	if controller.auxiliaryBufferFull {
		controller.statusRegister |= 0x20 // Set Auxiliary Output Buffer Full
	}

	// Bit 6: Time-out
	if controller.timeout {
		controller.statusRegister |= 0x40 // Set Time-out
	}

	// Bit 7: Parity error
	if controller.parityError {
		controller.statusRegister |= 0x80 // Set Parity error
	}
}

func (controller *Ps2Controller) ReadDataPort() uint8 {
	data := controller.outputBuffer
	controller.DisableDataPortReadyForRead()
	return data
}

func (controller *Ps2Controller) WriteCommandRegister(value uint8) {
	log.Printf("PS2 controller write command: [%#04x]", value)
	switch value {
	case 0x20: // Read Configuration Byte
		controller.BufferOutputData(controller.configurationByte)
		controller.EnableDataPortReadyForRead()
	case 0x60: // Write Configuration Byte
		controller.configurationByte = controller.ReadDataPort()
	case 0xA7: // Disable second PS/2 port
		controller.port2_enabled = false
	case 0xA8: // Enable second PS/2 port
		controller.port2_enabled = true
	case 0xA9: // Test second PS/2 port
		// TODO: Implement test logic for second PS/2 port
	case 0xAA: // Test PS/2 Controller
		controller.BufferOutputData(0x55)
		controller.EnableDataPortReadyForRead()
	case 0xAB: // Test first PS/2 port
		// TODO: Implement test logic for first PS/2 port
	case 0xAD: // Disable first PS/2 port
		controller.port1_enabled = false
	case 0xAE: // Enable first PS/2 port
		controller.port1_enabled = true
	case 0xC0: // Read input port
		// TODO: Implement input port reading
	case 0xD0: // Read output port
		// TODO: Implement output port reading
	case 0xD1: // Write output port
		// TODO: Implement output port writing
	case 0xD2: // Write keyboard output buffer
		// TODO: Implement keyboard output buffer writing
	case 0xD3: // Write mouse output buffer
		// TODO: Implement mouse output buffer writing
	case 0xD4: // Write to mouse
		// TODO: Implement writing to mouse
	case 0xF0: // Read test inputs
		// TODO: Implement reading test inputs
	default:
		log.Printf("Unknown PS2 controller write command: [%#04x]", value)
	}
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

func (controller *Ps2Controller) ConnectDevice(device Ps2Device) {
	device.Connect(controller)
	controller.endpoint = &device
}

func (controller *Ps2Controller) BufferInputData(data uint8) {
	controller.inputBuffer = data
	controller.DisableDataPortReadyForRead()
	controller.DisableDataPortReadyForWrite()
}

func (controller *Ps2Controller) BufferOutputData(data uint8) {
	controller.outputBuffer = data
	controller.EnableDataPortReadyForRead()
	controller.DisableDataPortReadyForWrite()

	if controller.endpoint != nil {
		device := *controller.endpoint
		device.ReceiveData(data)
	}
}

func (controller *Ps2Controller) ReadSystemControlPort() uint8 {
	return controller.systemControlPort
}

func (controller *Ps2Controller) WriteSystemControlPort(data uint8) {
	controller.systemControlPort = data

	if controller.systemControlPort&0x01 == 0x01 {
		controller.systemFlag = true
	} else {
		controller.systemFlag = false
	}

	if controller.systemControlPort&0x02 == 0x02 {
		controller.keyboardLocked = true
	} else {
		controller.keyboardLocked = false
	}

	if controller.systemControlPort&0x04 == 0x04 {
		controller.auxiliaryBufferFull = true
	} else {
		controller.auxiliaryBufferFull = false
	}

	if controller.systemControlPort&0x08 == 0x08 {
		controller.timeout = true
	} else {
		controller.timeout = false
	}

	if controller.systemControlPort&0x10 == 0x10 {
		controller.parityError = true
	} else {
		controller.parityError = false
	}

	controller.updateStatusRegister()
}
