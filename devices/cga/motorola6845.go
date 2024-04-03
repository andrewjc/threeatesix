package cga

import (
	"github.com/andrewjc/threeatesix/devices/bus"
	"log"
)

/*
   Simulated motorola cga device
*/

type Motorola6845 struct {
	busId uint32

	BlinkEnabled        bool
	GraphicsMode640x200 bool
	VideoEnabled        bool
	MonochromeSignal    bool
	TextMode            bool
	GraphicsMode320x200 bool
	TextMode80x25       bool

	// CGA video memory
	videoMemory [0x4000]uint8

	// CGA palette
	palette [16]uint8

	// CGA cursor position
	cursorPosition uint16

	// CGA status register
	statusRegister uint8
}

func NewMotorola6845() *Motorola6845 {
	chip := &Motorola6845{}
	// Initialize video memory
	for i := range chip.videoMemory {
		chip.videoMemory[i] = 0
	}
	// Initialize palette with default CGA colors
	chip.palette = [16]uint8{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
	}
	return chip
}

func (c *Motorola6845) SetDeviceBusId(id uint32) {
	c.busId = id
}

func (c *Motorola6845) OnReceiveMessage(message bus.BusMessage) {
	// Handle bus messages if needed
}

func (c *Motorola6845) WriteAddr8(port_addr uint16, value uint8) {
	switch port_addr {
	case 0x03D8:
		// CGA Controller Mode Control Register
		c.SetModeControlRegister(value)
	case 0x03D9:
		// CGA Controller Color Select Register
		c.SetColorSelectRegister(value)
	case 0x03DA:
		// CGA Status Register (read-only)
		log.Printf("Attempted write to read-only CGA Status Register")
	default:
		log.Printf("Unhandled CGA write to port %#04x with value %#02x", port_addr, value)
	}
}

func (c *Motorola6845) ReadAddr8(port_addr uint16) uint8 {
	switch port_addr {
	case 0x03DA:
		// CGA Status Register
		return c.ReadStatusRegister()
	default:
		log.Printf("Unhandled CGA read from port %#04x", port_addr)
		return 0
	}
}

func (c *Motorola6845) SetModeControlRegister(value uint8) {
	/*
	   03D8    r/w    CGA mode control register  (except PCjr)
	           bit 7-6      not used
	           bit 5  = 1  blink enabled
	           bit 4  = 1  640*200 graphics mode
	           bit 3  = 1  video enabled
	           bit 2  = 1  monochrome signal
	           bit 1  = 0  text mode
	              = 1  320*200 graphics mode
	           bit 0  = 0  40*25 text mode
	              = 1  80*25 text mode
	*/
	log.Printf("CGA Controller set mode control register: [%#04x]", value)

	c.BlinkEnabled = (value & 0x20) != 0
	c.GraphicsMode640x200 = (value & 0x10) != 0
	c.VideoEnabled = (value & 0x08) != 0
	c.MonochromeSignal = (value & 0x04) != 0
	c.TextMode = (value & 0x02) == 0
	c.GraphicsMode320x200 = (value & 0x02) != 0
	c.TextMode80x25 = (value & 0x01) != 0

	log.Printf("CGA Controller set mode control register: BlinkEnabled: %v, GraphicsMode640x200: %v, VideoEnabled: %v, MonochromeSignal: %v, TextMode: %v, GraphicsMode320x200: %v, TextMode80x25: %v", c.BlinkEnabled, c.GraphicsMode640x200, c.VideoEnabled, c.MonochromeSignal, c.TextMode, c.GraphicsMode320x200, c.TextMode80x25)
}

func (c *Motorola6845) SetColorSelectRegister(value uint8) {
	/*
	   03D9    r/w    CGA color select register
	           bit 7-6  background color
	           bit 5-4  palette
	           bit 3    high intensity
	           bit 2-0  foreground color
	*/
	log.Printf("CGA Controller set color select register: [%#04x]", value)

	backgroundColor := (value >> 6) & 0x03
	palette := (value >> 4) & 0x03
	highIntensity := (value & 0x08) != 0
	foregroundColor := value & 0x07

	log.Printf("CGA Controller set color select register: BackgroundColor: %d, Palette: %d, HighIntensity: %v, ForegroundColor: %d", backgroundColor, palette, highIntensity, foregroundColor)
}

func (c *Motorola6845) ReadStatusRegister() uint8 {
	/*
	   03DA    r      CGA status register
	           bit 7-4  not used
	           bit 3 = 1  in vertical retrace
	           bit 2-0  not used
	*/
	// Check if the CGA is in vertical retrace
	if c.isInVerticalRetrace() {
		c.statusRegister |= 0x08
	} else {
		c.statusRegister &^= 0x08
	}

	return c.statusRegister
}

func (c *Motorola6845) GetBus() *bus.Bus {
	// TODO: Implement GetBus method to return the bus associated with the CGA device
	return nil
}

func (c *Motorola6845) SetBus(bus *bus.Bus) {
	// TODO: Implement SetBus method to set the bus associated with the CGA device
}

func (c *Motorola6845) isInVerticalRetrace() bool {
	// Get the current scanline number
	scanline := c.getCurrentScanline()

	// Get the vertical sync and vertical total values from the respective registers
	verticalSync := c.getVerticalSyncRegister()
	verticalTotal := c.getVerticalTotalRegister()

	// Check if the current scanline is within the vertical retrace period
	if scanline >= verticalSync && scanline < verticalTotal {
		return true
	}

	return false
}

func (c *Motorola6845) getCurrentScanline() uint16 {
	// The current scanline number is determined by the value in the vertical sync position register
	// The vertical sync position register contains the scanline number where the vertical sync pulse starts

	// Read the vertical sync position register (register index 7)
	verticalSyncPosition := c.readRegister(7)

	// Consider the vertical sync position as the current scanline
	return verticalSyncPosition
}

func (c *Motorola6845) getVerticalSyncRegister() uint16 {
	// The vertical sync register (register index 7) contains the scanline number where the vertical sync pulse starts
	return c.readRegister(7)
}

func (c *Motorola6845) getVerticalTotalRegister() uint16 {
	// The vertical total register (register index 4) contains the total number of scanlines in a frame
	return c.readRegister(4)
}

func (c *Motorola6845) readRegister(index uint8) uint16 {
	// The Motorola 6845 has 18 registers, each 8 bits wide
	// The register index is written to port 0x03D4, and the register value is read from port 0x03D5

	// Write the register index to port 0x03D4
	c.WriteRegisterIndex(0x03D4, index)

	// Read the register value from port 0x03D5
	value := c.ReadRegisterValue(0x03D5)

	return uint16(value)
}

func (c *Motorola6845) WriteRegisterIndex(port uint16, index uint8) {
	// Write the register index to the specified port
	c.WritePort(port, index)
}

func (c *Motorola6845) ReadRegisterValue(port uint16) uint8 {
	// Read the register value from the specified port
	return c.ReadPort(port)
}

func (c *Motorola6845) WritePort(port uint16, value uint8) {
	// TODO: Implement the logic to write the value to the specified port
	// This would involve writing the value to the port using the appropriate I/O mechanism
}

func (c *Motorola6845) ReadPort(port uint16) uint8 {
	// TODO: Implement the logic to read the value from the specified port
	// This would involve reading the value from the port using the appropriate I/O mechanism
	return 0
}
