package cmos

import (
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

type Motorola146818 struct {
	bus      *bus.Bus
	busId    uint32
	cmosData [128]uint8
	index    uint8
}

func NewMotorola146818() *Motorola146818 {
	return &Motorola146818{}
}

func (d *Motorola146818) GetDeviceBusId() uint32 {
	return d.busId
}

func (d *Motorola146818) SetDeviceBusId(id uint32) {
	d.busId = id
}

func (d *Motorola146818) SetBus(bus *bus.Bus) {
	d.bus = bus
}

func (d *Motorola146818) OnReceiveMessage(message bus.BusMessage) {
}

func (d *Motorola146818) GetPortMap() *bus.DevicePortMap {
	return &bus.DevicePortMap{
		ReadPorts:  []uint16{0x71},
		WritePorts: []uint16{0x70, 71},
	}
	return nil
}

func (d *Motorola146818) ReadAddr8(addr uint16) uint8 {
	switch addr {
	case 0x71: // Data port read
		if d.index < 128 {
			friendlyCmosString := common.CmosRegisterWriteToFriendlyString(d.index, d.cmosData[d.index])
			log.Printf("CMOS RAM: %#02x -> %#02x (%s)", d.index, d.cmosData[d.index], friendlyCmosString)
			return d.cmosData[d.index]
		}
	default:
		log.Printf("Motorola6845: Unsupported read from address 0x%04X", addr)
	}
	return 0
}

func (d *Motorola146818) WriteAddr8(addr uint16, data uint8) {
	switch addr {
	case 0x70: // Index port
		d.index = data & 0x7F // Masking the top bit to avoid side effects like disabling NMI
	case 0x71: // Data port
		if d.index < 128 {
			d.cmosData[d.index] = data
		}
	default:
		log.Printf("Motorola6845: Unsupported write to address 0x%04X with data 0x%02X", addr, data)
	}
}
