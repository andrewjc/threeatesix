package intel8086

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"log"
	"math/bits"
)

func INSTR_ADC(core *CpuCore) {
	var term1, term2, result uint32
	var bitLength uint32
	var err error
	var modrm ModRm
	var bytesConsumed uint32
	var t1 *uint8
	var t1w *uint16
	var t2 *uint8
	var t2w *uint16
	var t1Name, t2Name string

	readOperands := func() {
		core.currentByteAddr++
		modrm, bytesConsumed, err = core.consumeModRm()
		if err != nil {
			log.Fatalf("Error reading operands for adc instruction: %v", err)
		}
		core.currentByteAddr += bytesConsumed
	}

	updateFlags := func() {
		bitLength = uint32(bits.Len32(result))
		core.registers.SetFlag(CarryFlag, result>>(bitLength) != 0)
		core.registers.SetFlag(ZeroFlag, result == 0)
		core.registers.SetFlag(SignFlag, (result>>(bitLength-1))&1 != 0)
		core.registers.SetFlag(OverFlowFlag, ((term1^term2)&0x80000000) == 0 && ((term1^result)&0x80000000) != 0)
	}

	switch core.currentOpCodeBeingExecuted {
	case 0x14, 0x15:
		// ADC AL, imm8 or ADC AX, imm16
		core.currentByteAddr++
		if core.currentOpCodeBeingExecuted == 0x14 {
			term1 = uint32(core.registers.AL)
			if imm8, err := core.readImm8(); err != nil {
				goto eof
			} else {
				term2 = uint32(imm8)
			}
			result = term1 + term2 + uint32(core.registers.GetFlagInt(CarryFlag))
			core.registers.AL = uint8(result)
			core.logInstruction(fmt.Sprintf("[%#04x] adc al, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term2))
		} else {
			term1 = uint32(core.registers.AX)
			if imm16, err := core.readImm16(); err != nil {
				goto eof
			} else {
				term2 = uint32(imm16)
			}
			result = term1 + term2 + uint32(core.registers.GetFlagInt(CarryFlag))
			core.registers.AX = uint16(result)
			core.logInstruction(fmt.Sprintf("[%#04x] adc ax, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term2))
		}
	case 0x80, 0x81, 0x83:
		// ADC r/m8, imm8 or ADC r/m16, imm16 or ADC r/m16, imm8
		readOperands()
		if core.currentOpCodeBeingExecuted == 0x80 {
			t1, t1Name, err = core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			term1 = uint32(*t1)
			if imm8, err := core.readImm8(); err != nil {
				goto eof
			} else {
				term2 = uint32(imm8)
			}
			result = term1 + term2 + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint8(result)
			if _, err = core.writeRm8(&modrm, &tmp); err != nil {
				goto eof
			}
			core.logInstruction(fmt.Sprintf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), t1Name, term2))
		} else {
			t1w, t1Name, err = core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			term1 = uint32(*t1w)
			if imm8, err := core.readImm8(); core.currentOpCodeBeingExecuted == 0x83 && err == nil {
				term2 = uint32(int8(imm8)) // sign-extend imm8
			} else if err != nil {
				if imm16, err := core.readImm16(); err != nil {
					goto eof
				} else {
					term2 = uint32(imm16)
				}
			}
			result = term1 + term2 + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint16(result)
			if _, err = core.writeRm16(&modrm, &tmp); err != nil {
				goto eof
			}
			core.logInstruction(fmt.Sprintf("[%#04x] adc %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), t1Name, term2))
		}
	case 0x10, 0x11, 0x12, 0x13:
		// ADC r/m8, r8 or ADC r/m16, r16 or ADC r8, r/m8 or ADC r16, r/m16
		readOperands()
		if core.currentOpCodeBeingExecuted == 0x10 || core.currentOpCodeBeingExecuted == 0x12 {
			t1, t1Name, err = core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			term1 = uint32(*t1)
			t2, t2Name = core.readR8(&modrm)
			term2 = uint32(*t2)
			result = term1 + term2 + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint8(result)
			if core.currentOpCodeBeingExecuted == 0x10 {
				if _, err = core.writeRm8(&modrm, &tmp); err != nil {
					goto eof
				}
			} else {
				core.writeR8(&modrm, &tmp)
			}
			core.logInstruction(fmt.Sprintf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
		} else {
			t1w, t1Name, err = core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			term1 = uint32(*t1w)
			t2w, t2Name = core.readR16(&modrm)
			term2 = uint32(*t2w)
			result = term1 + term2 + uint32(core.registers.GetFlagInt(CarryFlag))
			tmp := uint16(result)
			if core.currentOpCodeBeingExecuted == 0x11 {
				if _, err = core.writeRm16(&modrm, &tmp); err != nil {
					goto eof
				}
			} else {
				core.writeR16(&modrm, &tmp)
			}
			core.logInstruction(fmt.Sprintf("[%#04x] adc %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
		}
	default:
		log.Fatalf("[%#04x] Unimplemented adc instruction: %#04x", core.GetCurrentlyExecutingInstructionAddress(), core.currentOpCodeBeingExecuted)
	}

	updateFlags()

eof:
}

func INSTR_ADD(core *CpuCore) {
	var term1, term2, result uint32
	var size uint8
	var srcName, destName string
	var err error
	var mask uint32

	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0x04: // ADD AL, imm8
		term1 = uint32(core.registers.AL)
		term2, err = common.UInt32From8(core.readImm8())
		if err != nil {
			goto eof
		}
		result = term1 + term2
		core.registers.AL = uint8(result)
		size = 8
		destName = "AL"
		srcName = fmt.Sprintf("#%02x", term2)

	case 0x05: // ADD AX, imm16
		term1 = uint32(core.registers.AX)
		term2, err = common.UInt32From16(core.readImm16())
		if err != nil {
			goto eof
		}
		result = term1 + term2
		core.registers.AX = uint16(result)
		size = 16
		destName = "AX"
		srcName = fmt.Sprintf("#%04x", term2)

	case 0x80: // ADD r/m8, imm8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			goto eof
		}
		core.currentByteAddr += bytesConsumed

		t1, t1Name, err := core.readRm8(&modrm)
		if err != nil {
			goto eof
		}
		term1 = uint32(*t1)
		term2, err = common.UInt32From8(core.readImm8())
		if err != nil {
			goto eof
		}
		result = term1 + term2
		tmp := uint8(result)
		_, err = core.writeRm8(&modrm, &tmp)
		size = 8
		srcName = fmt.Sprintf("#%02x", term2)
		destName = t1Name

	case 0x81: // ADD r/m16, imm16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			goto eof
		}
		core.currentByteAddr += bytesConsumed

		t1, t1Name, err := core.readRm16(&modrm)
		if err != nil {
			goto eof
		}
		term1 = uint32(*t1)
		term2, err = common.UInt32From16(core.readImm16())
		if err != nil {
			goto eof
		}
		result = term1 + term2
		tmp := uint16(result)
		_, err = core.writeRm16(&modrm, &tmp)
		size = 16
		srcName = fmt.Sprintf("#%04x", term2)
		destName = t1Name

	case 0x83: // ADD r/m16, imm8 (sign-extended)
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			goto eof
		}
		core.currentByteAddr += bytesConsumed

		t1, t1Name, err := core.readRm16(&modrm)
		if err != nil {
			goto eof
		}
		term1 = uint32(*t1)
		imm8, err := core.readImm8()
		if err != nil {
			goto eof
		}
		term2 = uint32(int16(int8(imm8))) // Sign extend
		result = term1 + term2
		tmp := uint16(result)
		_, err = core.writeRm16(&modrm, &tmp)
		if err != nil {
			goto eof
		}
		size = 16
		srcName = fmt.Sprintf("#%02x", imm8)
		destName = t1Name

	case 0x00: // ADD r/m8, r8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			goto eof
		}
		core.currentByteAddr += bytesConsumed

		t1, t1Name, err := core.readRm8(&modrm)
		if err != nil {
			goto eof
		}
		t2, t2Name := core.readR8(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(*t2)
		tmp := uint8(term1 + term2)
		_, err = core.writeRm8(&modrm, &tmp)
		if err != nil {
			goto eof
		}
		size = 8
		destName = t1Name
		srcName = t2Name

	case 0x01: // ADD r/m16, r16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			goto eof
		}
		core.currentByteAddr += bytesConsumed

		t1, t1Name, err := core.readRm16(&modrm)
		if err != nil {
			goto eof
		}
		t2, t2Name := core.readR16(&modrm)
		term1 = uint32(*t1)
		term2 = uint32(*t2)
		result = term1 + term2
		tmp := uint16(result)
		_, err = core.writeRm16(&modrm, &tmp)
		if err != nil {
			goto eof
		}
		size = 16
		destName = t1Name
		srcName = t2Name

	case 0x02: // ADD r8, r/m8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			goto eof
		}
		core.currentByteAddr += bytesConsumed

		t1, t1Name := core.readR8(&modrm)
		t2, t2Name, err := core.readRm8(&modrm)
		if err != nil {
			goto eof
		}
		term1 = uint32(*t1)
		term2 = uint32(*t2)
		tmp := uint8(term1 + term2)
		_, err = core.writeRm8(&modrm, &tmp)
		if err != nil {
			goto eof
		}
		size = 8
		destName = t1Name
		srcName = t2Name

	case 0x03: // ADD r16, r/m16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			goto eof
		}
		core.currentByteAddr += bytesConsumed

		t1, t1Name := core.readR16(&modrm)
		t2, t2Name, err := core.readRm16(&modrm)
		if err != nil {
			goto eof
		}
		term1 = uint32(*t1)
		term2 = uint32(*t2)
		result = term1 + term2
		tmp := uint16(result)
		core.writeR16(&modrm, &tmp)
		size = 16
		destName = t1Name
		srcName = t2Name

	default:
		log.Fatalf("[%#04x] Unimplemented add instruction: %#04x", core.GetCurrentlyExecutingInstructionAddress(), core.currentOpCodeBeingExecuted)
	}

	// Update flags
	mask = uint32(1<<size - 1)
	core.registers.SetFlag(CarryFlag, result > mask)
	core.registers.SetFlag(ZeroFlag, (result&mask) == 0)
	core.registers.SetFlag(SignFlag, (result&(1<<(size-1))) != 0)
	core.registers.SetFlag(OverFlowFlag, ((term1^term2)&(1<<(size-1)) == 0) && ((term1^result)&(1<<(size-1)) != 0))

	// Update auxiliary carry flag (used for BCD arithmetic)
	core.registers.SetFlag(AdjustFlag, ((term1&0xF)+(term2&0xF)) > 0xF)

	core.logInstruction(fmt.Sprintf("[%#04x] ADD %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

eof:
}

func INSTR_AND(core *CpuCore) {

	var srcName string
	var term1 uint32
	var term2 uint32
	var result uint32

	var signr int16
	var sign1 int16
	var sign2 int16

	var dataSize uint32

	switch core.currentOpCodeBeingExecuted {
	case 0x24:
		{
			// 	and AL,imm8
			core.currentByteAddr++
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			term1 = uint32(core.registers.AL)
			result = uint32(term1) & uint32(term2)
			core.registers.AL = uint8(result)

			dataSize = 8

			core.logInstruction(fmt.Sprintf("[%#04x] and al, %#08x", core.GetCurrentlyExecutingInstructionAddress(), term2))
			goto success
		}
	case 0x25:
		{
			// and AX,imm16
			core.currentByteAddr++

			var result uint32
			var term1 uint32
			var term1_name string
			var dataSize uint8

			term1, term1_name, dataSize = core.GetRegister16(&core.registers.AX)
			term2, _, err := core.GetImmediate16()
			if err != nil {
				goto eof
			}
			result = term1 & term2
			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term1_name, term2))

			if dataSize == 32 {
				core.registers.EAX = result
			} else {
				core.registers.AX = uint16(result)
			}

			goto success
		}
	case 0x80:
		{
			// and r/m8,imm8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			result = uint32(term1) & uint32(term2)
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			dataSize = 8

			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), t1Name, term2))
			goto success
		}
	case 0x81:
		{
			// and r/m16,imm16
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			srcName, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			dataSize = 16
			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, srcName))
			goto success
		}
	case 0x83:
		{
			// and r/m16,imm8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			srcName, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			dataSize = 16
			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, srcName))
			goto success
		}
	case 0x20:
		{
			// and r/m8,r8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			t2, t2Name := core.readR8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			dataSize = 8
			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	case 0x21:
		{
			// and r/m16,r16
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			t2, t2Name := core.readR16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			_, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			dataSize = 16
			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	case 0x22:
		{
			// add r8,r/m8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name := core.readR8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			dataSize = 8
			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	case 0x23:
		{
			// add r16,r/m16
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name := core.readR16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term2 = uint32(*t2)
			result = uint32(term1) & uint32(term2)
			tmp := uint16(result)
			core.writeR16(&modrm, &tmp)
			dataSize = 16
			core.logInstruction(fmt.Sprintf("[%#04x] and %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	default:
		log.Fatalf("[%#04x] Unimplemented and instruction: %#04x", core.GetCurrentlyExecutingInstructionAddress(), core.currentOpCodeBeingExecuted)
	}

success:

	// update flags
	sign1 = int16(term1 >> (dataSize - 1))
	sign2 = int16(term2 >> (dataSize - 1))
	signr = int16((result >> (dataSize - 1)) & 0x01)

	core.registers.SetFlag(CarryFlag, (result>>dataSize) == 1)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr == 1)

	core.registers.SetFlag(OverFlowFlag, (sign1^sign2 == 0) && (sign1^signr == 1))

eof:
}

func INSTR_OR(core *CpuCore) {
	var term1, term2, result uint32
	var dataSize uint8 = 16
	var srcName, dstName string
	var err error

	var signr int16
	var sign1 int16
	var sign2 int16

	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0x0c: // OR AL, imm8
		imm8, err := core.readImm8()
		if err != nil {
			return
		}
		term1 = uint32(core.registers.AL)
		term2 = uint32(imm8)
		result = term1 | term2
		core.registers.AL = uint8(result)
		dataSize = 8
		dstName = "AL"
		srcName = fmt.Sprintf("#%#04x", imm8)

	case 0x0d: // OR AX, imm16
		term1, dstName, dataSize = core.GetRegister16(&core.registers.AX)
		term2, _, err = core.GetImmediate16()
		if err != nil {
			return
		}

		result = term1 | term2
		if dataSize == 32 {
			core.registers.EAX = result
			srcName = fmt.Sprintf("#%#04x", term2)
		} else {
			core.registers.AX = uint16(result)
			srcName = fmt.Sprintf("#%#04x", term2)
		}

	case 0x80, 0x82: // OR r/m8, imm8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		rm, _, err := core.readRm8(&modrm)
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed
		term1 = uint32(*rm)
		imm8, err := core.readImm8()
		if err != nil {
			return
		}
		term2 = uint32(imm8)
		result = term1 | term2
		*rm = uint8(result)
		dataSize = 8
		dstName = modrm.String()
		srcName = fmt.Sprintf("#%#04x", imm8)

	case 0x81, 0x83: // OR r/m16, imm16 or imm8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		rm, rmName, err := core.readRm16(&modrm)
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed
		term1 = uint32(*rm)
		if core.currentOpCodeBeingExecuted == 0x83 {
			imm8, err := core.readImm8()
			if err != nil {
				return
			}
			term2 = uint32(imm8) // Zero-extend the 8-bit immediate to 16 bits
		} else {
			imm16, err := core.readImm16()
			if err != nil {
				return
			}
			term2 = uint32(imm16)
		}
		result = term1 | term2
		*rm = uint16(result)
		dataSize = 16
		dstName = rmName
		srcName = fmt.Sprintf("#%#04x", term2)

	case 0x08: // OR r/m8, r8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		rm, rmName, err := core.readRm8(&modrm)
		if err != nil {
			return
		}
		r8, rName := core.readR8(&modrm)

		core.currentByteAddr += bytesConsumed
		term1 = uint32(*rm)
		term2 = uint32(*r8)
		result = term1 | term2
		*rm = uint8(result)
		dataSize = 8
		dstName = rmName
		srcName = rName

	case 0x09: // OR r/m16, r16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		rm, rmName, err := core.readRm16(&modrm)
		if err != nil {
			return
		}
		r16, rName := core.readR16(&modrm)
		core.currentByteAddr += bytesConsumed
		term1 = uint32(*rm)
		term2 = uint32(*r16)
		result = term1 | term2
		*rm = uint16(result)
		dataSize = 16
		dstName = rmName
		srcName = rName

	case 0x0A: // OR r8, r/m8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		r8, rName := core.readR8(&modrm)

		rm, rmName, err := core.readRm8(&modrm)
		if err != nil {

		}

		core.currentByteAddr += bytesConsumed
		term1 = uint32(*r8)
		term2 = uint32(*rm)
		result = term1 | term2
		*r8 = uint8(result)
		dataSize = 8
		srcName = rmName
		dstName = rName

	case 0x0B: // OR r16, r/m16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			return
		}
		r16, rName := core.readR16(&modrm)

		rm, rmName, err := core.readRm16(&modrm)
		if err != nil {
			return
		}
		core.currentByteAddr += bytesConsumed
		term1 = uint32(*r16)
		term2 = uint32(*rm)
		result = term1 | term2
		*r16 = uint16(result)
		dataSize = 16
		srcName = rmName
		dstName = rName

	default:
		log.Fatalf("[%#04x] Unimplemented or instruction: %#04x", core.GetCurrentlyExecutingInstructionAddress(), core.currentOpCodeBeingExecuted)
	}

	// Update flags
	sign1 = int16(term1 >> (dataSize - 1))
	sign2 = int16(term2 >> (dataSize - 1))
	signr = int16((result >> (dataSize - 1)) & 0x01)

	core.registers.SetFlag(CarryFlag, (result>>dataSize) == 1)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr == 1)

	core.registers.SetFlag(OverFlowFlag, (sign1^sign2 == 0) && (sign1^signr == 1))

	// Increment instruction pointer
	core.logInstruction(fmt.Sprintf("[%#04x] OR %s, %s", core.GetCurrentlyExecutingInstructionAddress(), dstName, srcName))

}

