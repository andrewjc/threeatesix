package main

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/cpu"
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

	isaBus ISABus
	ram    []byte
	rom    RomImages
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

type ISABus struct {
}

func NewPC() *PersonalComputer {
	pc := &PersonalComputer{}

	pc.isaBus = ISABus{}
	pc.ram = make([]byte, MAX_RAM_BYTES)
	pc.rom = RomImages{}
	pc.cpu = cpu.New80386CPU()
	pc.mathCoProcessor = cpu.New80287MathCoProcessor()

	return pc
}
