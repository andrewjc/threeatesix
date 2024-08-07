package pc

import (
	"fmt"
	"github.com/andrewjc/threeatesix/common"
	"github.com/andrewjc/threeatesix/devices/bus"
	"github.com/andrewjc/threeatesix/devices/cga"
	"github.com/andrewjc/threeatesix/devices/cmos"
	"github.com/andrewjc/threeatesix/devices/hid/kb"
	"github.com/andrewjc/threeatesix/devices/intel8086"
	"github.com/andrewjc/threeatesix/devices/intel82335"
	"github.com/andrewjc/threeatesix/devices/intel8237"
	"github.com/andrewjc/threeatesix/devices/intel8259a"
	"github.com/andrewjc/threeatesix/devices/intel82C54"
	"github.com/andrewjc/threeatesix/devices/io"
	"github.com/andrewjc/threeatesix/devices/memmap"
	"github.com/andrewjc/threeatesix/devices/monitor"
	"github.com/andrewjc/threeatesix/devices/ps2"
	"io/ioutil"
	"log"
	"os"
)

type romimages struct {
	bios []byte
	vga  []byte
}

// PersonalComputer represents the virtual PC being emulated
type PersonalComputer struct {
	cpu             *intel8086.CpuCore
	mathCoProcessor *intel8086.CpuCore

	bus *bus.Bus

	ram []byte
	rom romimages

	programmableInterruptController1 *intel8259a.Intel8259a
	programmableInterruptController2 *intel8259a.Intel8259a
	programmableIntervalTimer        *intel82C54.Intel82C54

	memController    *memmap.MemoryAccessController
	ioPortController *io.IOPortAccessController

	ps2Controller *ps2.Ps2Controller

	hardwareMonitor                *monitor.HardwareMonitor
	cgaController                  *cga.Motorola6845
	cmos                           *cmos.Motorola146818
	highIntegrationInterfaceDevice *intel82335.Intel82335
	dmaController                  *intel8237.Intel8237
	dmaController2                 *intel8237.Intel8237
}

// BiosFilename - name of the bios image the virtual machine will boot up
const BiosFilename = "bios/ami386.bin"
const VideoBiosFilename = "bios/vgabios.bin"

// MaxRAMBytes - the amount of ram installed in this virtual machine
const MaxRAMBytes = 0x1E84800 //32 million (32mb)
//const MaxRAMBytes = 0xF42400 //8mb
//const MaxRAMBytes = 0x100000000 //4GB

func (pc *PersonalComputer) Power() {
	// do stuff
	log.SetFlags(log.Lshortfile)

	pc.cpu.Init(pc.bus)
	pc.mathCoProcessor.Init(pc.bus)

	// hack a hard drive into cmos settings
	pc.cmos.WriteAddr8(0x70, 0x12)
	pc.cmos.WriteAddr8(0x71, 0x01)

	for {
		if pc.cpu.GetIP() == 0x0 {
			log.Printf("Instruction pointer is 0, halting")
			break
		} //loop until instruction pointer equals 0

		pc.cpu.Step()
		//pc.mathCoProcessor.Step()
		pc.programmableIntervalTimer.Step()
	}
}

func NewPc() *PersonalComputer {
	pc := &PersonalComputer{}

	pc.bus = bus.NewDeviceBus()
	pc.ram = make([]byte, MaxRAMBytes)
	pc.rom = romimages{}
	pc.cpu = intel8086.New80386CPU()
	pc.mathCoProcessor = intel8086.New80287MathCoProcessor()

	pc.programmableInterruptController1 = intel8259a.NewIntel8259a() //pic1
	pc.programmableInterruptController2 = intel8259a.NewIntel8259a() //pic2
	pc.programmableInterruptController1.IsPrimaryDevice(true)
	pc.programmableInterruptController2.IsSecondaryDevice(true)

	pc.programmableIntervalTimer = intel82C54.NewIntel82C54()      //pit
	pc.highIntegrationInterfaceDevice = intel82335.NewIntel82335() //hiid
	pc.dmaController = intel8237.NewIntel8237()                    //dma
	pc.dmaController2 = intel8237.NewIntel8237()                   //dma2
	pc.dmaController.IsPrimaryDevice(true)
	pc.dmaController2.IsSecondaryDevice(true)

	pc.cgaController = cga.NewMotorola6845()
	pc.cmos = cmos.NewMotorola146818()

	pc.memController = memmap.NewMemoryController(&pc.ram, &pc.rom.bios, &pc.rom.vga)

	pc.ioPortController = io.NewIOPortController()

	pc.memController.SetBus(pc.bus)
	pc.ioPortController.SetBus(pc.bus)

	pc.ps2Controller = ps2.CreatePS2Controller()
	pc.ps2Controller.SetBus(pc.bus)

	pc.hardwareMonitor = monitor.NewHardwareMonitor()

	pc.bus.RegisterDevice(pc.hardwareMonitor, common.MODULE_DEBUG_MONITOR)

	pc.bus.RegisterDevice(pc.cpu, common.MODULE_PRIMARY_PROCESSOR)
	pc.bus.RegisterDevice(pc.mathCoProcessor, common.MODULE_MATH_CO_PROCESSOR)
	pc.bus.RegisterDevice(pc.programmableInterruptController1, common.MODULE_INTERRUPT_CONTROLLER_1)
	pc.bus.RegisterDevice(pc.programmableInterruptController2, common.MODULE_INTERRUPT_CONTROLLER_2)
	pc.bus.RegisterDevice(pc.programmableIntervalTimer, common.MODULE_PIT)
	pc.bus.RegisterDevice(pc.highIntegrationInterfaceDevice, common.MODULE_INTEL_82335)
	pc.bus.RegisterDevice(pc.cgaController, common.MODULE_CGA)
	pc.bus.RegisterDevice(pc.cmos, common.MODULE_CMOS)
	pc.bus.RegisterDevice(pc.dmaController, common.MODULE_DMA_CONTROLLER)
	pc.bus.RegisterDevice(pc.dmaController2, common.MODULE_DMA_CONTROLLER_2)

	pc.bus.RegisterDevice(pc.memController, common.MODULE_MEMORY_ACCESS_CONTROLLER)
	pc.bus.RegisterDevice(pc.ioPortController, common.MODULE_IO_PORT_ACCESS_CONTROLLER)

	pc.bus.RegisterDevice(pc.ps2Controller, common.MODULE_PS2_CONTROLLER)

	pc.ps2Controller.ConnectDevice(kb.NewPs2Keyboard())

	return pc
}

func (pc *PersonalComputer) GetPrimaryCpu() *intel8086.CpuCore {
	return pc.cpu
}

func (pc *PersonalComputer) GetMemoryController() *memmap.MemoryAccessController {
	return pc.memController
}

func (pc *PersonalComputer) GetBus() *bus.Bus {
	return pc.bus
}

func (pc *PersonalComputer) LoadBios() {

	biosData, err := ioutil.ReadFile(BiosFilename)
	if err != nil {
		fmt.Printf("Failed to load BIOS: %s\n", err)
		os.Exit(1) // Consider handling this more softly, perhaps returning an error
	}

	pc.rom.bios = biosData

	videoBiosData, err := ioutil.ReadFile(VideoBiosFilename)
	if err != nil {
		fmt.Printf("Failed to load Video BIOS: %s\n", err)
		os.Exit(1) // Consider handling this more softly, perhaps returning an error
	}

	pc.rom.vga = videoBiosData
}
