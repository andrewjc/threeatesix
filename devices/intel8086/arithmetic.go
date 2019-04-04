package intel8086

import (
	"log"
	"math/bits"
)

func INSTR_ADC(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x14:
		{
			// 	adc AL,imm8
			core.IncrementIP()
			term2 = uint32(core.readImm8())
			term1 = uint32(core.registers.AL)
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			core.registers.AL = uint8(term1)

			log.Printf("[%#04x] adc al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x15:
		{
			// adc AX,imm16
			core.IncrementIP()
			term2 = uint32(core.readImm16())
			term1 = uint32(core.registers.AX)
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			core.registers.AX = uint16(term1)

			log.Printf("[%#04x] adc ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x80:
		{
			// adc r/m8,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x81:
		{
			// adc r/m16,imm16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm16())
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x83:
		{
			// adc r/m16,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x10:
		{
			// adc r/m8,r8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x11:
		{
			// adc r/m16,r16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x12:
		{
			// adc r8,r/m8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x13:
		{
			// adc r16,r/m16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint16(result)
			core.writeR16(&modrm, &tmp)

			log.Printf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.SetFlag(CarryFlag, result>>(bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr != 0)

	core.registers.SetFlag(OverFlowFlag,  (sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0))


}


func INSTR_ADD(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x04:
		{
			// 	add AL,imm8
			core.IncrementIP()
			term2 = uint32(core.readImm8())
			term1 = uint32(core.registers.AL)
			result = uint32(term1) + uint32(term2)
			core.registers.AL = uint8(term1)

			log.Printf("[%#04x] add al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x05:
		{
			// 		add AX,imm16
			core.IncrementIP()
			term2 = uint32(core.readImm16())
			term1 = uint32(core.registers.AX)
			result = uint32(term1) + uint32(term2)
			core.registers.AX = uint16(term1)

			log.Printf("[%#04x] add ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x80:
		{
			// add r/m8,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) + uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x81:
		{
			// add r/m16,imm16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm16())
			result = uint32(term1) + uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x83:
		{
			// add r/m16,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) + uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x00:
		{
			// add r/m8,r8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x01:
		{
			// add r/m16,r16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x02:
		{
			// add r8,r/m8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x03:
		{
			// add r16,r/m16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) + uint32(term2)
			tmp := uint16(result)
			core.writeR16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.SetFlag(CarryFlag, result>>(bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr != 0)

	core.registers.SetFlag(OverFlowFlag,  (sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0))

}

func INSTR_AND(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x24:
		{
			// 	and AL,imm8
			core.IncrementIP()
			term2 = uint32(core.readImm8())
			term1 = uint32(core.registers.AL)
			result = uint32(term1) & uint32(term2)
			core.registers.AL = uint8(term1)

			log.Printf("[%#04x] add al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x25:
		{
			// 	and AX,imm16
			core.IncrementIP()
			term2 = uint32(core.readImm16())
			term1 = uint32(core.registers.AX)
			result = uint32(term1) & uint32(term2)
			core.registers.AX = uint16(term1)

			log.Printf("[%#04x] add ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x80:
		{
			// and r/m8,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) & uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x81:
		{
			// and r/m16,imm16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm16())
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x83:
		{
			// and r/m16,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x20:
		{
			// and r/m8,r8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x21:
		{
			// and r/m16,r16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x22:
		{
			// add r8,r/m8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x23:
		{
			// add r16,r/m16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			core.writeR16(&modrm, &tmp)

			log.Printf("[%#04x] add %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.SetFlag(CarryFlag, result>>(bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr != 0)

	core.registers.SetFlag(OverFlowFlag,  (sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0))


}


func INSTR_OR(core *CpuCore) {
	core.IncrementIP()

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x0c:
		{
			// OR AL,imm8
			term2 = uint32(core.readImm8())
			term1 = uint32(core.registers.AL)
			result = uint32(term1) | uint32(term2)
			core.registers.AL = uint8(result)

			log.Printf("[%#04x] or al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x0d:
		{
			// OR AX,imm16
			term2 = uint32(core.readImm16())
			term1 = uint32(core.registers.AX)
			result = uint32(term1) | uint32(term2)
			core.registers.AX = uint16(result)

			log.Printf("[%#04x] or ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x80:
		{
			// OR r/m8,imm8
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm8(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm8())
			tmp := uint8(uint32(term1) | uint32(term2))

			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x81:
		{
			// OR r/m16,imm16
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm16())
			tmp := uint16(uint32(term1) | uint32(term2))

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x83:
		{
			// OR r/m16,imm8
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm8())
			tmp := uint16(uint32(term1) | uint32(term2))

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x08:
		{
			// OR r/m8,r8
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm8(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR8(&modrm)
			term2 = uint32(*rm2)
			tmp := uint8(uint32(term1) | uint32(term2))

			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x09:
		{
			// OR r/m16,r16
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR16(&modrm)
			term2 = uint32(*rm2)
			tmp := uint16(uint32(term1) | uint32(term2))

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x0A:
		{
			// OR r8,r/m8
			modrm := core.consumeModRm()
			rm, rmStr := core.readR8(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readRm8(&modrm)
			term2 = uint32(*rm2)
			tmp := uint8(uint32(term1) | uint32(term2))

			core.writeR8(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x0B:
		{
			// OR r16,r/m16
			modrm := core.consumeModRm()
			rm, rmStr := core.readR16(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readRm16(&modrm)
			term2 = uint32(*rm2)
			tmp := uint16(uint32(term1) | uint32(term2))

			core.writeR16(&modrm, &tmp)

			log.Printf("[%#04x] or %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags

	core.registers.SetFlag(OverFlowFlag,  false)

	core.registers.SetFlag(CarryFlag, false)

	core.registers.SetFlag(SignFlag, (result >> bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(ParityFlag, bits.OnesCount32(result)%2 == 0)

}

func INSTR_XOR(core *CpuCore) {
	core.IncrementIP()

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x34:
		{
			// XOR AL,imm8
			term2 = uint32(core.readImm8())
			term1 = uint32(core.registers.AL)
			result = uint32(term1) ^ uint32(term2)
			core.registers.AL = uint8(result)

			log.Printf("[%#04x] xor al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x35:
		{
			// XOR AX,imm16
			term2 = uint32(core.readImm16())
			term1 = uint32(core.registers.AX)
			result = uint32(term1) ^ uint32(term2)
			core.registers.AX = uint16(result)

			log.Printf("[%#04x] xor ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x80:
		{
			// XOR r/m8,imm8
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm8(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm8())
			tmp := uint8(uint32(term1) ^ uint32(term2))

			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x81:
		{
			// XOR r/m16,imm16
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm16())
			tmp := uint16(uint32(term1) ^ uint32(term2))

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x83:
		{
			// XOR r/m16,imm8
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			term2 = uint32(core.readImm8())
			tmp := uint16(uint32(term1) ^ uint32(term2))

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, term2)
		}
	case 0x30:
		{
			// XOR r/m8,r8
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm8(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR8(&modrm)
			term2 = uint32(*rm2)
			tmp := uint8(uint32(term1) ^ uint32(term2))

			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x31:
		{
			// XOR r/m16,r16
			modrm := core.consumeModRm()
			rm, rmStr := core.readRm16(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR16(&modrm)
			term2 = uint32(*rm2)
			tmp := uint16(uint32(term1) ^ uint32(term2))

			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x32:
		{
			// XOR r8,r/m8
			modrm := core.consumeModRm()
			rm, rmStr := core.readR8(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readRm8(&modrm)
			term2 = uint32(*rm2)
			tmp := uint8(uint32(term1) ^ uint32(term2))

			core.writeR8(&modrm, &tmp)

			log.Printf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	case 0x33:
		{
			// XOR r16,r/m16
			modrm := core.consumeModRm()
			rm, rmStr := core.readR16(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str := core.readRm16(&modrm)
			term2 = uint32(*rm2)
			tmp := uint16(uint32(term1) ^ uint32(term2))

			core.writeR16(&modrm, &tmp)

			log.Printf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), rmStr, rm2Str)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags

	core.registers.SetFlag(OverFlowFlag,  false)

	core.registers.SetFlag(CarryFlag, false)

	core.registers.SetFlag(SignFlag, (result >> bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(ParityFlag, bits.OnesCount32(result)%2 == 0)

}


func INSTR_SUB(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentByteAtCodePointer {
	case 0x2c:
		{
			// 	SUB AL,imm8
			core.IncrementIP()
			term2 = uint32(core.readImm8())
			term1 = uint32(core.registers.AL)
			result = uint32(term1) - uint32(term2)
			core.registers.AL = uint8(term1)

			log.Printf("[%#04x] sub al, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x2d:
		{
			// 		SUB AX,imm16
			core.IncrementIP()
			term2 = uint32(core.readImm16())
			term1 = uint32(core.registers.AX)
			result = uint32(term1) - uint32(term2)
			core.registers.AX = uint16(term1)

			log.Printf("[%#04x] sub ax, %#04x", core.GetCurrentlyExecutingInstructionPointer(), term2)
		}
	case 0x80:
		{
			// SUB r/m8,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) - uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] sub %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x81:
		{
			// SUB r/m16,imm16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm16())
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] sub %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x83:
		{
			// SUB r/m16,imm8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			term2 = uint32(core.readImm8())
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] sub %s, %#04x", core.GetCurrentlyExecutingInstructionPointer(), t1Name, term2)
		}
	case 0x28:
		{
			// SUB r/m8,r8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x29:
		{
			// SUB r/m16,r16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readRm16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readR16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			core.writeRm16(&modrm, &tmp)

			log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x2A:
		{
			// SUB r8,r/m8
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint8(result)
			core.writeRm8(&modrm, &tmp)

			log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	case 0x2B:
		{
			// SUB r16,r/m16
			core.IncrementIP()
			modrm := core.consumeModRm()
			t1, t1Name := core.readR16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name := core.readRm16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			core.writeR16(&modrm, &tmp)

			log.Printf("[%#04x] sub %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Name, t2Name)
		}
	}

	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 := (term1 >> (bitLength)) & 0x01
	sign2 := (term2 >> (bitLength)) & 0x01
	signr := uint8((result >> (bitLength)) & 0x01)

	core.registers.SetFlag(CarryFlag, result>>(bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr != 0)

	core.registers.SetFlag(OverFlowFlag,  (sign1 == 0 && sign2 == 1 && signr == 1) || (sign1 == 1 && sign2 == 0 && signr == 0))


}

func INSTR_SHIFT(core *CpuCore) {

	// SAL
	var destTerm interface{}
	var bitLength uint8
	var t1Str string
	var t2Str string
	var opCode uint8
	var countTerm uint8

	modrm := core.consumeModRm()
	opCode = core.currentByteAtCodePointer

	switch modrm.mod {
	case 4:
		{
			// sal
			switch opCode {
			// 8 bit versions
			case 0xd0: destTerm, t1Str = core.readRm8(&modrm); countTerm = 1; t2Str = "1"; bitLength = 8
			case 0xd2: destTerm, t1Str = core.readRm8(&modrm); countTerm = core.registers.CL; t2Str = "CL"; bitLength = 8
			case 0xC0: destTerm, t1Str = core.readRm8(&modrm); countTerm = core.readImm8(); t2Str = string(countTerm); bitLength = 8

			// 16 bit versions
			case 0xd1: destTerm, t1Str = core.readRm16(&modrm); countTerm = 1; t2Str = "1"; bitLength = 16
			case 0xd3: destTerm, t1Str = core.readRm16(&modrm); countTerm = core.registers.CL; t2Str = "CL"; bitLength = 16
			case 0xC1: destTerm, t1Str = core.readRm16(&modrm); countTerm = core.readImm8(); t2Str = string(countTerm); bitLength = 16

			default: {
				log.Printf("Unhandled shift sal variant")
				doCoreDump(core)
				panic(false)
			}
			}

			for i:=uint8(0);i<countTerm;i++ {
				if bitLength == 8 {

					if *destTerm.(*uint8) << 1 > (0 & 11111111) {
						// shift into carry
						core.registers.SetFlag(CarryFlag,  *destTerm.(*uint8) >> (bitLength-1) & 1 == 1)
					}

					*destTerm.(*uint8) <<= 1
				}

				if bitLength == 16 {

					if *destTerm.(*uint16) << 1 > (0 & 11111111) {
						// shift into carry
						core.registers.SetFlag(CarryFlag,  *destTerm.(*uint16) >> (bitLength-1) & 1 == 1)
					}

					*destTerm.(*uint16) <<= 1
				}
			}

			log.Printf("[%#04x] sal %s, %s", core.GetCurrentlyExecutingInstructionPointer(), t1Str, t2Str)

		}
	case 5, 7:
		{
			// shr - shifts to the right, and clears the most significant bit (unsigned)
			// sar - shifts to the right, preserves the most significant bit (signed)
			switch opCode {
			// 8 bit versions
			case 0xd0: destTerm, t1Str = core.readRm8(&modrm); countTerm = 1; t2Str = "1"; bitLength = 8
			case 0xd2: destTerm, t1Str = core.readRm8(&modrm); countTerm = core.registers.CL; t2Str = "CL"; bitLength = 8
			case 0xC0: destTerm, t1Str = core.readRm8(&modrm); countTerm = core.readImm8(); t2Str = string(countTerm); bitLength = 8

			// 16 bit versions
			case 0xd1: destTerm, t1Str = core.readRm16(&modrm); countTerm = 1; t2Str = "1"; bitLength = 16
			case 0xd3: destTerm, t1Str = core.readRm16(&modrm); countTerm = core.registers.CL; t2Str = "CL"; bitLength = 16
			case 0xC1: destTerm, t1Str = core.readRm16(&modrm); countTerm = core.readImm8(); t2Str = string(countTerm); bitLength = 16

			default: {
				log.Printf("Unhandled shift sal variant")
				doCoreDump(core)
				panic(false)
			}
			}

			var msbBit interface{}
			for i:=uint8(0);i<countTerm;i++ {

				if bitLength == 8 {

					core.registers.SetFlag(CarryFlag,  *destTerm.(*uint8) & 1 == 1)

					msbBit = (*destTerm.(*uint8) >> 8) & 1
					*destTerm.(*uint8) >>= 1
					if modrm.mod == 7 {
						*destTerm.(*uint8) = *destTerm.(*uint8) | (msbBit.(uint8) << (bitLength-1))
					}
				}

				if bitLength == 16 {

					core.registers.SetFlag(CarryFlag,  *destTerm.(*uint16) & 1 == 1)

					msbBit = (*destTerm.(*uint16) >> 8) & 1
					*destTerm.(*uint8) >>= 1
					if modrm.mod == 7 {
						*destTerm.(*uint16) = *destTerm.(*uint16) | (msbBit.(uint16) << (bitLength-1))
					}
				}
			}

			opCode := "SHR"
			if modrm.mod == 7 {
				opCode = "SAR"
			}

			log.Printf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionPointer(), opCode, t1Str, t2Str)

		}
	}

}


