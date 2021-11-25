package chip8

import (
	"fmt"
	"log"
	"math/rand"

	"chip8/conf"
)

type Instruction struct {
	val    uint16
	opcode uint16
	x      uint8
	y      uint8
	n      uint8
	kk     uint8
	nnn    uint16
}

func (i Instruction) DecodeOpcode() uint16 {
	return i.val & 0xF000
}

func (i Instruction) DecodeX() uint8 {
	return uint8((i.val & 0x0F00) >> 8)
}

func (i Instruction) DecodeY() uint8 {
	return uint8((i.val & 0x00F0) >> 4)
}

func (i Instruction) DecodeN() uint8 {
	return uint8(i.val & 0x000F)
}

func (i Instruction) DecodeFlag() uint8 {
	return uint8(i.val & 0x000F)
}

func (i Instruction) DecodeKK() uint8 {
	return uint8(i.val & 0x00FF)
}

func (i Instruction) DecodeNNN() uint16 {
	return i.val & 0x0FFF
}

func (i *Instruction) Decode() {
	i.opcode = i.DecodeOpcode()
	i.x = i.DecodeX()
	i.y = i.DecodeY()
	i.n = i.DecodeN()
	i.kk = i.DecodeKK()
	i.nnn = i.DecodeNNN()
}

func NewInstruction(val uint16) *Instruction {
	return &Instruction{val: val}
}

type CPU struct {
	memory         *Memory
	regV           [conf.RegsNum]uint8
	regPC          uint16
	regI           uint16
	regSP          uint16
	drawFlag       bool
	screenBuf      [conf.ScreenHeight][conf.ScreenWidth]uint8
	keysPressedBuf [conf.KeysNum]uint8
	delayTimer     uint8
	soundTimer     uint8
	cycleStartTime int64
	ir             *Instruction
}

func NewCPU() *CPU {
	memory := NewMemory()
	return &CPU{memory: memory, regPC: conf.StartAddr}
}

func (cpu *CPU) GetMemory() *Memory {
	return cpu.memory
}

func (cpu *CPU) GetScreenBuf() *[conf.ScreenHeight][conf.ScreenWidth]uint8 {
	return &cpu.screenBuf
}

func (cpu *CPU) GetKeysPressedBuf() *[conf.KeysNum]uint8 {
	return &cpu.keysPressedBuf
}

func (cpu *CPU) GetDrawFlag() bool {
	return cpu.drawFlag
}

func (cpu *CPU) SetDrawFlag(flag bool) {
	cpu.drawFlag = flag
}

func (cpu *CPU) LoadRom(data []uint8) {
	cpu.memory.write(data)
}

func (cpu *CPU) fetch() {
	highByte := cpu.memory.ram[cpu.regPC]
	lowByte := cpu.memory.ram[cpu.regPC+1]
	instruction := (uint16(highByte) << 8) | uint16(lowByte)
	pc := cpu.regPC + 2
	cpu.regPC = pc
	cpu.ir = NewInstruction(instruction)
}

func (cpu *CPU) decode() {
	cpu.ir.Decode()
}

func (cpu CPU) GetOpcode() uint16 {
	return cpu.ir.opcode
}

func (cpu CPU) GetX() uint8 {
	return cpu.ir.x
}

func (cpu CPU) GetY() uint8 {
	return cpu.ir.y
}

func (cpu CPU) GetN() uint8 {
	return cpu.ir.n
}

func (cpu CPU) GetKK() uint8 {
	return cpu.ir.kk
}

func (cpu CPU) GetAddr() uint16 {
	return cpu.ir.nnn
}

