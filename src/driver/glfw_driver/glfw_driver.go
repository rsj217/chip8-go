package glfw_driver

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"chip8/conf"
	"chip8/driver"
)

type GLFWDriver struct {
	window *glfw.Window
}

func (f *GLFWDriver) Loop() bool {
	closed := f.window.ShouldClose()
	if closed {
		glfw.Terminate()
	}
	return !closed
}

func (f *GLFWDriver) Display() driver.Displayer {
	return NewGFWDisplay(f.window)
}

func (f *GLFWDriver) Keyboard() driver.Keyboarder {
	return NewGFWLKeyboard(f.window)
}

func NewGFWDriver() *GLFWDriver {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(int(64*conf.Zoom), int(32*conf.Zoom), "chip8", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)

	gl.Ortho(0, float64(64*conf.Zoom), float64(32*conf.Zoom), 0, 0, 1)
	return &GLFWDriver{window: window}
}
