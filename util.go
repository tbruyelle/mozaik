package main

import (
	"golang.org/x/mobile/gl"
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

func NewProgram(shaders ...gl.Shader) gl.Program {
	prg := gl.CreateProgram()
	for _, shader := range shaders {
		prg.AttachShader(shader)
	}
	prg.Link()
	if prg.Get(gl.LINK_STATUS) != gl.TRUE {
		panic("linker error: " + prg.GetInfoLog())
	}
	prg.Validate()
	for _, shader := range shaders {
		prg.DetachShader(shader)
		shader.Delete()
	}
	return prg
}

func compileShader(type_ int, source string) gl.Shader {
	shader := gl.CreateShader(type_)
	shader.Source(source)
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("shader compile error for source " + source + "\n" + shader.GetInfoLog())
	}
	return shader
}

func loadShader(type_ int, file string) gl.Shader {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return compileShader(type_, string(b))
}
