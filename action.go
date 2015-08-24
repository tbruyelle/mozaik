package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"golang.org/x/mobile/exp/sprite/clock"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type ActionFunc func(o *Object, t clock.Time)

func (a ActionFunc) Do(o *Object, t clock.Time) {
	a(o, t)
}

// wait pauses the display of the current object
type wait struct {
	until clock.Time
	next  Action
}

func (w wait) Do(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
		o.Dead = true
		return
	}
	if t > o.Time+w.until {
		// Once the time is elapsed,
		// start the next Action
		o.Time = 0
		o.Dead = false
		o.Action = w.next
	}
}

var colorTexMap = map[Color]int{
	Empty:       texEmpty,
	Red:         texBlockRed,
	Yellow:      texBlockYellow,
	Blue:        texBlockBlue,
	Green:       texBlockGreen,
	Pink:        texBlockPink,
	Orange:      texBlockOrange,
	LightBlue:   texBlockLightBlue,
	Purple:      texBlockPurple,
	Brown:       texBlockBrown,
	LightGreen:  texBlockLightGreen,
	Cyan:        texBlockCyan,
	LightPink:   texBlockLightPink,
	White:       texBlockWhite,
	LightPurple: texBlockLightPurple,
	LightBrown:  texBlockLightBrown,
	OtherWhite:  texBlockOtherWhite,
}

func blockSprite(o *Object) {
	b, ok := o.Data.(*Block)
	if !ok {
		log.Println("Invalid type assertion", o.Data)
		return
	}
	o.Sprite = g.world.texs[colorTexMap[b.Color]]
}

func blockIdle(o *Object, t clock.Time) {
	if o.Time == 0 {
		// Ensure no transformation in the idle action
		o.Time = t
		o.Reset()
	}
	blockSprite(o)
	if g.level.Win() {
		o.Time = 0
		o.Action = wait{until: clock.Time((o.X + o.Y) / 20), next: ActionFunc(blockPopOut)}
		return
	}
	o.Angle, o.Sx, o.Sy = 0, 0, 0
}

func signatureBlockIdle(o *Object, t clock.Time) {
	blockSprite(o)
}

func blockRotate(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
	}
	f := clock.EaseOut(o.Time, o.Time+16, t)
	o.Angle = math.Pi / 2 * f
	if f == 1 {
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
		o.Time = 0
		o.Action = ActionFunc(blockIdle)
		return
	}
	blockSprite(o)
	// Update also the scaling
	if f > .5 {
		f = (f - .5) / .5
		o.Scale = .8 + .2*f
	} else {
		f = f / .5
		o.Scale = 1 - .2*f
	}
}

func blockRotateInverse(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
	}
	f := clock.EaseOut(o.Time, o.Time+12, t)
	o.Angle = -math.Pi / 2 * f
	if f == 1 {
		// The rotation is over
		// First apply the rotation to the level struct
		// Use a mutex because this must be done only one time
		g.level.Lock()
		if g.level.rotating != nil {
			g.level.RotateSwitchInverse(g.level.rotating)
		}
		g.level.rotating = nil
		g.level.Unlock()
		// Return to the idle action
		o.Action = ActionFunc(blockIdle)
		return
	}
	blockSprite(o)
	// Update also the scaling
	if f > .5 {
		f = (f - .5) / .5
		o.Scale = .8 + .2*f
	} else {
		f = f / .5
		o.Scale = 1 - .2*f
	}
}

func blockPopIn(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
	}
	blockSprite(o)
	o.Dead = false
	f := clock.EaseOut(o.Time, o.Time+20, t)
	o.Tx = -o.X - o.Width + (o.X+o.Width)*f
	o.Ty = -o.Y - o.Height + (o.Y+o.Height)*f
	if f == 1 {
		o.Reset()
		o.Action = ActionFunc(blockIdle)
	}
}

func blockPopOut(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
		o.Rx = o.X + blockSize/2
		o.Ry = o.Y + blockSize/2
	}
	f := clock.EaseIn(o.Time, o.Time+100, t)
	o.Angle += f
	if f >= .3 {
		// Start moving
		if o.X < windowWidth/2 {
			o.Tx -= f * 3
		} else {
			o.Tx += f * 3
		}
		if o.Y < yMin+blockSize*2 {
			o.Ty -= f * 3
		} else {
			o.Ty += f * 3
		}
	}
}

func switchPopIn(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
		o.Sx = o.X + o.Width/2
		o.Sy = o.Y + o.Height/2
	}
	switchSprite(o)
	o.Scale = clock.EaseOut(o.Time, o.Time+20, t)
	if o.Scale == 1 {
		o.Reset()
		o.Action = ActionFunc(switchIdle)
	}
}

func switchPopOut(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
		o.Sx = o.X + o.Width/2
		o.Sy = o.Y + o.Height/2
	}
	switchSprite(o)
	o.Scale = 1 - clock.EaseIn(o.Time, o.Time+20, t)
}

func switchIdle(o *Object, t clock.Time) {
	if g.level.Win() {
		o.Time = 0
		o.Action = ActionFunc(switchPopOut)
		return
	}
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

func looseTxtPop(o *Object, t clock.Time) {
	o.Dead = g.level.Win() || g.level.moves < g.level.maxMoves
	if !o.Dead {
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
			o.Action = ActionFunc(winTxtZoomIn)
			g.listen = true
		}
	}
}

func winTxtZoomIn(o *Object, t clock.Time) {
	if o.Time == 0 {
		// Start the animation
		o.Time = t
		o.Sx = o.X + o.Width/2
		o.Sy = o.Y + o.Height/2
	}
	o.Scale = 1 + clock.EaseIn(o.Time, o.Time+20, t)*.2
	if o.Scale == 1.2 {
		o.Time = 0
		o.Action = ActionFunc(winTxtZoomOut)
	}
}

func winTxtZoomOut(o *Object, t clock.Time) {
	if o.Time == 0 {
		o.Time = t
	}
	o.Scale = 1.2 - clock.EaseOut(o.Time, o.Time+25, t)*.2
	if o.Scale == 1 {
		o.Time = 0
		o.Action = ActionFunc(winTxtZoomIn)
	}
}
