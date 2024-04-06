//go:build !tinygo.riscv && !cortexm && !tinygo.riscve

package runtime

func waitForEvents() {
	runtimePanic("deadlocked: no event source")
}
