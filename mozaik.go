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
				n := w.newNode(int(b.Color), blockSize, blockSize, float32(j)*blockSize, float32(i)*blockSize)
				w.blocks[b] = n
			}
		}
	}
	v := switchSize / 2
	for _, sw := range g.level.switches {
		_ = sw
		w.switches = append(w.switches, w.newNode(texSwitch1, switchSize, switchSize, sw.X+v, sw.Y+v))
	}
	return w
}

func (w *World) newNode(tex int, width, height, x, y float32) *sprite.Node {
	n := &sprite.Node{}
	w.eng.Register(n)
	w.scene.AppendChild(n)
	w.eng.SetSubTex(n, w.texs[tex])
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
	texSwitch1
	texSwitch2
	texSwitch3
	texSwitch4
	texSwitch5
	texSwitch6
	texSwitch7
	texSwitch8
	texSwitch9
)

const (
	TexBlockSize  = 128
	TexSwitchSize = 48
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
		texBlockRed:       sprite.SubTex{t, image.Rect(0, 0, TexBlockSize, TexBlockSize)},
		texBlockYellow:    sprite.SubTex{t, image.Rect(TexBlockSize, 0, TexBlockSize*2, TexBlockSize)},
		texBlockBlue:      sprite.SubTex{t, image.Rect(TexBlockSize*2, 0, TexBlockSize*3, TexBlockSize)},
		texBlockGreen:     sprite.SubTex{t, image.Rect(TexBlockSize*3, 0, TexBlockSize*4, TexBlockSize)},
		texBlockPink:      sprite.SubTex{t, image.Rect(0, TexBlockSize, TexBlockSize, TexBlockSize*2)},
		texBlockOrange:    sprite.SubTex{t, image.Rect(TexBlockSize, TexBlockSize, TexBlockSize*2, TexBlockSize*2)},
		texBlockLightBlue: sprite.SubTex{t, image.Rect(TexBlockSize*2, TexBlockSize, TexBlockSize*3, TexBlockSize*2)},
		texSwitch1:        sprite.SubTex{t, image.Rect(0, TexBlockSize*2, TexSwitchSize, TexBlockSize*2+TexSwitchSize)},
		texSwitch2:        sprite.SubTex{t, image.Rect(0, TexBlockSize*2, TexSwitchSize, TexBlockSize*2+TexSwitchSize)},
		texSwitch3:        sprite.SubTex{t, image.Rect(TexSwitchSize*2, TexBlockSize*2, TexSwitchSize*3, TexBlockSize*2+TexSwitchSize)},
		texSwitch4:        sprite.SubTex{t, image.Rect(TexSwitchSize*3, TexBlockSize*2, TexSwitchSize*4, TexBlockSize*2+TexSwitchSize)},
		texSwitch5:        sprite.SubTex{t, image.Rect(TexSwitchSize*4, TexBlockSize*2, TexSwitchSize*5, TexBlockSize*2+TexSwitchSize)},
		texSwitch6:        sprite.SubTex{t, image.Rect(TexSwitchSize*5, TexBlockSize*2, TexSwitchSize*6, TexBlockSize*2+TexSwitchSize)},
		texSwitch7:        sprite.SubTex{t, image.Rect(TexSwitchSize*6, TexBlockSize*2, TexSwitchSize*7, TexBlockSize*2+TexSwitchSize)},
		texSwitch8:        sprite.SubTex{t, image.Rect(TexSwitchSize*7, TexBlockSize*2, TexSwitchSize*8, TexBlockSize*2+TexSwitchSize)},
		texSwitch9:        sprite.SubTex{t, image.Rect(TexSwitchSize*8, TexBlockSize*2, TexSwitchSize*9, TexBlockSize*2+TexSwitchSize)},
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
	g.level = LoadLevel(g.currentLevel)
	g.world = NewWorld()
	// Load first level
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
