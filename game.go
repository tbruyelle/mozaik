package main

import (
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"
	_ "image/png"
	"log"
	"math"
)

const (
	PortraitWidth        = 576
	PortraitHeight       = 704
	LandscapeWidth       = PortraitHeight
	LandscapeHeight      = PortraitWidth
	BlockSize            = 128
	BlockRadius          = 10
	BlockPadding         = 0
	BlockCornerSegments  = 6
	SwitchSize           = 48
	SwitchSegments       = 20
	DashboardSize        = 144
	Padding              = 24
	XMin                 = Padding
	YMin                 = Padding
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
	portrait                                 bool
	windowWidth, windowHeight, padding       float32
	blockSize, blockRadius, blockPadding     float32
	switchSize                               float32
	dashboardSize                            float32
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

func NewGame(glctx gl.Context) {
	g = &Game{currentLevel: 1, listen: true}
	g.level = LoadLevel(g.currentLevel)
}

func initWorld(glctx gl.Context) {
	g.world = NewWorld(glctx)
}

func computeSizes(sz size.Event) {
	// Compute dimensions according to current window size
	windowWidth, windowHeight = float32(sz.WidthPt), float32(sz.HeightPt)
	portrait = windowHeight >= windowWidth
	log.Println("window", windowWidth, windowHeight, "portrait", portrait)

	var widthFactor, heightFactor float32
	if portrait {
		widthFactor = windowWidth / PortraitWidth
		heightFactor = windowHeight / PortraitHeight
	} else {
		widthFactor = windowWidth / LandscapeWidth
		heightFactor = windowHeight / LandscapeHeight
	}
	log.Println("factors", widthFactor, heightFactor)
	var minFactor float32
	if widthFactor < heightFactor {
		minFactor = widthFactor
	} else {
		minFactor = heightFactor
	}

	windowRadius = math.Sqrt(math.Pow(float64(windowHeight), 2) + math.Pow(float64(windowWidth), 2))

	padding = compute(Padding, minFactor)
	blockSize = compute(BlockSize, minFactor)
	log.Println("block size", BlockSize, blockSize)
	blockRadius = compute(BlockRadius, minFactor)
	blockPadding = compute(BlockPadding, minFactor)
	switchSize = compute(SwitchSize, minFactor)
	log.Println("switch size", SwitchSize, switchSize)
	dashboardSize = compute(DashboardSize, minFactor)
	xMin = compute(XMin, minFactor)
	yMin = compute(YMin, minFactor)
	if portrait {
		xMax = compute(PortraitWidth-Padding, widthFactor)
		yMax = compute(PortraitHeight-Padding-DashboardSize, heightFactor)
	} else {
		xMax = compute(LandscapeWidth-Padding-DashboardSize, widthFactor)
		yMax = compute(LandscapeHeight-Padding, heightFactor)
	}
	signatureBlockSize = compute(SignatureBlockSize, minFactor)
	signatureBlockRadius = compute(SignatureBlockRadius, minFactor)
	signatureLineWidth = compute(SignatureLineWidth, minFactor)
	lineWidth = compute(LineWidth, minFactor)
	winTxtWidth = compute(WinTxtWidth, minFactor)
	winTxtHeight = compute(WinTxtHeight, minFactor)
	charWidth = compute(CharWidth, minFactor)
	charHeight = compute(CharHeight, minFactor)
	looseTxtWidth = compute(LooseTxtWidth, minFactor)
	looseTxtHeight = compute(LooseTxtHeight, minFactor)
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
