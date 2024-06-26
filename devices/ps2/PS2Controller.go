package ps2

import (
	"github.com/andrewjc/threeatesix/common"
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
	inputPort            uint8
	outputPort           uint8
	testInputs           uint8

	endpoint Ps2Device
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
		WritePorts: []uint16{0x60, 0x61, 0x64},
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
	// Handle received messages here
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
	controller.BufferInputData(value)
}

func (controller *Ps2Controller) WriteControlPort(value uint8) {
	controller.isCommand = true
	controller.WriteCommandRegister(value)
}

func (controller *Ps2Controller) ReadStatusRegister() uint8 {
	controller.updateStatusRegister()
	return controller.statusRegister
}

func (controller *Ps2Controller) updateStatusRegister() {
	controller.statusRegister = 0

	if controller.dataPortReadEnabled {
		controller.statusRegister |= 0x01
	}

	if !controller.dataPortWriteEnabled {
		controller.statusRegister |= 0x02
	}

	if controller.systemFlag {
		controller.statusRegister |= 0x04
	}

	if controller.isCommand {
		controller.statusRegister |= 0x08
	}

	if controller.keyboardLocked {
		controller.statusRegister |= 0x10
	}

	if controller.auxiliaryBufferFull {
		controller.statusRegister |= 0x20
	}

	if controller.timeout {
		controller.statusRegister |= 0x40
	}

	if controller.parityError {
		controller.statusRegister |= 0x80
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
	case 0x20:
		controller.BufferOutputData(controller.configurationByte)
	case 0x60:
		controller.configurationByte = controller.inputBuffer
	case 0xA7:
		controller.port2_enabled = false
	case 0xA8:
		controller.port2_enabled = true
	case 0xA9:
		controller.BufferOutputData(controller.TestPort2())
	case 0xAA:
		controller.BufferOutputData(0x55)
	case 0xAB:
		controller.BufferOutputData(controller.TestPort1())
	case 0xAD:
		controller.port1_enabled = false
	case 0xAE:
		controller.port1_enabled = true
	case 0xC0:
		controller.BufferOutputData(controller.ReadInputPort())
	case 0xD0:
		controller.BufferOutputData(controller.ReadOutputPort())
	case 0xD1:
		controller.WriteOutputPort(controller.inputBuffer)
	case 0xD2:
		controller.BufferOutputData(controller.inputBuffer)
	case 0xD3:
		controller.BufferOutputData(controller.inputBuffer)
	case 0xD4:
		controller.endpoint.SendData(controller.inputBuffer)
	case 0xF0:
		controller.BufferOutputData(controller.ReadTestInputs())
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
	controller.endpoint = device
}

func (controller *Ps2Controller) BufferInputData(data uint8) {
	controller.inputBuffer = data
	controller.DisableDataPortReadyForRead()
	controller.DisableDataPortReadyForWrite()
	interruptMessage := bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_RAISE,
		Sender:  controller.busId,
		Data:    []byte{byte(2)},
	}
	err := controller.bus.SendMessageSingle(common.MODULE_INTERRUPT_CONTROLLER_1, interruptMessage)
	if err != nil {
		log.Printf("8259A: Error sending interrupt request message: %v", err)
	}
}

func (controller *Ps2Controller) BufferOutputData(data uint8) {
	controller.outputBuffer = data
	controller.EnableDataPortReadyForRead()
	controller.DisableDataPortReadyForWrite()

	if controller.endpoint != nil {
		controller.endpoint.ReceiveData(data)
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

func (controller *Ps2Controller) TestPort1() uint8 {
	if controller.port1_enabled && controller.endpoint != nil {
		controller.endpoint.SendData(0xAA)
		return 0x00 // Success
	}
	return 0x01 // Failure
}

func (controller *Ps2Controller) TestPort2() uint8 {
	if controller.port2_enabled && controller.endpoint != nil {
		controller.endpoint.SendData(0xAA)
		return 0x00 // Success
	}
	return 0x01 // Failure
}

func (controller *Ps2Controller) ReadInputPort() uint8 {
	// Implement reading from the input port
	return controller.inputPort
}

func (controller *Ps2Controller) ReadOutputPort() uint8 {
	// Implement reading from the output port
	return controller.outputPort
}

func (controller *Ps2Controller) WriteOutputPort(data uint8) {
	// Implement writing to the output port
	controller.outputPort = data
}

func (controller *Ps2Controller) ReadTestInputs() uint8 {
	// Implement reading test inputs
	return controller.testInputs
}
