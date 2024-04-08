package intel82C54

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"github.com/andrewjc/threeatesix/devices/intel8259a"

	"log"
	"time"
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
	counterValue       [3]uint16

	previousUpdateTime int64 // used to simulate the 1.19 MHz clock
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

func (p *Intel82C54) GetPortMap() *bus.DevicePortMap {
	return nil
}

func (p *Intel82C54) ReadAddr8(addr uint16) uint8 {
	//TODO implement me
	panic("implement me")
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

	//log.Printf("PIT Command Register Write: Counter=%d, Access Mode=%d, Mode=%d", counterIndex, p.counterAccessMode[counterIndex], p.counterMode[counterIndex])
}

func (p *Intel82C54) CounterRegisterWrite(counterIndex uint8, value uint8) {
	if counterIndex >= 3 {
		log.Printf("Invalid PIT counter index: %d", counterIndex)
		return
	}

	switch p.counterAccessMode[counterIndex] {
	case 0x00: // Counter Latch
		// Latch the current counter value
		p.counterLatch[counterIndex] = p.counterValue[counterIndex]
		log.Printf("PIT Counter %d Latched: Value=%#04x", counterIndex, p.counterLatch[counterIndex])
	case 0x01: // Read/Write LSB
		// Update the lower 8 bits of the counter latch and value with the new value
		lowerBits := uint16(value)
		upperBits := p.counterLatch[counterIndex] & 0xFF00
		p.counterLatch[counterIndex] = upperBits | lowerBits
		p.counterValue[counterIndex] = upperBits | lowerBits
	case 0x02: // Read/Write MSB
		// Update the upper 8 bits of the counter latch and value with the new value
		lowerBits := p.counterLatch[counterIndex] & 0x00FF
		upperBits := uint16(value) << 8
		p.counterLatch[counterIndex] = upperBits | lowerBits
		p.counterValue[counterIndex] = upperBits | lowerBits
	case 0x03: // Read/Write LSB then MSB
		if !p.counterInitialized[counterIndex] {
			// Update the lower 8 bits of the counter latch and value with the new value
			lowerBits := uint16(value)
			upperBits := p.counterLatch[counterIndex] & 0xFF00
			p.counterLatch[counterIndex] = upperBits | lowerBits
			p.counterValue[counterIndex] = upperBits | lowerBits
			p.counterInitialized[counterIndex] = true
		} else {
			// Update the upper 8 bits of the counter latch and value with the new value
			lowerBits := p.counterLatch[counterIndex] & 0x00FF
			upperBits := uint16(value) << 8
			p.counterLatch[counterIndex] = upperBits | lowerBits
			p.counterValue[counterIndex] = upperBits | lowerBits
			p.counterInitialized[counterIndex] = false
		}
	}

	//log.Printf("PIT Counter %d Register Write: Value=%#04x", counterIndex, p.counterLatch[counterIndex])
}

func (p *Intel82C54) CounterRegisterRead(counterIndex uint8) uint8 {
	if counterIndex >= 3 {
		log.Printf("Invalid PIT counter index: %d", counterIndex)
		return 0
	}

	var value uint8

	switch p.counterAccessMode[counterIndex] {
	case 0x00: // Counter Latch
		// Return the latched value of the counter
		value = uint8(p.counterLatch[counterIndex] & 0xFF)
	case 0x01: // Read/Write LSB
		// Return the lower 8 bits of the counter value
		value = uint8(p.counterValue[counterIndex] & 0xFF)
	case 0x02: // Read/Write MSB
		// Return the upper 8 bits of the counter value
		value = uint8((p.counterValue[counterIndex] >> 8) & 0xFF)
	case 0x03: // Read/Write LSB then MSB
		if !p.counterInitialized[counterIndex] {
			// Return the lower 8 bits of the counter value
			value = uint8(p.counterValue[counterIndex] & 0xFF)
			p.counterInitialized[counterIndex] = true
		} else {
			// Return the upper 8 bits of the counter value
			value = uint8((p.counterValue[counterIndex] >> 8) & 0xFF)
			p.counterInitialized[counterIndex] = false
		}
	}

	//log.Printf("PIT Counter %d Register Read: Value=%#02x", counterIndex, value)
	return value
}

func (p *Intel82C54) ReadCounter0() uint8 {
	return p.CounterRegisterRead(0)
}

func (p *Intel82C54) WriteCounter0(value uint8) {
	p.CounterRegisterWrite(0, value)
}

func (p *Intel82C54) GetBus() *bus.Bus {
	return p.bus
}

func (p *Intel82C54) SetBus(bus *bus.Bus) {
	p.bus = bus
}

func (p *Intel82C54) Step() {
	// Step the PIT counters
	// ensure this only runs at 1.19318 MHz
	// 1193180 Hz / 65536 = 18.2065 Hz
	// 18.2065 Hz / 3 = 6.0688 Hz
	// 1 / 6.0688 Hz = 0.1648 seconds
	// 0.1648 seconds * 1000 = 164.8 milliseconds
	// 164.8 milliseconds * 1000 = 164800 microseconds

	// only run this every 164800 microseconds
	previousUpdateTime := p.previousUpdateTime
	currentTime := time.Now().UnixNano() / 1000
	if currentTime-previousUpdateTime < 164800 {
		return
	}

	p.previousUpdateTime = currentTime

	for i := 0; i < 3; i++ {
		if p.counterInitialized[i] {
			if p.counterValue[i] == 0 {
				// Counter has reached zero, generate an interrupt
				interruptController := p.bus.FindSingleDevice(common.MODULE_INTERRUPT_CONTROLLER_1).(*intel8259a.Intel8259a)
				interruptController.TriggerInterrupt(uint8(common.IRQ_PIT_COUNTER0 + i))

				// Reload the counter value
				p.counterValue[i] = p.counterLatch[i]

				log.Printf("PIT Counter %d Interrupt", i)

			} else {
				// Decrement the counter value
				p.counterValue[i]--
			}
		}
	}

}
