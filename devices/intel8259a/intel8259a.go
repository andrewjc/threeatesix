package intel8259a

import (
    "github.com/andrewjc/threeatesix/devices/bus"
    "log"
)

/*
   Simulated 8259A Interrupt Controller Chip
*/

type Intel8259a struct {
    busId           uint32
    irqMask         uint8 // Interrupt request mask
    irqService      uint8 // In-service interrupts
    requestIrr      uint8 // Interrupt request register
    serviceIrr      uint8 // In-service interrupt register
    interruptVector uint8 // Base interrupt vector
    isrPointer      uint8 // Pointer to the current position in ISR
    irrPointer      uint8 // Pointer to the current position in IRR
    priorityMode    uint8 // Interrupt priority mode (0: fixed, 1: rotating)
    autoEoi         bool  // Automatic End-of-Interrupt mode
    specificEoi     bool  // Specific End-of-Interrupt mode
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
    interruptVector := device.interruptVector + irq
    log.Printf("8259A: Triggering interrupt IRQ %d with vector 0x%02X", irq, interruptVector)
    // TODO: Send interrupt signal to the CPU with the appropriate interrupt vector
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
        // TODO: Handle initialization sequence
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
    } else {
        // OCW2: Operation Control Word 2
        if (value & 0xe0) == 0x00 {
            // Rotate in automatic EOI mode clear
            log.Printf("8259A: Rotate in automatic EOI mode clear")
        } else if (value & 0xe0) == 0x20 {
            // Non-specific EOI command
            device.EndOfInterrupt(0)
        } else if (value & 0xe0) == 0x40 {
            // No operation
            log.Printf("8259A: No operation command")
        } else if (value & 0xe0) == 0x60 {
            // Specific EOI command
            irq := value & 0x07
            device.EndOfInterrupt(irq)
        } else if (value & 0xe0) == 0x80 {
            // Rotate in automatic EOI mode set
            log.Printf("8259A: Rotate in automatic EOI mode set")
        } else if (value & 0xe0) == 0xa0 {
            // Rotate on non-specific EOI command
            log.Printf("8259A: Rotate on non-specific EOI command")
        } else if (value & 0xe0) == 0xc0 {
            // Set priority command
            log.Printf("8259A: Set priority command")
        } else if (value & 0xe0) == 0xe0 {
            // Rotate on specific EOI command
            log.Printf("8259A: Rotate on specific EOI command")
        }
    }
}

// This would handle data writes to the PIC's data port.
func (device *Intel8259a) dataWrite(value uint8) {
    // Interpret value according to the PIC's data word structure
    device.irqMask = value
    log.Printf("8259A: Interrupt mask set to 0x%02X", device.irqMask)
    device.updateInterrupts()
}
