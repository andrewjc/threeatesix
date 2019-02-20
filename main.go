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

	machine.loadBios()

	machine.power()

}

type PersonalComputer struct {
	cpu    cpu.CpuCore
	isaBus ISABus
	ram    []byte
}

const BIOS_FILENAME = "bios.bin"
const BIOS_MEMORY_LOCATION = 0x7c00
const MAX_RAM_BYTES = 0x1E84800 //32 million (32mb)

func (computer *PersonalComputer) loadBios() {
	biosData, err := ioutil.ReadFile(BIOS_FILENAME)

	if err != nil {
		fmt.Printf("Failed to load bios! - %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < len(biosData); i++ {
		computer.ram[i+BIOS_MEMORY_LOCATION] = biosData[i]
	}

}

func (computer *PersonalComputer) power() {
	// do stuff

	memInterconnect := common.CreateCpuMemoryInterconnect(&computer.ram)

	computer.cpu.Init(BIOS_MEMORY_LOCATION, memInterconnect)

	for {
		// stuff

		computer.cpu.Step()
	}
}


type ISABus struct {
}


func NewPC() *PersonalComputer {
	pc := &PersonalComputer{}

	pc.isaBus = ISABus{}
	pc.ram = make([]byte, MAX_RAM_BYTES)
	pc.cpu = cpu.New80386CPU()

	return pc
}
