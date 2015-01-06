package main

import (
	"math"

	"golang.org/x/mobile/gl"
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
	ModelBase
	block *Block
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
	t.modelView.Mul(t.modelView, rotate(t.angle))

	t.ModelBase.Draw()

	t.modelView = &modelViewBackup
}
