package main

import (
	"math"

	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/sprite"
	"golang.org/x/mobile/sprite/clock"
)

type Object struct {
	// Position
	X, Y float32
	// Speed
	Vx, Vy float32
	// Rotation
	Rx, Ry, Angle float32
	// Scale
	Sx, Sy        float32
	Width, Height float32
	Sprite        sprite.SubTex
	Action        func(o *Object)
	Dead          bool
	Tick          int
	// Data contains any relevant information needed about the object
	Data interface{}
}

func (o *Object) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	if o.Action != nil {
		// Invoke the action
		o.Action(o)
	}

	if o.Dead {
		// Do nothing if dead object
		return
	}

	// Set the texture
	e.SetSubTex(n, o.Sprite)

	// Compute affine transformations
	mv := &f32.Affine{}
	mv.Identity()

	if o.Angle == 0 {
		mv.Translate(mv, o.X, o.Y)
		mv.Mul(mv, &f32.Affine{
			{o.Width, 0, 0},
			{0, o.Height, 0},
		})
	} else {
		mv.Translate(mv, o.Rx, o.Ry)
		mv.Rotate(mv, -o.Angle)
		w := o.Width
		if o.X < o.Rx {
			w = -w
		}
		h := o.Height
		if o.Y < o.Ry {
			h = -h
		}
		mv.Mul(mv, &f32.Affine{
			{w, 0, 0},
			{0, h, 0},
		})
	}
	if o.Sx != 0 || o.Sy != 0 {
		mv.Scale(mv, o.Sx, o.Sy)
	}
	e.SetTransform(n, *mv)
}

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
