// +build ignore

package main

import (
	"fmt"
	"github.com/andrebq/gas"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/gltext"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	FRAME_RATE          = time.Second / 40
	BlockCornerSegments = 6
	SwitchSegments      = 20
)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

// Queue of work to run in main thread.
var mainfunc = make(chan func())

// Run all the functions that need to run in the main thread.
func Main() {
	var f func()
	for f = range mainfunc {
		f()
	}
}

// Put the function f on the main thread function queue.
func do(f func()) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		f()
		done <- true
	}
	<-done
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

var (
	window       *glfw.Window
	err          error
	g            *Game
	fonts        []*gltext.Font
	fontInd      int
)

func main() {
	_ = fmt.Sprint()

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	// antialiasing
	//glfw.WindowHint(glfw.Samples, 4)

	window, err = glfw.CreateWindow(WindowWidth, WindowHeight, "Mozaik", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Ensure thread context
	window.MakeContextCurrent()
	//glfw.SwapInterval(1)

	window.SetKeyCallback(keyCb)
	window.SetMouseButtonCallback(mouseCb)

	gl.Init()
	gl.ClearColor(0.9, 0.85, 0.46, 0.0)
	// useless in 2D
	gl.Disable(gl.DEPTH_TEST)
	// antialiasing
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.LINE_SMOOTH)

	for i := int32(32); i < 72; i++ {
		font := loadFonts(i)
		defer font.Release()
		fonts = append(fonts, font)
	}

	// Compute window radius
	windowRadius = math.Sqrt(math.Pow(WindowHeight, 2) + math.Pow(WindowWidth, 2))

	// Use window coordinates
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, WindowWidth, WindowHeight, 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	g = NewGame()

	g.Start()
	go eventLoop(g)
	go renderLoop(g)
	Main()
	g.Stop()
}

// Load the font
func loadFonts(scale int32) *gltext.Font {
	file, err := gas.Abs("code.google.com/p/freetype-go/testdata/luxisr.ttf")
	if err != nil {
		panic(err)
	}
	fd, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	font, err := gltext.LoadTruetype(fd, scale, 32, 127, gltext.LeftToRight)
	if err != nil {
		panic(err)
	}
	return font
}

func keyCb(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			g.Stop()
			os.Exit(0)

		case glfw.KeyR:
			g.Reset()

		case glfw.KeyW:
			g.Warp()

		case glfw.KeyC:
			g.UndoLastMove()

		case glfw.KeySpace:
			g.Continue()

		case glfw.Key1:
			keySwitch("1")
		case glfw.Key2:
			keySwitch("2")
		case glfw.Key3:
			keySwitch("3")
		case glfw.Key4:
			keySwitch("4")
		case glfw.Key5:
			keySwitch("5")
		case glfw.Key6:
			keySwitch("6")
		case glfw.Key7:
			keySwitch("7")
		case glfw.Key8:
			keySwitch("8")
		case glfw.Key9:
			keySwitch("9")

		}
	}
}

func keySwitch(name string) {
	if g.Listen() {
		g.level.TriggerSwitchName(name)
	}
}

func mouseCb(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		switch button {

		case glfw.MouseButtonLeft:
			x, y := window.GetCursorPosition()
			g.Click(int(x), int(y))
		}
	}
}

func eventLoop(g *Game) {
	eventTicker := time.NewTicker(time.Millisecond * 10)

	for {
		select {
		case <-eventTicker.C:
			do(func() {
				glfw.PollEvents()
			})
		}
	}
}

func renderLoop(g *Game) {
	mainTicker := time.NewTicker(FRAME_RATE)

	for {
		select {
		case <-mainTicker.C:
			g.Update()
			do(func() {
				render(g)
			})
		}
	}
}

func render(g *Game) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Background
	renderBackground()

	// Reinit blocks as not renderered
	for i := 0; i < len(g.level.blocks); i++ {
		for j := 0; j < len(g.level.blocks[i]); j++ {
			if g.level.blocks[i][j] != nil {
				g.level.blocks[i][j].Rendered = false
			}
		}
	}

	// Render first the rotating switch if there's one
	if g.level.rotating != nil {
		renderSwitchBlocks(g.level.rotating)
	}
	// Render the remaining switches
	for _, s := range g.level.switches {
		if s != g.level.rotating {
			renderSwitchBlocks(s)
		}
	}

	// render the switches
	for _, s := range g.level.switches {
		renderSwitch(s)
	}

	renderDashboard()

	if g.level.Win() {
		fontInd++
		if fontInd >= len(fonts) {
			fontInd = 0
		}
		gl.LoadIdentity()
		gl.Color3ub(107, 142, 35)
		fonts[fontInd].Printf(float32(WindowWidth/2)-100, float32(WindowHeight/2), "GAGNE !!")
	}

	window.SwapBuffers()
}

