package common

import (
	"fmt"
	"strings"
)

func CmosRegisterWriteToFriendlyString(registerSelect uint8, value uint8) string {
	switch registerSelect {
	case 0x00:
		return "Seconds: " + formatBCD(value)
	case 0x01:
		return "Seconds Alarm: " + formatBCD(value)
	case 0x02:
		return "Minutes: " + formatBCD(value)
	case 0x03:
		return "Minutes Alarm: " + formatBCD(value)
	case 0x04:
		return "Hours: " + formatBCD(value&0x7F) + ", 24-Hour Mode: " + formatBool(value&0x80 != 0)
	case 0x05:
		return "Hours Alarm: " + formatBCD(value&0x7F) + ", 24-Hour Mode: " + formatBool(value&0x80 != 0)
	case 0x06:
		return "Day of Week: " + formatDayOfWeek(value)
	case 0x07:
		return "Date of Month: " + formatBCD(value)
	case 0x08:
		return "Month: " + formatBCD(value)
	case 0x09:
		return "Year: " + formatBCD(value)
	case 0x0A:
		return "Register A: " + formatRegisterA(value)
	case 0x0B:
		return "Register B: " + formatRegisterB(value)
	case 0x0C:
		return "Register C: " + formatRegisterC(value)
	case 0x0D:
		return "Register D: " + formatRegisterD(value)
	case 0x0E:
		return "Diagnostic Status: " + fmt.Sprintf("%08b", value)
	case 0x0F:
		return "Shutdown Status: " + formatShutdownStatus(value)
	case 0x10:
		return "Floppy Disk Drive Type: " + formatFloppyDiskDriveType(value)
	case 0x11:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x12:
		return "Hard Disk Drive Type: " + formatHardDiskDriveType(value)
	case 0x13:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x14:
		return "Equipment Byte: " + formatEquipmentByte(value)
	case 0x15:
		return "Base Memory Low Byte: " + fmt.Sprintf("%02X", value)
	case 0x16:
		return "Base Memory High Byte: " + fmt.Sprintf("%02X", value)
	case 0x17:
		return "Extended Memory Low Byte: " + fmt.Sprintf("%02X", value)
	case 0x18:
		return "Extended Memory High Byte: " + fmt.Sprintf("%02X", value)
	case 0x19:
		return "Drive C Extended Type: " + formatDriveExtendedType(value)
	case 0x1A:
		return "Drive D Extended Type: " + formatDriveExtendedType(value)
	case 0x2E:
		return "CMOS Checksum High Byte: " + fmt.Sprintf("%02X", value)
	case 0x2F:
		return "CMOS Checksum Low Byte: " + fmt.Sprintf("%02X", value)
	case 0x30:
		return "Extended Memory Low Byte: " + fmt.Sprintf("%02X", value)
	case 0x31:
		return "Extended Memory High Byte: " + fmt.Sprintf("%02X", value)
	case 0x32:
		return "Date Century Byte: " + formatDateCentury(value)
	case 0x33:
		return "Information Flags: " + formatInformationFlags(value)
	case 0x34:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x35:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x38:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x3D:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x3E:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x3F:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x40:
		return "Floppy Disk Drive 0 Media Type: " + formatFloppyDiskMediaType(value)
	case 0x41:
		return "Floppy Disk Drive 1 Media Type: " + formatFloppyDiskMediaType(value)
	case 0x42:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x4E:
		return "Real-Time Clock Day of Month Alarm: " + formatBCD(value)
	case 0x4F:
		return "Real-Time Clock Month Alarm: " + formatBCD(value)
	case 0x50:
		return "Real-Time Clock Century: " + formatBCD(value)
	case 0x51:
		return "Real-Time Clock Century Alarm: " + formatBCD(value)
	case 0x52:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x53:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x5B:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x5C:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x5D:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x5E:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x5F:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x65:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x66:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x67:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x68:
		return "Extended CMOS Checksum Byte: " + fmt.Sprintf("%02X", value)
	case 0x69:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x6B:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x6C:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x6D:
		return "Reserved: " + fmt.Sprintf("%02X", value)
	case 0x71:
		return "RTC Address: " + fmt.Sprintf("%02X", value)
	case 0x72:
		return "RTC Data: " + fmt.Sprintf("%02X", value)
	case 0x73:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x74:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x75:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x76:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x77:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x78:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x79:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x7A:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x7B:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x7C:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x7D:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x7E:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0x7F:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xD0:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xD1:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xD2:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xD3:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xD4:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xD5:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xD6:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	case 0xDC:
		return "PS/2 Mouse Port Installed: " + formatBool(value&0x01 != 0)
	case 0xDD:
		return "INT 15h, E801h Memory Size: " + fmt.Sprintf("%02X", value)
	case 0xDE:
		return "INT 15h, E820h Memory Size: " + fmt.Sprintf("%02X", value)
	case 0xDF:
		return "Chipset-Specific Register: " + fmt.Sprintf("%02X", value)
	default:
		return "Unknown Register: " + fmt.Sprintf("%02X", registerSelect) + ", Value: " + fmt.Sprintf("%02X", value)
	}
}

