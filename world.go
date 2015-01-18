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

type World struct {
	background *Background
	scene      *sprite.Node
	eng        sprite.Engine
	texs       []sprite.SubTex
}

func compute(val float32, factor float32) float32 {
	return val * factor
}

func NewWorld() *World {

	// Clean
	// TODO
	w := &World{}

	w.background = NewBackground()

	w.eng = glsprite.Engine()
	w.loadTextures()
	w.LoadScene()
	return w
}

func (w *World) LoadScene() {
	w.scene = w.newNode()
	w.eng.SetTransform(w.scene, f32.Affine{
		{1, 0, 0},
		{0, 1, 0},
	})

	// Create the blocks
	for i := range g.level.blocks {
		for j := range g.level.blocks[i] {
			b := g.level.blocks[i][j]
			n := w.newNode()
			n.Arranger = &b.Object
			w.scene.AppendChild(n)
		}
	}
	// Create the switches
	for _, sw := range g.level.switches {
		n := w.newNode()
		n.Arranger = &sw.Object
		w.scene.AppendChild(n)
	}

	// The bottom dashboard
	dashboard := w.newNode()
	//w.scene.AppendChild(dashboard)
	w.eng.SetTransform(dashboard, f32.Affine{
		{1, 0, 0},
		{0, 1, windowHeight - dashboardHeight},
	})

	// Add the win block signature
	signature := w.newNode()
	w.scene.AppendChild(signature)
	signSize := signatureBlockSize * 4
	w.eng.SetTransform(signature, f32.Affine{
		{1, 0, windowWidth - signSize},
		{0, 1, windowHeight - signSize},
	})
	line, col := 0, 0
	for i := range g.level.winSignature {
		for j := range g.level.winSignature[i] {
			c := g.level.winSignature[i][j]
			if c != Empty {
				n := w.newNode()
				signature.AppendChild(n)
				b := newBlock(c, line, col, signatureBlockSize, 0)
				b.Action = ActionFunc(signatureBlockIdle)
				n.Arranger = &b.Object
			}
			col++
		}
		line++
		col = 0
	}

	// The move counter
	{
		n := w.newNode()
		w.scene.AppendChild(n)
		n.Arranger = &Object{
			X:      windowWidth / 2,
			Y:      windowHeight - 60,
			Width:  40,
			Height: 54,
			Sprite: w.texs[tex2],
		}
	}

	// Add the win text node
	{
		n := w.newNode()
		w.scene.AppendChild(n)
		n.Arranger = &Object{
			X:      windowWidth/2 - winTxtWidth/2,
			Y:      windowHeight/2 - winTxtHeight/2,
			Width:  winTxtWidth,
			Height: winTxtHeight,
			Sprite: w.texs[texWinTxt],
			Action: ActionFunc(winTxtPop),
		}
	}
}

