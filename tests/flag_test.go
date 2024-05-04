package tests

import (
	"github.com/andrewjc/threeatesix/devices/intel8086"
	"github.com/andrewjc/threeatesix/pc"
	"testing"
)

func Test_SetFlag(t *testing.T) {

	tests := []struct {
		name         string
		flagMask     uint16
		setTestValue bool
	}{
		// TODO: Add test cases.

		{"TestSettingFlagAndRetrieve", intel8086.ZeroFlag, true},
	}
	for _, tt := range tests {

		testPc := pc.NewPc() //build a new pc for each test run
		testPc.GetPrimaryCpu().Init(testPc.GetBus())
		testPc.GetMemoryController().UnlockBootVector()

		t.Run(tt.name, func(t *testing.T) {
			testPc.GetPrimaryCpu().SetFlag(intel8086.InterruptFlag, true)
			testPc.GetPrimaryCpu().SetFlag(intel8086.IoPrivilegeLevelFlag, true)
			testPc.GetPrimaryCpu().SetFlag(intel8086.DirectionFlag, true)

			testPc.GetPrimaryCpu().SetFlag(tt.flagMask, tt.setTestValue)

			zeroFlagStatus := testPc.GetPrimaryCpu().GetFlag(tt.flagMask)

			zeroFlagIntStatus := testPc.GetPrimaryCpu().GetFlagInt(tt.flagMask)

			if zeroFlagIntStatus != tt.flagMask {
				panic("Flag int doesn't match expected value")
			}

			if zeroFlagStatus != tt.setTestValue {
				panic("Flag bool doesn't match expected value")
			}
		})
	}
}
