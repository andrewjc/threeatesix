package common

import "fmt"

func BiosPostCodeToString(value uint8) string {
	switch value {
	case 0x0:
		return "Power-on reset"
	case 0x1:
		return "Register test about to start"
	case 0x2:
		return "Register test has passed"
	case 0x3:
		return "ROM BIOS checksum test passed"
	case 0x4:
		return "Passed keyboard controller test with and without mouse"
	case 0x5:
		return "Chipset initialized... DMA and interrupt controller disabled"
	case 0x6:
		return "Video system disabled and the system timer checks OK"
	case 0x7:
		return "8265 programmable interval timer initialized"
	case 0x8:
		return "Delta counter channel 2 initialized"
	case 0x9:
		return "Delta counter channel 1 initialized"
	case 0x0A:
		return "Delta counter channel 0 initialized"
	case 0x0B:
		return "Refresh started"
	case 0x0C:
		return "System timer started"
	case 0x0D:
		return "Refresh check OK"
	case 0x0E:
		return "Refresh period check OK"
	case 0x0F:
		return "System timer check OK"
	case 0x10:
		return "Ready to start 64KB base memory test"
	case 0x11:
		return "Address line test OK"
	case 0x12:
		return "64KB base memory test OK"
	case 0x13:
		return "Interrupt vectors initialized"

	case 0x14:
		return "8042 keyboard controller test OK"

	case 0x15:
		return "CMOS read/write test ok"

	case 0x16:
		return "CMOS checksum and battery ok"

	case 0x17:
		return "Monochrome video mode OK"
	case 0x18:
		return "CGA color mode set OK"
	case 0x19:
		return "Attempting to pass control to video ROM at C0000h"
	case 0x1A:
		return "Returned from video ROM"
	case 0x1B:
		return "Shadow RAM enabled"
	case 0x1C:
		return "Display memory read/write test OK"
	case 0x1D:
		return "Alternate display memory read/write test OK"
	case 0x1E:
		return "Global equipment byte set for proper"
	case 0x1F:
		return "Ready to initialize video system"
	case 0x20:
		return "Finished setting video mode"
	case 0x21:
		return "ROM type 27256 verified"
	case 0x22:
		return "The power-on message is displayed"
	case 0x30:
		return "Ready to start the virtual mode memory test"
	case 0x31:
		return "Virtual memory mode test started"
	case 0x32:
		return "CPU has switched to virtual mode"
	case 0x33:
		return "Testing the memory address lines"
	case 0x34:
		return "Testing the memory address lines"
	case 0x35:
		return "Lower 1MB of RAM found"
	case 0x36:
		return "Memory size computation checks OK"
	case 0x37:
		return "Memory test in progress"
	case 0x38:
		return "Memory below 1MB is initialized"
	case 0x39:
		return "Memory above 1MB is initialized"

	case 0x3A:
		return "Memory size is displayed"
	case 0x3B:
		return "Ready to test the lower 1MB of RAM"
	case 0x3C:
		return "Memory test of lower 1MB OK"
	case 0x3D:
		return "Memory test above 1MB OK"
	case 0x3E:
		return "Ready to shutdown for real-mode testing"
	case 0x3F:
		return "Shutdown OK - now in real mode"

	case 0x40:
		return "Cache memory now on... Ready to disable gate A 20"
	case 0x41:
		return "A20 line disabled successfully"
	case 0x42:
		return "i486 internal cache turned ok"
	case 0x43:
		return "Ready to start DMA controller test"

	case 0x50:
		return "Shutdown OK - now in real mode"
	case 0x51:
		return "Starting DMA controller 1 register test"
	case 0x52:
		return "DMA controller 1 test passed, starting DMA controller 2 register test"
	case 0x53:
		return "DMA controller 2 test passed"
	case 0x54:
		return "Ready to test latch on DMA controller 1 and 2"
	case 0x55:
		return "DMA controller and 2 latch test OK"
	case 0x56:
		return "DMA controller 1 and 2 configured OK"
	case 0x57:
		return "8259 programmable interrupt controller initialized OK"

	case 0x70:
		return "Start of keyboard test"
	case 0x71:
		return "Keyboard controller OK"
	case 0x72:
		return "Keyboard test OK... Starting mouse interface test"
	case 0x73:
		return "Keyboard and mouse global initialized OK"
	case 0x74:
		return "Display setup prompt... Floppy setup ready to start"
	case 0x75:
		return "Floppy controller test OK"

	case 0x76:
		return "Hard disk setup ready to start"
	case 0x77:
		return "Hard disk controller setup OK"
	case 0x79:
		return "Ready to initialize timer data"

	case 0x7A:
		return "Timer data area initialized"
	case 0x7B:
		return "CMOS battery verified OK"
	case 0x7E:
		return "CMOS memory size updated"
	case 0x7F:
		return "Enable setup routine if <DELETE> is pressed"
	default:
		return fmt.Sprintf("Unhandled post code: %#02x", value)
	}

}
