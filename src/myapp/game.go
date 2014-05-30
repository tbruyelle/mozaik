package main

import (
	"math"
	"runtime"
	"time"

	"git.tideland.biz/goas/loop"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/mandala"
)

const (
	FRAMES_PER_SECOND = 24
)

type viewportSize struct {
	width, height int
}

type renderLoopControl struct {
	resizeViewport chan viewportSize
	pause          chan mandala.PauseEvent
	resume         chan bool
	window         chan mandala.Window
}

var (
	window mandala.Window
	g      *Game
	// Background rotate state
	bgRotate     = float32(0)
	windowRadius float64
)

func newRenderLoopControl() *renderLoopControl {
	return &renderLoopControl{
		resizeViewport: make(chan viewportSize),
		pause:          make(chan mandala.PauseEvent),
		resume:         make(chan bool),
		window:         make(chan mandala.Window, 1),
	}
}

func draw() {
// Draw
	gl.ClearColor(0.9, 0.85, 0.46, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	w := g.world
	w.background.Draw()
}

// Run runs renderLoop. The loop renders a frame and swaps the buffer
// at each tick received.
func renderLoopFunc(control *renderLoopControl) loop.LoopFunc {

	return func(loop loop.Loop) error {
		var window mandala.Window
		// Lock/unlock the loop to the current OS thread. This is
		// necessary because OpenGL functions should be called from
		// the same thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// Create an instance of ticker and immediately stop
		// it because we don't want to swap buffers before
		// initializing a rendering state.
		ticker := time.NewTicker(time.Duration(1e9 / int(FRAMES_PER_SECOND)))
		ticker.Stop()

		for {
			select {
			case window = <-control.window:
				ticker.Stop()

				window.MakeContextCurrent()

				width, height := window.GetSize()
				gl.Viewport(0, 0, width, height)

				mandala.Logf("Restarting rendering loop...")
				ticker = time.NewTicker(time.Duration(1e9 / int(FRAMES_PER_SECOND)))

				// Compute window radius
				windowRadius = math.Sqrt(math.Pow(float64(height), 2) + math.Pow(float64(width), 2))

					gl.Init()
	gl.Disable(gl.DEPTH_TEST)
	// antialiasing
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.LINE_SMOOTH)

			// At each tick render a frame and swap buffers.
			case <-ticker.C:
				draw()
				window.SwapBuffers()

			case event := <-control.pause:
				ticker.Stop()
				event.Paused <- true

			case <-control.resume:

			case <-loop.ShallStop():
				ticker.Stop()
				return nil
			}
		}
	}
}

// eventLoopFunc listen to events originating from the
// framework.
func eventLoopFunc(renderLoopControl *renderLoopControl) loop.LoopFunc {
	return func(loop loop.Loop) error {

		for {
			select {

			// Receive an EGL state from the
			// framework and notify the render
			// loop about that.
			// case eglState := <-mandala.Init:
			// 	mandala.Logf("EGL surface initialized W:%d H:%d", eglState.SurfaceWidth, eglState.SurfaceHeight)
			// 	renderLoopControl.eglState <- eglState

			// Receive events from the framework.
			//
			// When the application starts the
			// typical events chain is:
			//
			// * onCreate
			// * onResume
			// * onInputQueueCreated
			// * onNativeWindowCreated
			// * onNativeWindowResized
			// * onWindowFocusChanged
			// * onNativeRedrawNeeded
			//
			// Pausing (i.e. clicking on the back
			// button) the application produces
			// following events chain:
			//
			// * onPause
			// * onWindowDestroy
			// * onWindowFocusChanged
			// * onInputQueueDestroy
			// * onDestroy

			case untypedEvent := <-mandala.Events():
				switch event := untypedEvent.(type) {

				// Receive a native window
				// from the framework and send
				// it to the render loop in
				// order to begin the
				// rendering process.
				case mandala.NativeWindowCreatedEvent:
					renderLoopControl.window <- event.Window

				// Finger down/up on the screen.
				case mandala.ActionUpDownEvent:
					if event.Down {
						mandala.Logf("Finger is DOWN at %f %f", event.X, event.Y)
					} else {
						mandala.Logf("Finger is now UP")
					}

					// Finger is moving on the screen.
				case mandala.ActionMoveEvent:
					mandala.Logf("Finger is moving at coord %f %f", event.X, event.Y)

				case mandala.DestroyEvent:
					mandala.Logf("Stop rendering...\n")
					mandala.Logf("Quitting from application...\n")
					return nil

				case mandala.NativeWindowRedrawNeededEvent:

				case mandala.PauseEvent:
					mandala.Logf("Application was paused. Stopping rendering ticker.")
					renderLoopControl.pause <- event

				case mandala.ResumeEvent:
					mandala.Logf("Application was resumed. Reactivating rendering ticker.")
					renderLoopControl.resume <- true

				}
			}
		}
	}
}
