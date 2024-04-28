package intel82C54

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
	"time"
)

type Intel82C54 struct {
	bus                *bus.Bus
	busId              uint32
	controlWord        uint8
	counterRegister    [3]uint16
	counterLatch       [3]uint16
	counterMode        [3]uint8
	counterAccessMode  [3]uint8
	counterInitialized [3]bool
	counterValue       [3]uint16
	counterOutput      [3]bool
	counterNullCount   [3]bool
	counterLatched     [3]bool
	statusLatched      bool
	previousUpdateTime int64
}

func NewIntel82C54() *Intel82C54 {
	return &Intel82C54{}
}

func (p *Intel82C54) GetDeviceBusId() uint32 {
	return p.busId
}

func (p *Intel82C54) SetDeviceBusId(id uint32) {
	p.busId = id
}

func (p *Intel82C54) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (p *Intel82C54) GetPortMap() *bus.DevicePortMap {
	return &bus.DevicePortMap{
		ReadPorts:  []uint16{0x0040, 0x0041, 0x0042, 0x0043},
		WritePorts: []uint16{0x0040, 0x0041, 0x0042, 0x0043},
	}
}

func (p *Intel82C54) ReadAddr8(addr uint16) uint8 {
	switch addr {
	case 0x0040:
		// Read from Counter 0
		return p.CounterRegisterRead(0)
	case 0x0041:
		// Read from Counter 1
		return p.CounterRegisterRead(1)
	case 0x0042:
		// Read from Counter 2
		return p.CounterRegisterRead(2)
	default:
		log.Printf("PIT: Invalid read address: %#04x", addr)
		return 0
	}
}

func (p *Intel82C54) WriteAddr8(addr uint16, data uint8) {
	switch addr {
	case 0x0040:
		// Write to Counter 0
		p.CounterRegisterWrite(0, data)
	case 0x0041:
		// Write to Counter 1
		p.CounterRegisterWrite(1, data)
	case 0x0042:
		// Write to Counter 2
		p.CounterRegisterWrite(2, data)
	case 0x0043:
		// Write to Control Word Register
		p.CommandRegisterWrite(data)
	default:
		log.Printf("PIT: Invalid write address: %#04x", addr)
	}
}

func (p *Intel82C54) CommandRegisterWrite(value uint8) {
	p.controlWord = value

	// Extract the counter index from bits 6-7 of the control word
	counterIndex := (value >> 6) & 0x03

	// Check if it's a counter latch command
	if (value>>6)&0x03 == 0x00 {
		// Counter latch command
		p.counterLatched[counterIndex] = true
		return
	}

	// Check if it's a read-back command
	if (value>>6)&0x03 == 0x03 {
		// Read-back command
		p.handleReadBackCommand(value)
		return
	}

	// Extract the counter access mode from bits 4-5 of the control word
	p.counterAccessMode[counterIndex] = (value >> 4) & 0x03

	// Extract the counter mode from bits 1-3 of the control word
	p.counterMode[counterIndex] = (value >> 1) & 0x07

	// Set the counter as initialized
	p.counterInitialized[counterIndex] = true

	// Reset the counter value and output
	p.counterValue[counterIndex] = 0
	p.counterOutput[counterIndex] = false
}

func (p *Intel82C54) CounterRegisterWrite(counterIndex uint8, value uint8) {
	if counterIndex >= 3 {
		log.Printf("Invalid PIT counter index: %d", counterIndex)
		return
	}

	switch p.counterAccessMode[counterIndex] {
	case 0x00:
		// Counter latch command (handled in CommandRegisterWrite)
		return
	case 0x01:
		// Read/Write LSB only
		// Update the counter register LSB with the new value
		p.counterRegister[counterIndex] = (p.counterRegister[counterIndex] & 0xFF00) | uint16(value)
	case 0x02:
		// Read/Write MSB only
		// Update the counter register MSB with the new value
		p.counterRegister[counterIndex] = (p.counterRegister[counterIndex] & 0x00FF) | (uint16(value) << 8)
	case 0x03:
		// Read/Write LSB then MSB
		if !p.counterInitialized[counterIndex] {
			// Update the counter register LSB with the new value
			p.counterRegister[counterIndex] = (p.counterRegister[counterIndex] & 0xFF00) | uint16(value)
			p.counterInitialized[counterIndex] = true
		} else {
			// Update the counter register MSB with the new value
			p.counterRegister[counterIndex] = (p.counterRegister[counterIndex] & 0x00FF) | (uint16(value) << 8)
			p.counterInitialized[counterIndex] = false

			// Load the new count into the counter
			p.counterValue[counterIndex] = p.counterRegister[counterIndex]
			p.counterNullCount[counterIndex] = false
		}
	}
}

