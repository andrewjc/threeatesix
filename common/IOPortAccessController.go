package common

import "log"

/*
	IO Port Access Controller
	Provides read/write functions for port mapped IO
*/

type IOPortAccessController struct {
	backingMemory []byte
}

func (r *IOPortAccessController) ReadAddr8(addr uint16) uint8 {
	var byteData uint8
	byteData = (r.backingMemory)[addr]

	return byteData
}

func (r *IOPortAccessController) WriteAddr8(addr uint16, value uint8) {
	r.backingMemory[addr] = value
}

func (r *IOPortAccessController) ReadAddr16(addr uint16) uint16 {
	b1 := uint16(r.ReadAddr8(addr))
	b2 := uint16(r.ReadAddr8(addr + 1))
	return b2<<8 | b1
}

func (r *IOPortAccessController) WriteAddr16(addr uint16, value uint16) {
	log.Fatal("TODO!")
}

func CreateIOPortController() *IOPortAccessController{
	return &IOPortAccessController{backingMemory:make([]byte, 0x10000)}
}
