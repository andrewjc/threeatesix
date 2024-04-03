package io

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"github.com/andrewjc/threeatesix/devices/cga"
	"github.com/andrewjc/threeatesix/devices/intel82335"
	"github.com/andrewjc/threeatesix/devices/intel8259a"
	"github.com/andrewjc/threeatesix/devices/intel82C54"
	"github.com/andrewjc/threeatesix/devices/ps2"
	"log"
)

/*
	IO Port Access Controller
	Provides read/write functions for port mapped IO
*/

type IOPortAccessController struct {
	bus                              *bus.Bus
	busId                            uint32
	highIntegrationInterfaceDevice   *intel82335.Intel82335
	programmableInterruptController1 *intel8259a.Intel8259a
	programmableInterruptController2 *intel8259a.Intel8259a
	programmableIntervalTimer        *intel82C54.Intel82C54
	cmosRegisterSelect               uint8
	cmosRegisterData                 []uint8

	cgaController *cga.Motorola6845
}

func CreateIOPortController() *IOPortAccessController {
	return &IOPortAccessController{
		cmosRegisterData:                 make([]byte, 0x10000),
		highIntegrationInterfaceDevice:   intel82335.NewIntel82335(),
		programmableInterruptController1: intel8259a.NewIntel8259a(),
		programmableInterruptController2: intel8259a.NewIntel8259a(),
		programmableIntervalTimer:        intel82C54.NewIntel82C54(),
		cgaController:                    cga.NewMotorola6845(),
	}
}

func (mem *IOPortAccessController) SetDeviceBusId(id uint32) {
	mem.busId = id
}

func (mem *IOPortAccessController) OnReceiveMessage(message bus.BusMessage) {

}

func (r *IOPortAccessController) ReadAddr8(addr uint16) uint8 {

	if addr == 0x64 {
		// Status Register READ
		sr := r.GetBus().FindSingleDevice(common.MODULE_PS2_CONTROLLER).(*ps2.Ps2Controller).ReadStatusRegister()
		return sr
	}

	if addr == 0x60 {
		sr := r.GetBus().FindSingleDevice(common.MODULE_PS2_CONTROLLER).(*ps2.Ps2Controller).ReadDataPort()
		return sr
	}

	if addr == 0x0042 {
		// Read from IO port 0x0042
		return r.programmableIntervalTimer.ReadCounter0()
	}

	if addr == 0x24 {
		// RC1 roll compare register???
		//log.Printf("RC1 roll compare register read")
		return r.highIntegrationInterfaceDevice.Rc1RegisterRead()
	}

	if addr == 0x71 {
		// CMOS RAM
		return r.cmosRegisterData[r.cmosRegisterSelect]
	}

	if addr == 0xe3 || addr == 0xe4 {
		return 0
	}

	log.Fatalf("Unhandled IO port read: PORT=[%#04x]", addr)

	return 0
}

func (r *IOPortAccessController) WriteAddr8(port_addr uint16, value uint8) {

	if port_addr == 0x00F1 {
		// 80287 math coprocessor
		r.GetBus().SendMessageSingle(common.MODULE_MATH_CO_PROCESSOR, bus.BusMessage{common.MESSAGE_REQUEST_CPU_MODESWITCH, []byte{common.REAL_MODE}})
		return
	}

	if port_addr == 0x64 {
		// Command Register Write
		r.GetBus().FindSingleDevice(common.MODULE_PS2_CONTROLLER).(*ps2.Ps2Controller).WriteCommandRegister(value)
		return
	}

	if port_addr == 0x60 {
		// Data Port Write
		r.GetBus().FindSingleDevice(common.MODULE_PS2_CONTROLLER).(*ps2.Ps2Controller).WriteDataPort(value)
		return
	}

	if port_addr == 0x61 {
		// Command Port Write
		r.GetBus().FindSingleDevice(common.MODULE_PS2_CONTROLLER).(*ps2.Ps2Controller).WriteControlPort(value)
		return
	}

	if port_addr == 0x80 {
		// bios post diag
		log.Printf("BIOS POST: %#02x - %s", value, common.BiosPostCodeToString(value))
		return
	}

	if port_addr == 0x22 {
		// MCR register setup
		r.highIntegrationInterfaceDevice.McrRegisterInitialize(value)
		return
	}

	if port_addr == 0x24 {
		// RC1 roll compare register???
		//log.Printf("RC1 roll compare register write")
		r.highIntegrationInterfaceDevice.Rc1RegisterWrite(value)
		return
	}

	if port_addr == 0x70 {
		// CMOS RAM
		r.cmosRegisterSelect = value
		return
	}

	if port_addr == 0x71 {
		// CMOS RAM
		log.Printf("CMOS RAM WRITE: %#02x, %#02x", r.cmosRegisterSelect, value)
		r.cmosRegisterData[r.cmosRegisterSelect] = value
		return
	}

	if port_addr == 0x92 {
		// A20 Gate
		// log.Printf("A20 GATE: %#02x", value)
		if value == 0x00 {
			r.GetBus().SendMessageSingle(common.MODULE_MEMORY_ACCESS_CONTROLLER, bus.BusMessage{common.MESSAGE_DISABLE_A20_GATE, []byte{value}})
		} else {
			r.GetBus().SendMessageSingle(common.MODULE_MEMORY_ACCESS_CONTROLLER, bus.BusMessage{common.MESSAGE_ENABLE_A20_GATE, []byte{value}})
		}

		return
	}

	if port_addr == 0x8 {
		// DMA command register write
		log.Printf("DMA COMMAND REGISTER WRITE: %#02x", value)
		r.highIntegrationInterfaceDevice.DmaCommandRegisterWrite(value)
		return
	}

	if port_addr == 0xd0 {
		log.Printf("Interrupt Request Level Priority Controller Configuration: %#02x", value)
		r.programmableInterruptController2.SetInterruptRequest(value)
		return
	}

	if port_addr == 0x20 || port_addr == 0x21 {
		r.programmableInterruptController1.CommandWordWrite(value)
		return
	}

	if port_addr == 0xA0 || port_addr == 0xA1 {
		r.programmableInterruptController2.CommandWordWrite(value)
		return
	}

	if port_addr == 0x03d8 {
		// CGA
		r.cgaController.WriteAddr8(port_addr, value)
		return
	}

	if port_addr == 0x0042 {
		// Write to IO port 0x0042
		r.programmableIntervalTimer.WriteCounter0(value)
		return
	}

	if port_addr == 0x0043 {
		// PIT
		r.programmableIntervalTimer.CommandRegisterWrite(value)
		return
	}

	log.Fatalf("Unhandled IO port write: PORT=[%#04x], value=%#02x", port_addr, value)
}

func (r *IOPortAccessController) ReadAddr16(addr uint16) uint16 {
	b1 := uint16(r.ReadAddr8(addr))
	b2 := uint16(r.ReadAddr8(addr + 1))
	return b2<<8 | b1
}

func (r *IOPortAccessController) WriteAddr16(addr uint16, value uint16) {

}

func (controller *IOPortAccessController) GetBus() *bus.Bus {
	return controller.bus
}

func (controller *IOPortAccessController) SetBus(bus *bus.Bus) {
	controller.bus = bus
	controller.highIntegrationInterfaceDevice.SetBus(bus)
}
