package glfw_driver

import (
	"github.com/go-gl/glfw/v3.2/glfw"

	"chip8/conf"
)

type GFWLKeyboard struct {
	window     *glfw.Window
	keycodeMap map[glfw.Key]uint8
}

func NewGFWLKeyboard(window *glfw.Window) *GFWLKeyboard {
	keycodeMap := map[glfw.Key]byte{
		glfw.Key1: 0x1,
		glfw.Key2: 0x2,
		glfw.Key3: 0x3,
		glfw.Key4: 0xc,
		glfw.KeyQ: 0x4,
		glfw.KeyW: 0x5,
		glfw.KeyE: 0x6,
		glfw.KeyR: 0xd,
		glfw.KeyA: 0x7,
		glfw.KeyS: 0x8,
		glfw.KeyD: 0x9,
		glfw.KeyF: 0xe,
		glfw.KeyZ: 0xa,
		glfw.KeyX: 0x0,
		glfw.KeyC: 0xb,
		glfw.KeyV: 0xf,
	}
	return &GFWLKeyboard{
		window:     window,
		keycodeMap: keycodeMap,
	}
}

func (key GFWLKeyboard) PollEvent(keysPressedBuf *[conf.KeysNum]uint8) {
	glfw.PollEvents()
	for k, v := range key.keycodeMap {
		switch key.window.GetKey(k) {
		case glfw.Press:
			keysPressedBuf[v] = 1
		case glfw.Release:
			keysPressedBuf[v] = 0
		}
	}
}
