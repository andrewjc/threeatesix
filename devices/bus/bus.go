package bus

import (
	"container/list"
	"github.com/google/uuid"
	"log"
)

type DeviceType uint8

type Bus struct {
	deviceMap map[DeviceType]*list.List
}

type BusMessage struct {
	Subject uint32
	Sender  uint32
	Data    []byte
}

type BusDevice interface {
	SetDeviceBusId(id uint32)
	OnReceiveMessage(message BusMessage)
	SetBus(bus *Bus)
}

func NewDeviceBus() *Bus {
	bus := &Bus{}

	bus.deviceMap = make(map[DeviceType]*list.List)

	return bus
}

func (bus *Bus) RegisterDevice(device BusDevice, deviceType DeviceType) {

	if _, ok := bus.deviceMap[deviceType]; !ok {
		bus.deviceMap[deviceType] = list.New()
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
