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
	blockSet   []*sprite.Node
}

func compute(val float32, factor float32) float32 {
	return val * factor
}

// Switch implement the Arranger interface
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

type BlockSetArranger struct {
	sw *Switch
}

func (bsa *BlockSetArranger) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	v := switchSize / 2
	mv := &f32.Affine{}
	mv.Identity()
	mv.Translate(mv, bsa.sw.X+v, bsa.sw.Y+v)
	if bsa.sw.rotate != 0 {
		mv.Rotate(mv, -bsa.sw.rotate)
	}
	e.SetTransform(n, *mv)
}

type BlockArranger struct {
	sw *Switch
	// the block position according to the switch
	x, y int
	// translation according to the switch
	tx, ty float32
}

func (ba *BlockArranger) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	// find the corresponding block
	b := g.level.blocks[ba.x][ba.y]
	if b.Rendered {
		e.SetTransform(n, f32.Affine{
			{0, 0, 0},
			{0, 0, 0},
		})
	} else {
		b.Rendered = true //FIXME put Rendered in the arranger
		e.SetSubTex(n, g.world.texs[b.Color])
		mv := &f32.Affine{}
		mv.Identity()
		mv.Translate(mv, ba.tx, ba.ty)

		mv.Mul(mv, &f32.Affine{
			{blockSize, 0, 0},
			{0, blockSize, 0},
		})

		e.SetTransform(n, *mv)
	}
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

	// Create the switches
	for _, sw := range g.level.switches {
		n := w.newNode()
		n.Arranger = sw
		w.scene.AppendChild(n)
		// for each switch add the corresponding block set
		bs := w.newNode()
		bs.Arranger = &BlockSetArranger{sw: sw}
		// blockSet are not added to the scene node because they'll be
		// rendered in a specific order, according to the presence of
		// a rotating switch.
		w.blockSet = append(w.blockSet, bs)
		// for each block set add the corresponding blocks
		w.addBlock(bs, sw, sw.line, sw.col, -blockSize, -blockSize)
		w.addBlock(bs, sw, sw.line, sw.col+1, 0, -blockSize)
		w.addBlock(bs, sw, sw.line+1, sw.col+1, 0, 0)
		w.addBlock(bs, sw, sw.line+1, sw.col, -blockSize, 0)
	}

	// Add the win block signature
	signatureNode := w.newNode()
	w.scene.AppendChild(signatureNode)
	w.eng.SetTransform(signatureNode, f32.Affine{
		{1, 0, xMax - signatureBlockSize*4},
		{0, 1, yMax - signatureBlockSize*4},
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
			col++
		}
	}

	return w
}

func (w *World) addBlock(parent *sprite.Node, sw *Switch, x, y int, tx, ty float32) {
	b := w.newNode()
	b.Arranger = &BlockArranger{sw: sw, x: x, y: y, tx: tx, ty: ty}
	parent.AppendChild(b)
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
