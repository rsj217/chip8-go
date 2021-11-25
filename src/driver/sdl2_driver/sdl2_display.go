package sdl2_driver

import (
	sdl2 "github.com/veandco/go-sdl2/sdl"

	"chip8/conf"
)

type SDL2Display struct {
	window  *sdl2.Window
	surface *sdl2.Surface
}

func (d *SDL2Display) DrawPixel(x, y uint8) {

	rect := &sdl2.Rect{
		X: int32(uint64(x) * conf.Zoom),
		Y: int32(uint64(y) * conf.Zoom),
		W: int32(conf.Zoom),
		H: int32(conf.Zoom),
	}
	_ = d.surface.FillRect(rect, conf.FgColor)

}

func (d *SDL2Display) DrawFrame(screenBuf *[conf.ScreenHeight][conf.ScreenWidth]uint8) {
	_ = d.surface.FillRect(nil, conf.BgColor)
	for y := uint8(0); y < conf.ScreenHeight; y++ {
		for x := uint8(0); x < conf.ScreenWidth; x++ {
			if screenBuf[y][x] == 1 {
				d.DrawPixel(x, y)
			}
		}
	}
}

func (d *SDL2Display) Render(screenBuf *[conf.ScreenHeight][conf.ScreenWidth]uint8) {
	d.DrawFrame(screenBuf)
	_ = d.window.UpdateSurface()
}

func NewDisplay(window *sdl2.Window, surface *sdl2.Surface) *SDL2Display {
	display := &SDL2Display{window: window, surface: surface}
	return display
}
