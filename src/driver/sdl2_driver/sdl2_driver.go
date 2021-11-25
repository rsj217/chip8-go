package sdl2_driver

import (
	"log"

	sdl2 "github.com/veandco/go-sdl2/sdl"

	"chip8/conf"
	"chip8/driver"
)

type SDL2Driver struct {
	window  *sdl2.Window
	surface *sdl2.Surface
}

func (s *SDL2Driver) Loop() bool {
	return true
}

func (s *SDL2Driver) Display() driver.Displayer {
	return NewDisplay(s.window, s.surface)
}

func (s *SDL2Driver) Keyboard() driver.Keyboarder {
	return NewSDL2Keyboard()
}

func NewSDL2Driver() *SDL2Driver {
	if err := sdl2.Init(sdl2.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	window, err := sdl2.CreateWindow(
		"chip8",
		sdl2.WINDOWPOS_UNDEFINED,
		sdl2.WINDOWPOS_UNDEFINED,
		int32(64*conf.Zoom), int32(32*conf.Zoom),
		sdl2.WINDOW_SHOWN,
	)
	if err != nil {
		sdl2.Quit()
		log.Fatal(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		_ = window.Destroy()
		sdl2.Quit()
		log.Fatal(err)
	}
	display := &SDL2Display{window: window, surface: surface}
	if err := display.surface.FillRect(nil, conf.BgColor); err != nil {
		_ = window.Destroy()
		sdl2.Quit()
		log.Fatal(err)
	}

	if err := display.window.UpdateSurface(); err != nil {
		_ = window.Destroy()
		sdl2.Quit()
		log.Fatal(err)
	}
	return &SDL2Driver{window: window, surface: surface}
}
