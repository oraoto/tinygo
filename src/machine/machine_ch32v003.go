//go:build ch32v003

package machine

import "device/wch"

const deviceName = wch.Device

const (
	portA Pin = iota * 16
	portC
	portD
)

const (
	PA0 = portA + iota
	PA1
	PA2
	PA3
	PA4
	PA5
	PA6
	PA7
)

const (
	PC0 = portC + iota
	PC1
	PC2
	PC3
	PC4
	PC5
	PC6
	PC7
)

const (
	PD0 = portD + iota
	PD1
	PD2
	PD3
	PD4
	PD5
	PD6
	PD7
)

const (
	PinInput       PinMode = 0b00 // Input Mode
	PinOutput10MHz PinMode = 0b01 // Output mode, max speed 10MHz
	PinOutput2MHz  PinMode = 0b10 // Output mode, max speed 2MHz
	PinOutput50MHz PinMode = 0b11 // Output mode, max speed 50MHz
	PinOutput      PinMode = PinOutput2MHz

	PinInputModeAnalog     PinMode = 0b00 << 2 // Analog input mode
	PinInputModeFloating   PinMode = 0b01 << 2 // Floating input mode
	PinInputModePullUpDown PinMode = 0b10 << 2 // Input pull up/down mode

	PinOutputModeGPPushPull   PinMode = 0b00 << 2 // Output mode general purpose push/pull
	PinOutputModeGPOpenDrain  PinMode = 0b01 << 2 // Output mode general purpose open drain
	PinOutputModeAltPushPull  PinMode = 0b10 << 2 // Output mode alt. purpose push/pull
	PinOutputModeAltOpenDrain PinMode = 0b11 << 2 // Output mode alt. purpose open drain

	PinInputPulldown PinMode = PinInputModePullUpDown
	PinInputPullup   PinMode = PinInputModePullUpDown | 0b10000
)

func (p Pin) Configure(config PinConfig) {
	p.enableClock()
	port := p.getPort()
	pin := uint8(p) % 16

	port.CFGLR.ReplaceBits(uint32(config.Mode), 0b1111, pin*4)

	if config.Mode&0b11 == PinInputModePullUpDown {
		if config.Mode&0b10000 == 0b10000 {
			port.OUTDR.ReplaceBits(1, 0b1, pin)
		}
	}
}

func (p Pin) Set(high bool) {
	port := p.getPort()
	pin := uint8(p) % 16
	if high {
		port.BSHR.Set(1 << pin)
	} else {
		port.BSHR.Set(1 << (pin + 16))
	}
}

func (p Pin) Get() bool {
	port := p.getPort()
	pin := uint8(p) % 16
	val := port.INDR.Get() & (1 << pin)
	return (val > 0)
}

func (p Pin) getPort() *wch.GPIO_Type {
	switch p / 16 {
	case 0:
		return wch.GPIOA
	case 1:
		return wch.GPIOC
	case 2:
		return wch.GPIOD
	default:
		panic("machine: unknown port")
	}
}

func (p Pin) enableClock() {
	switch p / 16 {
	case 0:
		wch.RCC.APB2PCENR.SetBits(wch.RCC_APB2PCENR_IOPAEN)
	case 1:
		wch.RCC.APB2PCENR.SetBits(wch.RCC_APB2PCENR_IOPCEN)
	case 2:
		wch.RCC.APB2PCENR.SetBits(wch.RCC_APB2PCENR_IOPDEN)
	default:
		panic("machine: unknown port")
	}
}
