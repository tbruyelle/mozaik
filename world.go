package main

import (
	"fmt"
	"image"
	"log"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/gl"
)

type World struct {
	background  *Background
	moveCounter MoveCounter
	scene       *sprite.Node
	eng         sprite.Engine
	texs        []sprite.SubTex
}

func compute(val float32, factor float32) float32 {
	return val * factor
}

func NewWorld(glctx gl.Context) *World {

	// Clean
	// TODO
	w := &World{}

	w.background = NewBackground(glctx)

	w.eng = eng
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
			b.Layout(i, j, blockSize, blockPadding, xMin, yMin)
			n.Arranger = &b.Object
			w.scene.AppendChild(n)
		}
	}
	// Create the switches
	for _, sw := range g.level.switches {
		n := w.newNode()
		sw.Layout(switchSize)
		n.Arranger = &sw.Object
		w.scene.AppendChild(n)
	}

	// The bottom dashboard
	dashboard := w.newNode()
	//w.scene.AppendChild(dashboard)
	var dashboardX, dashboardY float32
	if portrait {
		dashboardX = 0
		dashboardY = windowHeight - dashboardSize
	} else {
		dashboardX = windowWidth - dashboardSize
		dashboardY = 0
	}
	w.eng.SetTransform(dashboard, f32.Affine{
		{1, 0, dashboardX},
		{0, 1, dashboardY},
	})

	// Add the win block signature
	signature := w.newNode()
	w.scene.AppendChild(signature)
	signSize := signatureBlockSize * 4
	w.eng.SetTransform(signature, f32.Affine{
		{1, 0, windowWidth - signSize - padding},
		{0, 1, windowHeight - signSize - padding},
	})
	line, col := 0, 0
	for i := range g.level.winSignature {
		for j := range g.level.winSignature[i] {
			c := g.level.winSignature[i][j]
			if c != Empty {
				n := w.newNode()
				signature.AppendChild(n)
				b := &Block{Color: c}
				b.Action = ActionFunc(signatureBlockIdle)
				b.Layout(line, col, signatureBlockSize, 0, 0, 0)
				n.Arranger = &b.Object
			}
			col++
		}
		line++
		col = 0
	}

	// The move counter
	var counterX, counterY float32
	if portrait {
		counterX = padding * 2
		counterY = windowHeight - dashboardSize + (dashboardSize-charHeight)/2
	} else {
		counterX = windowWidth - dashboardSize + (dashboardSize-charWidth*2)/2
		counterY = padding * 2
	}
	for i := 0; i < 2; i++ {
		n := w.newNode()
		w.scene.AppendChild(n)
		c := new(Char)
		w.moveCounter[i] = c
		n.Arranger = c
		c.X = counterX
		c.Y = counterY
		c.Width = charWidth
		c.Height = charHeight
		counterX += charWidth
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

	// The loose text node
	{
		n := w.newNode()
		w.scene.AppendChild(n)
		n.Arranger = &Object{
			X:      windowWidth/2 - gameoverTxtWidth/2,
			Y:      windowHeight/2 - gameoverTxtHeight/2,
			Width:  gameoverTxtWidth,
			Height: gameoverTxtHeight,
			Sprite: w.texs[texLooseTxt],
			Action: ActionFunc(looseTxtPop),
		}
	}
}

func (w *World) Draw(glctx gl.Context, t clock.Time, sz size.Event) {
	// Background
	w.background.Draw()
	// the move counter
	w.printMoves(g.level)
	// The scene
	w.eng.Render(w.scene, t, sz)
}

func (w *World) newNode() *sprite.Node {
	n := &sprite.Node{}
	w.eng.Register(n)
	return n
}

type MoveCounter [2]*Char

func (w *World) printMoves(l Level) {
	remain := l.maxMoves - l.moves
	moves := fmt.Sprintf("%d", remain)
	if remain < 10 {
		if !w.moveCounter[0].Dead {
			// Center the unique number only once.
			if portrait {
				w.moveCounter[1].X -= charWidth / 2
			} else {
				w.moveCounter[1].Y -= charHeight / 2
			}
		}
		w.moveCounter[0].Sprite = w.texs[texEmpty]
		w.moveCounter[0].Dead = true
		w.moveCounter[1].Set(w, rune(moves[0]))
		return
	}
	if w.moveCounter[0].Dead {
		w.moveCounter[0].Dead = false
		// Center the unique number only once.
		if portrait {
			w.moveCounter[1].X += charWidth / 2
		} else {
			w.moveCounter[1].Y += charHeight / 2
		}
	}
	w.moveCounter[0].Set(w, rune(moves[0]))
	w.moveCounter[1].Set(w, rune(moves[1]))
}

func (w *World) decMoves() {
}

type Char struct {
	Object
	val string
}

