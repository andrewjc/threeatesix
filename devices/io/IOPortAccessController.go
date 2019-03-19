package io

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

/*
	IO Port Access Controller
	Provides read/write functions for port mapped IO
*/

type IOPortAccessController struct {
	backingMemory         []byte
	bus                   *bus.Bus
	busId                 uint32
}


func (mem *IOPortAccessController) SetDeviceBusId(id uint32) {
	mem.busId = id
}

func (mem *IOPortAccessController) OnReceiveMessage(message bus.BusMessage) {

}

func (r *IOPortAccessController) ReadAddr8(addr uint16) uint8 {
	var byteData uint8
	byteData = (r.backingMemory)[addr]

	return byteData
}

func (r *IOPortAccessController) WriteAddr8(addr uint16, value uint8) {

	if addr == 0x00F1 {
		// 80287 math coprocessor
		r.GetBus().SendMessageSingle(common.MODULE_MATH_CO_PROCESSOR, bus.BusMessage{common.MESSAGE_REQUEST_CPU_MODESWITCH, []byte{common.REAL_MODE}})
		//r.coProcessorController.EnterMode(REAL_MODE)
		return
	}

	/*if addr == 0x0A0H {
		// Interrupt controller 1
		r.cpuController.
	}

	if addr == 0x0A1H {
		// Interrupt controller 2
	}*/

	r.backingMemory[addr] = value
}

func (r *IOPortAccessController) ReadAddr16(addr uint16) uint16 {
	b1 := uint16(r.ReadAddr8(addr))
	b2 := uint16(r.ReadAddr8(addr + 1))
	return b2<<8 | b1
}

func (r *IOPortAccessController) WriteAddr16(addr uint16, value uint16) {
	log.Fatal("TODO!")
}

func (controller *IOPortAccessController) GetBus() *bus.Bus {
	return controller.bus
}

func (controller *IOPortAccessController) SetBus(bus *bus.Bus) {
	controller.bus = bus
}

func CreateIOPortController() *IOPortAccessController {
	return &IOPortAccessController{backingMemory: make([]byte, 0x10000)}
}
