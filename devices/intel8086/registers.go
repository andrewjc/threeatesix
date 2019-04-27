package intel8086

import "fmt"

type SegmentRegister struct {
	base uint16
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

	IP uint16 // instruction pointer
	SP uint16
	BP uint16
	SI uint16
	DI uint16

	// accumulator registers
	// used for I/O port access, arithmetic, interrupt calls
	AH uint8
	AL uint8
	AX uint16

	// base registers
	// used as a base pointer for memory access
	BX uint16
	BH uint8
	BL uint8

	// counter registers
	// used as loop counter and for shifts
	CX uint16
	CH uint8
	CL uint8

	// data registers
	// used for I/O port access, arithmetic, interrupt calls
	DX uint16
	DH uint8
	DL uint8

	// Flags
	FLAGS uint16

	// Control Flag
	CR0   uint32
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
	default:
		return fmt.Sprintf("Unrecognised segment register index %d", i)
	}
}

