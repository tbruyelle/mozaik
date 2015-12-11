package main

import (
	"log"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/exp/sprite/glsprite"
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
	images       *glutil.Images
	eng          sprite.Engine
	fps          *debug.FPS
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					glctx = nil
				}
			case size.Event:
				sz = e
			case paint.Event:
				if glctx == nil || e.External {
					// Not ready yet
					continue
				}
				draw(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				touch_(sz, e)
			}
		}
	})
}

func onStart(glctx gl.Context) {
	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
	eng = glsprite.Engine(images)
}

func onStop() {
	eng.Release()
	fps.Release()
	images.Release()
}

func initialize(glctx gl.Context, sz size.Event) {
	glctx.Disable(gl.DEPTH_TEST)
	// antialiasing
	glctx.Enable(gl.BLEND)
	glctx.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	NewGame(glctx, sz)
}

func draw(glctx gl.Context, sz size.Event) {
	// Keep until golang.org/x/mogile/x11.go handle Start callback
	if g == nil {
		initialize(glctx, sz)
	}

	now := clock.Time(time.Since(start) * FPS / time.Second)
	if now == lastClock {
		return
	}
	lastClock = now

	glctx.ClearColor(0.9, 0.09, 0.26, 0.0)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	g.world.Draw(glctx, now, sz)
	fps.Draw(sz)
}

func touch_(sz size.Event, t touch.Event) {
	log.Printf("TOUCH %+v", t)
	if t.Type == touch.TypeEnd {

		g.Click(float32(t.X)/sz.PixelsPerPt, float32(t.Y)/sz.PixelsPerPt)
	}
}
