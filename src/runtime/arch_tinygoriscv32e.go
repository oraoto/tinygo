//go:build tinygo.riscve

package runtime

import "device/riscv"

const GOARCH = "arm" // riscv pretends to be arm

// The bitness of the CPU (e.g. 8, 32, 64).
const TargetBits = 32

const callInstSize = 4

const deferExtraRegs = 0

// RISC-V has a maximum alignment of 16 bytes (both for RV32 and for RV64).
// Source: https://riscv.org/wp-content/uploads/2015/01/riscv-calling.pdf
func align(ptr uintptr) uintptr {
	return (ptr + 15) &^ 15
}

func getCurrentStackPointer() uintptr {
	return uintptr(stacksave())
}

// The safest thing to do here would just be to disable interrupts for
// procPin/procUnpin. Note that a global variable is safe in this case, as any
// access to procPinnedMask will happen with interrupts disabled.

var procPinnedMask uintptr

//go:linkname procPin sync/atomic.runtime_procPin
func procPin() {
	procPinnedMask = riscv.DisableInterrupts()
}

//go:linkname procUnpin sync/atomic.runtime_procUnpin
func procUnpin() {
	riscv.EnableInterrupts(procPinnedMask)
}

func waitForEvents() {
	mask := riscv.DisableInterrupts()
	if !runqueue.Empty() {
		riscv.Asm("wfi")
	}
	riscv.EnableInterrupts(mask)
}
