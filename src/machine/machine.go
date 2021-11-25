package machine

import (
	"time"

	"chip8/chip8"
	"chip8/conf"
	"chip8/driver"
	"chip8/driver/glfw_driver"
	"chip8/driver/sdl2_driver"
	"chip8/util"
)

type Machine struct {
	driver   driver.Driver
	display  driver.Displayer
	keyboard driver.Keyboarder
	cpu      *chip8.CPU
}

func NewSDL2Machine() *Machine {
	d := sdl2_driver.NewSDL2Driver()
	display := d.Display()
	keyboard := d.Keyboard()
	cpu := chip8.NewCPU()
	return &Machine{
		driver:   d,
		display:  display,
		keyboard: keyboard,
		cpu:      cpu,
	}

}

func NewGLFWMachine() *Machine {
	d := glfw_driver.NewGFWDriver()
	display := d.Display()
	keyboard := d.Keyboard()
	cpu := chip8.NewCPU()
	return &Machine{
		driver:   d,
		display:  display,
		keyboard: keyboard,
		cpu:      cpu,
	}
}

func NewMachine() (machine *Machine) {
	switch conf.Driver {
	case conf.SDL2Driver:
		machine = NewSDL2Machine()
	case conf.GLFWDriver:
		machine = NewGLFWMachine()
	default:
		panic("invalid driver")
	}
	return
}

func (m Machine) LoadRom(path string) {
	data := util.ReadROM(path)
	m.cpu.LoadRom(data)
}

func (m Machine) Run() {
	var cycle uint8
	for m.driver.Loop() {
		m.cpu.Cycle()
		m.keyboard.PollEvent(m.cpu.GetKeysPressedBuf())
		if m.cpu.GetDrawFlag() {
			m.display.Render(m.cpu.GetScreenBuf())
			m.cpu.SetDrawFlag(false)
		}
		time.Sleep(time.Duration(1e6/uint32(conf.ClockSpeed)) * time.Microsecond)
		if cycle >= uint8(conf.ClockSpeed/conf.TimerSpeed) { //timer start
			cycle = 0
			m.cpu.Ticker()
		} // timer end
		cycle++
	}
}
