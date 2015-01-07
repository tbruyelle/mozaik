package main

import (
	"log"
	"math"
)

func blockIdle(o *Object) {
	// Ensure no transformation in the idle action
	o.Angle, o.Sx, o.Sy = 0, 0, 0
	b, ok := o.Data.(*Block)
	if !ok {
		log.Println("Invalid type assertion", o.Data)
		return
	}
	o.Sprite = g.world.texs[b.Color]
}

const (
	rotateTicks         = 10
	rotateRevertTicks   = 6
	rotateComplete      = math.Pi / 2
	halfRotate          = rotateComplete / 2
	rotatePerTick       = rotateComplete / rotateTicks
	rotateRevertPerTick = rotateComplete / rotateRevertTicks
	scaleMin            = 0.9
)

func scaleStep(rotate float32) float32 {
	//return float32(math.Cos(float64(4*rotate))/12 + 1 - 1.0/12)
	return float32(math.Cos(float64(4*rotate))/12 + 0.91666)
}

func blockRotate(o *Object) {
	b, ok := o.Data.(*Block)
	if !ok {
		log.Println("Invalid type assertion", o.Data)
		return
	}
	o.Angle += rotatePerTick
	if o.Angle >= rotateComplete {
		// The rotation is over
		// First apply the rotation to the level struct
		// Use a mutex because this must be done only one time
		g.level.Lock()
		if g.level.rotating != nil {
			g.level.RotateSwitch(g.level.rotating)
		}
		g.level.rotating = nil
		g.level.Unlock()
		// Return to the idle action
		o.Action = blockIdle
		return
	}
	o.Sprite = g.world.texs[b.Color]
	// Update also the scaling
	scale := float32(math.Cos(float64(o.Angle*4))/12 + .91666)
	o.Sx, o.Sy = scale, scale
}

func switchIdle(o *Object) {
	sw, ok := o.Data.(*Switch)
	if !ok {
		log.Println("Invalid type assertion", o.Data)
		return
	}
	switch sw.name {
	case "1":
		o.Sprite = g.world.texs[texSwitch1]
	case "2":
		o.Sprite = g.world.texs[texSwitch2]
	case "3":
		o.Sprite = g.world.texs[texSwitch3]
	case "4":
		o.Sprite = g.world.texs[texSwitch4]
	case "5":
		o.Sprite = g.world.texs[texSwitch5]
	case "6":
		o.Sprite = g.world.texs[texSwitch6]
	case "7":
		o.Sprite = g.world.texs[texSwitch7]
	case "8":
		o.Sprite = g.world.texs[texSwitch8]
	case "9":
		o.Sprite = g.world.texs[texSwitch9]
	}
}

func winTxtPop(o *Object) {
	o.Dead = !g.level.Win()
}
