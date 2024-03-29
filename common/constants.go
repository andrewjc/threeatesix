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
	MODULE_PS2_CONTROLLER
	MODULE_INTEL_82335_MCR
	MODULE_DEBUG_MONITOR
)

const (
	SEGMENT_CS = iota
	SEGMENT_SS
	SEGMENT_DS
	SEGMENT_ES
	SEGMENT_FS
	SEGMENT_GS
)
