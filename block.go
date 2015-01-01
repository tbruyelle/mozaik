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
	Rendered bool
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

func (s *Switch) String() string {
	return fmt.Sprintf("sw{line:%d, col:%d}", s.line, s.col)
}
