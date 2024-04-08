package main

import (
	"fmt"
	"github.com/andrewjc/threeatesix/devices/intel8086"
	"github.com/andrewjc/threeatesix/pc"
	"testing"
)

const MaxRAMBytes = 0x1E84800 //32 million (32mb)

var testPc pc.PersonalComputer

func Test_INSTR_JCXZ_SHORT_REL8(t *testing.T) {

	tests := []struct {
		name        string
		instruction []uint8
		cxValue     uint16
		expectedIP  uint16
	}{
		// TODO: Add test cases.

		{"TestJCXZ_SHORT_REL8_CxZeroValue", []uint8{0xe3, 0x3b}, 0, 0x013d},
		{"TestJCXZ_SHORT_REL8_CxNonZeroValue", []uint8{0xe3, 0x3b}, 1, 0x0102},
	}
	for _, tt := range tests {

		testPc := pc.NewPc() //build a new pc for each test run
		testPc.GetPrimaryCpu().Init(testPc.GetBus())
		testPc.GetMemoryController().UnlockBootVector()

		t.Run(tt.name, func(t *testing.T) {
			testPc.GetPrimaryCpu().SetCS(0x0)
			testPc.GetPrimaryCpu().SetIP(0x100)

			for x := 0; x < len(tt.instruction); x++ {
				testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(x)), tt.instruction[x])
			}

			testPc.GetPrimaryCpu().GetRegisters().CX = tt.cxValue

			testPc.GetPrimaryCpu().Step()

			if testPc.GetPrimaryCpu().GetIP() != tt.expectedIP {
				panic(fmt.Errorf("Expected ip [%#04x] but got [%#04x]", tt.expectedIP, testPc.GetPrimaryCpu().GetIP()))
			}
		})
	}
}

func Test_INSTR_JZ_SHORT_REL8(t *testing.T) {

	tests := []struct {
		name        string
		instruction []uint8
		zeroFlag    bool
		expectedIP  uint16
	}{
		// TODO: Add test cases.
		{"TestJZ_SHORT_REL8_ZeroFlag", []uint8{0x74, 0xee}, true, 0x00f0},
		{"TestJZ_SHORT_REL8_NonZeroFlag", []uint8{0x74, 0xee}, false, 0x0102},
	}
	for _, tt := range tests {

		testPc := pc.NewPc() //build a new pc for each test run
		testPc.GetPrimaryCpu().Init(testPc.GetBus())
		testPc.GetMemoryController().UnlockBootVector()

		t.Run(tt.name, func(t *testing.T) {
			testPc.GetPrimaryCpu().SetCS(0x0)
			testPc.GetPrimaryCpu().SetIP(0x100)

			for x := 0; x < len(tt.instruction); x++ {
				testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(x)), tt.instruction[x])
			}

			testPc.GetPrimaryCpu().SetFlag(intel8086.ZeroFlag, tt.zeroFlag)

			testPc.GetPrimaryCpu().Step()

			if testPc.GetPrimaryCpu().GetIP() != tt.expectedIP {
				panic(fmt.Errorf("Expected ip [%#04x] but got [%#04x]", tt.expectedIP, testPc.GetPrimaryCpu().GetIP()))
			}
		})
	}
}

func Test_LOD_INSTRUCTIONS(t *testing.T) {

	tests := []struct {
		name        string
		instruction []uint8
		cxValue     uint16
	}{
		// TODO: Add test cases.
		{"Test_REP_LODSW", []uint8{0xf3, 0xad}, 2},
		{"Test_LODSB", []uint8{0x74, 0xac}, 0},
		{"Test_LODSD", []uint8{0x74, 0xad}, 0},
	}
	for _, tt := range tests {

		testPc := pc.NewPc() //build a new pc for each test run
		testPc.GetPrimaryCpu().Init(testPc.GetBus())
		testPc.GetMemoryController().UnlockBootVector()

		t.Run(tt.name, func(t *testing.T) {
			testPc.GetPrimaryCpu().SetCS(0x0)
			testPc.GetPrimaryCpu().SetIP(0x100)
			testPc.GetPrimaryCpu().GetRegisters().CX = tt.cxValue

			for x := 0; x < len(tt.instruction); x++ {
				testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(x)), tt.instruction[x])
			}

			testPc.GetPrimaryCpu().Step()

			/*if testPc.GetPrimaryCpu().GetIP() != tt.expectedIP {
				panic(fmt.Errorf("Expected ip [%#04x] but got [%#04x]", tt.expectedIP,testPc.GetPrimaryCpu().GetIP() ))
			}*/
		})
	}
}

func Test_BufferInitTest_INSTRUCTIONS(t *testing.T) {

	tests := []struct {
		name        string
		instruction []uint8
	}{
		// TODO: Add test cases.
		{"Test_Put1Into100LenBuffer", []uint8{0xf3, 0xad}},
	}
	for _, tt := range tests {

		testPc := pc.NewPc() //build a new pc for each test run
		testPc.GetPrimaryCpu().Init(testPc.GetBus())
		testPc.GetMemoryController().UnlockBootVector()

		t.Run(tt.name, func(t *testing.T) {
			testPc.GetPrimaryCpu().SetCS(0x0)
			testPc.GetPrimaryCpu().SetIP(0x100)

			for x := 0; x < len(tt.instruction); x++ {
				testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(x)), tt.instruction[x])
			}

			testPc.GetPrimaryCpu().Step()

			/*if testPc.GetPrimaryCpu().GetIP() != tt.expectedIP {
				panic(fmt.Errorf("Expected ip [%#04x] but got [%#04x]", tt.expectedIP,testPc.GetPrimaryCpu().GetIP() ))
			}*/
		})
	}
}
