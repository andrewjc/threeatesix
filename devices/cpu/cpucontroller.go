package cpu

type CpuController interface {
	GetIP() uint16
	GetCS() uint16
	IncrementIP()
	SetIP(addr uint16)
	SetCS(addr uint16)
	EnterMode(mode uint8)
}
