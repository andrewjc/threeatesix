package main

import (
	"github.com/andrewjc/threeatesix/pc"
)

/*
	ThreeAteSix - A 386 emulator
*/

func main() {

	machine := pc.NewPc()

	machine.LoadBios()
	machine.Power()

}

