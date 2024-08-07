package intel8237

import (
	"github.com/andrewjc/threeatesix/devices/bus"
	"github.com/andrewjc/threeatesix/devices/memmap"
	"log"
)

/*
   The Intel 8237 is a Direct Memory Access (DMA) controller chip used in IBM PC-compatible computers.
   It provides four DMA channels that can be used to transfer data between memory and peripheral devices
   without the intervention of the CPU.
*/

type Intel8237 struct {
	bus   *bus.Bus
	busId uint32

	isPrimaryDevice   bool
	isSecondaryDevice bool

	addressRegisters  [4]uint16 // Address registers for each DMA channel
	countRegisters    [4]uint16 // Count registers for each DMA channel
	statusRegister    uint8     // Status register
	commandRegister   uint8     // Command register
	requestRegister   uint8     // Request register
	maskRegister      uint8     // Mask register
	modeRegisters     [4]uint8  // Mode registers for each DMA channel
	flipFlop          bool      // Byte pointer flip-flop
	temporaryRegister uint8     // Temporary register

	pageRegisters [4]uint8
}

func NewIntel8237() *Intel8237 {
	return &Intel8237{}
}

func (d *Intel8237) GetDeviceBusId() uint32 {
	return d.busId
}

func (d *Intel8237) SetDeviceBusId(id uint32) {
	d.busId = id
}

func (d *Intel8237) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (d *Intel8237) GetPortMap() *bus.DevicePortMap {
	if d.isPrimaryDevice {
		return &bus.DevicePortMap{
			ReadPorts: []uint16{
				0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007, // Channel registers
				0x0008,                         // Status register
				0x000D,                         // Temporary register
				0x0087, 0x0083, 0x0081, 0x0082, // Page registers
			},
			WritePorts: []uint16{
				0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007, // Channel registers
				0x0008,                         // Command register
				0x0009,                         // Request register
				0x000A,                         // Single mask register bit
				0x000B,                         // Mode register
				0x000C,                         // Clear byte pointer flip-flop
				0x000D,                         // Master clear / Temporary register
				0x000E,                         // Clear mask register
				0x000F,                         // Write all mask register bits
				0x0087, 0x0083, 0x0081, 0x0082, // Page registers
			},
		}
	} else if d.isSecondaryDevice {
		return &bus.DevicePortMap{
			ReadPorts: []uint16{
				0x00C0, 0x00C2, 0x00C4, 0x00C6, 0x00C8, 0x00CA, 0x00CC, 0x00CE, // Channel registers
				0x00D0,                         // Status register
				0x00DA,                         // Temporary register
				0x008F, 0x008B, 0x0089, 0x008A, // Page registers
			},
			WritePorts: []uint16{
				0x00C0, 0x00C2, 0x00C4, 0x00C6, 0x00C8, 0x00CA, 0x00CC, 0x00CE, // Channel registers
				0x00D0,                         // Command register
				0x00D2,                         // Request register
				0x00D4,                         // Single mask register bit
				0x00D6,                         // Mode register
				0x00D8,                         // Clear byte pointer flip-flop
				0x00DA,                         // Master clear / Temporary register
				0x00DC,                         // Clear mask register
				0x00DE,                         // Write all mask register bits
				0x008F, 0x008B, 0x0089, 0x008A, // Page registers
			},
		}
	}
	return nil
}

func (d *Intel8237) ReadAddr8(addr uint16) uint8 {
	if d.isPrimaryDevice {
		switch addr {
		case 0x0008:
			return d.ReadStatusRegister()
		case 0x000D:
			return d.ReadTemporaryRegister()
		}
	} else if d.isSecondaryDevice {
		switch addr {
		case 0x00D0:
			return d.ReadStatusRegister()
		case 0x00DA:
			return d.ReadTemporaryRegister()
		}
	}

	log.Fatalf("Intel8237 Invalid read address: %#04x", addr)
	panic(0)
}