func (w *World) Draw(t clock.Time) {
	// Background
	w.background.Draw()
	// The scene
	w.eng.Render(w.scene, t)
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
	texBlockPurple
	texBlockBrown
	texBlockLightGreen
	texBlockCyan
	texBlockLightPink
	texBlockWhite
	texBlockLightPurple
	texBlockLightBrown
	texBlockOtherWhite
	texSwitch1
	texSwitch2
	texSwitch3
	texSwitch4
	texSwitch5
	texSwitch6
	texSwitch7
	texSwitch8
	texSwitch9
	texWinTxt
	tex0
	tex1
	tex2
	tex3
	tex4
	tex5
	tex6
	tex7
	tex8
	tex9
	texEmpty
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
		// Empty texture
		texEmpty: sprite.SubTex{},
		// Block textures
		texBlockRed:         sprite.SubTex{t, image.Rect(0, 0, TexBlockSize, TexBlockSize)},
		texBlockYellow:      sprite.SubTex{t, image.Rect(TexBlockSize, 0, TexBlockSize*2, TexBlockSize)},
		texBlockBlue:        sprite.SubTex{t, image.Rect(TexBlockSize*2, 0, TexBlockSize*3, TexBlockSize)},
		texBlockGreen:       sprite.SubTex{t, image.Rect(TexBlockSize*3, 0, TexBlockSize*4, TexBlockSize)},
		texBlockBrown:       sprite.SubTex{t, image.Rect(TexBlockSize*4, 0, TexBlockSize*5, TexBlockSize)},
		texBlockLightGreen:  sprite.SubTex{t, image.Rect(TexBlockSize*5, 0, TexBlockSize*6, TexBlockSize)},
		texBlockCyan:        sprite.SubTex{t, image.Rect(TexBlockSize*6, 0, TexBlockSize*7, TexBlockSize)},
		texBlockLightPink:   sprite.SubTex{t, image.Rect(TexBlockSize*7, 0, TexBlockSize*8, TexBlockSize)},
		texBlockPink:        sprite.SubTex{t, image.Rect(0, TexBlockSize, TexBlockSize, TexBlockSize*2)},
		texBlockOrange:      sprite.SubTex{t, image.Rect(TexBlockSize, TexBlockSize, TexBlockSize*2, TexBlockSize*2)},
		texBlockLightBlue:   sprite.SubTex{t, image.Rect(TexBlockSize*2, TexBlockSize, TexBlockSize*3, TexBlockSize*2)},
		texBlockPurple:      sprite.SubTex{t, image.Rect(TexBlockSize*3, TexBlockSize, TexBlockSize*4, TexBlockSize*2)},
		texBlockWhite:       sprite.SubTex{t, image.Rect(TexBlockSize*4, TexBlockSize, TexBlockSize*5, TexBlockSize*2)},
		texBlockLightPurple: sprite.SubTex{t, image.Rect(TexBlockSize*5, TexBlockSize, TexBlockSize*6, TexBlockSize*2)},
		texBlockLightBrown:  sprite.SubTex{t, image.Rect(TexBlockSize*6, TexBlockSize, TexBlockSize*7, TexBlockSize*2)},
		texBlockOtherWhite:  sprite.SubTex{t, image.Rect(TexBlockSize*7, TexBlockSize, TexBlockSize*8, TexBlockSize*2)},
		// Switches textures
		texSwitch1: sprite.SubTex{t, image.Rect(0, TexBlockSize*2, TexSwitchSize-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch2: sprite.SubTex{t, image.Rect(TexSwitchSize, TexBlockSize*2, TexSwitchSize*2-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch3: sprite.SubTex{t, image.Rect(TexSwitchSize*2, TexBlockSize*2, TexSwitchSize*3-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch4: sprite.SubTex{t, image.Rect(TexSwitchSize*3, TexBlockSize*2, TexSwitchSize*4-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch5: sprite.SubTex{t, image.Rect(TexSwitchSize*4, TexBlockSize*2, TexSwitchSize*5-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch6: sprite.SubTex{t, image.Rect(TexSwitchSize*5, TexBlockSize*2, TexSwitchSize*6-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch7: sprite.SubTex{t, image.Rect(TexSwitchSize*6, TexBlockSize*2, TexSwitchSize*7-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch8: sprite.SubTex{t, image.Rect(TexSwitchSize*7, TexBlockSize*2, TexSwitchSize*8-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch9: sprite.SubTex{t, image.Rect(TexSwitchSize*8, TexBlockSize*2, TexSwitchSize*9-1, TexBlockSize*2+TexSwitchSize)},
		// Win text texture
		texWinTxt: sprite.SubTex{t, image.Rect(0, TexBlockSize*2+TexSwitchSize, 300, TexBlockSize*2+TexSwitchSize+90)},
	}

	// Load the number textures
	numStartX := 320
	numStartY := 320
	numEndY := 374

	texId := tex0
	for i := 0; i < 10; i++ {
		w.texs[texId] = sprite.SubTex{t, image.Rect(numStartX, numStartY, numStartX+40, numEndY)}
		numStartX += 40
		texId++
	}
}
