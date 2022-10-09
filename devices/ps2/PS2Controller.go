package ps2

import (
    "github.com/andrewjc/threeatesix/common"
    "github.com/andrewjc/threeatesix/devices/bus"
    "log"
)

type Ps2Controller struct {
    bus                *bus.Bus
    busId              uint32
    statusRegister     uint8
    bufferedOutputData []uint8

    internalRam      uint8 // used for storing config bytes
    pendingOperation uint8
    port1_enabled    bool
    port2_enabled    bool
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

func (controller *Ps2Controller) ReadDataPort() uint8 {
    if len(controller.bufferedOutputData) > 0 {
        response := controller.bufferedOutputData[0]
        controller.bufferedOutputData = nil
        controller.DisableDataPortReadyForRead()
        return response
    }
    return 0x00
}

func (controller *Ps2Controller) WriteCommandRegister(value uint8) {
    log.Printf("PS2 controller write command: [%#04x]", value)
    if value == 0xAA {
        // Test PS2 Controller
        controller.SendBufferedResponse([]uint8{0x55}) //OK Message
        return
    }

    if value == 0xAD {
        // Disable first ps2 port
        controller.port1_enabled = false
        return
    }

    if value == 0xAE {
        // Enable first ps2 port
        controller.port1_enabled = true
        return
    }

    if value == 0x60 {
        // Write next byte to byte 0 of internal data
        controller.pendingOperation = 0x60
        controller.EnableDataPortReadyForWrite()
        return
    }

    log.Printf("Unknown PS2 controller write command: [%#04x]", value)

}

func (controller *Ps2Controller) SendBufferedResponse(response []uint8) {
    controller.bufferedOutputData = response
    controller.EnableDataPortReadyForRead()
}

func (controller *Ps2Controller) EnableDataPortReadyForWrite() {
    controller.statusRegister = common.Reset(controller.statusRegister, 1)
}

func (controller *Ps2Controller) DisableDataPortReadyForWrite() {
    controller.statusRegister = common.Set(controller.statusRegister, 1)
}
func (controller *Ps2Controller) EnableDataPortReadyForRead() {
    controller.statusRegister = common.Set(controller.statusRegister, 0)
}

func (controller *Ps2Controller) DisableDataPortReadyForRead() {
    controller.statusRegister = common.Reset(controller.statusRegister, 0)
}
