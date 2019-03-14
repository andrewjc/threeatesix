package cpu

import "github.com/andrewjc/threeatesix/common"

func New80287MathCoProcessor() CpuCore {

	cpuCore := New80386CPU()

	cpuCore.partId = common.MODULE_MATH_CO_PROCESSOR

	return cpuCore
}
