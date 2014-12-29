package main

import (
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"math"
)

const (
	VShaderBasic = `#version 100

attribute vec4 position;
attribute vec4 color;

varying vec4 theColor;

uniform mat4 modelViewProjection;

void main() {
	gl_Position = modelViewProjection * position;
	theColor = color;
}`

	FShaderBasic = `#version 100

precision mediump float;

varying vec4 theColor;

void main() {
	gl_FragColor = theColor;
}`
)

type BlockModel struct {
	ModelGroup
	block *Block
}

func blockColor(b *Block) Color {
	switch b.Color {
	case Red:
		return RedColor
	case Blue:
		return BlueColor
	case LightBlue:
		return LightBlueColor
	case Orange:
		return OrangeColor
	case Green:
		return GreenColor
	case Pink:
		return PinkColor
	case Yellow:
		return YellowColor
	}
	return WhiteColor
}

func NewBlockModel(b *Block, size, radius float32) *BlockModel {
	model := &BlockModel{block: b}

	c := blockColor(b)
	s := size - radius
	// Inner square
	innervs := []Vertex{
		NewVertex(radius, radius, 0, c),
		NewVertex(radius, s, 0, c),
		NewVertex(s, radius, 0, c),
		NewVertex(s, s, 0, c),
	}
	model.Add(gl.TRIANGLE_STRIP, innervs, VShaderBasic, FShaderBasic)
	// Bottom square
	topvs := []Vertex{
		NewVertex(radius, s, 0, c),
		NewVertex(radius, size, 0, c),
		NewVertex(s, s, 0, c),
		NewVertex(s, size, 0, c),
	}
	model.Add(gl.TRIANGLE_STRIP, topvs, VShaderBasic, FShaderBasic)
	// Top square
	bottomvs := []Vertex{
		NewVertex(radius, radius, 0, c),
		NewVertex(radius, 0, 0, c),
		NewVertex(s, radius, 0, c),
		NewVertex(s, 0, 0, c),
	}
	model.Add(gl.TRIANGLE_STRIP, bottomvs, VShaderBasic, FShaderBasic)
	// Right square
	leftvs := []Vertex{
		NewVertex(s, radius, 0, c),
		NewVertex(size, radius, 0, c),
		NewVertex(s, s, 0, c),
		NewVertex(size, s, 0, c),
	}
	model.Add(gl.TRIANGLE_STRIP, leftvs, VShaderBasic, FShaderBasic)
	// Left square
	rightvs := []Vertex{
		NewVertex(radius, radius, 0, c),
		NewVertex(0, radius, 0, c),
		NewVertex(radius, s, 0, c),
		NewVertex(0, s, 0, c),
	}
	model.Add(gl.TRIANGLE_STRIP, rightvs, VShaderBasic, FShaderBasic)
	// Bottom right corner
	addCorner(model, c, s, s, radius, 0)
	// Bottom left corner
	addCorner(model, c, radius, s, radius, 1)
	// Top left corner
	addCorner(model, c, radius, radius, radius, 2)
	// Top right corner
	addCorner(model, c, s, radius, radius, 3)

	return model
}

func addCorner(model *BlockModel, c Color, x, y, radius, start float32) {

	max := float64(BlockCornerSegments * (start + 1))
	vs := []Vertex{NewVertex(x, y, 0, c)}
	for i := float64(start * BlockCornerSegments); i <= max; i++ {
		a := math.Pi / 2 * i / BlockCornerSegments
		xr := math.Cos(a) * float64(radius)
		yr := math.Sin(a) * float64(radius)
		vs = append(vs, NewVertex(x+float32(xr), y+float32(yr), 0, c))
	}
	model.Add(gl.TRIANGLE_FAN, vs, VShaderBasic, FShaderBasic)
}

type SwitchModel struct {
	ModelBase
	sw *Switch
}

func NewSwitchModel(sw *Switch) *SwitchModel {
	model := &SwitchModel{sw: sw}

	vs := []Vertex{NewVertex(0, 0, 0, WhiteColor)}
	vv := float64(switchSize / 2)
	for i := float64(0); i <= SwitchSegments; i++ {
		a := 2 * math.Pi * i / SwitchSegments
		vs = append(vs, NewVertex(float32(math.Sin(a)*vv), float32(math.Cos(a)*vv), 0, WhiteColor))
	}
	model.Init(gl.TRIANGLE_FAN, vs, VShaderBasic, FShaderBasic)

	v := switchSize / 2

	//model.modelView = gl.Ortho2D(0, WindowWidth, WindowHeight, 0).Mul(f32.Translate(float32(sw.X+v), float32(sw.Y+v), 0))
	model.modelView = ortho2D(0, windowWidth, windowHeight, 0)
	model.modelView.Mul(model.modelView, translate(float32(sw.X)+v, float32(sw.Y)+v, 0))
	return model
}

var (
	topLeftModelView     = translate(-blockSize, -blockSize, 0)
	topRightModelView    = translate(0, -blockSize, 0)
	bottomRightModelView = identity()
	bottomLeftModelView  = translate(-blockSize, 0, 0)
)

// TODO the switch number
func (t *SwitchModel) Draw() {
	modelViewBackup := *t.modelView
	s := t.sw
	var rotatemv *f32.Mat4
	if s.rotate == 0 {
		rotatemv = identity()
	} else {
		rotatemv = rotate(s.rotate)
	}
	var scalemv *f32.Mat4
	if s.scale == 0 {
		scalemv = identity()
	} else {
		scalemv = scale(s.scale)
	}
	blockmv := &f32.Mat4{}
	blockmv.Mul(scalemv, rotatemv)

	// Draw the associated blocks
	// top left block
	t.drawBlock(g.level.blocks[s.line][s.col], mul(blockmv, topLeftModelView))
	// top right block
	t.drawBlock(g.level.blocks[s.line][s.col+1], mul(blockmv, topRightModelView))
	// bottom right block
	t.drawBlock(g.level.blocks[s.line+1][s.col+1], mul(blockmv, bottomRightModelView))
	// bottom left block
	t.drawBlock(g.level.blocks[s.line+1][s.col], mul(blockmv, bottomLeftModelView))

	t.ModelBase.Draw()

	t.modelView = &modelViewBackup
}

func (t *SwitchModel) drawBlock(b *Block, modelView *f32.Mat4) {
	if !b.Rendered {
		b.Rendered = true
		bm := g.world.blocks[b]
		bm.modelView.Mul(t.modelView, modelView)
		bm.Draw()
	}
}

type Background struct {
	ModelBase
	angle float32
}

func NewBackground() *Background {
	model := &Background{}
	vs := []Vertex{}

	for i := float64(0); i <= BgSegments; i++ {
		if math.Mod(i, 2) == 0 {
			vs = append(vs, NewVertex(0, 0, 0, BgColor))
		}
		a := 2 * math.Pi * i / BgSegments
		vs = append(vs, NewVertex(float32(math.Sin(a)*windowRadius), float32(math.Cos(a)*windowRadius), 0, BgColor))
	}
	model.Init(gl.TRIANGLES, vs, VShaderBasic, FShaderBasic)
	return model
}

func (t *Background) Draw() {
	if t.angle > math.Pi {
		t.angle = t.angle - math.Pi
	} else {
		t.angle += 0.02
	}
	modelViewBackup := *t.modelView
	t.modelView.Mul(t.modelView, rotate(-t.angle))

	t.ModelBase.Draw()

	t.modelView = &modelViewBackup
}
