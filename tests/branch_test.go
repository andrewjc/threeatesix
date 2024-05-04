package tests_test

import (
	"github.com/andrewjc/threeatesix/devices/intel8086"
	"github.com/andrewjc/threeatesix/devices/memmap"
	"github.com/andrewjc/threeatesix/pc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupCpuCore() (*intel8086.CpuCore, *memmap.MemoryAccessController) {
	testPc := pc.NewPc() //build a new pc for each test run
	testPc.GetPrimaryCpu().Init(testPc.GetBus())
	testPc.GetMemoryController().UnlockBootVector()

	testPc.GetPrimaryCpu().SetCS(0x0)
	testPc.GetPrimaryCpu().SetIP(1)
	return testPc.GetPrimaryCpu(), testPc.GetMemoryController()
}

func TestINSTR_JMP_NEAR_REL16(t *testing.T) {
	// Setup
	core, mem := setupCpuCore()

	// Set initial IP value
	core.GetRegisters().IP = 0x1000
	core.GetRegisters().CS = intel8086.SegmentRegister{
		Base:  0x1000,
		Limit: 0x1000,
	}

	// Set the instruction
	mem.WriteMemoryAddr8(core.GetCurrentCodePointer(), 0xE9)
	mem.WriteMemoryAddr16(uint32(core.GetCurrentCodePointer())+1, 0x05)

	// Call the function
	intel8086.INSTR_JMP_NEAR_REL16(core)

	// Assert the expected results
	assert.Equal(t, uint16(0x1008), core.GetRegisters().IP)
	assert.Equal(t, uint32(0x1000), core.GetRegisters().CS.Base)

}