func formatBCD(value uint8) string {
	return fmt.Sprintf("%02d", value>>4*10+value&0xF)
}

func formatBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func formatDayOfWeek(value uint8) string {
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	if value >= 1 && value <= 7 {
		return days[value-1]
	}
	return "Invalid Day"
}

func formatRegisterA(value uint8) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UIP: %d, ", value>>7&1))
	sb.WriteString(fmt.Sprintf("DV: %d, ", value>>6&1))
	sb.WriteString(fmt.Sprintf("RS: %02b", value&0x0F))
	return sb.String()
}

func formatRegisterB(value uint8) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SET: %d, ", value>>7&1))
	sb.WriteString(fmt.Sprintf("PIE: %d, ", value>>6&1))
	sb.WriteString(fmt.Sprintf("AIE: %d, ", value>>5&1))
	sb.WriteString(fmt.Sprintf("UIE: %d, ", value>>4&1))
	sb.WriteString(fmt.Sprintf("SQWE: %d, ", value>>3&1))
	sb.WriteString(fmt.Sprintf("DM: %d, ", value>>2&1))
	sb.WriteString(fmt.Sprintf("24/12: %d, ", value>>1&1))
	sb.WriteString(fmt.Sprintf("DSE: %d", value&1))
	return sb.String()
}

func formatRegisterC(value uint8) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("IRQF: %d, ", value>>7&1))
	sb.WriteString(fmt.Sprintf("PF: %d, ", value>>6&1))
	sb.WriteString(fmt.Sprintf("AF: %d, ", value>>5&1))
	sb.WriteString(fmt.Sprintf("UF: %d", value>>4&1))
	return sb.String()
}

func formatRegisterD(value uint8) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("VRT: %d, ", value>>7&1))
	sb.WriteString(fmt.Sprintf("Reserved: %07b", value&0x7F))
	return sb.String()
}

func formatShutdownStatus(value uint8) string {
	statuses := []string{
		"Normal",
		"Reset by Keyboard Controller",
		"Reset by CMOS Corrupted",
		"Reset by Memory Size Change",
		"Reset by Watchdog Timer",
		"Reset by Power Management Event",
		"Reset by Software Reset",
		"Reset by After Memory Test",
	}
	if value >= 0 && value <= 7 {
		return statuses[value]
	}
	return "Unknown Status"
}

func formatFloppyDiskDriveType(value uint8) string {
	types := []string{
		"None",
		"360KB 5.25\"",
		"1.2MB 5.25\"",
		"720KB 3.5\"",
		"1.44MB 3.5\"",
		"2.88MB 3.5\"",
	}
	if value >= 1 && value <= 5 {
		return types[value]
	}
	return "Unknown Type"
}

func formatHardDiskDriveType(value uint8) string {
	types := []string{
		"None",
		"Type 1",
		"Type 2",
		"Type 3",
		"Type 4",
		"Type 5",
		"Type 6",
		"Type 7",
		"Type 8",
		"Type 9",
		"Type 10",
		"Type 11",
		"Type 12",
		"Type 13",
		"Type 14",
		"Type 15",
	}
	if value >= 0 && value <= 15 {
		return types[value]
	}
	return "Unknown Type"
}

func formatEquipmentByte(value uint8) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Floppy Drive Count: %d, ", value>>6&3))
	sb.WriteString(fmt.Sprintf("Math Coprocessor: %s, ", formatBool(value>>2&1 != 0)))
	sb.WriteString(fmt.Sprintf("Mouse: %s", formatBool(value&1 != 0)))
	return sb.String()
}

func formatDriveExtendedType(value uint8) string {
	types := []string{
		"None",
		"1.44MB 3.5\"",
		"2.88MB 3.5\"",
		"Hard Disk",
	}
	if value >= 0 && value <= 3 {
		return types[value]
	}
	return "Unknown Type"
}

func formatDateCentury(value uint8) string {
	return fmt.Sprintf("%02d", value)
}

func formatInformationFlags(value uint8) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Reserved: %d, ", value>>7&1))
	sb.WriteString(fmt.Sprintf("Reserved: %d, ", value>>6&1))
	sb.WriteString(fmt.Sprintf("Reserved: %d, ", value>>5&1))
	sb.WriteString(fmt.Sprintf("Reserved: %d, ", value>>4&1))
	sb.WriteString(fmt.Sprintf("Reserved: %d, ", value>>3&1))
	sb.WriteString(fmt.Sprintf("Reserved: %d, ", value>>2&1))
	sb.WriteString(fmt.Sprintf("Reserved: %d, ", value>>1&1))
	sb.WriteString(fmt.Sprintf("Reserved: %d", value&1))
	return sb.String()
}

func formatFloppyDiskMediaType(value uint8) string {
	types := []string{
		"360KB 5.25\"",
		"1.2MB 5.25\"",
		"720KB 3.5\"",
		"1.44MB 3.5\"",
		"2.88MB 3.5\"",
	}
	if value >= 0 && value <= 4 {
		return types[value]
	}
	return "Unknown Type"
}
