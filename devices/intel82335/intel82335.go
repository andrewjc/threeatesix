package intel82335

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"

	"log"
	"strings"
)

/*
Simulated 82335 High Integration Interface Device for 386SX
*/
const (
	MCR_BIOS_ROM_ENABLE       = 0  // Bit position for BIOS ROM enable
	MCR_S640_BASE_MEMORY_SIZE = 3  // Bit position for base memory size (512KB or 640KB)
	MCR_DRAM_MEMORY_SIZE      = 4  // Bit position for DRAM memory size (1MBx1 or 256KBx1/256KBx4)
	MCR_MEMORY_INTERLEAVE     = 6  // Bit position for memory interleave mode
	MCR_ROM_SIZE              = 8  // Bit position for ROM size (256KB or 512KB)
	MCR_ADAPTER_ROM_ENABLE    = 9  // Bit position for adapter ROM enable
	MCR_VIDEO_RAM_ENABLE      = 10 // Bit position for video RAM enable
	MCR_VIDEO_READ_ONLY       = 11 // Bit position for video read-only mode

	MCR_MEMORY_INTERLEAVE_1_BANK = 1 // Memory interleave mode: 1 bank
	MCR_MEMORY_INTERLEAVE_2_BANK = 2 // Memory interleave mode: 2 banks
	MCR_MEMORY_INTERLEAVE_3_BANK = 3 // Memory interleave mode: 3 banks
	MCR_MEMORY_INTERLEAVE_4_BANK = 4 // Memory interleave mode: 4 banks
)

type Intel82335 struct {
	busId uint32   // Device bus ID
	bus   *bus.Bus // Reference to the bus

	biosRomAccessEnabled bool  // 0 = enable BIOS ROM access, 1 = disable BIOS ROM access and enable BIOS shadow
	s640BaseMemorySize   bool  // 0 = 512KB, 1 = 640KB
	dRamSize             bool  // 0 = 1MBx1 DRAM, 1 = 256KBx1 or 256KBx4 DRAM
	romSize              bool  // 0 = 256KB ROM, 1 = 512KB ROM
	adapterRomEnabled    bool  // 0 = enable adapter ROM access, 1 = disable adapter ROM access and enable shadow
	videoRamEnabled      bool  // 0 = enable video RAM access, 1 = disable video RAM access
	videoReadOnly        bool  // 0 = video read-write, 1 = video read-only
	memoryInterleaving   uint8 // Memory interleaving mode

	controlRegister uint8   // Control register value
	mcrRegisters    []uint8 // MCR registers

	// RC1 module
	rc1RollCompareRegister uint8 // RC1 roll compare register

	// DMA module
	dmaCommandRegister uint8 // DMA command register
}

func NewIntel82335() *Intel82335 {
	chip := &Intel82335{
		mcrRegisters:           make([]uint8, 0x100),
		rc1RollCompareRegister: 0x00,
	}
	return chip
}

func (controller *Intel82335) GetBus() *bus.Bus {
	return controller.bus
}

func (controller *Intel82335) SetBus(bus *bus.Bus) {
	controller.bus = bus
}

func (device *Intel82335) GetDeviceBusId() uint32 {
	return device.busId
}

func (device *Intel82335) SetDeviceBusId(id uint32) {
	device.busId = id
}

func (device *Intel82335) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (controller *Intel82335) GetPortMap() *bus.DevicePortMap {
	return &bus.DevicePortMap{
		ReadPorts:  []uint16{0x0022, 0x0024},
		WritePorts: []uint16{0x0022, 0x0024},
	}
}

func (controller *Intel82335) ReadAddr8(addr uint16) uint8 {
	switch addr {
	case 0x0022:
		// Example: Read a system configuration register or similar
		return controller.ReadSelectedRegister()
	case 0x0024:
		// Example: Read another important register
		return controller.Rc1RegisterRead()
	default:
		log.Printf("Intel82335: Invalid read address: %#04x", addr)
		return 0
	}
}

func (controller *Intel82335) WriteAddr8(addr uint16, data uint8) {
	switch addr {
	case 0x0022:
		// Example: Write to a system configuration register
		controller.McrRegisterInitialize(data)
	case 0x0024:
		// Example: Write to another control register
		controller.Rc1RegisterWrite(data)
	default:
		log.Printf("Intel82335: Invalid write address: %#04x", addr)
	}
}

