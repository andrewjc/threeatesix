package main

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"github.com/andrewjc/threeatesix/devices/intel8086"
	"github.com/andrewjc/threeatesix/devices/intel8259a"
	"github.com/andrewjc/threeatesix/devices/io"
	"github.com/andrewjc/threeatesix/devices/memmap"

	"io/ioutil"
	"os"
)

/*
	ThreeAteSix - A 386 emulator
*/

func main() {
	machine := newPc()

	machine.power()

}

type romimages struct {
	bios []byte
}

// PersonalComputer represents the virtual PC being emulated
type PersonalComputer struct {
	cpu             *intel8086.CpuCore
	mathCoProcessor *intel8086.CpuCore

	bus *bus.Bus

	ram []byte
	rom romimages

	masterInterruptController *intel8259a.Intel8259a
	slaveInterruptController  *intel8259a.Intel8259a

	memController *memmap.MemoryAccessController
	ioPortController *io.IOPortAccessController

}

// BiosFilename - name of the bios image the virtual machine will boot up
const BiosFilename = "bios.bin"

// MaxRAMBytes - the amount of ram installed in this virtual machine
const MaxRAMBytes = 0x1E84800 //32 million (32mb)

func (pc *PersonalComputer) power() {
	// do stuff

	pc.loadBios()

	pc.memController.SetBus(pc.bus)

	pc.ioPortController.SetBus(pc.bus)

	pc.bus.RegisterDevice(pc.cpu, common.MODULE_PRIMARY_PROCESSOR)
	pc.bus.RegisterDevice(pc.mathCoProcessor, common.MODULE_MATH_CO_PROCESSOR)
	pc.bus.RegisterDevice(pc.masterInterruptController, common.MODULE_MASTER_INTERRUPT_CONTROLLER)
	pc.bus.RegisterDevice(pc.slaveInterruptController, common.MODULE_SLAVE_INTERRUPT_CONTROLLER)

	pc.bus.RegisterDevice(pc.memController, common.MODULE_MEMORY_ACCESS_CONTROLLER)
	pc.bus.RegisterDevice(pc.ioPortController, common.MODULE_IO_PORT_ACCESS_CONTROLLER)

	pc.memController.HandleMemoryMapSwitch(common.REAL_MODE)

	pc.cpu.Init(pc.bus)
	pc.mathCoProcessor.Init(pc.bus)

	for {
		// stuff

		pc.cpu.Step()
	}
}

func (pc *PersonalComputer) loadBios() {
	biosData, err := ioutil.ReadFile(BiosFilename)

	if err != nil {
		fmt.Printf("Failed to load bios! - %s", err.Error())
		os.Exit(1)
	}

	pc.rom.bios = make([]byte, len(biosData))
	for i := 0; i < len(biosData); i++ {
		pc.rom.bios[i] = biosData[i]
	}

}

func newPc() *PersonalComputer {
	pc := &PersonalComputer{}

	pc.bus = bus.NewDeviceBus()
	pc.ram = make([]byte, MaxRAMBytes)
	pc.rom = romimages{}
	pc.cpu = intel8086.New80386CPU()
	pc.mathCoProcessor = intel8086.New80287MathCoProcessor()

	pc.masterInterruptController = intel8259a.NewIntel8259a() //pic1
	pc.slaveInterruptController = intel8259a.NewIntel8259a()  //pic2

	pc.memController = memmap.CreateMemoryController(&pc.ram, &pc.rom.bios)

	pc.ioPortController = io.CreateIOPortController()

	return pc
}
