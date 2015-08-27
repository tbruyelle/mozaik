package main

import (
	"log"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/gl"
)

const (
	FPS = 60
)

var (
	g            *Game
	windowRadius float64
	start        = time.Now()
	lastClock    = clock.Time(-1)
)

func main() {
	app.Main(func(a app.App) {
		var sz size.Event
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case size.Event:
				sz = e
			case paint.Event:
				draw(sz)
				a.EndPaint(e)
			case touch.Event:
				touch_(sz, e)
			}
		}
	})
}

func initialize(sz size.Event) {
	gl.Disable(gl.DEPTH_TEST)
	// antialiasing
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	NewGame(sz)
}

func draw(sz size.Event) {
	// Keep until golang.org/x/mogile/x11.go handle Start callback
	if g == nil {
		initialize(sz)
	}

	now := clock.Time(time.Since(start) * FPS / time.Second)
	if now == lastClock {
		return
	}
	lastClock = now

	gl.ClearColor(0.9, 0.09, 0.26, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	g.world.Draw(now, sz)
	debug.DrawFPS(sz)
}

func touch_(sz size.Event, t touch.Event) {
	log.Printf("TOUCH %+v", t)
	if t.Type == touch.TypeEnd {

		g.Click(float32(t.X)/sz.PixelsPerPt, float32(t.Y)/sz.PixelsPerPt)
	}
}
