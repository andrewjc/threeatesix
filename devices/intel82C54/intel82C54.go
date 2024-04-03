package intel82C54

import (
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type Intel82C54 struct {
	bus                *bus.Bus
	busId              uint32
	commandRegister    uint8
	counterRegister    [3]uint8
	counterLatch       [3]uint16
	counterMode        [3]uint8
	counterAccessMode  [3]uint8
	counterInitialized [3]bool
}

func NewIntel82C54() *Intel82C54 {
	return &Intel82C54{}
}

func (p *Intel82C54) SetDeviceBusId(id uint32) {
	p.busId = id
}

func (p *Intel82C54) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (p *Intel82C54) CommandRegisterWrite(value uint8) {
	p.commandRegister = value

	// Extract the counter index from bits 6-7 of the command register
	counterIndex := (value >> 6) & 0x03

	// Extract the counter access mode from bits 4-5 of the command register
	p.counterAccessMode[counterIndex] = (value >> 4) & 0x03

	// Extract the counter mode from bits 1-3 of the command register
	p.counterMode[counterIndex] = (value >> 1) & 0x07

	// Set the counter as initialized
	p.counterInitialized[counterIndex] = true

	log.Printf("PIT Command Register Write: Counter=%d, Access Mode=%d, Mode=%d", counterIndex, p.counterAccessMode[counterIndex], p.counterMode[counterIndex])
}

func (p *Intel82C54) CounterRegisterWrite(counterIndex uint8, value uint8) {
	if counterIndex >= 3 {
		log.Printf("Invalid PIT counter index: %d", counterIndex)
		return
	}

	switch p.counterAccessMode[counterIndex] {
	case 0x00: // Counter Latch
		// Latch not supported yet
		log.Printf("PIT Counter Latch not supported")
	case 0x01: // Read/Write LSB
		// Update the lower 8 bits of the counter latch with the new value
		lowerBits := uint16(value)
		upperBits := p.counterLatch[counterIndex] & 0xFF00
		p.counterLatch[counterIndex] = upperBits | lowerBits
	case 0x02: // Read/Write MSB
		// Update the upper 8 bits of the counter latch with the new value
		lowerBits := p.counterLatch[counterIndex] & 0x00FF
		upperBits := uint16(value) << 8
		p.counterLatch[counterIndex] = upperBits | lowerBits
	case 0x03: // Read/Write LSB then MSB
		if !p.counterInitialized[counterIndex] {
			// Update the lower 8 bits of the counter latch with the new value
			lowerBits := uint16(value)
			upperBits := p.counterLatch[counterIndex] & 0xFF00
			p.counterLatch[counterIndex] = upperBits | lowerBits
			p.counterInitialized[counterIndex] = true
		} else {
			// Update the upper 8 bits of the counter latch with the new value
			lowerBits := p.counterLatch[counterIndex] & 0x00FF
			upperBits := uint16(value) << 8
			p.counterLatch[counterIndex] = upperBits | lowerBits
			p.counterInitialized[counterIndex] = false
		}
	}

	log.Printf("PIT Counter %d Register Write: Value=%#04x", counterIndex, p.counterLatch[counterIndex])
}

func (p *Intel82C54) GetBus() *bus.Bus {
	return p.bus
}

func (p *Intel82C54) SetBus(bus *bus.Bus) {
	p.bus = bus
}
