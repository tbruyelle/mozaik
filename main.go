package main

import (
	"math"
	"time"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/geom"
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
	g = NewGame()
	width, height := geom.Width.Px(), geom.Height.Px()

	gl.Viewport(0, 0, int(width), int(height))

	// Compute window radius
	windowRadius = math.Sqrt(math.Pow(float64(height), 2) + math.Pow(float64(width), 2))

	//gl.Init()
	gl.Disable(gl.DEPTH_TEST)
	// antialiasing
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
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
		// Update all component states
		g.Update()

		// Reinit blocks as not renderered
		for i := 0; i < len(g.level.blocks); i++ {
			for j := 0; j < len(g.level.blocks[i]); j++ {
				if g.level.blocks[i][j] != nil {
					g.level.blocks[i][j].Rendered = false
				}
			}
		}

		gl.ClearColor(0.9, 0.85, 0.46, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		w := g.world
		w.background.Draw()
		if g.level.rotating != nil {
			// Start draw the rotating switch
			for _, swm := range w.switches {
				if swm.sw == g.level.rotating {
					swm.Draw()
				}
			}
		}
		// Draw the remaining switches
		for _, swm := range w.switches {
			swm.Draw()
		}
		debug.DrawFPS()
	}
}

func touch(t event.Touch) {
	if t.Type == event.TouchEnd {
		g.Click(int(t.Loc.X.Px()), int(t.Loc.Y.Px()))
	}
}
