//go:build ch32v003

package runtime

// import "device/wch"

import (
	"device/riscv"
	"device/wch"
	"machine"

	"runtime/interrupt"
)

type timeUnit int64

//export main
func main() {
	interrupt.New(3, func(interrupt.Interrupt) {
		println("interrupt called")
	}).Enable()

	enableUART()
	run()
}

//export __mulsi3
func __mulsi3(a, b int) int {
	var p int
	var i int
	absA := a
	if absA < 0 {
		absA = -a
	}
	for i = 0; i < absA; i++ {
		p += b
	}

	return p
}

//export __udivsi3
func __udivsi3(dividend, divisor int) int {
	if divisor == 0 {
		return 0
	}

	isNegative := (dividend < 0) != (divisor < 0)
	if dividend < 0 {
		dividend = -dividend
	}
	if divisor < 0 {
		divisor = -divisor
	}

	var quotient int

	for dividend >= divisor {
		dividend -= divisor
		quotient++
	}

	if isNegative {
		return -quotient
	}

	return quotient
}

func enableUART() {
	wch.RCC.APB2PCENR.SetBits(wch.RCC_APB2PCENR_USART1EN)

	// remap, default mapping 00
	wch.AFIO.PCFR.ClearBits(wch.AFIO_PCFR_USART1RM)     // low bit 0
	wch.AFIO.PCFR.ClearBits(wch.AFIO_PCFR_USART1REMAP1) // high bit 0

	// configure pin output
	tx := machine.PD5

	tx.Configure(machine.PinConfig{Mode: machine.PinOutput50MHz + machine.PinOutputModeAltPushPull})

	// 115200
	var divider uint32
	divider = 8 * 1000 * 1000 / 115200
	//gigadevice.USART0.BAUD.Set(divider)
	wch.USART1.BRR.Set(divider)

	// stopbit 1
	wch.USART1.SetCTLR2_STOP(0b00)

	// Enable USART port, tx
	wch.USART1.CTLR1.Set(
		wch.USART_CTLR1_TE |
			wch.USART_CTLR1_UE,
	)
}

// func waitForEvents() {
// 	arm.Asm("wfe")
// }

func putchar(c byte) {
	for {
		if wch.USART1.STATR.HasBits(wch.USART_STATR_TXE) {
			wch.USART1.DATAR.Set(uint32(c))
			break
		}
	}
}

func exit(_ int) {

}

func abort() {
	println("abort")
	for {
		riscv.Asm("wfi")
	}
}

func initTimer() {}

func sleepTicks(d timeUnit) {}

func ticks() timeUnit {
	return 0
}

func nanosecondsToTicks(ns int64) timeUnit {
	return timeUnit(ns)
}

func ticksToNanoseconds(ticks timeUnit) int64 {
	return int64(ticks)
}
