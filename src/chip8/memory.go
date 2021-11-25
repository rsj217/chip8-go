package chip8

import (
	"chip8/conf"
)

type Memory struct {
	ram   [conf.RamSize]uint8
	stack [conf.StackSize]uint16
}

func NewMemory() *Memory {
	memory := &Memory{}
	for i := range conf.FontSet {
		memory.ram[i] = conf.FontSet[i]
	}
	return memory
}

func (m *Memory) write(data []uint8) {
	for i := uint16(0); i < uint16(len(data)); i++ {
		m.ram[conf.StartAddr+i] = data[i]
	}
}

func (m *Memory) stackPush(SP, item uint16) {
	m.stack[SP] = item
}

func (m *Memory) stackPop(SP uint16) uint16 {
	return m.stack[SP]
}