func (device *Intel82335) McrRegisterInitialize(registerValue uint8) {
	device.controlRegister = registerValue

	// Extract configuration bits from the MCR register value
	device.biosRomAccessEnabled = getRegisterBit(registerValue, MCR_BIOS_ROM_ENABLE)
	device.s640BaseMemorySize = getRegisterBit(registerValue, MCR_S640_BASE_MEMORY_SIZE)
	device.dRamSize = getRegisterBit(registerValue, MCR_DRAM_MEMORY_SIZE)
	device.memoryInterleaving = getMemoryInterleaveMode(registerValue)
	device.romSize = getRegisterBit(registerValue, MCR_ROM_SIZE)
	device.adapterRomEnabled = getRegisterBit(registerValue, MCR_ADAPTER_ROM_ENABLE)
	device.videoRamEnabled = getRegisterBit(registerValue, MCR_VIDEO_RAM_ENABLE)
	device.videoReadOnly = getRegisterBit(registerValue, MCR_VIDEO_READ_ONLY)

	// Log the updated configuration
	log.Printf("MCR Set Config: %s", device.toString())

	if device.biosRomAccessEnabled {
		log.Printf("BIOS ROM Access Enabled")

		// Send a message to the bus to notify other devices of the change
		device.bus.SendMessage(bus.BusMessage{
			Subject: common.MESSAGE_BIOS_ROM_ACCESS_ENABLED,
			Sender:  device.busId,
			Data:    []byte{0},
		})
	} else {
		log.Printf("BIOS ROM Access Disabled, BIOS Shadow Enabled")

		// Send a message to the bus to notify other devices of the change
		device.bus.SendMessage(bus.BusMessage{
			Subject: common.MESSAGE_BIOS_ROM_ACCESS_DISABLED,
			Sender:  device.busId,
			Data:    []byte{0},
		})
	}
}

func (device *Intel82335) toString() string {
	var strs []string

	// Append configuration strings based on the current settings
	if device.biosRomAccessEnabled {
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

	// Join the configuration strings and return the result
	return strings.Join(strs, ",")
}

func (device *Intel82335) ReadSelectedRegister() uint8 {
	// Read the value of the selected MCR register based on the control register value
	return device.mcrRegisters[device.controlRegister]
}

func (controller *Intel82335) Rc1RegisterRead() uint8 {
	// Read the value of the RC1 roll compare register
	return controller.rc1RollCompareRegister
}

func (controller *Intel82335) Rc1RegisterWrite(value uint8) {
	// Write the value to the RC1 roll compare register
	controller.rc1RollCompareRegister = value

	// Log the updated value
	log.Printf("RC1 Register Write: %#02x", value)

	// Send a message to the bus to notify other devices of the change
	controller.bus.SendMessage(bus.BusMessage{
		Subject: common.MESSAGE_RC1_REGISTER_UPDATE,
		Sender:  controller.busId,
		Data:    []byte{value},
	})

}

func (controller *Intel82335) DmaCommandRegisterWrite(value uint8) {
	/*
	   Bitfields for DMA channel 0-3 command register:
	   Bit(s)	Description	(Table P002)
	    7	DACK sense active high
	    6	DREQ sense active high
	    5	=1 extended write selection
	       =0 late write selection
	    4	rotating priority instead of fixed priority
	    3	compressed timing
	    2	=1 enable controller
	       =0 enable memory-to-memory
	    1-0	channel number
	   SeeAlso: #P001,#P004,#P005,#P079
	*/
	// Write the value to the DMA command register
	controller.dmaCommandRegister = value
}

func getMemoryInterleaveMode(registervalue uint8) uint8 {
	// Extract the memory interleave mode bits from the register value
	o2 := getRegisterBit(registervalue, MCR_MEMORY_INTERLEAVE)
	o1 := getRegisterBit(registervalue, MCR_MEMORY_INTERLEAVE+1)

	// Determine the memory interleave mode based on the bit values
	switch {
	case !o1 && !o2:
		return MCR_MEMORY_INTERLEAVE_1_BANK
	case !o1 && o2:
		return MCR_MEMORY_INTERLEAVE_2_BANK
	case o1 && !o2:
		return MCR_MEMORY_INTERLEAVE_3_BANK
	case o1 && o2:
		return MCR_MEMORY_INTERLEAVE_4_BANK
	default:
		log.Fatalln("Invalid memory interleave mode configuration!")
		panic(0)
	}
}

func getRegisterBit(source uint8, position uint8) bool {
	// Extract the bit value at the specified position from the source value
	return (source>>position)&1 == 1
}
