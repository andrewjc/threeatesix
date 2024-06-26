package intel8259a

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type Intel8259a struct {
	bus                 *bus.Bus
	busId               uint32
	irqMask             uint8
	irqRequest          uint8
	inService           uint8
	interruptVectorBase uint8
	autoEOI             bool
	mode8086            bool
	slaveMode           bool
	masterMode          bool
	slaveID             uint8
	interruptOutput     bool
	readISR             bool
	readIRR             bool
}

func NewIntel8259a() *Intel8259a {
	return &Intel8259a{
		irqMask: 0xff,
	}
}

func (d *Intel8259a) GetDeviceBusId() uint32 {
	return d.busId
}

func (d *Intel8259a) SetDeviceBusId(id uint32) {
	d.busId = id
}

func (d *Intel8259a) SetBus(bus *bus.Bus) {
	d.bus = bus
}

func (d *Intel8259a) OnReceiveMessage(message bus.BusMessage) {
	if message.Subject == common.MESSAGE_INTERRUPT_RAISE {
		d.assertInterrupt(message.Data[0])
	} else if message.Subject == common.MESSAGE_INTERRUPT_ACKNOWLEDGE {
		d.acknowledgeInterrupt()
	} else if message.Subject == common.MESSAGE_INTERRUPT_COMPLETE {
		d.completeInterrupt(message.Data[0])
	}
}

func (d *Intel8259a) GetPortMap() *bus.DevicePortMap {
	if d.masterMode {
		return &bus.DevicePortMap{
			ReadPorts:  []uint16{0x20, 0x21},
			WritePorts: []uint16{0x20, 0x21},
		}
	} else if d.slaveMode {
		return &bus.DevicePortMap{
			ReadPorts:  []uint16{0xA0, 0xA1},
			WritePorts: []uint16{0xA0, 0xA1},
		}
	}
	return nil
}

func (d *Intel8259a) ReadAddr8(addr uint16) uint8 {
	switch addr {
	case 0x20, 0xA0:
		if d.readISR {
			d.readISR = false
			return d.inService
		} else if d.readIRR {
			d.readIRR = false
			return d.irqRequest
		} else {
			return d.irqRequest
		}
	case 0x21, 0xA1:
		return d.irqMask
	default:
		log.Printf("8259A: Unsupported read from address 0x%04X", addr)
		return 0
	}
}

func (d *Intel8259a) WriteAddr8(addr uint16, data uint8) {
	switch addr {
	case 0x20, 0xA0:
		if data&0x10 != 0 {
			d.initialize(data)
		} else if data&0x08 != 0 {
			d.operationCommand3(data)
		} else if data&0x04 != 0 {
			d.operationCommand2(data)
		} else {
			d.interruptVectorBase = data & 0xF8
		}
	case 0x21, 0xA1:
		if d.slaveMode && data&0x04 != 0 {
			d.slaveID = data & 0x07
		} else {
			d.irqMask = data
		}
	default:
		log.Printf("8259A: Unsupported write to address 0x%04X with data 0x%02X", addr, data)
	}
}

func (d *Intel8259a) initialize(data uint8) {
	d.irqMask = 0xFF
	d.irqRequest = 0
	d.inService = 0
	d.interruptVectorBase = data & 0xF8
	d.autoEOI = data&0x02 != 0
	d.mode8086 = data&0x01 != 0
	d.slaveMode = data&0x08 == 0
	d.masterMode = !d.slaveMode
	d.readISR = false
	d.readIRR = false
	d.interruptOutput = false
}

func (d *Intel8259a) operationCommand2(data uint8) {
	rotate := data&0x80 != 0
	autoEOI := data&0x40 != 0
	specificEOI := data&0x20 != 0
	nonSpecificEOI := data&0x10 != 0
	rotateIRQ := data & 0x07

	if nonSpecificEOI {
		d.completeInterrupt(d.findHighestPriorityIRQ())
	} else if specificEOI {
		d.completeInterrupt(rotateIRQ)
	}

	if autoEOI {
		d.autoEOI = true
	} else if rotate {
		d.rotateInterrupt(rotateIRQ)
	}
}

func (d *Intel8259a) operationCommand3(data uint8) {
	if data&0x01 != 0 {
		d.readISR = true
	} else if data&0x02 != 0 {
		d.readIRR = true
	}
}

func (d *Intel8259a) assertInterrupt(irq uint8) {
	d.irqRequest |= 1 << irq
	d.bus.SendMessageSingle(common.MODULE_PRIMARY_PROCESSOR, bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_RAISE,
		Data:    []byte{d.irqRequest},
		Sender:  d.busId,
	})
}

func (d *Intel8259a) acknowledgeInterrupt() {
	irq := d.findHighestPriorityIRQ()
	if irq != 0xFF {
		d.inService |= 1 << irq
		d.irqRequest &^= 1 << irq
		interruptVector := d.interruptVectorBase + irq
		d.sendInterrupt(interruptVector)
	}
}

func (d *Intel8259a) completeInterrupt(irq uint8) {
	d.inService &^= 1 << irq
	d.updateInterruptOutput()
}

func (d *Intel8259a) updateInterruptOutput() {
	if d.irqRequest&^d.irqMask != 0 {
		d.interruptOutput = true
		d.bus.SendMessageSingle(common.MODULE_PRIMARY_PROCESSOR, bus.BusMessage{
			Subject: common.MESSAGE_INTERRUPT_RAISE,
			Data:    []byte{0},
			Sender:  d.busId,
		})
	} else {
		d.interruptOutput = false
	}
}

func (d *Intel8259a) findHighestPriorityIRQ() uint8 {
	for irq := uint8(0); irq < 8; irq++ {
		if d.irqRequest&(1<<irq) != 0 && d.inService&(1<<irq) == 0 {
			return irq
		}
	}
	return 0xFF
}

func (d *Intel8259a) rotateInterrupt(irq uint8) {
	d.irqRequest = (d.irqRequest << irq) | (d.irqRequest >> (8 - irq))
	d.inService = (d.inService << irq) | (d.inService >> (8 - irq))
}

func (d *Intel8259a) sendInterrupt(interruptVector uint8) {
	d.bus.SendMessageSingle(common.MODULE_PRIMARY_PROCESSOR, bus.BusMessage{
		Subject: common.MESSAGE_INTERRUPT_EXECUTE,
		Data:    []byte{interruptVector},
		Sender:  d.busId,
	})
}

func (d *Intel8259a) IsPrimaryDevice(b bool) {
	d.masterMode = b
}

func (d *Intel8259a) IsSecondaryDevice(b bool) {
	d.slaveMode = b
}
