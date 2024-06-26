package tests

import (
	"github.com/andrewjc/threeatesix/pc"
	"testing"
)

func Test_MovTests(t *testing.T) {

	testPc := pc.NewPc() //build a new pc for each test run
	testPc.GetPrimaryCpu().Init(testPc.GetBus())
	testPc.GetMemoryController().UnlockBootVector()

	testPc.GetPrimaryCpu().SetCS(0x0)
	testPc.GetPrimaryCpu().SetIP(1)

	instructions := []uint8{0xb8, 0x00, 0x90, 0x8E, 0xD8}

	for x := 0; x < len(instructions); x++ {
		testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(x)), instructions[x])
	}

	for {
		if testPc.GetPrimaryCpu().GetIP() == 0 {
			break
		}
		testPc.GetPrimaryCpu().Step()
	}
}

func Test_MovDecodeOperandSizeOverridePrefix(t *testing.T) {

	testPc := pc.NewPc() //build a new pc for each test run
	testPc.GetPrimaryCpu().Init(testPc.GetBus())
	testPc.GetMemoryController().UnlockBootVector()

	testPc.GetPrimaryCpu().SetCS(0x0)
	testPc.GetPrimaryCpu().SetIP(1)

	instructions := []uint8{0x66, 0x25, 0xFF, 0xFF, 0xFF, 0x9F}

	for x := 0; x < len(instructions); x++ {
		testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(x)), instructions[x])
	}

	for {
		if testPc.GetPrimaryCpu().GetIP() == 0 {
			break
		}
		testPc.GetPrimaryCpu().Step()
	}
}
