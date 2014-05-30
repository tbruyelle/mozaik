package main

import (
	"github.com/go-gl/gl"
	"github.com/remogatto/mathgl"
	"math"
)

const (
	VShaderBasic = `#version 330 core

	layout(location=0) in vec4 position;
	layout(location=1) in vec4 color;

	smooth out vec4 theColor;

	uniform mat4 modelViewProjection;

	void main() {
		theColor = color;
		gl_Position = modelViewProjection * position;
	}`

	FShaderBasic = `#version 330 core

	smooth in vec4 theColor;

	out vec4 outputColor;

	void main() {
		outputColor = theColor;
	}`
)

type BlockModel struct {
	ModelGroup
	block *Block
}

func getBlockColor(b *Block) Color {
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

	c := getBlockColor(b)
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
	vv := float64(SwitchSize / 2)
	for i := float64(0); i <= SwitchSegments; i++ {
		a := 2 * math.Pi * i / SwitchSegments
		vs = append(vs, NewVertex(float32(math.Sin(a)*vv), float32(math.Cos(a)*vv), 0, WhiteColor))
	}
	model.Init(gl.TRIANGLE_FAN, vs, VShaderBasic, FShaderBasic)

	v := SwitchSize / 2
	model.modelView = mathgl.Ortho2D(0, WindowWidth, WindowHeight, 0).Mul4(mathgl.Translate3D(float32(sw.X+v), float32(sw.Y+v), 0))
	return model
}

var (
	topLeftModelView     = mathgl.Translate3D(-BlockSize, -BlockSize, 0)
	topRightModelView    = mathgl.Translate3D(0, -BlockSize, 0)
	bottomRightModelView = mathgl.Ident4f()
	bottomLeftModelView  = mathgl.Translate3D(-BlockSize, 0, 0)
)

// TODO the switch number
func (t *SwitchModel) Draw() {
	modelViewBackup := t.modelView
	s := t.sw
	var rotatemv mathgl.Mat4f
	if s.rotate == 0 {
		rotatemv = mathgl.Ident4f()
	} else {
		rotatemv = mathgl.HomogRotate3D(t.sw.rotate, [3]float32{0, 0, 1})
	}
	var scalemv mathgl.Mat4f
	if s.scale == 0 {
		scalemv = mathgl.Ident4f()
	} else {
		scalemv = mathgl.Scale3D(s.scale, s.scale, 0)
	}
	blockmv := scalemv.Mul4(rotatemv)

	// Draw the associated blocks
	// top left block
	t.drawBlock(g.level.blocks[s.line][s.col], blockmv.Mul4(topLeftModelView))
	// top right block
	t.drawBlock(g.level.blocks[s.line][s.col+1], blockmv.Mul4(topRightModelView))
	// bottom right block
	t.drawBlock(g.level.blocks[s.line+1][s.col+1], blockmv.Mul4(bottomRightModelView))
	// bottom left block
	t.drawBlock(g.level.blocks[s.line+1][s.col], blockmv.Mul4(bottomLeftModelView))

	t.ModelBase.Draw()

	t.modelView = modelViewBackup
}

func (t *SwitchModel) drawBlock(b *Block, modelView mathgl.Mat4f) {
	if !b.Rendered {
		b.Rendered = true
		bm := g.world.blocks[b]
		bm.modelView = t.modelView.Mul4(modelView)
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
		t.angle += 0.03
	}
	modelViewBackup := t.modelView
	t.modelView = t.modelView.Mul4(mathgl.HomogRotate3D(-t.angle, [3]float32{0, 0, 1}))

	t.ModelBase.Draw()

	t.modelView = modelViewBackup
}
