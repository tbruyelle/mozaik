package main

import (
	"fmt"
	_ "image/png"
	"math"

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
	WinTxtWidth          = 300
	WinTxtHeight         = 90
)

type Game struct {
	currentLevel int
	level        Level
	listen       bool
	world        *World
}

func NewGame() *Game {
	game := &Game{currentLevel: 1, listen: true}
	return game
}

func (g *Game) Start() {
	g.ComputeSizes()
	g.level = LoadLevel(g.currentLevel)
	g.world = NewWorld()
}

func (g *Game) ComputeSizes() {
	// Compute dimensions according to current window size
	windowWidth, windowHeight = float32(geom.Width), float32(geom.Height)

	fmt.Println("window", windowWidth, windowHeight)
	widthFactor := windowWidth / WindowWidth
	heightFactor := windowHeight / WindowHeight

	windowRadius = math.Sqrt(math.Pow(float64(windowHeight), 2) + math.Pow(float64(windowWidth), 2))

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
	winTxtWidth = compute(WinTxtWidth, widthFactor)
	winTxtHeight = compute(WinTxtHeight, widthFactor)
}

func (g *Game) Stop() {
}

func (g *Game) Click(x, y float32) {
	if g.listen {
		if g.level.Win() {
			// Next level
			g.Warp()
		} else {
			g.level.PressSwitch(x, y)
		}
	}
}

func (g *Game) Listen() bool {
	return g.listen && g.level.rotating == nil
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
		//FIXME clean resources
		g.world.LoadScene()
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
		// TODO
		//sw.ChangeState(NewResetState())
	}
}