func (p *Intel82C54) CounterRegisterRead(counterIndex uint8) uint8 {
	if counterIndex >= 3 {
		log.Printf("Invalid PIT counter index: %d", counterIndex)
		return 0
	}

	var value uint8

	switch p.counterAccessMode[counterIndex] {
	case 0x00:
		// Counter latch command
		if p.counterLatched[counterIndex] {
			// Return the latched value of the counter
			value = uint8(p.counterLatch[counterIndex] & 0xFF)
			p.counterLatched[counterIndex] = false
		} else {
			// Return the current value of the counter
			value = uint8(p.counterValue[counterIndex] & 0xFF)
		}
	case 0x01:
		// Read/Write LSB only
		// Return the LSB of the counter value
		value = uint8(p.counterValue[counterIndex] & 0xFF)
	case 0x02:
		// Read/Write MSB only
		// Return the MSB of the counter value
		value = uint8((p.counterValue[counterIndex] >> 8) & 0xFF)
	case 0x03:
		// Read/Write LSB then MSB
		if !p.counterInitialized[counterIndex] {
			// Return the LSB of the counter value
			value = uint8(p.counterValue[counterIndex] & 0xFF)
			p.counterInitialized[counterIndex] = true
		} else {
			// Return the MSB of the counter value
			value = uint8((p.counterValue[counterIndex] >> 8) & 0xFF)
			p.counterInitialized[counterIndex] = false
		}
	}

	return value
}

func (p *Intel82C54) handleReadBackCommand(value uint8) {
	// Check if count and/or status should be latched
	latchCount := (value & 0x20) == 0
	latchStatus := (value & 0x10) == 0

	// Check which counters are selected
	for i := 0; i < 3; i++ {
		if (value>>(i+1))&0x01 == 0x01 {
			if latchCount {
				// Latch the count value
				p.counterLatch[i] = p.counterValue[i]
				p.counterLatched[i] = true
			}
			if latchStatus {
				// Latch the status
				p.statusLatched = true
			}
		}
	}
}

func (p *Intel82C54) GetBus() *bus.Bus {
	return p.bus
}

func (p *Intel82C54) SetBus(bus *bus.Bus) {
	p.bus = bus
}

func (p *Intel82C54) Step() {
	// Step the PIT counters at 1.19318 MHz
	previousUpdateTime := p.previousUpdateTime
	currentTime := time.Now().UnixNano() / 1000
	if currentTime-previousUpdateTime < 164800 {
		return
	}

	p.previousUpdateTime = currentTime

	for i := 0; i < 3; i++ {
		if p.counterInitialized[i] {
			switch p.counterMode[i] {
			case 0:
				// Mode 0: Interrupt on Terminal Count
				if p.counterValue[i] == 0 {
					// Counter reached zero, set output high
					p.counterOutput[i] = true
				} else {
					// Decrement the counter value
					p.counterValue[i]--
				}
			case 1:
				// Mode 1: Hardware Retriggerable One-Shot
				if p.counterValue[i] == 0 {
					// Counter reached zero, set output high
					p.counterOutput[i] = true
				} else {
					// Decrement the counter value
					p.counterValue[i]--
				}
			case 2:
				// Mode 2: Rate Generator
				if p.counterValue[i] == 1 {
					// Counter reached one, set output low for one clock pulse
					p.counterOutput[i] = false
				} else if p.counterValue[i] == 0 {
					// Counter reached zero, reload initial count and set output high
					p.counterValue[i] = p.counterRegister[i]
					p.counterOutput[i] = true
				} else {
					// Decrement the counter value
					p.counterValue[i]--
				}
			case 3:
				// Mode 3: Square Wave Mode
				if p.counterValue[i] == 0 {
					// Counter reached zero, reload initial count
					p.counterValue[i] = p.counterRegister[i]
					// Toggle the output state
					p.counterOutput[i] = !p.counterOutput[i]
				} else {
					// Decrement the counter value
					p.counterValue[i]--
				}
			case 4:
				// Mode 4: Software Triggered Mode
				if p.counterValue[i] == 0 {
					// Counter reached zero, set output low for one clock pulse
					p.counterOutput[i] = false
				} else {
					// Decrement the counter value
					p.counterValue[i]--
				}
			case 5:
				// Mode 5: Hardware Triggered Mode
				if p.counterValue[i] == 0 {
					// Counter reached zero, set output low for one clock pulse
					p.counterOutput[i] = false
				} else {
					// Decrement the counter value
					p.counterValue[i]--
				}
			}

			// Check if counter reached zero
			if p.counterValue[i] == 0 {
				// Send interrupt request to the bus
				if p.bus != nil {

					// Reset the counter value
					p.counterValue[i] = p.counterRegister[i]

					// Set the output high
					p.counterOutput[i] = true

					// Set the null count flag
					p.counterNullCount[i] = true

					// Set the status latched flag
					p.statusLatched = true

					// Set the counter latched flag
					p.counterLatched[i] = true

					// Set the counter initialized flag
					p.counterInitialized[i] = false

					// Set the counter mode to mode 0
					p.counterMode[i] = 0

					// Set the counter access mode to mode 3
					p.counterAccessMode[i] = 3

					// Set the control word to 0x36
					p.controlWord = 0x36

					// Set the previous update time to the current time
					p.previousUpdateTime = currentTime

					interruptMessage := bus.BusMessage{
						Subject: common.MESSAGE_INTERRUPT_RAISE,
						Sender:  p.busId,
						Data:    []byte{byte(i)},
					}
					err := p.bus.SendMessageSingle(common.MODULE_INTERRUPT_CONTROLLER_1, interruptMessage)
					if err != nil {
						log.Printf("8259A: Error sending interrupt request message: %v", err)
					}
				}
			}
		}
	}
}