func renderSwitchBlocks(s *Switch) {
	// TODO constant
	v := SwitchSize / 2
	x, y := float32(s.X+v), float32(s.Y+v)
	gl.LoadIdentity()
	gl.Translatef(x, y, 0)
	if s.rotate != 0 {
		gl.Rotatef(float32(s.rotate), 0, 0, 1)
	}
	bsf := float32(BlockSize - s.Z)
	padding := float32(BlockPadding)

	var b *Block
	// Render block top left
	b = g.level.blocks[s.line][s.col]
	if !b.Rendered {
		gl.PushMatrix()
		gl.Translatef(-bsf-padding, -bsf-padding, 0)
		renderBlock(b, bsf)
		gl.PopMatrix()
		b.Rendered = true
	}

	// Render block top right
	b = g.level.blocks[s.line][s.col+1]
	if !b.Rendered {
		gl.PushMatrix()
		gl.Translatef(padding, -bsf-padding, 0)
		renderBlock(b, bsf)
		gl.PopMatrix()
		b.Rendered = true
	}

	// Render block bottom right
	b = g.level.blocks[s.line+1][s.col+1]
	if !b.Rendered {
		gl.PushMatrix()
		gl.Translatef(padding, padding, 0)
		renderBlock(b, bsf)
		gl.PopMatrix()
		b.Rendered = true
	}

	// render block bottom left
	b = g.level.blocks[s.line+1][s.col]
	if !b.Rendered {
		gl.PushMatrix()
		gl.Translatef(-bsf-padding, padding, 0)
		renderBlock(b, bsf)
		gl.PopMatrix()
		b.Rendered = true
	}
}

func renderBlock(b *Block, s float32) {
	renderBlock_(b.Color, s, s, BlockRadius, LineWidth)
}

func renderBlockSignature(color ColorDef) {
	renderBlock_(color, SignatureBlockSize, SignatureBlockSize, SignatureBlockRadius, SignatureLineWidth)
}

func renderBlock_(color ColorDef, w, h, radius, lineWidth float32) {
	setColor(color)
	gl.Begin(gl.QUADS)
	// Render inner square
	gl.Vertex2f(radius, radius)
	gl.Vertex2f(radius, h-radius)
	gl.Vertex2f(w-radius, h-radius)
	gl.Vertex2f(w-radius, radius)
	// Render top square
	gl.Vertex2f(radius, h-radius)
	gl.Vertex2f(radius, h)
	gl.Vertex2f(w-radius, h)
	gl.Vertex2f(w-radius, h-radius)
	// Render bottom square
	gl.Vertex2f(radius, radius)
	gl.Vertex2f(radius, 0)
	gl.Vertex2f(w-radius, 0)
	gl.Vertex2f(w-radius, radius)
	// Render left square
	gl.Vertex2f(w-radius, radius)
	gl.Vertex2f(w, radius)
	gl.Vertex2f(w, h-radius)
	gl.Vertex2f(w-radius, h-radius)
	// Render right square
	gl.Vertex2f(radius, radius)
	gl.Vertex2f(0, radius)
	gl.Vertex2f(0, h-radius)
	gl.Vertex2f(radius, h-radius)
	gl.End()
	// Render bottom right corner
	ww, hh := float64(w-radius), float64(h-radius)
	renderCorner(color, ww, hh, 0, radius, lineWidth)
	// Render bottom left corner
	ww, hh = float64(radius), float64(h-radius)
	renderCorner(color, ww, hh, 1, radius, lineWidth)
	// Render top left corner
	ww, hh = float64(radius), float64(radius)
	renderCorner(color, ww, hh, 2, radius, lineWidth)
	// Render top right corner
	ww, hh = float64(w-radius), float64(radius)
	renderCorner(color, ww, hh, 3, radius, lineWidth)

	if lineWidth != 0 {
		// Render the shape
		gl.LineWidth(lineWidth)
		gl.Color3i(0, 0, 0)
		gl.Begin(gl.LINES)

		gl.Vertex2f(radius, 0)
		gl.Vertex2f(w-radius, 0)

		gl.Vertex2f(0, radius)
		gl.Vertex2f(0, h-radius)

		gl.Vertex2f(w, radius)
		gl.Vertex2f(w, h-radius)

		gl.Vertex2f(radius, h)
		gl.Vertex2f(w-radius, h)

		gl.End()
	}
	gl.PopMatrix()
}

