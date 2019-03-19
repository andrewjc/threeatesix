package intel8259a

import "github.com/andrewjc/threeatesix/devices/bus"

/*
	Simulated 8259A Interrupt Controller Chip

	IRQ0 through IRQ7 are the master 8259's interrupt lines, while IRQ8 through IRQ15 are the slave 8259's interrupt lines.
*/

type Intel8259a struct {
	busId uint32
}

func NewIntel8259a() *Intel8259a {
	chip := &Intel8259a{}

	return chip
}

func (device *Intel8259a) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *Intel8259a) OnReceiveMessage(message bus.BusMessage) {

}