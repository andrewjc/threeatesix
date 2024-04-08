package io

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"github.com/andrewjc/threeatesix/devices/cga"
	"github.com/andrewjc/threeatesix/devices/intel82335"
	"github.com/andrewjc/threeatesix/devices/intel8237"
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
	bus                *bus.Bus
	busId              uint32
	cmosRegisterSelect uint8
	cmosRegisterData   []uint8
}

func CreateIOPortController() *IOPortAccessController {
	return &IOPortAccessController{
		cmosRegisterData: make([]byte, 0x10000),
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
		return r.GetBus().FindSingleDevice(common.MODULE_PIT).(*intel82C54.Intel82C54).ReadCounter0()
	}

	if addr == 0x24 {
		// RC1 roll compare register???
		//log.Printf("RC1 roll compare register read")
		sr := r.GetBus().FindSingleDevice(common.MODULE_INTEL_82335).(*intel82335.Intel82335).Rc1RegisterRead()
		return sr
	}

	if addr == 0x71 {
		// CMOS RAM
		return r.cmosRegisterData[r.cmosRegisterSelect]
	}

	if addr == 0xe3 {
		// 8237 DMA controller
		return r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).ReadTemporaryRegister()
	}
	if addr == 0xe4 {
		// 8237 DMA controller
		return r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).ReadStatusRegister()
	}

	if addr == 0x00D3 {
		// 8237 DMA controller
		return r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).ReadTemporaryRegister()
	}
	if addr == 0x00D0 {
		// 8237 DMA controller
		return r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).ReadStatusRegister()
	}

	log.Fatalf("Unhandled IO port read: PORT=[%#04x]", addr)

	return 0
}

func (r *IOPortAccessController) WriteAddr8(port_addr uint16, value uint8) {

	if port_addr == 0x00F1 {
		// 80287 math coprocessor
		err := r.GetBus().SendMessageSingle(common.MODULE_MATH_CO_PROCESSOR, bus.BusMessage{Subject: common.MESSAGE_REQUEST_CPU_MODESWITCH, Data: []byte{common.REAL_MODE}})
		if err != nil {
			log.Fatalf("Failed to send message to math coprocessor: %s", err)
			return
		}
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
		r.GetBus().FindSingleDevice(common.MODULE_INTEL_82335).(*intel82335.Intel82335).McrRegisterInitialize(value)
		return
	}

	if port_addr == 0x24 {
		// RC1 roll compare register???
		//log.Printf("RC1 roll compare register write")
		r.GetBus().FindSingleDevice(common.MODULE_INTEL_82335).(*intel82335.Intel82335).Rc1RegisterWrite(value)
		return
	}

	if port_addr == 0x26 {
		// DMA command register write
		r.GetBus().FindSingleDevice(common.MODULE_INTEL_82335).(*intel82335.Intel82335).DmaCommandRegisterWrite(value)
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
			err := r.GetBus().SendMessageSingle(common.MODULE_MEMORY_ACCESS_CONTROLLER, bus.BusMessage{Subject: common.MESSAGE_DISABLE_A20_GATE, Data: []byte{value}})
			if err != nil {
				log.Fatalf("Failed to send message to memory access controller: %s", err)
				return
			}
		} else {
			err := r.GetBus().SendMessageSingle(common.MODULE_MEMORY_ACCESS_CONTROLLER, bus.BusMessage{Subject: common.MESSAGE_ENABLE_A20_GATE, Data: []byte{value}})
			if err != nil {
				log.Fatalf("Failed to send message to memory access controller: %s", err)
				return
			}
		}

		return
	}

	if port_addr == 0x08 {
		// Write command register to DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).WriteCommandRegister(value)
		return
	}

	if port_addr == 0x09 {
		// Write request register to DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).WriteRequestRegister(value)
		return
	}

	if port_addr == 0x0A {
		// Write single mask register to DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).WriteSingleMaskRegister(value)
		return
	}

	if port_addr == 0x0B {
		// Write mode register to DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).WriteModeRegister(value)
		return
	}

	if port_addr == 0x0C {
		// Clear byte pointer flip-flop in DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).ClearBytePointerFlipFlop()
		return
	}

	if port_addr == 0x0D {
		// Read temporary register from DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).ReadTemporaryRegister()
		return
	}

	if port_addr == 0x0D {
		// Master clear DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).MasterClear()
		return
	}

	if port_addr == 0x0E {
		// Clear mask register in DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).ClearMaskRegister()
		return
	}

	if port_addr == 0x0F {
		// Write mask register to DMA controller
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER).(*intel8237.Intel8237).WriteMaskRegister(value)
		return
	}

	if port_addr == 0x00D0 {
		// Write command register to DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).WriteCommandRegister(value)
		return
	}

	if port_addr == 0x00D2 {
		// Write request register to DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).WriteRequestRegister(value)
		return
	}

	if port_addr == 0x00D4 {
		// Write single mask register to DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).WriteSingleMaskRegister(value)
		return
	}

	if port_addr == 0x00D6 {
		// Write mode register to DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).WriteModeRegister(value)
		return
	}

	if port_addr == 0x00D8 {
		// Clear byte pointer flip-flop in DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).ClearBytePointerFlipFlop()
		return
	}

	if port_addr == 0x00DA {
		// Read temporary register from DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).ReadTemporaryRegister()
		return
	}

	if port_addr == 0x00DA {
		// Master clear DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).MasterClear()
		return
	}

	if port_addr == 0x00DC {
		// Clear mask register in DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).ClearMaskRegister()
		return
	}

	if port_addr == 0x00DE {
		// Write mask register to DMA controller (channels 4-7)
		r.GetBus().FindSingleDevice(common.MODULE_DMA_CONTROLLER_2).(*intel8237.Intel8237).WriteMaskRegister(value)
		return
	}

	if port_addr == 0x20 || port_addr == 0xA0 {
		r.GetBus().FindSingleDevice(common.MODULE_INTERRUPT_CONTROLLER_1).(*intel8259a.Intel8259a).DataWrite(value)
		return
	}

	if port_addr == 0x21 || port_addr == 0xA1 {
		r.GetBus().FindSingleDevice(common.MODULE_INTERRUPT_CONTROLLER_2).(*intel8259a.Intel8259a).DataWrite(value)
		return
	}

	if port_addr == 0x03d8 {
		// CGA
		r.GetBus().FindSingleDevice(common.MODULE_CGA).(*cga.Motorola6845).WriteAddr8(port_addr, value)
		return
	}

	if port_addr == 0x0042 {
		// Write to IO port 0x0042
		r.GetBus().FindSingleDevice(common.MODULE_PIT).(*intel82C54.Intel82C54).WriteCounter0(value)
		return
	}

	if port_addr == 0x0043 {
		// PIT
		r.GetBus().FindSingleDevice(common.MODULE_PIT).(*intel82C54.Intel82C54).CommandRegisterWrite(value)
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
}
