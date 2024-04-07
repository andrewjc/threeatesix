package intel8259a

import (
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
	"os"
)

/*
   Simulated 8259A Interrupt Controller Chip
*/

type Intel8259a struct {
	busId               uint32
	irqMask             uint8 // Interrupt request mask
	irqService          uint8 // In-service interrupts
	requestIrr          uint8 // Interrupt request register
	serviceIrr          uint8 // In-service interrupt register
	interruptVectorBase uint8 // Interrupt vector base
	interruptVector     uint8 // Interrupt vector
	isrPointer          uint8 // Pointer to the current position in ISR
	irrPointer          uint8 // Pointer to the current position in IRR
	priorityMode        uint8 // Interrupt priority mode (0: fixed, 1: rotating)
	autoEoi             bool  // Automatic End-of-Interrupt mode
	specificEoi         bool  // Specific End-of-Interrupt mode
	isSlavePIC          bool  // Is this a slave PIC
	slaveIRQ            uint8 // Slave PIC's interrupt input
	operatingMode       bool  // Operating mode (0: MCS-80/85, 1: x86)
	bufferedMode        bool  // Buffered mode
}

func NewIntel8259a() *Intel8259a {
	chip := &Intel8259a{}
	chip.irqMask = 0xff         // All interrupts masked initially
	chip.interruptVector = 0x08 // Default interrupt vector base
	return chip
}

func (device *Intel8259a) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *Intel8259a) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (device *Intel8259a) SetInterruptRequest(irq uint8) {
	// Set the specified interrupt request bit in the request IRR
	interruptRequestBit := uint8(1 << irq)
	device.requestIrr |= interruptRequestBit
	log.Printf("8259A: Interrupt request set for IRQ %d", irq)
	device.updateInterrupts()
}

func (device *Intel8259a) clearInterruptRequest(irq uint8) {
	// Clear the specified interrupt request bit in the request IRR
	interruptRequestBit := uint8(1 << irq)
	device.requestIrr &= ^interruptRequestBit
	log.Printf("8259A: Interrupt request cleared for IRQ %d", irq)
	device.updateInterrupts()
}

func (device *Intel8259a) setServiceRequest(irq uint8) {
	// Set the specified interrupt service bit in the IRQ service register
	interruptServiceBit := uint8(1 << irq)
	device.irqService |= interruptServiceBit
	log.Printf("8259A: Interrupt service set for IRQ %d", irq)
}

func (device *Intel8259a) clearServiceRequest(irq uint8) {
	// Clear the specified interrupt service bit in the IRQ service register
	interruptServiceBit := uint8(1 << irq)
	device.irqService &= ^interruptServiceBit
	log.Printf("8259A: Interrupt service cleared for IRQ %d", irq)
}

func (device *Intel8259a) updateInterrupts() {
	// Check for pending interrupts and trigger the highest priority one
	pendingInterrupts := device.requestIrr & ^device.irqMask & ^device.irqService
	if pendingInterrupts != 0 {
		highestPriorityIrq := device.getHighestPriorityIrq(pendingInterrupts)
		device.triggerInterrupt(highestPriorityIrq)
	}
}

func (device *Intel8259a) getHighestPriorityIrq(interrupts uint8) uint8 {
	// Find the highest priority interrupt among the pending ones
	for irq := uint8(0); irq < 8; irq++ {
		if (interrupts & (1 << irq)) != 0 {
			return irq
		}
	}
	return 0xff // No pending interrupts
}

func (device *Intel8259a) triggerInterrupt(irq uint8) {
	// Trigger the specified interrupt
	device.setServiceRequest(irq)
	device.clearInterruptRequest(irq)
	interruptVector := device.interruptVectorBase + irq

	// If this is a slave PIC, send the interrupt to the master PIC
	if device.isSlavePIC {
		device.setServiceRequest(device.slaveIRQ)
		device.clearInterruptRequest(device.slaveIRQ)
		interruptVector = device.interruptVectorBase + device.slaveIRQ
	} else {
		// If this is a master PIC, check if there is a slave PIC
		if device.slaveIRQ != 0xff {
			// Send the interrupt to the slave PIC
			device.setServiceRequest(device.slaveIRQ)
			device.clearInterruptRequest(device.slaveIRQ)
			interruptVector = device.interruptVectorBase + device.slaveIRQ
		} else {
			// No slave PIC, so just trigger the interrupt
			interruptVector = device.interruptVectorBase + irq
		}
	}

	log.Printf("8259A: Triggering interrupt IRQ %d with vector 0x%02X", irq, interruptVector)

}

