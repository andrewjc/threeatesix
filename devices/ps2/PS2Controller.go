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

	commandByte        uint8
	lastCommand        uint8
	expectingParameter bool
	selfTestPassed     bool
	keyboardEnabled    bool
	mouseEnabled       bool
	keyboardIRQEnabled bool
	mouseIRQEnabled    bool

	refreshCycleToggle   bool
	ioChannelCheck       bool
	ioChannelCheckStatus bool
	keyboardA20          bool
	outputRegisterFull   bool
	inputRegisterFull    bool
	clockGate2           bool
	speakerData          bool

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
	controller := &Ps2Controller{
		port1_enabled:        true,
		port2_enabled:        true,
		dataPortWriteEnabled: true,
		dataPortReadEnabled:  false,
		selfTestPassed:       true,
		keyboardEnabled:      false,
		mouseEnabled:         false,
		keyboardIRQEnabled:   false,
		mouseIRQEnabled:      false,
	}
	controller.resetController()
	return controller
}

func (controller *Ps2Controller) resetController() {
	controller.commandByte = 0x5D // Default command byte value
	controller.lastCommand = 0
	controller.expectingParameter = false
	controller.systemFlag = false
	controller.isCommand = false
	controller.ioChannelCheck = false
	controller.ioChannelCheckStatus = false
	controller.keyboardLocked = false
	controller.auxiliaryBufferFull = false
	controller.timeout = false
	controller.parityError = false
	controller.inputBuffer = 0
	controller.outputBuffer = 0
	controller.systemControlPort = 0
	controller.updateSystemControlPort()

	log.Println("PS/2 Controller: System Reset requested")
}

func (controller *Ps2Controller) GetBus() *bus.Bus {
	return controller.bus
}

func (controller *Ps2Controller) SetBus(bus *bus.Bus) {
	controller.bus = bus
}

func (controller *Ps2Controller) WriteDataPort(value uint8) {
	if controller.expectingParameter {
		controller.handleCommandParameter(value)
	} else {
		// Handle data port writes (e.g., sending commands to keyboard/mouse)
		if controller.port1_enabled && controller.endpoint != nil {
			controller.endpoint.SendData(value)
		}
	}
}

func (controller *Ps2Controller) handleCommandParameter(value uint8) {
	switch controller.lastCommand {
	case 0x60: // Write command byte
		controller.commandByte = value
		controller.updateControllerState()
	case 0xD1: // Write output port
		controller.WriteOutputPort(value)
	case 0xD4: // Write to second PS/2 port
		if controller.port2_enabled {
			// Send to mouse device if implemented
		}
	default:
		log.Printf("Unexpected parameter for command [%#02x]: [%#02x]", controller.lastCommand, value)
	}
	controller.expectingParameter = false
}

func (controller *Ps2Controller) updateControllerState() {
	controller.keyboardEnabled = controller.commandByte&0x10 != 0
	controller.mouseEnabled = controller.commandByte&0x20 != 0
	controller.keyboardIRQEnabled = controller.commandByte&0x01 != 0
	controller.mouseIRQEnabled = controller.commandByte&0x02 != 0
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
	controller.outputRegisterFull = false
	controller.updateSystemControlPort()
	return data
}

func (controller *Ps2Controller) WriteCommandRegister(value uint8) {
	controller.lastCommand = value
	controller.expectingParameter = false

	switch value {
	case 0x20: // Read command byte
		controller.BufferOutputData(controller.commandByte)
	case 0x60: // Write command byte
		controller.expectingParameter = true
	case 0xA7: // Disable second PS/2 port
		controller.port2_enabled = false
	case 0xA8: // Enable second PS/2 port
		controller.port2_enabled = true
	case 0xA9: // Test second PS/2 port
		controller.BufferOutputData(controller.TestPort2())
	case 0xAA: // Test PS/2 Controller
		controller.performSelfTest()
	case 0xAB: // Test first PS/2 port
		controller.BufferOutputData(controller.TestPort1())
	case 0xAD: // Disable first PS/2 port
		controller.port1_enabled = false
	case 0xAE: // Enable first PS/2 port
		controller.port1_enabled = true
	case 0xF0: // Pulse output line
		// Implement pulse output line logic if needed
	case 0xFE: // Resend
		controller.resetController()
	case 0xFF: // Reset
		controller.resetController()
		controller.BufferOutputData(0xFA) // ACK
	default:
		log.Printf("Unknown PS2 controller command: [%#02x]", value)
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
	controller.DisableDataPortReadyForWrite()
	controller.EnableDataPortReadyForRead()
	controller.outputRegisterFull = true
	controller.updateSystemControlPort()
	if controller.keyboardIRQEnabled {
		controller.triggerInterrupt()
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
	controller.updateSystemControlPort()
	return controller.systemControlPort
}

func (controller *Ps2Controller) updateSystemControlPort() {
	controller.systemControlPort &= 0x0F // Clear the upper 4 bits

	if controller.ioChannelCheck {
		controller.systemControlPort |= 0x10
	}
	if controller.ioChannelCheckStatus {
		controller.systemControlPort |= 0x20
	}
	if controller.refreshCycleToggle {
		controller.systemControlPort |= 0x40
	}
	controller.refreshCycleToggle = !controller.refreshCycleToggle // Toggle for next read

	if controller.outputRegisterFull {
		controller.systemControlPort |= 0x80
	}
}

func (controller *Ps2Controller) WriteSystemControlPort(data uint8) {
	controller.systemControlPort = data
	controller.speakerData = data&0x02 != 0
	controller.clockGate2 = data&0x01 != 0

	if data&0x80 != 0 {
		controller.ioChannelCheckStatus = false
	}
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

func (controller *Ps2Controller) performSelfTest() {
	if controller.selfTestPassed {
		controller.BufferOutputData(0x55) // Test passed
	} else {
		controller.BufferOutputData(0xFC) // Test failed
	}
}

func (controller *Ps2Controller) triggerInterrupt() {
	interruptMessage := bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_RAISE,
		Sender:  controller.busId,
		Data:    []byte{byte(1)}, // IRQ 1 for keyboard
	}
	err := controller.bus.SendMessageSingle(common.MODULE_INTERRUPT_CONTROLLER_1, interruptMessage)
	if err != nil {
		log.Printf("PS/2 Controller: Error sending interrupt request message: %v", err)
	}
}

func (controller *Ps2Controller) ResetKeyboard() {
	if controller.endpoint != nil {
		controller.endpoint.SendData(0xFF) // Send reset command to keyboard
	}
}

func (controller *Ps2Controller) SetTypematicRateDelay(data uint8) {
	if controller.endpoint != nil {
		controller.endpoint.SendData(0xF3) // Set typematic rate/delay command
		controller.endpoint.SendData(data) // Send the rate/delay value
	}
}

func (controller *Ps2Controller) EnableKeyboardScanning() {
	if controller.endpoint != nil {
		controller.endpoint.SendData(0xF4) // Enable scanning command
	}
}

func (controller *Ps2Controller) EnableIRQ1() {
	controller.commandByte |= 0x01
	controller.keyboardIRQEnabled = true
}
