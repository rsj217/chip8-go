package main

import "chip8/machine"

func start(rom string) {
	m := machine.NewMachine()
	m.LoadRom(rom)
	m.Run()
}

func main() {
	rom := "../roms/BC"
	start(rom)
}
