package intel8259a

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
	"os"
)

/*
   Simulated 8259A Interrupt Controller Chip
*/

type Intel8259a struct {
	bus                 *bus.Bus
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
	debugMode           bool  // Debug mode
}

func NewIntel8259a() *Intel8259a {
	chip := &Intel8259a{}
	chip.irqMask = 0xff         // All interrupts masked initially
	chip.interruptVector = 0x08 // Default interrupt vector base
	chip.debugMode = true
	return chip
}

func (device *Intel8259a) GetDeviceBusId() uint32 {
	return device.busId
}

func (device *Intel8259a) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *Intel8259a) SetBus(bus *bus.Bus) {
	device.bus = bus
}

func (device *Intel8259a) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (device *Intel8259a) SetInterruptRequest(irq uint8) {
	// Set the specified interrupt request bit in the request IRR
	interruptRequestBit := uint8(1 << irq)
	device.requestIrr |= interruptRequestBit

	device.updateInterrupts()
	device.log("8259A: Interrupt request set for IRQ %d", irq)
}

func (device *Intel8259a) clearInterruptRequest(irq uint8) {
	// Clear the specified interrupt request bit in the request IRR
	interruptRequestBit := uint8(1 << irq)
	device.requestIrr &= ^interruptRequestBit
	device.log("8259A: Interrupt request cleared for IRQ %d", irq)
	device.updateInterrupts()
}

func (device *Intel8259a) setServiceRequest(irq uint8) {
	// Set the specified interrupt service bit in the IRQ service register
	interruptServiceBit := uint8(1 << irq)
	device.irqService |= interruptServiceBit
	device.log("8259A: Interrupt service set for IRQ %d", irq)
}

func (device *Intel8259a) clearServiceRequest(irq uint8) {
	// Clear the specified interrupt service bit in the IRQ service register
	interruptServiceBit := uint8(1 << irq)
	device.irqService &= ^interruptServiceBit
	device.log("8259A: Interrupt service cleared for IRQ %d", irq)
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
	// Calculate the interrupt vector based on the vector base and the IRQ
	device.interruptVector = device.interruptVectorBase + irq

	// Set the service request for the IRQ
	device.setServiceRequest(irq)

	// Clear the interrupt request for the IRQ
	device.clearInterruptRequest(irq)

	// Trigger the interrupt on the bus
	if device.bus != nil {
		interruptMessage := bus.BusMessage{
			Subject: common.MESSAGE_INTERRUPT_REQUEST,
			Sender:  device.busId,
			Data:    []byte{device.interruptVector},
		}
		device.bus.SendMessage(interruptMessage)
	}

	device.log("8259A: Triggering interrupt IRQ %d with vector 0x%02X", irq, device.interruptVector)
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
	device.updateInterrupts()
	device.log("8259A: End-of-Interrupt for IRQ %d", irq)
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
		device.log("8259A: Initialization Command Word 1 received")
	} else if (value & 0x08) != 0 {
		// OCW3: Operation Control Word 3
		if (value & 0x02) != 0 {
			// Set priority mode
			device.priorityMode = (value >> 5) & 0x01
			device.log("8259A: Priority mode set to %d", device.priorityMode)
		}
		if (value & 0x01) != 0 {
			// Set EOI mode
			device.autoEoi = (value & 0x02) != 0
			device.specificEoi = (value & 0x40) != 0
			device.log("8259A: EOI mode set to Auto: %t, Specific: %t", device.autoEoi, device.specificEoi)
		}
	} else if (value & 0x04) != 0 {
		// ICW3: Initialization Command Word 3
		device.log("8259A: Initialization Command Word 3 received")

		// Set the master/slave relationship and the slave PIC's interrupt input
		device.isSlavePIC = (value & 0x08) != 0
		device.slaveIRQ = value & 0x07

		device.log("8259A: Slave PIC configured, isSlavePIC: %t, slaveIRQ: %d", device.isSlavePIC, device.slaveIRQ)

		// Enable all interrupts
		device.irqMask = 0x00
		device.log("8259A: Interrupt mask set to 0x%02X", device.irqMask)

		device.updateInterrupts()

	} else if (value & 0x02) != 0 {
		// ICW2: Initialization Command Word 2
		device.interruptVectorBase = value & 0xF8
		device.log("8259A: Initialization Command Word 2 received with vector base 0x%02X", device.interruptVectorBase)
	} else if (value & 0x01) != 0 {
		// ICW4: Initialization Command Word 4
		device.autoEoi = (value & 0x02) != 0
		device.specificEoi = (value & 0x01) != 0
		device.operatingMode = (value & 0x01) != 0 // 0: MCS-80/85 mode, 1: x86 mode
		device.bufferedMode = (value & 0x08) != 0

		device.log("8259A: Initialization Command Word 4 received with EOI mode Auto: %t, Specific: %t", device.autoEoi, device.specificEoi)
		device.log("8259A: Operating mode set to %s, Buffered mode set to %t", device.operatingMode, device.bufferedMode)

		// Enable all interrupts
		device.irqMask = 0x00
		device.log("8259A: Interrupt mask set to 0x%02X", device.irqMask)
		device.updateInterrupts()

	} else if (value == 0x20) || (value == 0x60) {
		// Non-specific EOI command
		device.EndOfInterrupt(0)
	} else if value == 0x60 {
		// Specific EOI command
		device.EndOfInterrupt(0)
	} else if value == 0x0A {
		// Rotate in Automatic EOI mode
		device.autoEoi = true
		device.log("8259A: Automatic EOI mode set")
	} else if value == 0x0B {
		// Rotate in Automatic EOI mode with a special case for EOI
		device.autoEoi = true
		device.specificEoi = true
		device.log("8259A: Automatic EOI mode with specific EOI set")
	} else if value == 0x0C {
		// Set Interrupt Mask
		device.irqMask = 0x01
		device.updateInterrupts()
		device.log("8259A: Interrupt mask set to 0x%02X", device.irqMask)
	} else if value == 0x0D {
		// Read Interrupt Request Register
		device.log("8259A: Interrupt Request Register read")
	} else if value == 0x0E {
		// Read In-Service Register
		device.log("8259A: In-Service Register read")
	} else if value == 0x0F {
		// Read Interrupt Mask Register
		device.log("8259A: Interrupt Mask Register read")
	} else if value == 0x70 {
		// Read Initialization Command Word 1
		device.log("8259A: Initialization Command Word 1 read")
	} else if value == 0x71 {
		// Read Initialization Command Word 2
		device.log("8259A: Initialization Command Word 2 read")
	} else if value == 0x72 {
		// Read Initialization Command Word 3
		device.log("8259A: Initialization Command Word 3 read")
	} else if value == 0x73 {
		// Read Initialization Command Word 4
		device.log("8259A: Initialization Command Word 4 read")
	} else if value == 0x74 {
		// Read Initialization Command Word 4
		device.log("8259A: Read Interrupt Mask Register")
	} else if value == 0 {
		// NOP: No Operation
		device.log("8259A: No Operation command received")
	} else {
		log.Printf("8259A: Unhandled command word 0x%02X", value)
		os.Exit(0)
	}
}

// This would handle data writes to the PIC's data port.
func (device *Intel8259a) dataWrite(value uint8) {
	// Interpret value according to the PIC's data word structure
	device.irqMask = value
	device.updateInterrupts()
	device.log("8259A: Interrupt mask set to 0x%02X", device.irqMask)
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

func (device *Intel8259a) log(format string, args ...interface{}) {
	if device.debugMode {
		log.Printf(format, args...)
	}
}
