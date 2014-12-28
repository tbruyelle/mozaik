package main

import (
	"fmt"
	"math"
	"reflect"
)

type State interface {
	Enter(g *Game, sw *Switch)
	Exit(g *Game, sw *Switch)
	Update(g *Game, sw *Switch)
	AllowChange(state State) bool
}

type IdleState struct {
}

func NewIdleState() State {
	return IdleState{}
}

func (s IdleState) Enter(g *Game, sw *Switch) {
}

func (s IdleState) Exit(g *Game, sw *Switch) {}

func (s IdleState) Update(g *Game, sw *Switch) {
}

func (s IdleState) AllowChange(state State) bool {
	switch state.(type) {
	case IdleState:
		return false
	}
	return true
}

func StateName(s State) string {
	return reflect.TypeOf(s).String()
}

// RotateState performs a 90d rotation
type RotateState struct {
	IdleState
}

func NewRotateState() State {
	return RotateState{}
}

const (
	rotateTicks         = 12
	rotateRevertTicks   = 6
	rotateDegree        = math.Pi / 2
	halfRotate          = rotateDegree / 2
	rotatePerTick       = rotateDegree / rotateTicks
	rotateRevertPerTick = rotateDegree / rotateRevertTicks
	scaleMin            = 0.9
)

func smoothstep(step float64, goal int) (r float64) {
	x := float64(step) / float64(goal)
	defer func() { fmt.Println("smooth", x, r) }()
	return 3*math.Pow(x, 2) - 2*math.Pow(x, 3)
}

func scaleStep(rotate float32) float32 {
	//return float32(math.Cos(float64(4*rotate))/12 + 1 - 1.0/12)
	return float32(math.Cos(float64(4*rotate))/12 + 0.91666)
}

func (s RotateState) Enter(g *Game, sw *Switch) {
	g.level.rotating = sw
	sw.rotate = 0
	sw.scale = 1
}

func (s RotateState) Exit(g *Game, sw *Switch) {
	g.level.RotateSwitch(sw)
	g.level.rotating = nil
	sw.rotate = 0
	sw.scale = 1
}

func (s RotateState) Update(g *Game, sw *Switch) {
	// Update the rotation
	sw.rotate += rotatePerTick
	if sw.rotate >= rotateDegree {
		// End of rotation
		sw.ChangeState(NewIdleState())
	} else {
		sw.scale = scaleStep(sw.rotate)
	}
}

func (s RotateState) AllowChange(state State) bool {
	switch state.(type) {
	case RotateState:
		return false
	}
	return true
}

// rotateStateReverse is used to cancel a previous rotate
type RotateStateReverse struct {
	RotateState
}

func NewRotateStateReverse() State {
	return RotateStateReverse{}
}

func (s RotateStateReverse) Exit(g *Game, sw *Switch) {
	g.level.RotateSwitchInverse(sw)
	sw.rotate = 0
	sw.scale = 1
	g.level.rotating = nil
}

func (s RotateStateReverse) Update(g *Game, sw *Switch) {
	sw.rotate -= rotateRevertPerTick
	if sw.rotate <= -rotateDegree {
		sw.ChangeState(NewIdleState())
	} else {
		sw.scale = scaleStep(sw.rotate)
	}
}

// ResetState cancels all moves
type ResetState struct {
	RotateStateReverse
}

func NewResetState() State {
	return ResetState{}
}

func (s ResetState) Update(g *Game, sw *Switch) {
	sw.rotate -= rotateRevertPerTick
	if sw.rotate <= -rotateDegree {
		// Process next switch
		last := g.level.PopLastRotated()
		if last != nil {
			if last != sw {
				sw.ChangeState(NewIdleState())
			}
			last.ChangeState(NewResetState())
		} else {
			g.listen = true
			sw.ChangeState(NewIdleState())
		}
	} else {
		sw.scale = scaleStep(sw.rotate)
	}
}