func renderCorner(color ColorDef, ww, hh, start float64, radius, lineWidth float32) {
	setColor(color)
	max := BlockCornerSegments * (start + 1)
	// Render the corner
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2d(ww, hh)
	for i := start * BlockCornerSegments; i <= max; i++ {
		a := math.Pi / 2 * i / BlockCornerSegments
		x := math.Cos(a) * float64(radius)
		y := math.Sin(a) * float64(radius)
		gl.Vertex2d(ww+x, hh+y)
	}
	gl.End()

	if lineWidth != 0 {
		// Render the shape
		gl.LineWidth(lineWidth)
		gl.Color3i(0, 0, 0)
		gl.Begin(gl.LINE_STRIP)
		for i := start * BlockCornerSegments; i <= max; i++ {
			a := math.Pi / 2 * i / BlockCornerSegments
			x := math.Cos(a) * float64(radius)
			y := math.Sin(a) * float64(radius)
			gl.Vertex2d(ww+x, hh+y)
		}
		gl.End()
	}
}

func renderSwitch(s *Switch) {
	gl.LoadIdentity()
	// TODO constant
	v := SwitchSize / 2
	x, y := float32(s.X+v), float32(s.Y+v)
	gl.Translatef(x, y, 0)
	// Render the switch
	gl.Color3f(1, 1, 1)
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Vertex2d(0, 0)
	vv := float64(v)
	for i := float64(0); i <= SwitchSegments; i++ {
		a := 2 * math.Pi * i / SwitchSegments
		gl.Vertex2d(math.Sin(a)*vv, math.Cos(a)*vv)
	}
	gl.End()

	if LineWidth != 0 {
		// Render the shape
		gl.Color3i(0, 0, 0)
		gl.LineWidth(LineWidth)
		gl.Begin(gl.LINE_LOOP)
		for i := float64(0); i <= SwitchSegments; i++ {
			a := 2 * math.Pi * i / SwitchSegments
			gl.Vertex2d(math.Sin(a)*vv, math.Cos(a)*vv)
		}
		gl.End()
	}

	// Write the switch name
	gl.LoadIdentity()
	w, h := fonts[6].Metrics(s.name)
	gl.Color3i(0, 0, 0)
	fonts[6].Printf(x-float32(w)/2, y-float32(h)/2+2, s.name)
}

func setColor(color ColorDef) {
	switch color {
	case Red:
		gl.Color3ub(239, 14, 84)
	case Yellow:
		gl.Color3ub(255, 218, 58)
	case Blue:
		gl.Color3ub(100, 149, 237)
	case Green:
		gl.Color3ub(88, 164, 0)
	case Pink:
		gl.Color3ub(255, 178, 255)
	case Orange:
		gl.Color3ub(243, 122, 17)
	case LightBlue:
		gl.Color3ub(98, 222, 255)
	}
}

func renderDashboard() {
	gl.LoadIdentity()
	setColor(Blue)
	fonts[32].Printf(float32(XMin), float32(WindowHeight-DashboardHeight),
		"Level %d", g.currentLevel)
	gl.Translatef(float32(XMin)+300, float32(WindowHeight-DashboardHeight), 0)
	line := 0
	for _, c := range g.level.winSignature {
		if c == '\n' {
			line++
			gl.LoadIdentity()
			gl.Translatef(float32(XMin)+300, float32(WindowHeight-DashboardHeight+SignatureBlockSize*line+BlockPadding), 0)
			continue
		}
		if c != '-' {
			renderBlockSignature(atoc(string(c)))
		}
		gl.Translated(SignatureBlockSize+BlockPadding, 0, 0)
	}
}