func (cpu *CPU) execute() {
	switch cpu.GetOpcode() {

	case 0x0000:
		switch cpu.GetKK() {
		case 0x00E0:
			cpu.screenBuf = [conf.ScreenHeight][conf.ScreenWidth]uint8{}
			cpu.drawFlag = true
		case 0x00EE:
			if cpu.regSP <= 0 {
				log.Fatal(fmt.Errorf("cpu.stack pointer below 0"))
			}
			cpu.regSP--
			cpu.regPC = cpu.memory.stackPop(cpu.regSP)
		}

	case 0x1000:
		cpu.regPC = cpu.GetAddr()

	case 0x2000:
		cpu.memory.stackPush(cpu.regSP, cpu.regPC)
		cpu.regSP++
		cpu.regPC = cpu.GetAddr()

	case 0x3000:
		x := cpu.GetX()
		kk := cpu.GetKK()
		if cpu.regV[x] == kk {
			cpu.regPC += 2
		}

	case 0x4000:
		x := cpu.GetX()
		kk := cpu.GetKK()
		if cpu.regV[x] != kk {
			cpu.regPC += 2
		}

	case 0x5000:
		x := cpu.GetX()
		y := cpu.GetY()
		if cpu.regV[x] == cpu.regV[y] {
			cpu.regPC += 2
		}

	case 0x6000:
		x := cpu.GetX()
		kk := cpu.GetKK()
		cpu.regV[x] = kk

	case 0x7000:
		x := cpu.GetX()
		kk := cpu.GetKK()
		cpu.regV[x] += kk
		cpu.regV[x] &= 0xFF

	case 0x8000:
		switch cpu.GetN() {
		case 0x0000:
			x := cpu.GetX()
			y := cpu.GetY()
			cpu.regV[x] = cpu.regV[y]

		case 0x0001:
			x := cpu.GetX()
			y := cpu.GetY()
			cpu.regV[x] |= cpu.regV[y]

		case 0x0002:
			x := cpu.GetX()
			y := cpu.GetY()
			cpu.regV[x] &= cpu.regV[y]

		case 0x0003:
			x := cpu.GetX()
			y := cpu.GetY()
			cpu.regV[x] ^= cpu.regV[y]

		case 0x0004:
			x := cpu.GetX()
			y := cpu.GetY()
			cpu.regV[x] += cpu.regV[y]

			if cpu.regV[x] > 0xFF {
				cpu.regV[0x0F] = 0x01
			} else {
				cpu.regV[0x0F] = 0x00
			}
			cpu.regV[x] &= 0xFF

		case 0x0005:
			x := cpu.GetX()
			y := cpu.GetY()
			if cpu.regV[x] < cpu.regV[y] {
				cpu.regV[0x0F] = 0x00
			} else {
				cpu.regV[0x0F] = 0x01

			}
			cpu.regV[x] -= cpu.regV[y]
			cpu.regV[x] &= 0xFF

		case 0x0006:
			x := cpu.GetX()
			cpu.regV[0x0F] = cpu.regV[x] & 0x01
			cpu.regV[x] >>= 1

		case 0x0007:
			x := cpu.GetX()
			y := cpu.GetY()
			cpu.regV[0x0F] = 0x01
			if cpu.regV[x] < cpu.regV[y] {
				cpu.regV[0x0F] = 0x01
			} else {
				cpu.regV[0x0F] = 0x00
			}
			cpu.regV[x] = cpu.regV[y] - cpu.regV[x]
			cpu.regV[x] &= 0xFF

		case 0x000E:
			x := cpu.GetX()
			cpu.regV[0x0F] = (cpu.regV[x] >> 7) & 0x01
			cpu.regV[x] = cpu.regV[x] << 1
			cpu.regV[x] &= 0xFF
		}

	case 0x9000:
		x := cpu.GetX()
		y := cpu.GetY()
		if cpu.regV[x] != cpu.regV[y] {
			cpu.regPC += 2
		}

	case 0xA000:
		cpu.regI = cpu.GetAddr()

	case 0xB000:
		addr := cpu.GetAddr()
		cpu.regPC = uint16(cpu.regV[0]) + addr

	case 0xC000:
		x := cpu.GetX()
		kk := cpu.GetKK()
		cpu.regV[x] = uint8(rand.Intn(256)) & kk

	case 0xD000:
		n := cpu.GetN()
		x := cpu.GetX()
		y := cpu.GetY()

		vx := cpu.regV[x]
		vy := cpu.regV[y]

		cpu.regV[0x0F] = 0
		for yy := uint8(0); yy < n; yy++ {
			sysByte := cpu.memory.ram[cpu.regI+uint16(yy)]
			for xx := uint8(0); xx < 8; xx++ {
				xCord := vx + xx
				yCord := vy + yy

				if xCord < conf.ScreenWidth && yCord < conf.ScreenHeight {
					sysBit := (sysByte >> (7 - xx)) & 0x01
					if (cpu.screenBuf[yCord][xCord] & sysBit) == 1 {
						cpu.regV[0x0F] = 1
					}
					cpu.screenBuf[yCord][xCord] ^= sysBit
				}
			}
		}
		cpu.drawFlag = true

	case 0xE000:
		switch cpu.GetKK() {
		case 0x009E:
			x := cpu.GetX()
			if cpu.keysPressedBuf[cpu.regV[x]] == 1 {
				cpu.regPC += 2
			}

		case 0x00A1:
			x := cpu.GetX()
			if cpu.keysPressedBuf[cpu.regV[x]] == 0 {
				cpu.regPC += 2
			}
		}

	case 0xF000:
		switch cpu.GetKK() {

		case 0x0007:
			x := cpu.GetX()
			cpu.regV[x] = cpu.delayTimer

		case 0x000A:
			x := cpu.GetX()
			pressed := false
			for i := uint8(0); i < 16; i++ {
				if cpu.keysPressedBuf[i] == 1 {
					cpu.regV[x] = i
					pressed = true
					break
				}
			}

			if !pressed {
				cpu.regPC -= 2
			}

		case 0x0015:
			x := cpu.GetX()
			cpu.delayTimer = cpu.regV[x]

		case 0x0018:
			x := cpu.GetX()
			cpu.soundTimer = cpu.regV[x]

		case 0x001E:
			x := cpu.GetX()
			cpu.regI = cpu.regI + uint16(cpu.regV[x])

		case 0x0029:
			x := cpu.GetX()
			cpu.regI = uint16(cpu.regV[x]) * 5

		case 0x0033:
			x := cpu.GetX()
			cpu.memory.ram[cpu.regI] = cpu.regV[x] / 100
			cpu.memory.ram[cpu.regI+1] = (cpu.regV[x] % 100) / 10
			cpu.memory.ram[cpu.regI+2] = (cpu.regV[x] % 100) % 10

		case 0x0055:
			x := cpu.GetX()
			for i := uint8(0); i <= x; i++ {
				cpu.memory.ram[cpu.regI+uint16(i)] = cpu.regV[i]
			}

		case 0x0065:
			x := cpu.GetX()
			for i := uint8(0); i <= x; i++ {
				cpu.regV[i] = cpu.memory.ram[cpu.regI+uint16(i)]
			}

		}
	}
}

func (cpu *CPU) Cycle() {
	cpu.fetch()
	cpu.decode()
	cpu.execute()
}

func (cpu *CPU) Ticker() {
	if cpu.delayTimer > 0 {
		cpu.delayTimer--
	}
	if cpu.soundTimer > 0 {
		// TODO: continuous
		cpu.soundTimer--
	}
}
