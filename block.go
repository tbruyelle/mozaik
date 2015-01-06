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
	Color    ColorDef
	X, Y     float32
	rotateSW *Switch
}

type Switch struct {
	state     State
	line, col int
	X, Y      float32
	scale     float32
	rotate    float32
	name      string
}

func (s *Switch) Rotate() {
	s.ChangeState(NewRotateState())
}

func (s *Switch) ChangeState(state State) {
	if s.state != nil {
		s.state.Exit(g, s)
		if !s.state.AllowChange(state) {
			fmt.Println("Change state not allowed")
			return
		}
	}
	s.state = state
	s.state.Enter(g, s)
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
