//go:build scheduler.tasks && tinygo.riscve

package task

import "unsafe"

var systemStack uintptr

// calleeSavedRegs is the list of registers that must be saved and restored when
// switching between tasks. Also see scheduler_riscv.S that relies on the
// exact layout of this struct.
type calleeSavedRegs struct {
	s0 uintptr // x8 (fp)
	s1 uintptr // x9

	pc uintptr
}

// archInit runs architecture-specific setup for the goroutine startup.
func (s *state) archInit(r *calleeSavedRegs, fn uintptr, args unsafe.Pointer) {
	// Store the initial sp for the startTask function (implemented in assembly).
	s.sp = uintptr(unsafe.Pointer(r))

	// Initialize the registers.
	// These will be popped off of the stack on the first resume of the goroutine.

	// Start the function at tinygo_startTask (defined in src/internal/task/task_stack_riscv.S).
	// This assembly code calls a function (passed in s0) with a single argument
	// (passed in s1). After the function returns, it calls Pause().
	r.pc = uintptr(unsafe.Pointer(&startTask))

	// Pass the function to call in s0.
	// This function is a compiler-generated wrapper which loads arguments out
	// of a struct pointer. See createGoroutineStartWrapper (defined in
	// compiler/goroutine.go) for more information.
	r.s0 = fn

	// Pass the pointer to the arguments struct in s1.
	r.s1 = uintptr(args)
	println("finish arch init")
}

func (s *state) resume() {
	swapTask(s.sp, &systemStack)
}

func (s *state) pause() {
	newStack := systemStack
	systemStack = 0
	swapTask(newStack, &s.sp)
}

// SystemStack returns the system stack pointer when called from a task stack.
// When called from the system stack, it returns 0.
func SystemStack() uintptr {
	return systemStack
}
