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

const (
	IRQ_PIT_COUNTER0 = iota + 0x20
	IRQ_PIT_COUNTER1
	IRQ_PIT_COUNTER2
	IRQ_KEYBOARD
	IRQ_COM2
	IRQ_COM1
	IRQ_LPT2
	IRQ_FLOPPY
	IRQ_LPT1
	IRQ_CMOS_RTC
	IRQ_FREE1
	IRQ_FREE2
	IRQ_FREE3
	IRQ_MOUSE
	IRQ_FPU
	IRQ_PRIMARY_ATA
	IRQ_SECONDARY_ATA
)
