package kb

import (
	"github.com/andrewjc/threeatesix/devices/ps2"
	"log"
)

type Ps2Keyboard struct {
	controller *ps2.Ps2Controller
	scanCodes  []uint8
}

func NewPs2Keyboard() ps2.Ps2Device {
	return &Ps2Keyboard{
		scanCodes: make([]uint8, 0),
	}
}

func (kb *Ps2Keyboard) Connect(controller *ps2.Ps2Controller) {
	kb.controller = controller
	log.Printf("PS/2 Keyboard connected")
	// Perform any necessary initialization or setup
}

func (kb *Ps2Keyboard) Disconnect() {
	kb.controller = nil
	log.Printf("PS/2 Keyboard disconnected")
	// Perform any necessary cleanup or teardown
}

func (kb *Ps2Keyboard) SendData(data uint8) {
	// Send the data to the PS/2 controller's input buffer
	kb.controller.BufferInputData(data)
	kb.controller.DisableDataPortReadyForRead()
	log.Printf("PS/2 Keyboard sent data: %#02x", data)
}

func (kb *Ps2Keyboard) ReceiveData(data uint8) {
	// Process the received data from the PS/2 controller

	// is the command byte set
	if kb.controller.ReadStatusRegister()&0x08 == 0x08 {
		// if so, process the data as a command
		kb.processCommand(data)
		log.Printf("PS/2 Keyboard received command: %#02x", data)
	} else {
		// otherwise, process the data as a scan code
		kb.processScanCode(data)
	}
}

func (kb *Ps2Keyboard) processScanCode(scanCode uint8) {
	// Process the received scan code
	kb.scanCodes = append(kb.scanCodes, scanCode)
	log.Printf("PS/2 Keyboard received scan code: %#02x", scanCode)
	// Perform any necessary actions based on the received scan code
}

func (kb *Ps2Keyboard) GetScanCodes() []uint8 {
	// Return the accumulated scan codes
	return kb.scanCodes
}

func (kb *Ps2Keyboard) ClearScanCodes() {
	// Clear the accumulated scan codes
	kb.scanCodes = make([]uint8, 0)
}

func (kb *Ps2Keyboard) processCommand(data uint8) {
	if data == 0x55 {
		// BAT (Basic Assurance Test) command
		// Perform the BAT and return the result
		kb.SendData(0xAA)
	} else {
		// Unrecognised command
		log.Printf("PS/2 Keyboard received unknown command: %#02x", data)
	}
}
