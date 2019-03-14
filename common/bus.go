package common

type Bus struct {
}

type BusDevice interface {
}

func (bus *Bus) registerDevice(device *BusDevice, deviceType uint8) {

}