func (d *Intel8237) WriteAddr8(addr uint16, data uint8) {

	// For address and count registers
	if addr >= 0x0000 && addr <= 0x0007 {
		channel := (addr >> 1) & 0x03
		if d.flipFlop {
			// High byte
			if addr&1 == 0 {
				d.addressRegisters[channel] = (d.addressRegisters[channel] & 0x00FF) | (uint16(data) << 8)
			} else {
				d.countRegisters[channel] = (d.countRegisters[channel] & 0x00FF) | (uint16(data) << 8)
			}
		} else {
			// Low byte
			if addr&1 == 0 {
				d.addressRegisters[channel] = (d.addressRegisters[channel] & 0xFF00) | uint16(data)
			} else {
				d.countRegisters[channel] = (d.countRegisters[channel] & 0xFF00) | uint16(data)
			}
		}
		d.flipFlop = !d.flipFlop
	}

	// Page registers for primary DMA controller
	if addr >= 0x0087 && addr <= 0x008F {
		channel := (addr - 0x0087) >> 1
		if channel < 4 {
			d.pageRegisters[channel] = data
		}
	}

	if d.isPrimaryDevice {
		switch addr {
		case 0x0000, 0x0002, 0x0004, 0x0006:
			// DMA channel memory address bytes 1 and 0 (low)
			channel := uint8(addr / 2)
			d.addressRegisters[channel] = (d.addressRegisters[channel] & 0xFF00) | uint16(data)
		case 0x0001, 0x0003, 0x0005, 0x0007:
			// DMA channel memory address bytes 1 and 0 (high)
			channel := uint8((addr - 1) / 2)
			d.addressRegisters[channel] = (uint16(data) << 8) | (d.addressRegisters[channel] & 0x00FF)
		case 0x0008:
			d.WriteCommandRegister(data)
		case 0x0009:
			d.WriteRequestRegister(data)
		case 0x000A:
			d.WriteSingleMaskRegister(data)
		case 0x000B:
			d.WriteModeRegister(data)
		case 0x000C:
			d.ClearBytePointerFlipFlop()
		case 0x000D:
			// ignored, write to temporary register
		case 0x000E:
			d.ClearMaskRegister()
		case 0x000F:
			d.WriteMaskRegister(data)
		}
	} else if d.isSecondaryDevice {
		switch addr {
		case 0x00C0, 0x00C4, 0x00C8, 0x00CC:
			// DMA channel memory address bytes 1 and 0 (low)
			channel := uint8((addr - 0x00C0) / 4)
			d.addressRegisters[channel] = (d.addressRegisters[channel] & 0xFF00) | uint16(data)
		case 0x00C2, 0x00C6, 0x00CA, 0x00CE:
			// DMA channel memory address bytes 1 and 0 (high)
			channel := uint8((addr - 0x00C2) / 4)
			d.addressRegisters[channel] = (uint16(data) << 8) | (d.addressRegisters[channel] & 0x00FF)
		case 0x00D0:
			d.WriteCommandRegister(data)
		case 0x00D2:
			d.WriteRequestRegister(data)
		case 0x00D4:
			d.WriteSingleMaskRegister(data)
		case 0x00D6:
			d.WriteModeRegister(data)
		case 0x00D8:
			d.ClearBytePointerFlipFlop()
		case 0x00DA:
			// ignored, write to temporary register
		case 0x00DC:
			d.ClearMaskRegister()
		case 0x00DE:
			d.WriteMaskRegister(data)
		}
	}
}

func (d *Intel8237) ReadStatusRegister() uint8 {
	return d.statusRegister
}

func (d *Intel8237) WriteCommandRegister(value uint8) {
	d.commandRegister = value
	// Implement command register functionality
	// The command register is used to control the overall operation of the 8237 chip.
	// It allows setting the DACK and DREQ sense, extended write selection, priority mode,
	// timing mode, controller enable, and memory-to-memory transfer enable.

	// DACK sense: bit 7
	// 1 = DACK active high, 0 = DACK active low
	dackSense := (value & 0x80) != 0

	// DREQ sense: bit 6
	// 1 = DREQ active high, 0 = DREQ active low
	dreqSense := (value & 0x40) != 0

	// Extended write selection: bit 5
	// 1 = extended write selection, 0 = late write selection
	extendedWrite := (value & 0x20) != 0

	// Priority mode: bit 4
	// 1 = rotating priority, 0 = fixed priority
	priorityMode := (value & 0x10) != 0

	// Timing mode: bit 3
	// 1 = compressed timing, 0 = normal timing
	timingMode := (value & 0x08) != 0

	// Controller enable: bit 2
	// 1 = enable controller, 0 = enable memory-to-memory
	controllerEnable := (value & 0x04) != 0

	// Implement the specific functionality based on the command register settings.
	// TODO: Implement the specific functionality based on the command register settings.
	_ = dackSense
	_ = dreqSense
	_ = extendedWrite
	_ = priorityMode
	_ = timingMode
	_ = controllerEnable
}

func (d *Intel8237) WriteRequestRegister(value uint8) {
	d.requestRegister = value
	// Implement request register functionality
	// The request register is used to initiate DMA transfers on specific channels.
	// Writing a '1' to a bit in this register initiates a DMA request for the corresponding channel.

	for i := 0; i < 4; i++ {
		if (value & (1 << i)) != 0 {
			// Initiate DMA request for channel i
			d.HandleDMARequest(uint8(i))
		}
	}
}

