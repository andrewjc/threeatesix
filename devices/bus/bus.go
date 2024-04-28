package bus

import (
	"container/list"
	"github.com/google/uuid"
	"log"
)

type DeviceType uint8

type Bus struct {
	deviceMap     map[DeviceType]*list.List
	devicePortMap map[uint16]*DevicePort
}

type BusMessage struct {
	Subject uint32
	Sender  uint32
	Data    []byte
}

type DevicePortMap struct {
	ReadPorts  []uint16
	WritePorts []uint16
}

type BusDevice interface {
	GetDeviceBusId() uint32
	SetDeviceBusId(id uint32)
	OnReceiveMessage(message BusMessage)
	SetBus(bus *Bus)
	GetPortMap() *DevicePortMap

	ReadAddr8(addr uint16) uint8
	WriteAddr8(addr uint16, data uint8)
}

func NewDeviceBus() *Bus {
	bus := &Bus{}

	bus.deviceMap = make(map[DeviceType]*list.List)
	bus.devicePortMap = make(map[uint16]*DevicePort)

	return bus
}

type DevicePortMode uint8

const DEVICE_PORT_READ = DevicePortMode(0)
const DEVICE_PORT_WRITE = DevicePortMode(1)

type DevicePort struct {
	Device BusDevice
	Port   uint16
	Mode   DevicePortMode
}

func (bus *Bus) GetDeviceOnPort(addr uint16) *DevicePort {
	if dev, ok := bus.devicePortMap[addr]; ok {
		return dev
	}

	return nil

}

func (bus *Bus) RegisterDevice(device BusDevice, deviceType DeviceType) {
	if _, ok := bus.deviceMap[deviceType]; !ok {
		bus.deviceMap[deviceType] = list.New()
	}

	// register the device ports
	portMap := device.GetPortMap()
	if portMap != nil {
		for _, port := range portMap.ReadPorts {
			bus.devicePortMap[port] = &DevicePort{Device: device, Port: port, Mode: DEVICE_PORT_READ}
		}

		for _, port := range portMap.WritePorts {
			bus.devicePortMap[port] = &DevicePort{Device: device, Port: port, Mode: DEVICE_PORT_WRITE}
		}
	}

	deviceList := bus.deviceMap[deviceType]
	device.SetDeviceBusId(getRandomUUID())
	device.SetBus(bus)

	deviceList.PushBack(device)
}

func getRandomUUID() uint32 {
	uuid, _ := uuid.NewRandom()
	return uuid.ID()
}

func (bus *Bus) FindDevice(deviceType DeviceType) *list.List {
	if device, ok := bus.deviceMap[deviceType]; ok {
		return device
	} else {
		log.Fatalf("Could not find device on bus of type %v", deviceType)
		return nil
	}
}

func (bus *Bus) FindSingleDevice(deviceType DeviceType) BusDevice {
	deviceList := bus.FindDevice(deviceType)
	return deviceList.Front().Value.(BusDevice)
}

// Sends a message to all devices on the bus
func (bus *Bus) SendMessage(message BusMessage) {
	for e, _ := range bus.deviceMap {
		devList := bus.deviceMap[e]
		for dev := devList.Front(); dev != nil; dev = dev.Next() {
			dev.Value.(BusDevice).OnReceiveMessage(message)
		}
	}
}

func (bus *Bus) SendMessageToAll(deviceType DeviceType, message BusMessage) error {
	if devList, ok := bus.deviceMap[deviceType]; ok {
		for dev := devList.Front(); dev != nil; dev = dev.Next() {
			dev.Value.(BusDevice).OnReceiveMessage(message)
		}
	} else {
		log.Fatalf("Could not find device on bus of type %v", deviceType)
	}

	return nil
}

func (bus *Bus) SendMessageSingle(deviceType DeviceType, message BusMessage) error {
	bus.FindSingleDevice(deviceType).OnReceiveMessage(message)

	return nil
}
func (bus *Bus) SendMessageToDeviceById(deviceId uint32, message BusMessage) error {
	for e, _ := range bus.deviceMap {
		devList := bus.deviceMap[e]
		for dev := devList.Front(); dev != nil; dev = dev.Next() {
			if dev.Value.(BusDevice).GetDeviceBusId() == deviceId {
				dev.Value.(BusDevice).OnReceiveMessage(message)
				return nil
			}
		}
	}

	return nil

}
