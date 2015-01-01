package main

import (
	"fmt"
	"image"
	"log"

	_ "image/png"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/sprite"
	"golang.org/x/mobile/sprite/glsprite"
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
	scene      *sprite.Node
	eng        sprite.Engine
	texs       []sprite.SubTex
	blocks     map[*Block]*sprite.Node
	switches   []*sprite.Node
}

func compute(val float32, factor float32) float32 {
	return val * factor
}

func NewWorld() *World {

	// Clean
	// TODO
	w := &World{}

	w.eng = glsprite.Engine()
	w.loadTextures()
	w.scene = &sprite.Node{}
	w.eng.Register(w.scene)
	w.eng.SetTransform(w.scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	// Create the blocks
	w.blocks = make(map[*Block]*sprite.Node)
	for i := 0; i < len(g.level.blocks); i++ {
		for j := 0; j < len(g.level.blocks[i]); j++ {
			b := g.level.blocks[i][j]
			if b != nil {
				n := w.newNode(w.texs[0], blockSize, blockSize, 0, 0)
				w.blocks[b] = n
			}
		}
	}
	for _, sw := range g.level.switches {
		_ = sw
		w.switches = append(w.switches, w.newNode(w.texs[0], switchSize, switchSize, 256, 256))
	}
	return w
}

func (w *World) newNode(t sprite.SubTex, width, height, x, y float32) *sprite.Node {
	n := &sprite.Node{}
	w.eng.Register(n)
	w.scene.AppendChild(n)
	w.eng.SetSubTex(n, t)
	w.eng.SetTransform(n, f32.Affine{
		{width, 0, x},
		{0, height, y},
	})
	return n
}

const (
	texBlockRed = iota
	texBlockYellow
	texBlockBlue
	texBlockGreen
	texBlockPink
	texBlockOrange
	texBlockLightBlue
)

const (
	texBlockSize = 128
)

func (w *World) loadTextures() {
	a, err := app.Open("textures/tiles.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	img, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := w.eng.LoadTexture(img)
	if err != nil {
		log.Fatal(err)
	}

	w.texs = []sprite.SubTex{
		texBlockRed:       sprite.SubTex{t, image.Rect(0, 0, texBlockSize, texBlockSize)},
		texBlockYellow:    sprite.SubTex{t, image.Rect(texBlockSize, 0, texBlockSize*2, texBlockSize)},
		texBlockBlue:      sprite.SubTex{t, image.Rect(texBlockSize*2, 0, texBlockSize*3, texBlockSize)},
		texBlockGreen:     sprite.SubTex{t, image.Rect(texBlockSize*3, 0, texBlockSize*4, texBlockSize)},
		texBlockPink:      sprite.SubTex{t, image.Rect(0, texBlockSize, texBlockSize, texBlockSize*2)},
		texBlockOrange:    sprite.SubTex{t, image.Rect(texBlockSize, texBlockSize, texBlockSize*2, texBlockSize*3)},
		texBlockLightBlue: sprite.SubTex{t, image.Rect(texBlockSize*2, texBlockSize, texBlockSize*3, texBlockSize*2)},
	}
}

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
	g.world = NewWorld()
	// Load first level
	g.level = LoadLevel(g.currentLevel)
	g.world.background = NewBackground()
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
		//FIXME g.world.Reset()
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