func (d *Intel8237) WriteSingleMaskRegister(value uint8) {
	channelSelect := value & 0x03
	setMask := (value & 0x04) != 0

	// The single mask register is used to mask or unmask individual DMA channels.
	// Bits 0-1 select the DMA channel, and bit 2 determines whether to set or clear the mask.
	if setMask {
		d.maskRegister |= 1 << channelSelect // Set the mask bit for the selected channel
	} else {
		d.maskRegister &= ^(1 << channelSelect) // Clear the mask bit for the selected channel
	}
}

func (d *Intel8237) WriteModeRegister(value uint8) {
	channelSelect := value & 0x03
	d.modeRegisters[channelSelect] = value
	// Implement mode register functionality
	// The mode registers are used to configure the transfer mode, address increment/decrement,
	// autoinitialization, and transfer type for each DMA channel.

	// Transfer mode: bits 7-6
	// 00 = demand mode, 01 = single mode, 10 = block mode, 11 = cascade mode
	transferMode := (value >> 6) & 0x03

	// Address increment/decrement: bit 5
	// 0 = increment, 1 = decrement
	addressMode := (value & 0x20) != 0

	// Autoinitialization: bit 4
	// 0 = disable, 1 = enable
	autoInitialization := (value & 0x10) != 0

	// Transfer type: bits 3-2
	// 00 = verify, 01 = write, 10 = read, 11 = illegal
	transferType := (value >> 2) & 0x03

	// Implement the specific functionality based on the mode register settings.
	// TODO: Implement the specific functionality based on the mode register settings.
	_ = transferMode
	_ = addressMode
	_ = autoInitialization
	_ = transferType
}

func (d *Intel8237) ClearBytePointerFlipFlop() {
	// Implement byte pointer flip-flop functionality
	// The byte pointer flip-flop is used to determine whether the high or low byte of the
	// address and count registers should be accessed.
	// Clearing the flip-flop resets it to point to the low byte.
	d.flipFlop = false
}

func (d *Intel8237) ReadTemporaryRegister() uint8 {
	// Implement temporary register functionality
	// The temporary register is used to store temporary data during DMA transfers.
	return d.temporaryRegister
}

func (d *Intel8237) MasterClear() {
	// Implement master clear functionality
	// The master clear operation resets the 8237 chip to its initial state.
	// It clears all registers and stops any ongoing DMA transfers.

	// Clear address registers
	for i := 0; i < 4; i++ {
		d.addressRegisters[i] = 0
	}

	// Clear count registers
	for i := 0; i < 4; i++ {
		d.countRegisters[i] = 0
	}

	// Clear status register
	d.statusRegister = 0

	// Clear command register
	d.commandRegister = 0

	// Clear request register
	d.requestRegister = 0

	// Clear mask register
	d.maskRegister = 0

	// Clear mode registers
	for i := 0; i < 4; i++ {
		d.modeRegisters[i] = 0
	}

	// Clear byte pointer flip-flop
	d.flipFlop = false

	// Clear temporary register
	d.temporaryRegister = 0

	// Stop ongoing DMA transfers
	for i := 0; i < 4; i++ {
		d.StopDMATransfer(uint8(i))
	}
}

func (d *Intel8237) ClearMaskRegister() {
	d.maskRegister = 0
	// The clear mask register operation clears all bits in the mask register,
	// unmasking all DMA channels.
}

func (d *Intel8237) WriteMaskRegister(value uint8) {
	d.maskRegister = value
	// The write mask register operation sets the entire mask register to the specified value,
	// masking or unmasking DMA channels based on the corresponding bits.
}

func (d *Intel8237) GetBus() *bus.Bus {
	return d.bus
}

func (d *Intel8237) SetBus(bus *bus.Bus) {
	d.bus = bus
}

func (d *Intel8237) StopDMATransfer(channel uint8) {
	// Implement the logic to stop the DMA transfer on the specified channel.

	// Clear the request flag for the specified channel
	d.requestRegister &= ^(1 << channel)

	// Reset the address and count registers for the specified channel
	d.addressRegisters[channel] = 0
	d.countRegisters[channel] = 0
}

func (d *Intel8237) IsDMATransferComplete(channel uint8) bool {
	// Implement the logic to check if the DMA transfer on the specified channel is complete.

	// Check if the count register for the specified channel has reached zero
	return d.countRegisters[channel] == 0
}

func (d *Intel8237) HandleDMARequest(channel uint8) {
	// Implement the logic to handle a DMA request on the specified channel.

	// Check if the channel is masked
	if (d.maskRegister & (1 << channel)) != 0 {
		return
	}

	// Check if a DMA transfer is already in progress for the specified channel
	if (d.statusRegister & (1 << channel)) != 0 {
		return
	}

	// Set the status flag for the specified channel
	d.statusRegister |= 1 << channel

	// Start the DMA transfer for the specified channel
	sourceAddr := uint32(d.addressRegisters[channel])
	destinationAddr := sourceAddr + 1 // Assumes destination address is next to source address
	count := d.countRegisters[channel]
	d.StartDMATransfer(channel, sourceAddr, destinationAddr, count)
}

