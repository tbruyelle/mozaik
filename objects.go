package main

import (
	"encoding/binary"
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
	// Translation
	Tx, Ty float32
	// Scale
	Sx, Sy        float32
	Width, Height float32
	Sprite        sprite.SubTex
	Action        Action
	Dead          bool
	Time          clock.Time
	// Data contains any relevant information needed about the object
	Data interface{}
}

type Action interface {
	Do(o *Object, t clock.Time)
}

func (o *Object) Reset() {
	o.Tx, o.Ty, o.Sx, o.Sy, o.Rx, o.Ry, o.Angle = 0, 0, 0, 0, 0, 0, 0
	o.Time = 0
}

func (o *Object) ZoomIn(f, start float32) {
	s := start + f
	o.Sx, o.Sy = s, s
	if start < 1 {
		o.Tx = o.Width / 2 * (1 - f)
		o.Ty = o.Height / 2 * (1 - f)
	} else {
		o.Tx = -o.Width / 2 * f
		o.Ty = -o.Height / 2 * f
	}
}

func (o *Object) ZoomOut(f, start float32) {
	s := start - f
	o.Sx, o.Sy = s, s
	mw, mh := o.Width/2, o.Height/2
	o.Tx = -mw*(start-1) + mw*f
	o.Ty = -mh*(start-1) + mh*f
}

func (o *Object) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	if o.Action != nil {
		// Invoke the action
		o.Action.Do(o, t)
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

	// Apply translations
	x, y := o.X+o.Tx, o.Y+o.Ty

	if o.Angle == 0 {
		mv.Translate(mv, x, y)
		mv.Mul(mv, &f32.Affine{
			{o.Width, 0, 0},
			{0, o.Height, 0},
		})
	} else {
		mv.Translate(mv, o.Rx, o.Ry)
		mv.Rotate(mv, -o.Angle)
		w := o.Width
		if x < o.Rx {
			w = -w
		}
		h := o.Height
		if y < o.Ry {
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
	b := &Background{}

	vertices := make([]float32, 0)
	for i := float64(0); i <= BgSegments; i++ {
		if math.Mod(i, 2) == 0 {
			// position
			vertices = append(vertices, 0, 0, 0, 1)
			// color
			vertices = append(vertices, .11, .03, .81, 1)
			b.vertexCount++
		}
		a := 2 * math.Pi * i / BgSegments
		// position
		vertices = append(vertices, float32(math.Sin(a)*windowRadius), float32(math.Cos(a)*windowRadius), 0, 1)
		// color
		vertices = append(vertices, .11, .03, .81, 1)
		b.vertexCount++
	}
	data := f32.Bytes(binary.LittleEndian, vertices...)

	b.Init(gl.TRIANGLES, data, VShaderBasic, FShaderBasic)
	return b
}

func (t *Background) Draw() {
	if t.angle > math.Pi {
		t.angle = t.angle - math.Pi
	} else {
		t.angle += 0.01
	}
	modelViewBackup := *t.modelView
	t.modelView.Mul(t.modelView, rotate(t.angle))

	t.ModelBase.Draw()

	t.modelView = &modelViewBackup
}
