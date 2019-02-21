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
	memoryAccessController *common.MemoryAccessController
	registers              *CpuRegisters
	opCodeMap              []OpCodeImpl
	mode                   uint8
}

func (core *CpuCore) SetIP(addr uint16) {
	core.registers.IP = addr
}

func (core *CpuCore) GetIP() uint16 {
	return core.registers.IP
}

func (core *CpuCore) GetCS() uint16 {
	return core.registers.CS
}

func (core *CpuCore) IncrementIP() {
	core.registers.IP++
}

func (core *CpuCore) Step() {
	instrByte := core.memoryAccessController.ReadHalfWord() //read 8 bit value

	instructionImpl := core.opCodeMap[instrByte.(uint8)]
	if instructionImpl != nil {
		instructionImpl(core)
	} else {
		log.Printf("CPU CORE ERROR!!!")
		core.dumpCoreInfo()
		log.Fatalf("CPU core failure. Unrecognised opcode: %#2x\n", instrByte)
	}

	fmt.Printf("CPU Stepped...\n")
}

func (core *CpuCore) Init(memController *common.MemoryAccessController) {
	core.memoryAccessController = memController

	core.EnterMode(common.REAL_MODE)

	core.memoryAccessController.SetCpuRegisterController(core)

	core.Reset()
}

func (core *CpuCore) Reset() {
	core.registers.CS = 0xF000
	core.registers.IP = 0xFFF0
	core.memoryAccessController.LockBootVector()
}

type CpuRegisters struct {

	// 16bit registers (real mode)
	CS uint16 // code segment
	DS uint16 // data segment
	IP uint16 // instruction pointer

	// 32bit registers (protected mode)
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

func (core *CpuCore) EnterMode(mode uint8) {
	core.mode = mode
	core.memoryAccessController.EnterMode(mode)
}

func (core *CpuCore) dumpCoreInfo() {
	if core.mode == common.REAL_MODE {
		log.Println("Cpu core in real mode")
		log.Printf("Registers: IP: %016X, CS: %016X, DS: %016X", core.registers.IP, core.registers.CS, core.registers.DS)
	}
}

/* CPU OPCODE IMPLEMENTATIONS */

func mapOpCodes(c *CpuCore) {

	c.opCodeMap[0xEA] = JMP_FAR_PTR16
	c.opCodeMap[0xE9] = JMP_NEAR_REL16

	c.opCodeMap[0x30] = XOR_rm8_r8

	// 0x40…0x47, 0xFE/0, 0xFF/0
	for i := 61; i <= 68; i++ {
		c.opCodeMap[uint8(i)] = INC_WD_REG
	}

}

type OpCodeImpl func(*CpuCore)

func JMP_FAR_PTR16(core *CpuCore) {
	destAddr := core.memoryAccessController.ReadNextWord().(uint16)
	log.Printf("INSTR JMP FAR PTR %016X", destAddr)

	core.memoryAccessController.SetIP(destAddr)
}

func JMP_NEAR_REL16(core *CpuCore) {
	destAddr := core.memoryAccessController.ReadNextWord().(uint16)
	log.Printf("INSTR JMP NEAR REL 16 %016X", destAddr)

	destAddr = core.registers.IP + destAddr

	core.memoryAccessController.SetIP(destAddr)
}

func XOR_rm8_r8(core *CpuCore) {
	// https://c9x.me/x86/html/file_module_x86_id_330.html
	// Destination = Destination ^ Source;

	destAddr := core.memoryAccessController.ReadNextWord()
	srcAddr := core.memoryAccessController.ReadNextWord()

	destVal := core.memoryAccessController.ReadAddr(destAddr.(uint16))
	srcVal := core.memoryAccessController.ReadAddr(srcAddr.(uint16))

	destVal = destVal ^ srcVal

	core.memoryAccessController.WriteAddr(destAddr.(uint16), destVal)
}

func INC_WD_REG(core *CpuCore) {
	// https://c9x.me/x86/html/file_module_x86_id_140.html

	// Destination = Destination + 1;

	destAddr := core.memoryAccessController.ReadNextWord()

	log.Printf("REG: %#2x", destAddr)
}
