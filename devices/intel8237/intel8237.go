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

	addressRegisters  [4]uint16 // Address registers for each DMA channel
	countRegisters    [4]uint16 // Count registers for each DMA channel
	statusRegister    uint8     // Status register
	commandRegister   uint8     // Command register
	requestRegister   uint8     // Request register
	maskRegister      uint8     // Mask register
	modeRegisters     [4]uint8  // Mode registers for each DMA channel
	flipFlop          bool      // Byte pointer flip-flop
	temporaryRegister uint8     // Temporary register
}

func NewIntel8237() *Intel8237 {
	return &Intel8237{}
}

func (d *Intel8237) SetDeviceBusId(id uint32) {
	d.busId = id
}

func (d *Intel8237) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (d *Intel8237) GetPortMap() *bus.DevicePortMap {
	return nil
}

func (d *Intel8237) ReadAddr8(addr uint16) uint8 {
	//TODO implement me
	panic("implement me")
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

// Additional functions for DMA transfer functionality

func (d *Intel8237) StartDMATransfer(channel uint8, sourceAddr uint32, destinationAddr uint32, count uint16) {
	// Implement the logic to start a DMA transfer on the specified channel
	// with the given source address, destination address, and transfer count.

	// Set the address registers for the specified channel
	d.addressRegisters[channel] = uint16(sourceAddr)

	// Set the count registers for the specified channel
	d.countRegisters[channel] = count

	// Set the request flag for the specified channel
	d.requestRegister |= 1 << channel

	// Start the DMA transfer
	d.HandleDMATransfer(channel)
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

	// Perform the DMA transfer
	sourceAddr := uint32(d.addressRegisters[channel])
	destinationAddr := sourceAddr + 1 // Assumes destination address is next to source address
	count := d.countRegisters[channel]

	memoryController := d.bus.FindSingleDevice(0).(*memmap.MemoryAccessController)
	mode := d.modeRegisters[channel]
	transferType := (mode >> 2) & 0x03
	if transferType == 0x01 {
		// Write transfer
		for i := uint16(0); i < count; i++ {
			data, _ := memoryController.ReadMemoryAddr8(uint32(sourceAddr))
			err := memoryController.WriteMemoryAddr8(uint32(destinationAddr), data)
			if err != nil {
				log.Printf("DMA Write Error: %v", err)
				return
			}
			sourceAddr++
			destinationAddr++
		}
	} else if transferType == 0x02 {
		// Read transfer
		for i := uint16(0); i < count; i++ {
			data, _ := memoryController.ReadMemoryAddr8(uint32(sourceAddr))
			err := memoryController.WriteMemoryAddr8(uint32(destinationAddr), data)
			if err != nil {
				log.Printf("DMA Read Error: %v", err)
				return
			}
			sourceAddr++
			destinationAddr++
		}
	}

	// Update the address and count registers
	d.addressRegisters[channel] = uint16(sourceAddr)
	d.countRegisters[channel] -= count

	// Check if the DMA transfer is complete
	if d.IsDMATransferComplete(channel) {
		// Clear the request and status flags for the specified channel
		d.requestRegister &= ^(1 << channel)
		d.statusRegister &= ^(1 << channel)
	}
}
