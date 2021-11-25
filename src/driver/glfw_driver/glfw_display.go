package glfw_driver

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"chip8/conf"
)

type GFWDisplay struct {
	window *glfw.Window
}

func (d *GFWDisplay) DrawPixel(x, y uint8) {
	gl.Rectf(
		float32(uint64(x)*conf.Zoom),
		float32(uint64(y)*conf.Zoom),
		float32(conf.Zoom+uint64(x)*conf.Zoom),
		float32(conf.Zoom+uint64(y)*conf.Zoom),
	)
}

func (d *GFWDisplay) DrawFrame(screenBuf *[conf.ScreenHeight][conf.ScreenWidth]uint8) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for y := uint8(0); y < conf.ScreenHeight; y++ {
		for x := uint8(0); x < conf.ScreenWidth; x++ {
			if screenBuf[y][x] == 1 {
				d.DrawPixel(x, y)
			}
		}
	}
}

func (d *GFWDisplay) Render(screenBuf *[conf.ScreenHeight][conf.ScreenWidth]uint8) {
	d.DrawFrame(screenBuf)
	d.window.SwapBuffers()
}

func NewGFWDisplay(window *glfw.Window) *GFWDisplay {
	return &GFWDisplay{
		window: window,
	}
}
