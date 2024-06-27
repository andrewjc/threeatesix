package intel8086

import "fmt"

func INSTR_ROUTER_2BYTE_01(core *CpuCore) {

	core.currentByteAddr++
	modrm, _, err := core.consumeModRm()
	if err != nil {
		core.logInstruction("Error in INSTR_ROUTER_2BYTE_01: %s\n", err)
		return
	}
	core.currentByteAddr-- // We need to re-read the modrm byte

	switch modrm.reg {
	/*case 0x00:
		return core.handleSGDT
	case 0x01:
		return core.handleSIDT
	case 0x02:
		return core.handleLGDT
	case 0x03:
		return core.handleLIDT*/
	case 0x04:
		INSTR_SMSW(core) //SMSW r/m16
	case 0x05:
		INSTR_SMSW(core) //SMSW r32/m16
	/*case 0x06:
		return core.handleLMSW
	case 0x07:
		return core.handleINVLPG*/
	default:
		core.logInstruction("[%#04x] Unrecognized 2-byte opcode: 0x0F 0x%02X", core.GetCurrentCodePointer(), modrm.reg)
		doCoreDump(core)
		panic(fmt.Sprintf("Unrecognized 2-byte opcode: 0x0F %#02x", modrm.reg))
	}
}
