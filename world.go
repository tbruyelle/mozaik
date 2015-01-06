package main

import (
	"image"
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/sprite"
	"golang.org/x/mobile/sprite/clock"
	"golang.org/x/mobile/sprite/glsprite"
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
}

func compute(val float32, factor float32) float32 {
	return val * factor
}

// Switch implement the Arranger interface.
func (sw *Switch) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	switch sw.name {
	case "1":
		e.SetSubTex(n, g.world.texs[texSwitch1])
	case "2":
		e.SetSubTex(n, g.world.texs[texSwitch2])
	case "3":
		e.SetSubTex(n, g.world.texs[texSwitch3])
	case "4":
		e.SetSubTex(n, g.world.texs[texSwitch4])
	case "5":
		e.SetSubTex(n, g.world.texs[texSwitch5])
	case "6":
		e.SetSubTex(n, g.world.texs[texSwitch6])
	case "7":
		e.SetSubTex(n, g.world.texs[texSwitch7])
	case "8":
		e.SetSubTex(n, g.world.texs[texSwitch8])
	case "9":
		e.SetSubTex(n, g.world.texs[texSwitch9])
	}
	mv := &f32.Affine{}
	mv.Identity()
	mv.Translate(mv, sw.X, sw.Y)
	mv.Mul(mv, &f32.Affine{
		{switchSize, 0, 0},
		{0, switchSize, 0},
	})
	e.SetTransform(n, *mv)
}

// Block implement the Arranger interface.
func (b *Block) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	e.SetSubTex(n, g.world.texs[b.Color])
	mv := &f32.Affine{}
	mv.Identity()
	if b.rotateSW == nil {
		// The block is not attached to a rotating switch
		mv.Translate(mv, b.X, b.Y)
		mv.Mul(mv, &f32.Affine{
			{blockSize, 0, 0},
			{0, blockSize, 0},
		})

	} else {
		// The block is attached to a rotating switch
		// we need to draw it according to the switch
		// to apply the correct affine transformations
		v := switchSize / 2
		mv.Translate(mv, b.rotateSW.X+v, b.rotateSW.Y+v)
		mv.Rotate(mv, -b.rotateSW.rotate)
		mv.Scale(mv, b.rotateSW.scale, b.rotateSW.scale)
		tx := blockSize
		if b.X < b.rotateSW.X {
			tx = -tx
		}
		ty := blockSize
		if b.Y < b.rotateSW.Y {
			ty = -ty
		}
		mv.Mul(mv, &f32.Affine{
			{tx, 0, 0},
			{0, ty, 0},
		})
	}
	e.SetTransform(n, *mv)
}

func NewWorld() *World {

	// Clean
	// TODO
	w := &World{}

	w.eng = glsprite.Engine()
	w.loadTextures()
	w.scene = w.newNode()
	w.eng.SetTransform(w.scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	// Create the blocks
	for i := range g.level.blocks {
		for j := range g.level.blocks[i] {
			b := g.level.blocks[i][j]
			if b != nil {
				n := w.newNode()
				n.Arranger = b
				w.scene.AppendChild(n)
			}
		}
	}
	// Create the switches
	for _, sw := range g.level.switches {
		n := w.newNode()
		n.Arranger = sw
		w.scene.AppendChild(n)
	}

	// Add the win block signature
	signatureNode := w.newNode()
	w.scene.AppendChild(signatureNode)
	w.eng.SetTransform(signatureNode, f32.Affine{
		{1, 0, windowWidth - signatureBlockSize*4},
		{0, 1, windowHeight - signatureBlockSize*4},
	})
	line, col := float32(0), float32(0)
	for _, c := range g.level.winSignature {
		if c == '\n' {
			//next line
			line++
			col = 0
			continue
		}
		if c != '-' {
			n := w.newNode()
			signatureNode.AppendChild(n)
			w.eng.SetSubTex(n, w.texs[atoi(string(c))])
			w.eng.SetTransform(n, f32.Affine{
				{signatureBlockSize, 0, col * signatureBlockSize},
				{0, signatureBlockSize, line * signatureBlockSize},
			})
		}
		col++
	}

	return w
}

func (w *World) Draw() {
	// Background
	w.background.Draw()
	// The scene
	w.eng.Render(w.scene, 0)
}

func (w *World) newNode() *sprite.Node {
	n := &sprite.Node{}
	w.eng.Register(n)
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
		texSwitch1:        sprite.SubTex{t, image.Rect(0, TexBlockSize*2, TexSwitchSize-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch2:        sprite.SubTex{t, image.Rect(TexSwitchSize, TexBlockSize*2, TexSwitchSize*2-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch3:        sprite.SubTex{t, image.Rect(TexSwitchSize*2, TexBlockSize*2, TexSwitchSize*3-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch4:        sprite.SubTex{t, image.Rect(TexSwitchSize*3, TexBlockSize*2, TexSwitchSize*4-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch5:        sprite.SubTex{t, image.Rect(TexSwitchSize*4, TexBlockSize*2, TexSwitchSize*5-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch6:        sprite.SubTex{t, image.Rect(TexSwitchSize*5, TexBlockSize*2, TexSwitchSize*6-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch7:        sprite.SubTex{t, image.Rect(TexSwitchSize*6, TexBlockSize*2, TexSwitchSize*7-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch8:        sprite.SubTex{t, image.Rect(TexSwitchSize*7, TexBlockSize*2, TexSwitchSize*8-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch9:        sprite.SubTex{t, image.Rect(TexSwitchSize*8, TexBlockSize*2, TexSwitchSize*9-1, TexBlockSize*2+TexSwitchSize)},
	}
}
