package main

import (
	"golang.org/x/mobile/f32"
	"io/ioutil"
	"strconv"
	"strings"
	"unsafe"
)

type Vertex struct {
	Coords Coords
	Color  Color
}

func NewVertex(X, Y, Z float32, color Color) Vertex {
	return Vertex{Coords: Coords{X, Y, Z, 1.0}, Color: color}
}

var (
	WhiteColor     = Color{1, 1, 1, 1}
	RedColor       = Color{0.93, 0.05, 0.33, 1}
	GreenColor     = Color{0.34, 0.64, 0, 1}
	BlueColor      = Color{0.39, 0.58, 0.93, 1}
	YellowColor    = Color{1, 0.85, 0.23, 1}
	PinkColor      = Color{1, 0.70, 1, 1}
	OrangeColor    = Color{0.95, 0.48, 0.07, 1}
	LightBlueColor = Color{0.38, 0.87, 1, 1}
	BgColor        = Color{1.0, 0.85, 0.23, 1.0}
)

type Coords struct{ X, Y, Z, W float32 }
type Color struct{ R, G, B, A float32 }

func Sequence(seqSize, ind int) int {
	r := ind / seqSize
	for r >= seqSize {
		r -= seqSize
	}
	return r

}

func identity() *f32.Mat4 {
	id := &f32.Mat4{}
	id.Identity()
	return id
}

func translate(Tx, Ty, Tz float32) *f32.Mat4 {
	ret := &f32.Mat4{}
	ret.Translate(identity(), Tx, Ty, Tx)
	return ret
}

func rotate(angle float32) *f32.Mat4 {
	ret := &f32.Mat4{}
	ret.Rotate(identity(), f32.Radian(angle), &f32.Vec3{0, 0, 1})
	return ret
}

func scale(scale float32) *f32.Mat4 {
	ret := &f32.Mat4{}
	ret.Scale(identity(), scale, scale, 0)
	return ret
}

func mul(m1 *f32.Mat4, m2 *f32.Mat4) *f32.Mat4 {
	ret := &f32.Mat4{}
	ret.Mul(m1, m2)
	return ret
}

func readVertexFile(file string) []Vertex {
	vertexes := make([]Vertex, 0)
	b, err := ioutil.ReadFile(file + ".coords")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		coords := strings.Split(line, ",")
		if len(coords) >= 4 {
			v := Vertex{}
			v.Coords.X = atof(coords[0])
			v.Coords.Y = atof(coords[1])
			v.Coords.Z = atof(coords[2])
			v.Coords.W = atof(coords[3])
			vertexes = append(vertexes, v)
		}
	}
	b, err = ioutil.ReadFile(file + ".colors")
	if err != nil {
		panic(err)
	}
	vind := 0
	lines = strings.Split(string(b), "\n")
	for _, line := range lines {
		colors := strings.Split(line, ",")
		if len(colors) >= 4 {
			v := &vertexes[vind]
			v.Color.R = atof(colors[0])
			v.Color.G = atof(colors[1])
			v.Color.B = atof(colors[2])
			v.Color.A = atof(colors[3])
			vind++
		}
	}
	return vertexes
}

func atof(s string) float32 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 10)
	if err != nil {
		panic(err)
	}
	return float32(f)
}

var (
	sizeFloat  = int(unsafe.Sizeof(float32(0)))
	sizeCoords = sizeFloat * 4
	sizeVertex = int(unsafe.Sizeof(Vertex{}))
)