func (device *Intel8259a) EndOfInterrupt(irq uint8) {
	// End-of-Interrupt handler
	if device.autoEoi {
		// Automatic EOI mode
		device.clearServiceRequest(device.isrPointer)
		device.isrPointer = (device.isrPointer + 1) % 8
	} else if device.specificEoi {
		// Specific EOI mode
		device.clearServiceRequest(irq)
	} else {
		// Normal EOI mode
		device.clearServiceRequest(device.isrPointer)
	}
	log.Printf("8259A: End-of-Interrupt for IRQ %d", irq)
	device.updateInterrupts()
}

// This would handle command words from the CPU.
func (device *Intel8259a) CommandWordWrite(value uint8) {
	// Interpret value according to the PIC's command word structure
	if (value & 0x10) != 0 {
		// ICW1: Initialization Command Word 1
		device.isrPointer = 0
		device.irrPointer = 0
		device.irqMask = 0xff
		device.requestIrr = 0
		device.irqService = 0
		device.priorityMode = 0
		device.autoEoi = false
		device.specificEoi = false
		log.Printf("8259A: Initialization Command Word 1 received")
	} else if (value & 0x08) != 0 {
		// OCW3: Operation Control Word 3
		if (value & 0x02) != 0 {
			// Set priority mode
			device.priorityMode = (value >> 5) & 0x01
			log.Printf("8259A: Priority mode set to %d", device.priorityMode)
		}
		if (value & 0x01) != 0 {
			// Set EOI mode
			device.autoEoi = (value & 0x02) != 0
			device.specificEoi = (value & 0x40) != 0
			log.Printf("8259A: EOI mode set to Auto: %t, Specific: %t", device.autoEoi, device.specificEoi)
		}
	} else if (value & 0x04) != 0 {
		// ICW3: Initialization Command Word 3
		log.Printf("8259A: Initialization Command Word 3 received")

		// Set the master/slave relationship and the slave PIC's interrupt input
		device.isSlavePIC = (value & 0x08) != 0
		device.slaveIRQ = value & 0x07

		log.Printf("8259A: Slave PIC configured, isSlavePIC: %t, slaveIRQ: %d", device.isSlavePIC, device.slaveIRQ)

		// Enable all interrupts
		device.irqMask = 0x00
		log.Printf("8259A: Interrupt mask set to 0x%02X", device.irqMask)

		device.updateInterrupts()

	} else if (value & 0x02) != 0 {
		// ICW2: Initialization Command Word 2
		device.interruptVectorBase = value & 0xF8
		log.Printf("8259A: Initialization Command Word 2 received with vector base 0x%02X", device.interruptVectorBase)
	} else if (value & 0x01) != 0 {
		// ICW4: Initialization Command Word 4
		device.autoEoi = (value & 0x02) != 0
		device.specificEoi = (value & 0x01) != 0
		device.operatingMode = (value & 0x01) != 0 // 0: MCS-80/85 mode, 1: x86 mode
		device.bufferedMode = (value & 0x08) != 0

		log.Printf("8259A: Initialization Command Word 4 received with EOI mode Auto: %t, Specific: %t", device.autoEoi, device.specificEoi)
		log.Printf("8259A: Operating mode set to %s, Buffered mode set to %t", device.operatingMode, device.bufferedMode)

		// Enable all interrupts
		device.irqMask = 0x00
		log.Printf("8259A: Interrupt mask set to 0x%02X", device.irqMask)
		device.updateInterrupts()

	} else if value == 0 {
		// NOP: No Operation
		log.Printf("8259A: No Operation command received")

	} else {
		log.Printf("8259A: Unhandled command word 0x%02X", value)
		os.Exit(0)
	}
}

// This would handle data writes to the PIC's data port.
func (device *Intel8259a) dataWrite(value uint8) {
	// Interpret value according to the PIC's data word structure
	device.irqMask = value
	log.Printf("8259A: Interrupt mask set to 0x%02X", device.irqMask)
	device.updateInterrupts()
}

func (device *Intel8259a) DataWrite(value uint8) {
	// Write data to the PIC's data port
	if (device.isrPointer & 0x01) == 0 {
		// Write to the command register
		device.CommandWordWrite(value)
	} else {
		// Write to the data register
		device.dataWrite(value)
	}
}

func (device *Intel8259a) HasPendingInterrupts() bool {
	// Check if there are any pending interrupts
	return device.requestIrr&^device.irqMask&^device.irqService != 0
}

func (device *Intel8259a) TriggerInterrupt(u uint8) {
	// Trigger the specified interrupt
	device.triggerInterrupt(u)
}