func (d *Intel8237) HandleDMATransfer(channel uint8) {

	// Check if the channel is masked
	if (d.maskRegister & (1 << channel)) != 0 {
		log.Printf("DMA Channel %d is masked", channel)
		return
	}

	// Check if a DMA transfer is already in progress for the specified channel
	if (d.statusRegister & (1 << channel)) == 0 {
		log.Printf("DMA Transfer not in progress for channel %d", channel)
		return
	}

	mode := d.modeRegisters[channel]
	transferMode := (mode >> 6) & 0x03
	addressMode := (mode & 0x20) != 0
	autoInitialization := (mode & 0x10) != 0
	transferType := (mode >> 2) & 0x03

	sourceAddr := uint32(d.addressRegisters[channel])
	count := d.countRegisters[channel]

	memoryController := d.bus.FindSingleDevice(0).(*memmap.MemoryAccessController)

	for i := uint16(0); i < count; i++ {
		switch transferType {
		case 0x01: // Write transfer
			data, _ := memoryController.ReadMemoryValue8(sourceAddr)
			err := memoryController.WriteMemoryAddr8(sourceAddr+1, data)
			if err != nil {
				log.Printf("DMA Write Error: %v", err)
				return
			}
		case 0x02: // Read transfer
			data, _ := memoryController.ReadMemoryValue8(sourceAddr)
			err := memoryController.WriteMemoryAddr8(sourceAddr+1, data)
			if err != nil {
				log.Printf("DMA Read Error: %v", err)
				return
			}
		}

		if addressMode {
			sourceAddr--
		} else {
			sourceAddr++
		}

		d.countRegisters[channel]--

		if transferMode == 0x01 { // Single mode
			break
		}
	}

	d.addressRegisters[channel] = uint16(sourceAddr)

	if d.countRegisters[channel] == 0 {
		if autoInitialization {
			// Reload original values
			// You need to store original values somewhere
		} else {
			d.StopDMATransfer(channel)
		}
	}

	if d.countRegisters[channel] == 0 {
		// Set terminal count bit in status register
		d.statusRegister |= 0x80

		if autoInitialization {
			// Reload original values
			// You need to store original values somewhere
		} else {
			d.StopDMATransfer(channel)
		}

		// Generate TC signal (this would typically trigger an interrupt)
		d.GenerateTerminalCount(channel)
	}
}

func (d *Intel8237) GenerateTerminalCount(channel uint8) {
	log.Printf("DMA channel %d reached terminal count", channel)
	// TODO: raise an interrupt here
}

func (d *Intel8237) IsPrimaryDevice(isPrimaryDevice bool) {
	d.isPrimaryDevice = isPrimaryDevice
	d.isSecondaryDevice = false
}

func (d *Intel8237) IsSecondaryDevice(isSecondaryDevice bool) {
	d.isPrimaryDevice = false
	d.isSecondaryDevice = isSecondaryDevice
}

func (d *Intel8237) WriteTemporaryRegister(data uint8) {
	d.temporaryRegister = data
}

func (d *Intel8237) AcknowledgeDMA(channel uint8) {
	// Clear the request bit
	d.requestRegister &= ^(1 << channel)

	// Set the corresponding bit in the status register
	d.statusRegister |= 1 << channel

	// Start the DMA transfer
	d.HandleDMATransfer(channel)
}

func (d *Intel8237) StartDMATransfer(channel uint8, sourceAddr uint32, destAddr uint32, count uint16) {
	// Ensure the channel is valid
	if channel > 3 {
		log.Printf("Error: Invalid DMA channel %d", channel)
		return
	}

	// Set up the address registers
	d.addressRegisters[channel] = uint16(sourceAddr & 0xFFFF)
	d.pageRegisters[channel] = uint8((sourceAddr >> 16) & 0xFF)

	// Set up the count register
	d.countRegisters[channel] = count

	// Set up the mode register
	// Assume single transfer mode, address increment, and write transfer for this example
	// You may want to make these configurable based on the specific needs of the transfer
	d.modeRegisters[channel] = 0x44 | channel // 01000100 | channel

	// Clear the byte pointer flip-flop
	d.ClearBytePointerFlipFlop()

	// Unmask the channel
	d.WriteSingleMaskRegister(channel)

	// Set the request bit for the channel
	d.WriteRequestRegister(0x04 | channel)

	// Log the start of the transfer
	log.Printf("Starting DMA transfer on channel %d", channel)
	log.Printf("Source address: %#08x", sourceAddr)
	log.Printf("Destination address: %#08x", destAddr)
	log.Printf("Count: %d", count)

	// Start the transfer
	d.HandleDMATransfer(channel)
}
