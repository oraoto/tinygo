//go:build ch32v003

package interrupt

import (
	"device/riscv"
)

func (i Interrupt) Enable() error {
	riscv.Asm("fence")
	return nil
}

// Adding pseudo function calls that is replaced by the compiler with the actual
// functions registered through interrupt.New.
//
//go:linkname callHandlers runtime/interrupt.callHandlers
func callHandlers(num int)

//go:inline
func callHandler(n int) {
	for {
		riscv.Asm("wfi")
	}
}

//export handleInterrupt
func handleInterrupt() {
	mcause := riscv.MCAUSE.Get()
	exception := mcause&(1<<31) == 0
	interruptNumber := uint32(mcause & 0x1f)

	if !exception && interruptNumber > 0 {
		// save MSTATUS & MEPC, which could be overwritten by another CPU interrupt
		mstatus := riscv.MSTATUS.Get()
		mepc := riscv.MEPC.Get()
		// Useing threshold to temporary disable this interrupts.
		// FYI: using CPU interrupt enable bit make runtime to loose interrupts.
		/*
			reg := (*volatile.Register32)(unsafe.Add(unsafe.Pointer(&esp.INTERRUPT_CORE0.CPU_INT_PRI_0), interruptNumber*4))
			thresholdSave := reg.Get()
			reg.Set(disableThreshold)
		*/
		riscv.Asm("fence")

		//interruptBit := uint32(1 << interruptNumber)

		/*
			// reset pending status interrupt
			if esp.INTERRUPT_CORE0.CPU_INT_TYPE.Get()&interruptBit != 0 {
				// this is edge type interrupt
				esp.INTERRUPT_CORE0.CPU_INT_CLEAR.SetBits(interruptBit)
				esp.INTERRUPT_CORE0.CPU_INT_CLEAR.ClearBits(interruptBit)
			} else {
				// this is level type interrupt
				esp.INTERRUPT_CORE0.CPU_INT_CLEAR.ClearBits(interruptBit)
			}
		*/

		// enable CPU interrupts
		riscv.MSTATUS.SetBits(1 << 3)

		// Call registered interrupt handler(s)
		callHandler(int(interruptNumber))

		// disable CPU interrupts
		riscv.MSTATUS.ClearBits(1 << 3)

		// restore interrupt threshold to enable interrupt again
		//reg.Set(thresholdSave)
		riscv.Asm("fence")

		// restore MSTATUS & MEPC
		riscv.MSTATUS.Set(mstatus)
		riscv.MEPC.Set(mepc)

		// do not enable CPU interrupts now
		// the 'MRET' in src/device/riscv/handleinterrupt.S will copies the state of MPIE back into MIE, and subsequently clears MPIE.
		// riscv.MSTATUS.SetBits(0x8)
	} else {
		// Topmost bit is clear, so it is an exception of some sort.
		// We could implement support for unsupported instructions here (such as
		// misaligned loads). However, for now we'll just print a fatal error.
		handleException(mcause)
	}
}

func handleException(mcause uintptr) {
	println("*** Exception:     pc:", riscv.MEPC.Get())
	println("*** Exception:   code:", uint32(mcause&0x1f))
	println("*** Exception: mcause:", mcause)
	switch uint32(mcause & 0x1f) {
	case 1:
		println("***    virtual addess:", riscv.MTVAL.Get())
	case 2:
		println("***            opcode:", riscv.MTVAL.Get())
	case 5:
		println("***      read address:", riscv.MTVAL.Get())
	case 7:
		println("***     write address:", riscv.MTVAL.Get())
	}
	for {
		riscv.Asm("wfi")
	}
}
