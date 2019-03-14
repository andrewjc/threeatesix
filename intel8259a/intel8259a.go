package intel8259a

/*
	Simulated 8259A Interrupt Controller Chip

	IRQ0 through IRQ7 are the master 8259's interrupt lines, while IRQ8 through IRQ15 are the slave 8259's interrupt lines.
*/

type Intel8259a struct {
}

func NewIntel8259a() Intel8259a {
	chip := Intel8259a{}

	return chip
}
