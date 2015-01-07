package main

import (
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/gl"
)

const (
	FRAMES_PER_SECOND = 30
)

var (
	g            *Game
	windowRadius float64
	ticker       *time.Ticker
)

func main() {
	app.Run(app.Callbacks{
		Draw:  draw,
		Touch: touch,
	})
}

func initialize() {
	gl.Disable(gl.DEPTH_TEST)
	// antialiasing
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	g = NewGame()
	g.Start()
}

func draw() {
	// Keep until golang.org/x/mogile/x11.go handle Start callback
	if g == nil {
		initialize()
		ticker = time.NewTicker(time.Duration(1e9 / int(FRAMES_PER_SECOND)))
	}

	select {
	case <-ticker.C:
		gl.ClearColor(0.9, 0.85, 0.46, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		g.world.Draw()
		debug.DrawFPS()
	}
}

func touch(t event.Touch) {
	if t.Type == event.TouchEnd {
		g.Click(float32(t.Loc.X), float32(t.Loc.Y))
	}
}
