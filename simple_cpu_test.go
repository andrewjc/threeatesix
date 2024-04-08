package main

import (
	"github.com/andrewjc/threeatesix/pc"
	"log"
	"testing"
)

func Test_DoCpuTest(t *testing.T) {

	testPc := pc.NewPc() //build a new pc for each test run
	testPc.GetPrimaryCpu().Init(testPc.GetBus())
	testPc.GetMemoryController().UnlockBootVector()

	testPc.GetPrimaryCpu().SetCS(0x0)
	testPc.GetPrimaryCpu().SetIP(1)

	instructions := []uint8{0xf, 0x20, 0xc0, 0x66, 0x25, 0xff, 0xff, 0xff, 0x9f, 0xf, 0x22, 0xc0, 0xff, 0xe7, 0xf, 0x1, 0xe0}

	for x := 0; x < len(instructions); x++ {
		testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(x)), instructions[x])
	}

	var tmp uint16
	for {
		if testPc.GetPrimaryCpu().GetIP() == 0 {
			break
		}
		testPc.GetPrimaryCpu().Step()

		if tmp == testPc.GetPrimaryCpu().GetIP() {
			log.Panic("CPU Appears Stuck!")
		} else {
			tmp = testPc.GetPrimaryCpu().GetIP()
		}
	}

}
