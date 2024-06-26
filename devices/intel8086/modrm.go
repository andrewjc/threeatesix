package intel8086

import "fmt"

type ModRm struct {
	mod uint8 // Mod field specifies the addressing mode.
	reg uint8 // Reg field typically used for register opcode extension.
	rm  uint8 // R/M field specifies register/memory operand addressing.

	sib   uint8 // SIB byte used in protected mode for more complex addressing.
	base  uint8 // Base register index (part of SIB byte decoding).
	index uint8 // Index register index (part of SIB byte decoding).
	scale uint8 // Scale factor (part of SIB byte decoding).

	disp8  uint8  // 8-bit displacement (used with certain addressing modes).
	disp16 uint16 // 16-bit displacement (used with certain addressing modes).
	disp32 uint32 // 32-bit displacement (used with certain addressing modes).
}

func (m *ModRm) String() string {
	return fmt.Sprintf("ModRm{mod=%d, reg=%d, rm=%d, sib=%d, base=%d, index=%d, scale=%d, disp8=%d, disp16=%d, disp32=%d}",
		m.mod, m.reg, m.rm, m.sib, m.base, m.index, m.scale, m.disp8, m.disp16, m.disp32)
}

func (core *CpuCore) consumeModRm() (ModRm, uint32, error) {
	var bytesConsumed uint32
	var modrmByte uint8
	m := ModRm{}

	modrmByte, err := core.memoryAccessController.ReadMemoryValue8(uint32(core.currentByteAddr))
	if err != nil {
		return m, bytesConsumed, err
	}

	bytesConsumed++
	m.mod = (modrmByte >> 6) & 0x03
	m.reg = (modrmByte >> 3) & 0x07
	m.rm = modrmByte & 0x07

	// Handle real mode and protected mode addressing
	if core.registers.CR0&1 == 0 { // Real mode
		bytesConsumed, err = handleRealModeAddressing(core, &m, bytesConsumed)
	} else { // Protected mode
		bytesConsumed, err = handleProtectedModeAddressing(core, &m, bytesConsumed)
	}
	if err != nil {
		return m, bytesConsumed, err
	}

	return m, bytesConsumed, nil
}
func handleRealModeAddressing(core *CpuCore, m *ModRm, bytesConsumed uint32) (uint32, error) {
	var err error

	switch m.mod {
	case 0:
		if m.rm == 6 { // Special case: direct addressing with a 16-bit displacement
			m.disp16, err = core.memoryAccessController.ReadMemoryValue16(uint32(core.currentByteAddr + bytesConsumed))
			if err != nil {
				return bytesConsumed, err
			}
			bytesConsumed += 2
		}
		// No displacement for other RM values when mod is 0, except for RM = 6

	case 1: // 8-bit displacement
		var disp8 uint8
		disp8, err = core.memoryAccessController.ReadMemoryValue8(uint32(core.currentByteAddr + bytesConsumed))
		if err != nil {
			return bytesConsumed, err
		}
		m.disp8 = disp8
		bytesConsumed++

	case 2: // 16-bit displacement
		m.disp16, err = core.memoryAccessController.ReadMemoryValue16(uint32(core.currentByteAddr + bytesConsumed))
		if err != nil {
			return bytesConsumed, err
		}
		bytesConsumed += 2
	}

	return bytesConsumed, nil
}

func handleProtectedModeAddressing(core *CpuCore, m *ModRm, bytesConsumed uint32) (uint32, error) {
	var err error

	// Check if there's an SIB byte to decode
	if m.mod != 3 && m.rm == 4 { // SIB byte is present when mod != 3 and rm = 4
		m.sib, err = core.memoryAccessController.ReadMemoryValue8(uint32(core.currentByteAddr + bytesConsumed))
		if err != nil {
			return bytesConsumed, err
		}
		bytesConsumed++
		// Decode the SIB byte
		m.scale = (m.sib >> 6) & 0x03
		m.index = (m.sib >> 3) & 0x07
		m.base = m.sib & 0x07
	}

	// Handle displacements based on the mod value
	switch m.mod {
	case 0:
		if m.rm == 5 { // Only when rm is 5, it means disp32 with no base register in mod 0
			m.disp32, err = core.memoryAccessController.ReadMemoryValue32(uint32(core.currentByteAddr + bytesConsumed))
			if err != nil {
				return bytesConsumed, err
			}
			bytesConsumed += 4
		} // else no displacement unless there's an SIB byte which changes handling

	case 1: // 8-bit displacement, applies to all rm values including when an SIB byte is present
		var disp8 uint8
		disp8, err = core.memoryAccessController.ReadMemoryValue8(uint32(core.currentByteAddr + bytesConsumed))
		if err != nil {
			return bytesConsumed, err
		}
		m.disp8 = disp8
		bytesConsumed++

	case 2: // 32-bit displacement, applies to all rm values including when an SIB byte is present
		m.disp32, err = core.memoryAccessController.ReadMemoryValue32(uint32(core.currentByteAddr + bytesConsumed))
		if err != nil {
			return bytesConsumed, err
		}
		bytesConsumed += 4
	}

	return bytesConsumed, nil
}
