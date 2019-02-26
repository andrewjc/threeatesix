package common

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

func CreateIOPortController() *IOPortAccessController{
	return &IOPortAccessController{backingMemory:make([]byte, 0x10000)}
}
