package ps2

import (
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type Ps2Controller struct {
	bus            *bus.Bus
	busId          uint32
	statusRegister uint8
}

func (controller *Ps2Controller) SetDeviceBusId(id uint32) {
	controller.busId = id
}

func (controller *Ps2Controller) OnReceiveMessage(message bus.BusMessage) {

}

func CreatePS2Controller() *Ps2Controller {
	return &Ps2Controller{}
}

func (controller *Ps2Controller) GetBus() *bus.Bus {
	return controller.bus
}

func (controller *Ps2Controller) SetBus(bus *bus.Bus) {
	controller.bus = bus
}

func (controller *Ps2Controller) ReadStatusRegister() uint8 {
	return controller.statusRegister
}

func (controller *Ps2Controller) WriteCommandRegister(value uint8) {
	log.Printf("PS2 controller write command: [%#04x]", value)
}