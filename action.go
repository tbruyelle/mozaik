package main

import (
	"golang.org/x/mobile/sprite/clock"
	"log"
	"math"
)

func blockIdle(o *Object, t clock.Time) {
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

func blockRotate(o *Object, t clock.Time) {
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

func switchIdle(o *Object, t clock.Time) {
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

func winTxtPop(o *Object, t clock.Time) {
	o.Dead = !g.level.Win()
	if !o.Dead {
		// Wait until animation done
		g.listen = false
		if o.Time == 0 {
			// Set time for the first pass
			o.Time = t
		}
		// Compute a translation animation
		f := clock.EaseIn(o.Time, o.Time+20, t)
		x := o.X + o.Width
		o.Tx = -x + x*f
		if f == 1 {
			// First animation is over
			o.Time = 0
			o.Action = winTxtBreethOut
			g.listen = true
		}
	}
}

func winTxtBreethOut(o *Object, t clock.Time) {
	if o.Time == 0 {
		// Start the animation
		o.Time = t
	}
	f := clock.EaseIn(o.Time, o.Time+20, t) * .2
	// Scale
	s := f + 1
	o.Sx, o.Sy = s, s
	// Translate to keep centered
	o.Tx = -o.Width * f / 2
	o.Ty = -o.Height * f / 2
	if f == .2 {
		o.Time = 0
		o.Action = winTxtBreethIn
	}
}

func winTxtBreethIn(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
	}
	f := clock.EaseOut(o.Time, o.Time+25, t) * .2
	s := 1.2 - f
	// Scale
	o.Sx, o.Sy = s, s
	// Translate to keep centered
	mw, mh := o.Width/2, o.Height/2
	o.Tx = -mw*.2 + mw*f
	o.Ty = -mh*.2 + mh*f
	if f == .2 {
		o.Time = 0
		o.Action = winTxtBreethOut
	}
}
