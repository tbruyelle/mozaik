package main

import (
	"fmt"

	"golang.org/x/mobile/geom"
)

const (
	WindowWidth          = 576
	WindowHeight         = 704
	BlockSize            = 128
	BlockRadius          = 10
	BlockPadding         = 0
	BlockCornerSegments  = 6
	SwitchSize           = 48
	SwitchSegments       = 20
	DashboardHeight      = 128
	XMin                 = 32
	YMin                 = 32
	XMax                 = WindowHeight - 32
	YMax                 = WindowWidth - 32 - DashboardHeight
	SignatureBlockSize   = 32
	SignatureBlockRadius = 6
	LineWidth            = 2
	SignatureLineWidth   = 1
	BgSegments           = 24
)

var (
	windowWidth, windowHeight                float32
	blockSize, blockRadius, blockPadding     float32
	switchSize                               float32
	dashboardHeight                          float32
	xMin, yMin, xMax, yMax                   float32
	signatureBlockSize, signatureBlockRadius float32
	lineWidth, signatureLineWidth            float32
)

type World struct {
	background *Background
	switches   []*SwitchModel
	blocks     map[*Block]*BlockModel
}

func compute(val float32, factor float32) float32 {
	return val * factor
}

func (w *World) Reset() {

	// Clean
	if len(w.switches) > 0 {
		for _, s := range w.switches {
			s.Destroy()
		}
	}
	w.switches = nil
	if len(w.blocks) > 0 {
		for k, m := range w.blocks {
			m.Destroy()
			delete(w.blocks, k)

		}
	}
	// Create the blocks
	w.blocks = make(map[*Block]*BlockModel)
	for i := 0; i < len(g.level.blocks); i++ {
		for j := 0; j < len(g.level.blocks[i]); j++ {
			b := g.level.blocks[i][j]
			if b != nil {
				w.blocks[b] = NewBlockModel(b, blockSize, blockRadius)
			}
		}
	}
	for _, sw := range g.level.switches {
		w.switches = append(w.switches, NewSwitchModel(sw))
	}
}

type Game struct {
	currentLevel int
	level        Level
	listen       bool
	world        *World
}

func NewGame() *Game {
	return &Game{currentLevel: 1, listen: true, world: &World{}}
}

func (g *Game) Start() {
	g.ComputeSizes()
	// Load first level
	g.level = LoadLevel(g.currentLevel)
	g.world.background = NewBackground()
	g.world.Reset()
}

func (g *Game) ComputeSizes() {
	// Compute dimensions according to current window size
	windowWidth, windowHeight = geom.Width.Px(), geom.Height.Px()
	fmt.Println("window", windowWidth, windowHeight)
	widthFactor := windowWidth / WindowWidth
	heightFactor := windowHeight / WindowHeight

	blockSize = compute(BlockSize, widthFactor)
	fmt.Println("size", BlockSize, blockSize)
	blockRadius = compute(BlockRadius, widthFactor)
	blockPadding = compute(BlockPadding, widthFactor)
	switchSize = compute(SwitchSize, widthFactor)
	fmt.Println("switch", SwitchSize, switchSize)
	dashboardHeight = compute(DashboardHeight, heightFactor)
	xMin = compute(XMin, widthFactor)
	yMin = compute(YMin, heightFactor)
	xMax = compute(XMax, widthFactor)
	yMax = compute(YMax, heightFactor)
	signatureBlockSize = compute(SignatureBlockSize, widthFactor)
	signatureBlockRadius = compute(SignatureBlockRadius, widthFactor)
	signatureLineWidth = compute(SignatureLineWidth, widthFactor)
	lineWidth = compute(LineWidth, widthFactor)
}

func (g *Game) Stop() {
}

func (g *Game) Click(x, y int) {
	if g.listen {
		g.level.PressSwitch(x, y)
	}
}

func (g *Game) Listen() bool {
	return g.listen && g.level.rotating == nil
}

func (g *Game) Update() {
	for _, s := range g.level.switches {
		s.state.Update(g, s)
	}
}

func (g *Game) Continue() {
	if g.level.Win() {
		g.Warp()
	}
}

func (g *Game) Warp() {
	if g.listen {
		// Next level
		g.currentLevel++
		g.level = LoadLevel(g.currentLevel)
		g.world.Reset()
	}
}

func (g *Game) UndoLastMove() {
	if g.listen {
		g.level.UndoLastMove()
	}
}

func (g *Game) Reset() {
	sw := g.level.PopLastRotated()
	if sw != nil {
		g.listen = false
		sw.ChangeState(NewResetState())
	}
}
