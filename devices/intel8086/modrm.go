package intel8086

type ModRm struct {
	mod uint8
	reg uint8
	rm  uint8

	sib uint8
	base   uint8
	index   uint8
	scale   uint8


	disp8  uint8
	disp16 uint16
	disp32 uint32

}


func (core *CpuCore) consumeModRm() (ModRm, uint32) {

	var bytesConsumed uint32
	m := ModRm{}

	modrmByte := core.memoryAccessController.ReadAddr8(uint32(core.currentByteAddr))
	bytesConsumed++

	m.mod = (modrmByte >> 6) & 0x03
	m.reg = (modrmByte >> 3) & 0x07
	m.rm = modrmByte & 0x07


	if core.registers.CR0 >> 0 & 1 == 0 {
		// real mode
		if m.mod == 1 {
			m.disp8 = core.memoryAccessController.ReadAddr8(uint32(core.currentByteAddr+bytesConsumed))
			bytesConsumed++
		} else if (m.mod == 0 && m.rm == 6) || m.mod == 2 {
			m.disp16 = core.memoryAccessController.ReadAddr16(core.currentByteAddr+bytesConsumed)
			bytesConsumed += 2
		}
	} else {
		// protected mode
		if m.mod != 3 && m.rm == 4 {
			m.sib = core.memoryAccessController.ReadAddr8(core.currentByteAddr+bytesConsumed)
			bytesConsumed++
		}

		if (m.mod == 0 && m.rm == 5) || m.mod == 2 {
			m.disp32 = core.memoryAccessController.ReadAddr32(core.currentByteAddr+bytesConsumed)

			bytesConsumed += 4
		} else if m.mod == 1 {
			m.disp8 = core.memoryAccessController.ReadAddr8(core.currentByteAddr+bytesConsumed)

			bytesConsumed++
		}
	}

	return m, bytesConsumed
}

// derived from:
// https://www.intel.com.au/content/www/au/en/architecture-and-technology/64-ia-32-architectures-software-developer-instruction-set-reference-manual-325383.html
// table 2.1
func (m *ModRm) getAddressMode16(core *CpuCore) uint16 {
	if m.mod == 0 {
		switch m.rm {
		case 0:
			return core.registers.BX + core.registers.SI
		case 1:
			return core.registers.BX + core.registers.DI
		case 2:
			return core.registers.BP + core.registers.SI
		case 3:
			return core.registers.BP + core.registers.DI
		case 4:
			return core.registers.SI
		case 5:
			return core.registers.DI
		case 6:
			return uint16(m.disp16)
		case 7:
			return core.registers.BX
		}
	} else if m.mod == 1 {
		if m.rm == 6 {
			return uint16(int32(core.registers.BP) + int32(m.disp8))
		}
		m.mod = 0
		return uint16(int32(m.getAddressMode16(core)) + int32(m.disp8))
	} else if m.mod == 2 {
		if m.rm == 6 {
			return uint16(int32(core.registers.BP) + int32(m.disp16))
		}
		m.mod = 0
		return uint16(int32(m.getAddressMode16(core)) + int32(m.disp16))
	}
	return uint16(0)
}


func (m *ModRm) getAddressMode32(core *CpuCore) uint32 {
	if m.mod == 0 {
		if m.rm == 5 {
			return m.disp32 // Is this a EBP?
		} else if m.rm == 4 {

			return m.regFromSib(core)
		}

		return *core.registers.registers32Bit[m.rm]
	} else if m.mod == 1 {
		var result uint32
		if m.rm == 4 {
			result = m.regFromSib(core)
		} else {
			result = *core.registers.registers32Bit[m.rm]
		}

		disp8 := m.disp8
		if disp8 < 0 {
			result -= uint32(-disp8)
		} else {
			result += uint32(disp8)
		}

		return result
	} else if m.mod == 2 {
		var result uint32
		if m.rm == 4 {
			result = m.regFromSib(core)
		} else {
			result = *core.registers.registers32Bit[m.rm]
		}
		result += m.disp32
		return result
	}

	return uint32(0)
}

func (m *ModRm) regFromSib(core *CpuCore) uint32 {

	// decode sip byte
	m.sib = core.memoryAccessController.ReadAddr8(core.currentByteAddr)
	if m.mod < 3 && m.rm == 4 {
		m.base = uint8(m.sib & 0x7)
		m.index = uint8((m.sib >> 3) & 0x7)
		m.scale = uint8((m.sib >> 6) & 0x3)
	}

	// calc base value
	var result uint32
	if m.base == 5 && m.mod == 0 {
		result = m.disp32
	} else {
		result = *core.registers.registers32Bit[m.base]
	}

	// index
	if m.index != 4 {
		result += (*core.registers.registers32Bit[m.index]) * uint32(1<<m.scale)
	}

	return result
}