func INSTR_XOR(core *CpuCore) {

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentOpCodeBeingExecuted {
	case 0x34:
		{
			// XOR AL,imm8
			core.currentByteAddr++
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}

			term1 = uint32(core.registers.AL)
			result = uint32(term1) ^ uint32(term2)
			core.registers.AL = uint8(result)

			core.logInstruction(fmt.Sprintf("[%#04x] xor al, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term2))
			goto success
		}
	case 0x35:
		{
			// XOR AX,imm16
			core.currentByteAddr++
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}

			term1 = uint32(core.registers.AX)
			result = uint32(term1) ^ uint32(term2)
			core.registers.AX = uint16(result)

			core.logInstruction(fmt.Sprintf("[%#04x] xor ax, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term2))
			goto success
		}
	case 0x80:
	case 0x82:
		{
			// XOR r/m8,imm8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm, rmStr, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			tmp := uint8(uint32(term1) ^ uint32(term2))

			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), rmStr, term2))
			goto success
		}
	case 0x81:
		{
			// XOR r/m16,imm16
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm, rmStr, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm)
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}
			tmp := uint16(uint32(term1) ^ uint32(term2))

			srcName, err := core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), rmStr, srcName))
			goto success
		}
	case 0x83:
		{
			// XOR r/m16,imm8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm, rmStr, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			tmp := uint16(uint32(term1) ^ uint32(term2))

			srcName, err := core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] xor %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), rmStr, srcName))
			goto success
		}
	case 0x30:
		{
			// XOR r/m8,r8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm, rmStr, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR8(&modrm)
			term2 = uint32(*rm2)
			tmp := uint8(uint32(term1) ^ uint32(term2))

			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] xor %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rmStr, rm2Str))
			goto success
		}
	case 0x31:
		{
			// XOR r/m16,r16
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm, rmStr, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*rm)
			rm2, rm2Str := core.readR16(&modrm)
			term2 = uint32(*rm2)
			tmp := uint16(uint32(term1) ^ uint32(term2))

			_, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			core.logInstruction(fmt.Sprintf("[%#04x] xor %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rmStr, rm2Str))
			goto success
		}
	case 0x32:
		{
			// XOR r8,r/m8
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm, rmStr := core.readR8(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term2 = uint32(*rm2)
			tmp := uint8(uint32(term1) ^ uint32(term2))

			core.writeR8(&modrm, &tmp)

			core.logInstruction(fmt.Sprintf("[%#04x] xor %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rmStr, rm2Str))
			goto success
		}
	case 0x33:
		{
			// XOR r16,r/m16
			core.currentByteAddr++
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			rm, rmStr := core.readR16(&modrm)
			term1 = uint32(*rm)
			rm2, rm2Str, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term2 = uint32(*rm2)
			tmp := uint16(uint32(term1) ^ uint32(term2))

			core.writeR16(&modrm, &tmp)

			core.logInstruction(fmt.Sprintf("[%#04x] xor %s, %s", core.GetCurrentlyExecutingInstructionAddress(), rmStr, rm2Str))
			goto success
		}
	default:
		log.Fatalf("[%#04x] Unimplemented xor instruction: %#04x", core.GetCurrentlyExecutingInstructionAddress(), core.currentOpCodeBeingExecuted)
	}

success:
	bitLength = uint32(bits.Len32(result))

	// update flags

	core.registers.SetFlag(OverFlowFlag, false)

	core.registers.SetFlag(CarryFlag, false)

	core.registers.SetFlag(SignFlag, (result>>bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(ParityFlag, bits.OnesCount32(result)%2 == 0)
eof:
}

func INSTR_SUB(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var term2 uint32
	var result uint32

	var size uint8
	var mask uint32

	var destName string
	var srcName string

	switch core.currentOpCodeBeingExecuted {
	case 0x2c:
		{
			// 	SUB AL,imm8
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			term1 = uint32(core.registers.AL)
			result = uint32(term1) - uint32(term2)
			core.registers.AL = uint8(result)
			size = 8

			srcName = fmt.Sprintf("#%#04x", term2)
			destName = "AL"

			goto success
		}
	case 0x2d:
		{
			// 		SUB AX,imm16
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}
			term1 = uint32(core.registers.AX)
			result = uint32(term1) - uint32(term2)
			core.registers.AX = uint16(term1)
			size = 16

			srcName = fmt.Sprintf("#%#04x", term2)
			destName = "AX"

			goto success
		}
	case 0x80:
		{
			// SUB r/m8,imm8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			result = uint32(term1) - uint32(term2)
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			size = 8

			srcName = fmt.Sprintf("#%#04x", term2)
			destName = t1Name

			goto success
		}
	case 0x81:
		{
			// SUB r/m16,imm16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, _, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			destName, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			size = 16

			srcName = fmt.Sprintf("#%#04x", term2)

			goto success
		}
	case 0x83:
		{
			// SUB r/m16,imm8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, _, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			destName, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			size = 16

			srcName = fmt.Sprintf("#%#04x", term2)

			goto success
		}
	case 0x28:
		{
			// SUB r/m8,r8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, _, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			t2, t2Name := core.readR8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint8(result)
			destName, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			size = 8

			srcName = t2Name

			goto success
		}
	case 0x29:
		{
			// SUB r/m16,r16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, _, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			t2, t2Name := core.readR16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			destName, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			size = 16

			srcName = t2Name

			goto success
		}
	case 0x2A:
		{
			// SUB r8,r/m8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, _ := core.readR8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint8(result)
			destName, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}
			size = 8

			srcName = t2Name

			goto success
		}
	case 0x2B:
		{
			// SUB r16,r/m16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name := core.readR16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term2 = uint32(*t2)
			result = uint32(term1) - uint32(term2)
			tmp := uint16(result)
			core.writeR16(&modrm, &tmp)
			size = 16

			srcName = t2Name
			destName = t1Name

			goto success
		}
	default:
		log.Fatal(fmt.Sprintf("Unrecognised SUB instruction: %#04x!", core.currentOpCodeBeingExecuted))
	}
success:

	// Create a mask for the operand size
	mask = uint32(1<<size - 1)

	// Update flags
	core.registers.SetFlag(CarryFlag, term2 > term1)
	core.registers.SetFlag(ZeroFlag, (result&mask) == 0)
	core.registers.SetFlag(SignFlag, (result&(1<<(size-1))) != 0)

	// Overflow occurs if the sign of the two operands is different and
	// the sign of the result is different from the sign of the first operand
	core.registers.SetFlag(OverFlowFlag,
		((term1^term2)&(term1^result)&(1<<(size-1))) != 0)

	// Update auxiliary carry flag (used for BCD arithmetic)
	core.registers.SetFlag(AdjustFlag, ((term1&0xF)-(term2&0xF))&0x10 != 0)

	core.logInstruction(fmt.Sprintf("[%#04x] SUB %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

eof:
}

func INSTR_INC(core *CpuCore) {
	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47:
		{
			// INC r16
			val, valName := core.registers.registers16Bit[core.currentOpCodeBeingExecuted-0x40], core.registers.index16ToString(core.currentOpCodeBeingExecuted-0x40)

			*val = *val + 1

			core.registers.SetFlag(ZeroFlag, *val == 0)

			core.logInstruction(fmt.Sprintf("[%#04x] inc %s", core.GetCurrentlyExecutingInstructionAddress(), valName))

		}
	default:
		log.Println(fmt.Printf("Unhandled INC instruction:  %#04x", core.currentOpCodeBeingExecuted))
		doCoreDump(core)
	}

}

func INSTR_INC_RM16(core *CpuCore) {
	var dest *uint16
	var destName string

	core.currentByteAddr++

	modrm, bytesConsumed, err := core.consumeModRm()
	if err != nil {
		goto eof
	}
	core.currentByteAddr += bytesConsumed

	dest, destName, err = core.readRm16(&modrm)
	if err != nil {
		goto eof
	}

	*dest = *dest + 1

	core.registers.SetFlag(ZeroFlag, *dest == 0)
	core.registers.SetFlag(SignFlag, (*dest>>15)&0x01 == 1)
	core.registers.SetFlag(OverFlowFlag, false) // Assume no overflow for INC

eof:
	core.logInstruction(fmt.Sprintf("[%#04x] %s %s", core.GetCurrentlyExecutingInstructionAddress(), "INC", destName))

}

func INSTR_INC_SHORT_REL8(core *CpuCore) {

	var dest *uint16
	var destName string
	var result uint16

	core.currentByteAddr++

	modrm, bytesConsumed, err := core.consumeModRm()
	if err != nil {
		goto eof
	}
	core.currentByteAddr += bytesConsumed

	dest, destName, err = core.readRm16(&modrm)
	if err != nil {
		goto eof
	}

	result = *dest + 1

	// Update flags
	core.registers.SetFlag(CarryFlag, result < *dest) // Set if carry occurs

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, (result>>15)&0x01 == 1) // Set if the most significant bit is 1

	core.registers.SetFlag(OverFlowFlag, (*dest>>15)&0x01 == 0 && (result>>15)&0x01 == 1) // Set if result exceeds maximum positive value

	*dest = result

eof:
	core.logInstruction(fmt.Sprintf("[%#04x] %s %s", core.GetCurrentlyExecutingInstructionAddress(), "INC", destName))

}

func INSTR_DEC(core *CpuCore) {
	core.currentByteAddr++

	switch core.currentOpCodeBeingExecuted {
	case 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F:
		{
			// DEC r16
			val, valName := core.registers.registers16Bit[core.currentOpCodeBeingExecuted-0x48], core.registers.index16ToString(core.currentOpCodeBeingExecuted-0x48)

			originalVal := *val
			result := *val - 1

			// Update flags
			core.registers.SetFlag(ZeroFlag, result == 0)
			core.registers.SetFlag(SignFlag, (result>>15)&0x01 == 1)
			core.registers.SetFlag(OverFlowFlag, originalVal == 0x8000)
			core.registers.SetFlag(ParityFlag, calculateParity(result))

			*val = result

			core.logInstruction(fmt.Sprintf("[%#04x] DEC %s", core.GetCurrentlyExecutingInstructionAddress(), valName))
		}
	default:
		log.Println(fmt.Printf("Unhandled DEC instruction:  %#04x", core.currentOpCodeBeingExecuted))
		doCoreDump(core)
	}
}

func INSTR_DEC_RM16(core *CpuCore) {
	var dest *uint16
	var destName string
	var result uint16

	core.currentByteAddr++

	modrm, bytesConsumed, err := core.consumeModRm()
	if err != nil {
		log.Fatalf("error decoding modrm: %s", err)
	}
	core.currentByteAddr += bytesConsumed

	dest, destName, err = core.readRm16(&modrm)
	if err != nil {
		log.Fatalf("error reading rm16: %s", err)
	}

	originalDest := *dest
	result = *dest - 1

	// Update flags
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(SignFlag, (result>>15)&0x01 == 1)
	core.registers.SetFlag(OverFlowFlag, originalDest == 0x8000 && result == 0x7FFF)
	core.registers.SetFlag(ParityFlag, calculateParity(result))

	*dest = result

	core.logInstruction(fmt.Sprintf("[%#04x] %s %s", core.GetCurrentlyExecutingInstructionAddress(), "DEC", destName))

}

func INSTR_DAS(core *CpuCore) {
	core.currentByteAddr++

	// DAS - Decimal Adjust after Subtraction
	// AL = AL - 6 if low nibble > 9 or AF = 1
	// AL = AL - 0x60 if high nibble > 9 or CF = 1

	if core.registers.GetFlag(AdjustFlag) || core.registers.AL&0xf > 9 {
		core.registers.AL -= 6
	}

	if core.registers.GetFlag(CarryFlag) || core.registers.AL > 0x9f {
		core.registers.AL -= 0x60
		core.registers.SetFlag(CarryFlag, true)
	}

	core.registers.SetFlag(ZeroFlag, core.registers.AL == 0)

	core.logInstruction(fmt.Sprintf("[%#04x] das", core.GetCurrentlyExecutingInstructionAddress()))

}

func INSTR_DAA(core *CpuCore) {
	core.currentByteAddr++

	// DAA - Decimal Adjust after Addition
	// AL = AL + 6 if low nibble > 9 or AF = 1
	// AL = AL + 0x60 if high nibble > 9 or CF = 1

	if core.registers.GetFlag(AdjustFlag) || core.registers.AL&0xf > 9 {
		core.registers.AL += 6
	}

	if core.registers.GetFlag(CarryFlag) || core.registers.AL > 0x9f {
		core.registers.AL += 0x60
		core.registers.SetFlag(CarryFlag, true)
	}

	core.registers.SetFlag(ZeroFlag, core.registers.AL == 0)

	core.logInstruction(fmt.Sprintf("[%#04x] daa", core.GetCurrentlyExecutingInstructionAddress()))

}

func INSTR_AAA(core *CpuCore) {
	core.currentByteAddr++

	// AAA - ASCII Adjust after Addition
	// AL = AL + 6 if low nibble > 9 or AF = 1
	// AH = AH + 1
	// AF = 1
	// CF = 1

	if core.registers.GetFlag(AdjustFlag) || core.registers.AL&0xf > 9 {
		core.registers.AL += 6
		core.registers.AH += 1
		core.registers.SetFlag(CarryFlag, true)
		core.registers.SetFlag(AdjustFlag, true)
	}

	core.logInstruction(fmt.Sprintf("[%#04x] aaa", core.GetCurrentlyExecutingInstructionAddress()))

}

func INSTR_AAS(core *CpuCore) {
	core.currentByteAddr++

	// AAS - ASCII Adjust after Subtraction
	// AL = AL - 6 if low nibble > 9 or AF = 1
	// AH = AH - 1
	// AF = 1
	// CF = 1

	if core.registers.GetFlag(AdjustFlag) || core.registers.AL&0xf > 9 {
		core.registers.AL -= 6
		core.registers.AH -= 1
		core.registers.SetFlag(CarryFlag, true)
		core.registers.SetFlag(AdjustFlag, true)
	}

	core.logInstruction(fmt.Sprintf("[%#04x] aas", core.GetCurrentlyExecutingInstructionAddress()))

}

func INSTR_SHIFT(core *CpuCore) {
	core.currentByteAddr++

	// SAL
	var destTerm interface{}
	var bitLength uint8
	var t1Str string
	var t2Str string
	var opCode uint8
	var countTerm uint8
	var err error

	modrm, bytesConsumed, err := core.consumeModRm()
	if err != nil {
		goto eof
	}

	opCode = core.currentOpCodeBeingExecuted

	switch modrm.mod {
	case 4:
		{
			// sal
			switch opCode {
			// 8 bit versions
			case 0xd0:
				destTerm, t1Str, err = core.readRm8(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = 1
				t2Str = "1"
				bitLength = 8
			case 0xd2:
				destTerm, t1Str, err = core.readRm8(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = core.registers.CL
				t2Str = "CL"
				bitLength = 8
			case 0xC0:
				destTerm, t1Str, err = core.readRm8(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm, err := core.readImm8()
				if err != nil {
					goto eof
				}
				t2Str = string(countTerm)
				bitLength = 8

			// 16 bit versions
			case 0xd1:
				destTerm, t1Str, err = core.readRm16(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = 1
				t2Str = "1"
				bitLength = 16
			case 0xd3:
				destTerm, t1Str, err = core.readRm16(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = core.registers.CL
				t2Str = "CL"
				bitLength = 16
			case 0xC1:
				destTerm, t1Str, err = core.readRm16(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm, err := core.readImm8()
				if err != nil {
					goto eof
				}
				t2Str = string(countTerm)
				bitLength = 16

			default:
				{
					core.logInstruction("Unhandled shift sal variant")
					doCoreDump(core)
					panic(false)
				}
			}

			for i := uint8(0); i < countTerm; i++ {
				if bitLength == 8 {

					if *destTerm.(*uint8)&0x80 != 0 {
						// shift into carry
						core.registers.SetFlag(CarryFlag, *destTerm.(*uint8)>>(bitLength-1)&1 == 1)
					}

					*destTerm.(*uint8) <<= 1
				}

				if bitLength == 16 {

					if *destTerm.(*uint16)&0x80 != 0 {
						// shift into carry
						core.registers.SetFlag(CarryFlag, *destTerm.(*uint16)>>(bitLength-1)&1 == 1)
					}

					*destTerm.(*uint16) <<= 1
				}
			}

			core.logInstruction(fmt.Sprintf("[%#04x] sal %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Str, t2Str))

		}
	case 5, 7:
		{
			// shr - shifts to the right, and clears the most significant bit (unsigned)
			// sar - shifts to the right, preserves the most significant bit (signed)
			switch opCode {
			// 8 bit versions
			case 0xd0:
				destTerm, t1Str, err = core.readRm8(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = 1
				t2Str = "1"
				bitLength = 8
			case 0xd2:
				destTerm, t1Str, err = core.readRm8(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = core.registers.CL
				t2Str = "CL"
				bitLength = 8
			case 0xC0:
				destTerm, t1Str, err = core.readRm8(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm, err := core.readImm8()
				if err != nil {
					goto eof
				}
				t2Str = string(countTerm)
				bitLength = 8

			// 16 bit versions
			case 0xd1:
				destTerm, t1Str, err = core.readRm16(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = 1
				t2Str = "1"
				bitLength = 16
			case 0xd3:
				destTerm, t1Str, err = core.readRm16(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm = core.registers.CL
				t2Str = "CL"
				bitLength = 16
			case 0xC1:
				destTerm, t1Str, err = core.readRm16(&modrm)
				if err != nil {
					goto eof
				}
				core.currentByteAddr += bytesConsumed
				countTerm, err := core.readImm8()
				if err != nil {
					goto eof
				}
				t2Str = string(countTerm)
				bitLength = 16

			default:
				{
					core.logInstruction("Unhandled shift sal variant")
					doCoreDump(core)
					panic(false)
				}
			}

			var msbBit interface{}
			for i := uint8(0); i < countTerm; i++ {

				if bitLength == 8 {

					core.registers.SetFlag(CarryFlag, *destTerm.(*uint8)&1 == 1)

					msbBit = (*destTerm.(*uint8) >> 8) & 1
					*destTerm.(*uint8) >>= 1
					if modrm.mod == 7 {
						*destTerm.(*uint8) = *destTerm.(*uint8) | (msbBit.(uint8) << (bitLength - 1))
					}
				}

				if bitLength == 16 {

					core.registers.SetFlag(CarryFlag, *destTerm.(*uint16)&1 == 1)

					msbBit = (*destTerm.(*uint16) >> 8) & 1
					*destTerm.(*uint16) >>= 1
					if modrm.mod == 7 {
						*destTerm.(*uint16) = *destTerm.(*uint16) | (msbBit.(uint16) << (bitLength - 1))
					}
				}
			}

			opCode := "SHR"
			if modrm.mod == 7 {
				opCode = "SAR"
			}

			core.logInstruction(fmt.Sprintf("[%#04x] %s %s, %s", core.GetCurrentlyExecutingInstructionAddress(), opCode, t1Str, t2Str))

		}
	}

eof:
}

func INSTR_STC(core *CpuCore) {
	core.registers.SetFlag(CarryFlag, true)
	core.logInstruction(fmt.Sprintf("[%#04x] STC", core.GetCurrentlyExecutingInstructionAddress()))

	core.registers.IP += 1
}

func INSTR_SBB(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var term2 uint32
	var result uint32

	var signr uint8
	var sign1 uint32
	var sign2 uint32
	var bitLength uint32

	switch core.currentOpCodeBeingExecuted {
	case 0x1C:
		{
			// SBB AL, imm8
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			term1 = uint32(core.registers.AL)
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			core.registers.AL = uint8(result)

			core.logInstruction(fmt.Sprintf("[%#04x] sbb al, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term2))
			goto success
		}
	case 0x1D:
		{
			// SBB AX, imm16
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}
			term1 = uint32(core.registers.AX)
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			core.registers.AX = uint16(result)

			core.logInstruction(fmt.Sprintf("[%#04x] sbb ax, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term2))
			goto success
		}
	case 0x1F:
		{
			// SBB BX, imm16
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}
			term1 = uint32(core.registers.BX)
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			core.registers.BX = uint16(result)

			core.logInstruction(fmt.Sprintf("[%#04x] sbb bx, %#04x", core.GetCurrentlyExecutingInstructionAddress(), term2))
			goto success
		}
	case 0x80:
		{
			// SBB r/m8, imm8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] sbb %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), t1Name, term2))
			goto success
		}
	case 0x81:
		{
			// SBB r/m16, imm16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm16()
			if err != nil {
				goto eof
			}
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			tmp := uint16(result)
			srcName, err := core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] sbb %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), t1Name, srcName))
			goto success
		}
	case 0x83:
		{
			// SBB r/m16, imm8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2, err := core.readImm8()
			if err != nil {
				goto eof
			}
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			tmp := uint16(result)
			srcName, err := core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] sbb %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), t1Name, srcName))
			goto success
		}
	case 0x18:
		{
			// SBB r/m8, r8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			t2, t2Name := core.readR8(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] sbb %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	case 0x19:
		{
			// SBB r/m16, r16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			t2, t2Name := core.readR16(&modrm)
			term2 = uint32(*t2)
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			tmp := uint16(result)
			_, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] sbb %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	case 0x1A:
		{
			// SBB r8, r/m8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name := core.readR8(&modrm)
			term1 = uint32(*t1)
			t2, t2Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term2 = uint32(*t2)
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] sbb %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	case 0x1B:
		{
			// SBB r16, r/m16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name := core.readR16(&modrm)
			term1 = uint32(*t1)
			t2, t2Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}

			core.currentByteAddr += bytesConsumed
			term2 = uint32(*t2)
			result = uint32(term1) - (uint32(term2) + uint32(core.registers.GetFlagInt(CarryFlag)))
			tmp := uint16(result)
			core.writeR16(&modrm, &tmp)

			core.logInstruction(fmt.Sprintf("[%#04x] sbb %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, t2Name))
			goto success
		}
	default:
		log.Fatal(fmt.Sprintf("Unrecognised SBB instruction: %#04x!", core.currentOpCodeBeingExecuted))
	}

success:
	bitLength = uint32(bits.Len32(result))

	// update flags
	sign1 = (term1 >> (bitLength)) & 0x01
	sign2 = (term2 >> (bitLength)) & 0x01
	signr = uint8((result >> (bitLength)) & 0x01)

	core.registers.SetFlag(CarryFlag, result>>(bitLength) != 0)

	core.registers.SetFlag(ZeroFlag, result == 0)

	core.registers.SetFlag(SignFlag, signr != 0)

	core.registers.SetFlag(OverFlowFlag, (sign1 == 1 && sign2 == 0 && signr == 0) || (sign1 == 0 && sign2 == 1 && signr == 1))

eof:
}

func INSTR_IMUL(core *CpuCore) {
	core.currentByteAddr++

	var result int32
	var overflow bool

	switch core.currentOpCodeBeingExecuted {
	case 0x69:
		// IMUL r16, r/m16, imm16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error consuming ModR/M byte: %v", err))
			return
		}
		core.currentByteAddr += bytesConsumed

		dest, destName := core.readR16(&modrm)
		src, srcName, err := core.readRm16(&modrm)
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error reading r/m16: %v", err))
			return
		}

		imm16, err := core.readImm16()
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error reading imm16: %v", err))
			return
		}

		term1 := int32(int16(*src))
		term2 := int32(int16(imm16))
		result = term1 * term2

		if result > 32767 || result < -32768 {
			overflow = true
		}

		*dest = uint16(result)

		core.logInstruction(fmt.Sprintf("[%#04x] imul %s, %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName, imm16))

	case 0x6B:
		// IMUL r16, r/m16, imm8
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error consuming ModR/M byte: %v", err))
			return
		}
		core.currentByteAddr += bytesConsumed

		dest, destName := core.readR16(&modrm)
		src, srcName, err := core.readRm16(&modrm)
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error reading r/m16: %v", err))
			return
		}

		imm8, err := core.readImm8()
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error reading imm8: %v", err))
			return
		}

		term1 := int32(int16(*src))
		term2 := int32(int8(imm8))
		result = term1 * term2

		if result > 32767 || result < -32768 {
			overflow = true
		}

		*dest = uint16(result)

		core.logInstruction(fmt.Sprintf("[%#04x] imul %s, %s, %#04x", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName, imm8))

	case 0x0F, 0xAF:
		// IMUL r16, r/m16
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error consuming ModR/M byte: %v", err))
			return
		}
		core.currentByteAddr += bytesConsumed

		dest, destName := core.readR16(&modrm)
		src, srcName, err := core.readRm16(&modrm)
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error reading r/m16: %v", err))
			return
		}

		term1 := int32(int16(*dest))
		term2 := int32(int16(*src))
		result = term1 * term2

		if result > 32767 || result < -32768 {
			overflow = true
		}

		*dest = uint16(result)

		core.logInstruction(fmt.Sprintf("[%#04x] imul %s, %s", core.GetCurrentlyExecutingInstructionAddress(), destName, srcName))

	case 0xF6:
		// IMUL r/m8 (AX <- AL * r/m8)
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error consuming ModR/M byte: %v", err))
			return
		}
		core.currentByteAddr += bytesConsumed

		src, srcName, err := core.readRm8(&modrm)
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error reading r/m8: %v", err))
			return
		}

		term1 := int16(int8(core.registers.AL))
		term2 := int16(int8(*src))
		result16 := term1 * term2

		if result16 > 127 || result16 < -128 {
			overflow = true
		}

		core.registers.AX = uint16(result16)

		core.logInstruction(fmt.Sprintf("[%#04x] imul %s", core.GetCurrentlyExecutingInstructionAddress(), srcName))

	case 0xF7:
		// IMUL r/m16 (DX:AX <- AX * r/m16)
		modrm, bytesConsumed, err := core.consumeModRm()
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error consuming ModR/M byte: %v", err))
			return
		}
		core.currentByteAddr += bytesConsumed

		src, srcName, err := core.readRm16(&modrm)
		if err != nil {
			core.logInstruction(fmt.Sprintf("Error reading r/m16: %v", err))
			return
		}

		term1 := int32(int16(core.registers.AX))
		term2 := int32(int16(*src))
		result = term1 * term2

		if result > 32767 || result < -32768 {
			overflow = true
		}

		core.registers.AX = uint16(result)
		core.registers.DX = uint16(result >> 16)

		core.logInstruction(fmt.Sprintf("[%#04x] imul %s", core.GetCurrentlyExecutingInstructionAddress(), srcName))

	default:
		core.logInstruction(fmt.Sprintf("Unrecognised IMUL instruction: %#04x!", core.currentOpCodeBeingExecuted))
		return
	}

	// Update flags
	core.registers.SetFlag(OverFlowFlag, overflow)
	core.registers.SetFlag(CarryFlag, overflow)
	core.registers.SetFlag(SignFlag, result < 0)
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(ParityFlag, bits.OnesCount32(uint32(result))%2 == 0)

}

func INSTR_MUL(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentOpCodeBeingExecuted {
	case 0xF6:
		{
			// MUL r/m8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2 = uint32(core.registers.AL)
			result = term1 * term2
			core.registers.AX = uint16(result)
			core.registers.SetFlag(CarryFlag, result>>16 != 0)

			core.logInstruction(fmt.Sprintf("[%#04x] mul %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, "AL"))
			goto success
		}
	case 0xF7:
		{
			// MUL r/m16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2 = uint32(core.registers.AX)
			result = term1 * term2
			core.registers.DX = uint16(result >> 16)
			core.registers.AX = uint16(result)

			core.logInstruction(fmt.Sprintf("[%#04x] mul %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, "AX"))
			goto success
		}
	default:
		log.Fatal(fmt.Sprintf("Unrecognised MUL instruction: %#04x!", core.currentOpCodeBeingExecuted))
	}

success:
	bitLength = uint32(bits.Len32(result))

	// update flags
	core.registers.SetFlag(OverFlowFlag, false)
	core.registers.SetFlag(CarryFlag, result>>bitLength != 0)
	core.registers.SetFlag(SignFlag, result>>bitLength != 0)
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(ParityFlag, bits.OnesCount32(result)%2 == 0)

eof:
}

func INSTR_DIV(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var term2 uint32
	var result uint32

	var bitLength uint32

	switch core.currentOpCodeBeingExecuted {
	case 0xF6:
		{
			// DIV r/m8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2 = uint32(core.registers.AL)
			result = term2 / term1
			core.registers.AX = uint16(result)
			core.registers.AL = uint8(term2 % term1)

			core.logInstruction(fmt.Sprintf("[%#04x] div %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, "AL"))
			goto success
		}
	case 0xF7:
		{
			// DIV r/m16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			term2 = uint32(core.registers.AX)
			result = term2 / term1
			core.registers.DX = uint16(result % term1)
			core.registers.AX = uint16(result)

			core.logInstruction(fmt.Sprintf("[%#04x] div %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, "AX"))
			goto success
		}
	default:
		log.Fatal(fmt.Sprintf("Unrecognised DIV instruction: %#04x!", core.currentOpCodeBeingExecuted))
	}

success:
	bitLength = uint32(bits.Len32(result))

	// update flags
	core.registers.SetFlag(OverFlowFlag, false)
	core.registers.SetFlag(CarryFlag, result>>bitLength != 0)
	core.registers.SetFlag(SignFlag, result>>bitLength != 0)
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(ParityFlag, bits.OnesCount32(result)%2 == 0)

eof:
}

func INSTR_IDIV(core *CpuCore) {
	core.currentByteAddr++

	var term1 int32
	var term2 int32
	var result int32

	var bitLength uint32

	switch core.currentOpCodeBeingExecuted {
	case 0xF6:
		{
			// IDIV r/m8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = int32(*t1)
			term2 = int32(core.registers.AL)
			result = term2 / term1
			core.registers.AX = uint16(result)
			core.registers.AL = uint8(term2 % term1)

			core.logInstruction(fmt.Sprintf("[%#04x] idiv %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, "AL"))
			goto success
		}
	case 0xF7:
		{
			// IDIV r/m16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = int32(*t1)
			term2 = int32(core.registers.AX)
			result = term2 / term1
			core.registers.DX = uint16(result % term1)
			core.registers.AX = uint16(result)

			core.logInstruction(fmt.Sprintf("[%#04x] idiv %s, %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name, "AX"))
			goto success
		}
	default:
		log.Fatal(fmt.Sprintf("Unrecognised IDIV instruction: %#04x!", core.currentOpCodeBeingExecuted))
	}

success:
	bitLength = uint32(bits.Len32(uint32(result)))

	// update flags
	core.registers.SetFlag(OverFlowFlag, false)
	core.registers.SetFlag(CarryFlag, result>>bitLength != 0)
	core.registers.SetFlag(SignFlag, result>>bitLength != 0)
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(ParityFlag, bits.OnesCount32(uint32(result))%2 == 0)

eof:
}

func INSTR_NEG(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var result uint32

	switch core.currentOpCodeBeingExecuted {
	case 0xF6:
		{
			// NEG r/m8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			result = uint32(0 - int32(term1))
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] neg %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name))
			goto success
		}
	case 0xF7:
		{
			// NEG r/m16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			result = uint32(0 - int32(term1))
			tmp := uint16(result)
			_, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] neg %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name))
			goto success
		}
	default:
		log.Fatal(fmt.Sprintf("Unrecognised NEG instruction: %#04x!", core.currentOpCodeBeingExecuted))
	}

success:
	// update flags
	core.registers.SetFlag(OverFlowFlag, false)
	core.registers.SetFlag(CarryFlag, result != 0)
	core.registers.SetFlag(SignFlag, (result>>15)&0x01 == 1)
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(ParityFlag, bits.OnesCount32(uint32(result))%2 == 0)

eof:
}

func INSTR_NOT(core *CpuCore) {
	core.currentByteAddr++

	var term1 uint32
	var result uint32

	var bitLength uint32

	switch core.currentOpCodeBeingExecuted {
	case 0xF6:
		{
			// NOT r/m8
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm8(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			result = ^term1
			tmp := uint8(result)
			_, err = core.writeRm8(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] not %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name))
			goto success
		}
	case 0xF7:
		{
			// NOT r/m16
			modrm, bytesConsumed, err := core.consumeModRm()
			if err != nil {
				goto eof
			}
			t1, t1Name, err := core.readRm16(&modrm)
			if err != nil {
				goto eof
			}
			core.currentByteAddr += bytesConsumed
			term1 = uint32(*t1)
			result = ^term1
			tmp := uint16(result)
			_, err = core.writeRm16(&modrm, &tmp)
			if err != nil {
				goto eof
			}

			core.logInstruction(fmt.Sprintf("[%#04x] not %s", core.GetCurrentlyExecutingInstructionAddress(), t1Name))
			goto success
		}
	default:
		log.Fatal(fmt.Sprintf("Unrecognised NOT instruction: %#04x!", core.currentOpCodeBeingExecuted))
	}
success:
	bitLength = uint32(bits.Len32(uint32(result)))

	// update flags
	core.registers.SetFlag(OverFlowFlag, false)
	core.registers.SetFlag(CarryFlag, result>>bitLength != 0)
	core.registers.SetFlag(SignFlag, result>>bitLength != 0)
	core.registers.SetFlag(ZeroFlag, result == 0)
	core.registers.SetFlag(ParityFlag, bits.OnesCount32(uint32(result))%2 == 0)

eof:
}

func calculateParity(value uint16) bool {
	count := 0
	for i := 0; i < 16; i++ {
		if (value>>i)&1 == 1 {
			count++
		}
	}
	return count%2 == 0
}
