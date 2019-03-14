package cpu

import "unsafe"

type ModRm struct {
	mod uint8
	reg uint8
	rm  uint8

	sib    uint8

	disp8  uint8
	disp32 uint32
}

func consumeModRm(core *CpuCore) ModRm {

	m := ModRm{}

	modrmByte := core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer())

	m.mod = (modrmByte >> 6) & 0x03
	m.reg = (modrmByte >> 3) & 0x07
	m.rm = modrmByte & 0x07

	core.IncrementIP()

	if m.mod != 3 && m.rm == 4 {
		m.sib = core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer())
		core.IncrementIP()
	}

	if (m.mod == 0 && m.rm == 5) || m.mod == 2 {
		m.disp32 = core.memoryAccessController.ReadAddr32(core.GetCurrentCodePointer())

		core.IncrementIP()
		core.IncrementIP()
		core.IncrementIP()
		core.IncrementIP()
	} else if m.mod == 1 {
		m.disp8 = core.memoryAccessController.ReadAddr8(core.GetCurrentCodePointer())

		core.IncrementIP()
	}

	return m
}


func getRm16(core *CpuCore) interface{} {
	modrm := consumeModRm(core)
	if modrm.mod == 3 {
		// reg
		return core.registers.registers16Bit[modrm.rm] //pntr
	} else {
		// mem

		return modrm // value

	}
}

func getReg16(core *CpuCore) *uint16 {
	modrm := consumeModRm(core)

	reg := getReg16FromModRm(core, modrm)
	return reg
}

func getReg16FromModRm(core *CpuCore, rm ModRm) *uint16 {
	// get the register reference from lookup table
	return core.registers.registers16Bit[rm.reg]
}

func getAddress16FromModRm(core *CpuCore, rm ModRm) uintptr {

	if rm.mod == 0 {
		if rm.rm == 5 {
			return uintptr(rm.disp32)
		}
		return uintptr(unsafe.Pointer(core.registers.registers16Bit[rm.rm]))
	}

	if rm.mod == 1 {

		v := uintptr(unsafe.Pointer(&core.registers.registers16Bit[rm.rm]))
		v = v + uintptr(rm.disp8)
		return v
	}

	if rm.mod == 2 {

		v := uintptr(unsafe.Pointer(&core.registers.registers16Bit[rm.rm]))
		v = v + uintptr(rm.disp32)
		return v
	}

	return 0
}
