package intel8086

import "fmt"

type SegmentRegister struct {
	base  uint16
	limit uint32
	//selector uint16
	access_information uint16
}

type CpuRegisters struct {
	registers8Bit  []*uint8
	registers16Bit []*uint16
	registers32Bit []*uint32

	registersSegmentRegisters []*SegmentRegister

	// 16bit registers (real mode)
	CS SegmentRegister // code segment
	DS SegmentRegister // data segment
	SS SegmentRegister // stack segment
	ES SegmentRegister // extended segment
	FS SegmentRegister // ?? segment
	GS SegmentRegister // ?? segment

	IP  uint16 // 16 bit instruction pointer
	SP  uint16
	BP  uint16
	SI  uint16
	DI  uint16
	EIP uint32 // 32 bit instruction pointer
	ESP uint32
	EBP uint32
	ESI uint32
	EDI uint32

	// accumulator registers
	// used for I/O port access, arithmetic, interrupt calls
	AH  uint8
	AL  uint8
	AX  uint16
	EAX uint32

	// base registers
	// used as a base pointer for memory access
	BX  uint16
	EBX uint32
	BH  uint8
	BL  uint8

	// counter registers
	// used as loop counter and for shifts
	CX  uint16
	ECX uint32
	CH  uint8
	CL  uint8

	// data registers
	// used for I/O port access, arithmetic, interrupt calls
	DX  uint16
	EDX uint32
	DH  uint8
	DL  uint8

	// Flags
	FLAGS uint16

	// Control Flag
	CR0 uint32
	CR1 uint32
	CR2 uint32
	CR3 uint32
	CR4 uint32
}

func (c *CpuRegisters) index8ToString(i uint8) string {

	switch {
	case i == 0:
		return "AL"
	case i == 1:
		return "CL"
	case i == 2:
		return "DL"
	case i == 3:
		return "BL"
	case i == 4:
		return "AH"
	case i == 5:
		return "CH"
	case i == 6:
		return "DH"
	case i == 7:
		return "BH"
	default:
		return fmt.Sprintf("Unrecognised 8 bit register index %d", i)
	}
}

func (c *CpuRegisters) index16ToString(i uint8) string {
	switch {
	case i == 0:
		return "AX"
	case i == 1:
		return "CX"
	case i == 2:
		return "DX"
	case i == 3:
		return "BX"
	case i == 4:
		return "SP"
	case i == 5:
		return "BP"
	case i == 6:
		return "SI"
	case i == 7:
		return "DI"
	default:
		return fmt.Sprintf("Unrecognised 16 bit register index %d", i)
	}
}

func (c *CpuRegisters) index32ToString(i uint8) string {
	switch {
	case i == 0:
		return "EAX"
	case i == 1:
		return "ECX"
	case i == 2:
		return "EDX"
	case i == 3:
		return "EBX"
	case i == 4:
		return "ESP"
	case i == 5:
		return "EBP"
	case i == 6:
		return "ESI"
	case i == 7:
		return "EDI"
	default:
		return fmt.Sprintf("Unrecognised 32 bit register index %d", i)
	}
}

func (core *CpuRegisters) indexSegmentToString(i uint8) string {
	switch {
	case i == 0:
		return "ES"
	case i == 1:
		return "CS"
	case i == 2:
		return "SS"
	case i == 3:
		return "DS"
	case i == 4:
		return "FS"
	case i == 5:
		return "GS"
	default:
		return fmt.Sprintf("Unrecognised segment register index %d", i)
	}
}
