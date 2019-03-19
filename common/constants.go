package common

const (
	REAL_MODE = iota
	PROTECTED_MODE
)

const (
	MODULE_PRIMARY_PROCESSOR = iota
	MODULE_MATH_CO_PROCESSOR
	MODULE_MASTER_INTERRUPT_CONTROLLER
	MODULE_SLAVE_INTERRUPT_CONTROLLER
	MODULE_MEMORY_ACCESS_CONTROLLER
	MODULE_IO_PORT_ACCESS_CONTROLLER
)
