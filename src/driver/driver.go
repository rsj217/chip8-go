package driver

import (
	"chip8/conf"
)

type Displayer interface {
	DrawPixel(x, y uint8)
	DrawFrame(screenBuf *[conf.ScreenHeight][conf.ScreenWidth]uint8)
	Render(screenBuf *[conf.ScreenHeight][conf.ScreenWidth]uint8)
}

type Keyboarder interface {
	PollEvent(keysPressedBuf *[conf.KeysNum]uint8)
}

type Driver interface {
	Loop() bool
	Display() Displayer
	Keyboard() Keyboarder
}
