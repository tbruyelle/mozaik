package main

import (
	"encoding/binary"
	"math"

	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/gl"
)

type Object struct {
	// Position
	X, Y float32
	// Speed
	Vx, Vy float32
	// Rotation
	Rx, Ry, Angle, AngleCenter float32
	// Translation
	Tx, Ty float32
	// Scale
	Sx, Sy, Scale float32
	Width, Height float32
	Sprite        sprite.SubTex
	Action        Action
	Dead          bool
	Time          clock.Time
	// Data contains any relevant information needed about the object
	Data interface{}
	// Object center
	cx, cy float32
}

type Action interface {
	Do(o *Object, t clock.Time)
}

func (o *Object) Reset() {
	o.Tx, o.Ty = 0, 0
	o.Sx, o.Sy, o.Scale = 0, 0, 0
	o.Rx, o.Ry, o.Angle, o.AngleCenter = 0, 0, 0, 0
	o.Time = 0
}

func (o *Object) Center() (float32, float32) {
	if o.cx == 0 {
		o.cx = o.X + o.Width/2
	}
	if o.cy == 0 {
		o.cy = o.Y + o.Height/2
	}
	return o.cx, o.cy
}

func (o *Object) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) {
	if o.Action != nil {
		// Invoke the action
		o.Action.Do(o, t)
	}

	// Set the texture
	e.SetSubTex(n, o.Sprite)

	if o.Dead {
		// Do nothing if dead object
		return
	}

	// Compute affine transformations
	mv := &f32.Affine{}
	mv.Identity()

	// Apply translations
	if o.Angle != 0 && o.Rx == o.Sx && o.Ry == o.Sy {
		// Optim when angle and scale use the same transformation
		mv.Translate(mv, o.Rx+o.Tx, o.Ry+o.Ty)
		mv.Rotate(mv, -o.Angle)
		mv.Scale(mv, o.Scale, o.Scale)
		mv.Translate(mv, -o.Rx-o.Tx, -o.Ry-o.Ty)
	} else {
		if o.Angle != 0 {
			mv.Translate(mv, o.Rx+o.Tx, o.Ry+o.Ty)
			mv.Rotate(mv, -o.Angle)
			mv.Translate(mv, -o.Rx-o.Tx, -o.Ry-o.Ty)
		}
		if o.Sx > 0 || o.Sy > 0 {
			mv.Translate(mv, o.Sx+o.Tx, o.Sy+o.Ty)
			mv.Scale(mv, o.Scale, o.Scale)
			mv.Translate(mv, -o.Sx-o.Tx, -o.Sy-o.Ty)
		}
	}
	if o.AngleCenter != 0 {
		cx, cy := o.Center()
		mv.Translate(mv, cx+o.Tx, cy+o.Ty)
		mv.Rotate(mv, -o.AngleCenter)
		mv.Translate(mv, -cx-o.Tx, -cy-o.Ty)
	}
	mv.Translate(mv, o.X+o.Tx, o.Y+o.Ty)
	mv.Mul(mv, &f32.Affine{
		{o.Width, 0, 0},
		{0, o.Height, 0},
	})
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

func NewBackground(glctx gl.Context) *Background {
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

	b.Init(glctx, gl.TRIANGLES, data, VShaderBasic, FShaderBasic)
	return b
}

func (t *Background) Draw() {
	if t.angle > math.Pi {
		t.angle = t.angle - math.Pi
	} else {
		if g.level.Win() {
			t.angle += 0.03
		} else {
			t.angle += 0.01
		}
	}
	modelViewBackup := *t.modelView
	t.modelView.Mul(t.modelView, rotate(t.angle))

	t.ModelBase.Draw()

	t.modelView = &modelViewBackup
}
