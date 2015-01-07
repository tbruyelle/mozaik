package main

import (
	"fmt"
)

type ColorDef int

const (
	Red ColorDef = iota
	Yellow
	Blue
	Green
	Pink
	Orange
	LightBlue
)

type Block struct {
	Object
	Color ColorDef
}

type Switch struct {
	Object
	line, col int
	name      string
}

func (s *Switch) Rotate() {
	blocks := s.Blocks()
	for i := range blocks {
		b := blocks[i]
		v := switchSize / 2
		b.Rx, b.Ry = s.X+v, s.Y+v
		b.Action = blockRotate
	}
}

// Blocks returns the block arround the switch in parameter.
func (sw *Switch) Blocks() []*Block {
	topLeft := g.level.blocks[sw.line][sw.col]
	topRight := g.level.blocks[sw.line][sw.col+1]
	bottomLeft := g.level.blocks[sw.line+1][sw.col]
	bottomRight := g.level.blocks[sw.line+1][sw.col+1]
	return []*Block{topLeft, topRight, bottomLeft, bottomRight}
}

func (s *Switch) String() string {
	return fmt.Sprintf("sw{line:%d, col:%d}", s.line, s.col)
}
