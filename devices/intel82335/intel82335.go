package intel82335

import (
	"fmt"
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
	"strings"
)

/*
	Simulated 82335 High Integration Interface Device for 386SX
*/
const (
	MCR_BIOS_ROM_ENABLE = 0
	MCR_S640_BASE_MEMORY_SIZE = 3
	MCR_DRAM_MEMORY_SIZE = 4
	MCR_MEMORY_INTERLEAVE = 6
	MCR_ROM_SIZE = 8
	MCR_ADAPTER_ROM_ENABLE = 9
	MCR_VIDEO_RAM_ENABLE = 10
	MCR_VIDEO_READ_ONLY = 11


	MCR_MEMORY_INTERLEAVE_1_BANK = 1
	MCR_MEMORY_INTERLEAVE_2_BANK = 2
	MCR_MEMORY_INTERLEAVE_3_BANK = 3
	MCR_MEMORY_INTERLEAVE_4_BANK = 4

)

type Intel82335 struct {
	busId uint32

	biosRomAccessEnabled bool // 0 = enable BIOS rom access. 1 = disable BIOS rom access and enable bios shadow
	s640BaseMemorySize   bool // 0 = 512KB, 1 = 640KB
	dRamSize             bool // 0 = 1MBx1DRAM, 1 = 256KBx1 or 256KBx4 DRAM
	romSize              bool // 0 = 256KB rom, 1 = 512KB rom
	adapterRomEnabled    bool // 0 = enable adapter ROM access. 1 = disable adapter ROM access and enable shadow
	videoRamEnabled      bool // 0 = enable video RAM access. 1 = disable video RAM access
	videoReadOnly        bool // 0 = video read write. 1 = video read only
	memoryInterleaving   uint8
	mcrRegisterLastValue uint8 //last value set by register initialize. don't manually write to this value.
}

func NewIntel82335() *Intel82335 {
	chip := &Intel82335{}

	return chip
}

func (device *Intel82335) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *Intel82335) OnReceiveMessage(message bus.BusMessage) {

}


func (device *Intel82335) McrRegisterInitialize(registerValue uint8) {

	device.mcrRegisterLastValue = registerValue

	// bits
	device.biosRomAccessEnabled = getRegisterBit(registerValue, MCR_BIOS_ROM_ENABLE)
	device.s640BaseMemorySize = getRegisterBit(registerValue, MCR_S640_BASE_MEMORY_SIZE)
	device.dRamSize = getRegisterBit(registerValue, MCR_DRAM_MEMORY_SIZE)
	device.memoryInterleaving = getMemoryInterleaveMode(registerValue)
	device.romSize = getRegisterBit(registerValue, MCR_ROM_SIZE)
	device.adapterRomEnabled = getRegisterBit(registerValue, MCR_ADAPTER_ROM_ENABLE)
	device.videoRamEnabled = getRegisterBit(registerValue, MCR_VIDEO_RAM_ENABLE)
	device.videoReadOnly = getRegisterBit(registerValue, MCR_VIDEO_READ_ONLY)


	log.Printf(fmt.Sprintf("MCR Set Config: %s", device.toString()))
}

func (device *Intel82335) toString() string {
	var strs []string

	if !device.biosRomAccessEnabled {
		strs = append(strs, "BiosRomAccessEnabled")
	} else {
		strs = append(strs, "BiosRomAccessDisabled,BiosRomShadowEnabled")
	}

	if !device.s640BaseMemorySize {
		strs = append(strs, "S640BaseMemorySize=512KB")
	} else {
		strs = append(strs, "S640BaseMemorySize=640KB")
	}

	if !device.dRamSize {
		strs = append(strs, "DSIZE=1MBx1")
	} else {
		strs = append(strs, "DSIZE=256KBx1/256KBx4")
	}

	strs = append(strs, fmt.Sprintf("MEMBANKS:%d", device.memoryInterleaving))

	if !device.romSize {
		strs = append(strs, "ROMSIZE=256KB")
	} else {
		strs = append(strs, "ROMSIZE=512KB")
	}

	if !device.adapterRomEnabled {
		strs = append(strs, "AdapterRomAccessEnabled")
	} else {
		strs = append(strs, "AdapterRomAccessDisabled,AdapterRomShadowEnabled")
	}

	if !device.videoRamEnabled {
		strs = append(strs, "VideoRamAccessEnabled")
	} else {
		strs = append(strs, "VideoRamAccessDisabled")
	}

	if !device.videoReadOnly {
		strs = append(strs, "VideoReadWrite")
	} else {
		strs = append(strs, "VideoReadOnly")
	}

	return strings.Join(strs, ",")
}

func (device *Intel82335) GetMcrRegister() uint8 {
	return device.mcrRegisterLastValue
}

func getMemoryInterleaveMode(registervalue uint8) uint8 {
	o2 := getRegisterBit(registervalue, MCR_MEMORY_INTERLEAVE)
	o1 := getRegisterBit(registervalue, MCR_MEMORY_INTERLEAVE+1)

	switch {
	case !o1 && !o2: return MCR_MEMORY_INTERLEAVE_1_BANK
	case !o1 && o2: return MCR_MEMORY_INTERLEAVE_2_BANK
	case o1 && !o2: return MCR_MEMORY_INTERLEAVE_3_BANK
	case o1 && o2: return MCR_MEMORY_INTERLEAVE_4_BANK
	default:
		panic("Invalid memory interleave mode configuration!")
	}
}

func getRegisterBit(source uint8, position uint8) bool {
	return (source>>position)&1 == 1
}