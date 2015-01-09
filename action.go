package main

import (
	"golang.org/x/mobile/sprite/clock"
	"log"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

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

func blockSprite(o *Object) {
	b, ok := o.Data.(*Block)
	if !ok {
		log.Println("Invalid type assertion", o.Data)
		return
	}
	o.Sprite = g.world.texs[b.Color]
}

func blockRotate(o *Object, t clock.Time) {
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
	blockSprite(o)
	// Update also the scaling
	scale := float32(math.Cos(float64(o.Angle*4))/12 + .91666)
	o.Sx, o.Sy = scale, scale
}

func blockPopStart(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t + clock.Time(rand.Intn(15))
		o.Dead = true
		return
	}
	if t > o.Time {
		o.Time = 0
		o.Action = blockPop
	}
}

func blockPop(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
		return
	}
	blockSprite(o)
	o.Dead = false
	f := clock.EaseIn(o.Time, o.Time+150, t)
	o.Tx = -o.X + o.Y/o.X*o.Y*f
	o.Ty = -o.Y + o.X/o.Y*o.X*f
	if f == 1 {
		o.Reset()
		o.Action = blockIdle
	}
}

func switchPop(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
		o.Dead = true
		return
	}
	if t <= o.Time+35 {
		// Wait until all the blocks have popped
		return
	}
	o.Dead = false
	switchSprite(o)
	f := clock.EaseIn(o.Time+35, o.Time+55, t)
	o.ZoomIn(f, 0)
	if f == 1 {
		o.Reset()
		o.Action = switchIdle
	}
}

func switchIdle(o *Object, t clock.Time) {
	switchSprite(o)
}

func switchSprite(o *Object) {
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
			o.Reset()
			o.Action = winTxtZoomIn
			g.listen = true
		}
	}
}

func winTxtZoomIn(o *Object, t clock.Time) {
	if o.Time == 0 {
		// Start the animation
		o.Time = t
	}
	f := clock.EaseIn(o.Time, o.Time+20, t) * .2
	o.ZoomIn(f, 1)
	if f == .2 {
		o.Time = 0
		o.Action = winTxtZoomOut
	}
}

func winTxtZoomOut(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
	}
	f := clock.EaseOut(o.Time, o.Time+25, t) * .2
	o.ZoomOut(f, 1.2)
	if f == .2 {
		o.Time = 0
		o.Action = winTxtZoomIn
	}
}
