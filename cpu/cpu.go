package cpu

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
)


func New80386CPU() CpuCore {

	cpuCore := CpuCore{}
	cpuCore.registers = &CpuRegisters{}

	cpuCore.opCodeMap = make([]OpCodeImpl, 256)

	mapOpCodes(&cpuCore)

	return cpuCore
}

type CpuCore struct {
	ram       *common.CpuMemInterconnect
	registers *CpuRegisters
	opCodeMap []OpCodeImpl
}

func (core *CpuCore) Step() {
	instrByte := core.ram.ReadHalfWord() //read 8 bit value

	instructionImpl := core.opCodeMap[instrByte]
	if instructionImpl != nil {
		instructionImpl(core)
	} else {
		log.Fatalf("CPU core failure. EIP: %08x - Unrecognised opcode: %#2x\n", core.registers.EIP, instrByte)
	}


	fmt.Printf("CPU Stepped...\n")
}

func (core *CpuCore) Init(bootEip uint32, ramInterconnect *common.CpuMemInterconnect) {
	core.ram = ramInterconnect
	core.registers.EIP = bootEip
	core.ram.SetCpuRegisterController(core)
}

type CpuRegisters struct {
	EIP uint32

	EAX uint32
	ECX uint32
	EDX uint32
	EBX uint32
	ESP uint32
	EBP uint32
	ESI uint32
	EDI uint32
}

func (core *CpuCore) IncrementEIP() {
	core.registers.EIP++
}

func (core *CpuCore) GetEIP() uint32 {
	return core.registers.EIP
}

/* CPU OPCODE IMPLEMENTATIONS */

func mapOpCodes(c *CpuCore) {
	c.opCodeMap[0x30] = XOR_rm8_r8


	// 0x40…0x47, 0xFE/0, 0xFF/0
	for i:=61;i<=68;i++ {
		c.opCodeMap[uint8(i)] = INC_WD_REG
	}

}

type OpCodeImpl func(*CpuCore)

func XOR_rm8_r8(core *CpuCore) {
	// https://c9x.me/x86/html/file_module_x86_id_330.html
	// Destination = Destination ^ Source;

	destAddr := core.ram.ReadNextWord()
	srcAddr := core.ram.ReadNextWord()

	destVal := core.ram.ReadAddr(destAddr)
	srcVal := core.ram.ReadAddr(srcAddr)

	destVal = destVal ^ srcVal

	core.ram.WriteAddr(destAddr, destVal)
}

func INC_WD_REG(core *CpuCore) {
	// https://c9x.me/x86/html/file_module_x86_id_140.html

	// Destination = Destination + 1;

	destAddr := core.ram.ReadNextWord()

	log.Printf("REG: %#2x", destAddr)
}