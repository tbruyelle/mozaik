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
	DashboardHeight      = 144
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
	CharWidth            = 40
	CharHeight           = 54
	LooseTxtWidth        = 338
	LooseTxtHeight       = 307
)

var (
	windowWidth, windowHeight                float32
	blockSize, blockRadius, blockPadding     float32
	switchSize                               float32
	dashboardHeight                          float32
	xMin, yMin, xMax, yMax                   float32
	signatureBlockSize, signatureBlockRadius float32
	lineWidth, signatureLineWidth            float32
	winTxtWidth, winTxtHeight                float32
	charWidth, charHeight                    float32
	looseTxtWidth, looseTxtHeight            float32
)

type Game struct {
	currentLevel int
	level        Level
	listen       bool
	world        *World
}

func NewGame() {
	g = &Game{currentLevel: 1, listen: true}
	computeSizes()
	g.level = LoadLevel(g.currentLevel)
	g.world = NewWorld()
}

func computeSizes() {
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
	charWidth = compute(CharWidth, widthFactor)
	charHeight = compute(CharHeight, widthFactor)
	looseTxtWidth = compute(LooseTxtWidth, widthFactor)
	looseTxtHeight = compute(LooseTxtHeight, widthFactor)
}

func (g *Game) Stop() {
}

func (g *Game) Click(x, y float32) {
	if g.Listen() {
		switch {
		case g.level.Win():
			// Next level
			g.Warp()

		case g.level.moves >= g.level.maxMoves:
			// Loose, restart
			g.currentLevel = 1
			g.level = LoadLevel(g.currentLevel)
			//FIXME clean resources
			g.world.LoadScene()

		case x < 30 && y < 30:
			// Trick to warp level
			// FIXME remove me
			g.Warp()

		case x > windowWidth-30 && y < 30:
			// Trick to undo moves
			// FIXME remove me
			g.level.UndoLastMove()

		default:
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
	if g.Listen() {
		// Next level
		g.currentLevel++
		g.level = LoadLevel(g.currentLevel)
		//FIXME clean resources
		g.world.LoadScene()
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
