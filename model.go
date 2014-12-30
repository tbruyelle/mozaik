package main

import (
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

type Model interface {
	Draw()
	Destroy()
	//pushModelView(modelView mathf32.Mat4)
	//popModelView()
}

type ModelBase struct {
	mode         gl.Enum
	buf          gl.Buffer
	vertices     []Vertex
	sizeVertices int
	vertexCount  int
	prg          gl.Program
	position     gl.Attrib
	color        gl.Attrib
	//	vao                   gl.VertexArray
	uniformMVP            gl.Uniform
	modelView, projection *f32.Mat4
	modelViewBackup       *f32.Mat4
}

type ModelGroup struct {
	models                []*ModelBase
	modelView, projection *f32.Mat4
}

func (t *ModelGroup) Add(mode gl.Enum, vertices []Vertex, vshaderf, fshaderf string) {
	m := &ModelBase{}
	m.Init(mode, vertices, vshaderf, fshaderf)
	t.models = append(t.models, m)
}

func (t *ModelBase) Init(mode gl.Enum, vertices []Vertex, vshaderf, fshaderf string) {
	t.mode = mode
	t.vertices = vertices
	t.vertexCount = len(t.vertices)
	t.sizeVertices = t.vertexCount * sizeVertex

	// Shaders
	var err error
	t.prg, err = glutil.CreateProgram(vshaderf, fshaderf)
	if err != nil {
		panic(err)
	}

	t.buf = gl.GenBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, t.buf)
	gl.BufferData2(gl.ARRAY_BUFFER, gl.STATIC_DRAW, t.sizeVertices, vertices)

	t.position = gl.GetAttribLocation(t.prg, "position")
	t.color = gl.GetAttribLocation(t.prg, "color")
	t.uniformMVP = gl.GetUniformLocation(t.prg, "modelViewProjection")

	// the projection matrix
	t.projection = &f32.Mat4{}
	t.projection.Identity()

	// the model view
	t.modelView = &f32.Mat4{}
	t.modelView.Identity()

	// Create VBO
	//	t.buffer = gl.GenBuffer()
	//	t.buffer.Bind(gl.ARRAY_BUFFER)
	//	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	//	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)
	//	t.buffer.Unbind(gl.ARRAY_BUFFER)
	//
	//	// Create VAO
	//	//t.vao = gl.GenVertexArray()
	//	//t.vao.Bind()
	//	t.buffer.Bind(gl.ARRAY_BUFFER)
	//
	//	// Attrib vertex data to VAO
	//	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	//	t.posLoc.EnableArray()
	//	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))
	//	t.colLoc.EnableArray()
	//
	//	t.buffer.Unbind(gl.ARRAY_BUFFER)
	//	//t.vao.Unbind()
}

// flatten returns column based flatten matrix.
func flatten(m *f32.Mat4) []float32 {
	f := make([]float32, 16)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			f[i*4+j] = m[i][j]
		}
	}
	return f
}

func (t *ModelBase) Draw() {
	gl.UseProgram(t.prg)

	mvp := &f32.Mat4{}
	mvp.Mul(t.modelView, t.projection)
	gl.UniformMatrix4fv(t.uniformMVP, flatten(mvp))
	//t.uniformMVP.UniformMatrix4f(false, (*[16]float32)(&mvp))

	gl.BindBuffer(gl.ARRAY_BUFFER, t.buf)

	gl.EnableVertexAttribArray(t.position)
	gl.VertexAttribPointer(t.position, 4, gl.FLOAT, false, sizeVertex, 0)
	gl.EnableVertexAttribArray(t.color)
	gl.VertexAttribPointer(t.color, 4, gl.FLOAT, false, sizeVertex, sizeCoords)

	gl.DrawArrays(t.mode, 0, t.vertexCount)

	gl.DisableVertexAttribArray(t.position)
	gl.DisableVertexAttribArray(t.color)
}

func (t *ModelGroup) Draw() {
	for _, m := range t.models {
		m.modelView = t.modelView
		m.Draw()
	}
}

func (t *ModelBase) Destroy() {
	gl.DeleteBuffer(t.buf)
	//t.vao.Delete()
	gl.DeleteProgram(t.prg)
}

func (t *ModelGroup) Destroy() {
	for _, m := range t.models {
		m.Destroy()
	}
}
