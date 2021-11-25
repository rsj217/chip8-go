package sdl2_driver

import (
	"log"

	sdl2 "github.com/veandco/go-sdl2/sdl"

	"chip8/conf"
)

type SDL2Keyboard struct {
	keycodeMap map[uint16]uint8
}

func NewSDL2Keyboard() *SDL2Keyboard {
	keycodeMap := make(map[uint16]uint8, 16)
	keycodeMap[uint16(sdl2.GetKeyFromName("1"))] = 0x1
	keycodeMap[uint16(sdl2.GetKeyFromName("2"))] = 0x2
	keycodeMap[uint16(sdl2.GetKeyFromName("3"))] = 0x3
	keycodeMap[uint16(sdl2.GetKeyFromName("4"))] = 0xC
	keycodeMap[uint16(sdl2.GetKeyFromName("Q"))] = 0x4
	keycodeMap[uint16(sdl2.GetKeyFromName("W"))] = 0x5
	keycodeMap[uint16(sdl2.GetKeyFromName("E"))] = 0x6
	keycodeMap[uint16(sdl2.GetKeyFromName("R"))] = 0xD
	keycodeMap[uint16(sdl2.GetKeyFromName("A"))] = 0x7
	keycodeMap[uint16(sdl2.GetKeyFromName("S"))] = 0x8
	keycodeMap[uint16(sdl2.GetKeyFromName("D"))] = 0x9
	keycodeMap[uint16(sdl2.GetKeyFromName("F"))] = 0xE
	keycodeMap[uint16(sdl2.GetKeyFromName("Z"))] = 0xA
	keycodeMap[uint16(sdl2.GetKeyFromName("X"))] = 0x0
	keycodeMap[uint16(sdl2.GetKeyFromName("C"))] = 0xB
	keycodeMap[uint16(sdl2.GetKeyFromName("V"))] = 0xF
	return &SDL2Keyboard{
		keycodeMap: keycodeMap,
	}
}

func (k SDL2Keyboard) PollEvent(keysPressedBuf *[conf.KeysNum]uint8) {
	for event := sdl2.PollEvent(); event != nil; event = sdl2.PollEvent() {
		switch t := event.(type) {
		case *sdl2.QuitEvent:
			sdl2.Quit()
			log.Fatal()
		case *sdl2.KeyboardEvent:

			switch t.Type {
			case sdl2.KEYDOWN:
				sym := k.keycodeMap[uint16(t.Keysym.Sym)]
				keysPressedBuf[sym] = 1
			case sdl2.KEYUP:
				sym := k.keycodeMap[uint16(t.Keysym.Sym)]
				keysPressedBuf[sym] = 0
			}
		}
	}
}
