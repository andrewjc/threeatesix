package main

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/cpu"
	"github.com/andrewjc/threeatesix/intel8259a"
	"io/ioutil"
	"os"
)

/*
	ThreeAteSix - A 386 emulator
*/

func main() {
	machine := NewPC()

	machine.power()

}

type RomImages struct {
	bios []byte
}

type PersonalComputer struct {
	cpu             cpu.CpuCore
	mathCoProcessor cpu.CpuCore

	bus common.Bus
	ram []byte
	rom RomImages

	masterInterruptController intel8259a.Intel8259a
	slaveInterruptController  intel8259a.Intel8259a
}

const BIOS_FILENAME = "bios.bin"

const MAX_RAM_BYTES = 0x1E84800 //32 million (32mb)

func (computer *PersonalComputer) power() {
	// do stuff

	memController := common.CreateMemoryController(&computer.ram, &computer.rom.bios)

	ioPortController := common.CreateIOPortController()

	computer.loadBios()

	memController.SetCpuController(&computer.cpu)
	memController.SetCoProcessorController(&computer.mathCoProcessor)
	ioPortController.SetCpuController(&computer.cpu)
	ioPortController.SetCoProcessorController(&computer.mathCoProcessor)

	computer.cpu.Init(memController, ioPortController)
	computer.mathCoProcessor.Init(memController, ioPortController)

	for {
		// stuff

		computer.cpu.Step()
	}
}

func (computer *PersonalComputer) loadBios() {
	biosData, err := ioutil.ReadFile(BIOS_FILENAME)

	if err != nil {
		fmt.Printf("Failed to load bios! - %s", err.Error())
		os.Exit(1)
	}

	computer.rom.bios = make([]byte, len(biosData))
	for i := 0; i < len(biosData); i++ {
		computer.rom.bios[i] = biosData[i]
	}

}

func NewPC() *PersonalComputer {
	pc := &PersonalComputer{}

	pc.bus = common.Bus{}
	pc.ram = make([]byte, MAX_RAM_BYTES)
	pc.rom = RomImages{}
	pc.cpu = cpu.New80386CPU()
	pc.mathCoProcessor = cpu.New80287MathCoProcessor()

	pc.masterInterruptController = intel8259a.NewIntel8259a() //pic1
	pc.slaveInterruptController = intel8259a.NewIntel8259a()  //pic2

	pc.bus.registerDevice(pc.cpu, common.MODULE_PRIMARY_PROCESSOR)
	pc.bus.registerDevice(pc.mathCoProcessor, common.MODULE_MATH_CO_PROCESSOR)
	pc.bus.registerDevice(pc.masterInterruptController, common.MODULE_MASTER_INTERRUPT_CONTROLLER)
	pc.bus.registerDevice(pc.slaveInterruptController, common.MODULE_SLAVE_INTERRUPT_CONTROLLER)

	return pc
}
