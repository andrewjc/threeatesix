package tests_test

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

func Test_INSTR_JMP_CALL_CLI_STC(t *testing.T) {

	tests := []struct {
		name         string
		instructions []uint8
		setupFunc    func(*pc.PersonalComputer)
		checkFunc    func(*pc.PersonalComputer) error
	}{

		{
			"Test_CALL_NEAR_REL16",
			[]uint8{0xE8, 0x08, 0x00}, // CALL 0008 (Relative call)
			func(pc *pc.PersonalComputer) {
				pc.GetPrimaryCpu().SetIP(0x2000)
				pc.GetPrimaryCpu().GetRegisters().SP = 0xFFFE
			},
			func(pc *pc.PersonalComputer) error {
				if pc.GetPrimaryCpu().GetIP() != 0x200B || pc.GetPrimaryCpu().GetRegisters().SP != 0xFFFC {
					return fmt.Errorf("CALL NEAR REL16 failed: expected IP=0x200B and SP=0xFFFC, got IP=0x%04X and SP=0x%04X", pc.GetPrimaryCpu().GetIP(), pc.GetPrimaryCpu().GetRegisters().SP)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		testPc := pc.NewPc() // build a new pc for each test run
		testPc.GetPrimaryCpu().Init(testPc.GetBus())
		testPc.GetMemoryController().UnlockBootVector()
		tt.setupFunc(testPc)

		t.Run(tt.name, func(t *testing.T) {
			testPc.GetPrimaryCpu().SetCS(0x0)

			for i, instr := range tt.instructions {
				testPc.GetMemoryController().WriteMemoryAddr8(uint32(testPc.GetPrimaryCpu().GetIP()+uint16(i)), instr)
			}

			testPc.GetPrimaryCpu().Step()

			if err := tt.checkFunc(testPc); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_INSTR_JMP_FAR_m16(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(*pc.PersonalComputer)
		expectedCS uint32
		expectedIP uint16
	}{
		{
			name: "Test_JMP_FAR_m16_SimpleJump",
			setupFunc: func(pc *pc.PersonalComputer) {
				// Setting up memory to contain the JMP FAR instruction followed by segment and offset
				// Assuming IP starts at 0x100
				pc.GetPrimaryCpu().SetCS(0x0000)
				pc.GetPrimaryCpu().SetIP(0x0100)
				// Writing JMP FAR instruction 0xEA and the immediate segment:offset data
				pc.GetMemoryController().WriteMemoryAddr8(0x0100, 0xEA)    // Opcode for JMP FAR
				pc.GetMemoryController().WriteMemoryAddr16(0x0101, 0x1234) // Offset
				pc.GetMemoryController().WriteMemoryAddr16(0x0103, 0x0002) // Segment
			},
			expectedCS: 0x0002,
			expectedIP: 0x1234,
		},
	}

	for _, tt := range tests {
		testPc := pc.NewPc() // Build a new PC for each test run
		testPc.GetPrimaryCpu().Init(testPc.GetBus())
		testPc.GetMemoryController().UnlockBootVector()
		tt.setupFunc(testPc)

		t.Run(tt.name, func(t *testing.T) {
			// Execute the JMP FAR instruction
			testPc.GetPrimaryCpu().Step()

			// Check if CS and IP registers are set to the expected values
			actualCS := testPc.GetPrimaryCpu().GetCS()
			actualIP := testPc.GetPrimaryCpu().GetIP()
			if actualCS != tt.expectedCS || actualIP != tt.expectedIP {
				t.Errorf("Test %s failed: expected CS:IP = %04X:%04X, got CS:IP = %04X:%04X",
					tt.name, tt.expectedCS, tt.expectedIP, actualCS, actualIP)
			}
		})
	}
}
