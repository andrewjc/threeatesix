package main

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/cpu"
	"io/ioutil"
	"log"
	"os"
	"unsafe"
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
	cpu    cpu.CpuCore
	isaBus ISABus
	ram    []byte
	rom    RomImages
}

const BIOS_FILENAME = "bios.bin"

const MAX_RAM_BYTES = 0x1E84800 //32 million (32mb)

func (computer *PersonalComputer) power() {
	// do stuff


	solid := uint32(123456789)

	pntr := uintptr(unsafe.Pointer(&solid))

	d1 := &pntr
	d3 := (*uint32)(unsafe.Pointer(pntr))

	log.Printf("%d %d %d", solid, *d1, *d3)

	pntr += 2
	d3 = (*uint32)(unsafe.Pointer(pntr))

	*d3 = 66

	log.Printf("%d %d %d", solid, *d1, *d3)



	memController := common.CreateMemoryController(&computer.ram, &computer.rom.bios)

	ioPortController := common.CreateIOPortController()

	computer.loadBios()

	computer.cpu.Init(memController, ioPortController)

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

	return pc
}