func (c *Char) Set(w *World, val rune) {
	// convert the rune to int
	c.Sprite = w.texs[tex0+val-48]
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
	texLooseTxt
	texEmpty
)

const (
	TexBlockSize      = 128
	TexSwitchSize     = 50
	TexCharWidth      = 38
	TexCharHeight     = 60
	TexWinWidth       = 470
	TexWinHeight      = 106
	TexGameoverWidth  = 560
	TexGameoverHeight = 106
)

func (w *World) loadTextures() {
	a, err := asset.Open("textures/tiles.png")
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
		texEmpty: {},
		// Block textures
		texBlockRed:         {t, image.Rect(0, 0, TexBlockSize, TexBlockSize)},
		texBlockYellow:      {t, image.Rect(TexBlockSize, 0, TexBlockSize*2, TexBlockSize)},
		texBlockBlue:        {t, image.Rect(TexBlockSize*2, 0, TexBlockSize*3, TexBlockSize)},
		texBlockGreen:       {t, image.Rect(TexBlockSize*3, 0, TexBlockSize*4, TexBlockSize)},
		texBlockBrown:       {t, image.Rect(TexBlockSize*4, 0, TexBlockSize*5, TexBlockSize)},
		texBlockLightGreen:  {t, image.Rect(TexBlockSize*5, 0, TexBlockSize*6, TexBlockSize)},
		texBlockCyan:        {t, image.Rect(TexBlockSize*6, 0, TexBlockSize*7, TexBlockSize)},
		texBlockLightPink:   {t, image.Rect(TexBlockSize*7, 0, TexBlockSize*8, TexBlockSize)},
		texBlockPink:        {t, image.Rect(0, TexBlockSize, TexBlockSize, TexBlockSize*2)},
		texBlockOrange:      {t, image.Rect(TexBlockSize, TexBlockSize, TexBlockSize*2, TexBlockSize*2)},
		texBlockLightBlue:   {t, image.Rect(TexBlockSize*2, TexBlockSize, TexBlockSize*3, TexBlockSize*2)},
		texBlockPurple:      {t, image.Rect(TexBlockSize*3, TexBlockSize, TexBlockSize*4, TexBlockSize*2)},
		texBlockWhite:       {t, image.Rect(TexBlockSize*4, TexBlockSize, TexBlockSize*5, TexBlockSize*2)},
		texBlockLightPurple: {t, image.Rect(TexBlockSize*5, TexBlockSize, TexBlockSize*6, TexBlockSize*2)},
		texBlockLightBrown:  {t, image.Rect(TexBlockSize*6, TexBlockSize, TexBlockSize*7, TexBlockSize*2)},
		texBlockOtherWhite:  {t, image.Rect(TexBlockSize*7, TexBlockSize, TexBlockSize*8, TexBlockSize*2)},
		// Switches textures
		texSwitch1: {t, image.Rect(0, TexBlockSize*2, TexSwitchSize-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch2: {t, image.Rect(TexSwitchSize, TexBlockSize*2, TexSwitchSize*2-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch3: {t, image.Rect(TexSwitchSize*2, TexBlockSize*2, TexSwitchSize*3-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch4: {t, image.Rect(TexSwitchSize*3, TexBlockSize*2, TexSwitchSize*4-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch5: {t, image.Rect(TexSwitchSize*4, TexBlockSize*2, TexSwitchSize*5-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch6: {t, image.Rect(TexSwitchSize*5, TexBlockSize*2, TexSwitchSize*6-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch7: {t, image.Rect(TexSwitchSize*6, TexBlockSize*2, TexSwitchSize*7-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch8: {t, image.Rect(TexSwitchSize*7, TexBlockSize*2, TexSwitchSize*8-1, TexBlockSize*2+TexSwitchSize)},
		texSwitch9: {t, image.Rect(TexSwitchSize*8, TexBlockSize*2, TexSwitchSize*9-1, TexBlockSize*2+TexSwitchSize)},
		// Win text texture
		texWinTxt: {t, image.Rect(0, TexBlockSize*2+TexSwitchSize, TexWinWidth, TexBlockSize*2+TexSwitchSize+TexWinHeight)},
		// Gameover text texture
		texLooseTxt: {t, image.Rect(0, TexBlockSize*2+TexSwitchSize+TexWinHeight, TexGameoverWidth, TexBlockSize*2+TexSwitchSize+TexWinHeight+TexGameoverHeight)},
	}

	// Load the number textures
	numStartX := 0
	numStartY := TexBlockSize*2 + TexSwitchSize + TexWinHeight + TexGameoverHeight
	numEndY := numStartY + TexCharHeight

	texId := tex0
	for i := 0; i < 10; i++ {
		w.texs[texId] = sprite.SubTex{t, image.Rect(numStartX, numStartY, numStartX+TexCharWidth, numEndY)}
		numStartX += TexCharWidth
		texId++
	}
}
