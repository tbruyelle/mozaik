package main

import (
	"github.com/tbruyelle/gl"
	"github.com/remogatto/mathgl"
)

type Model interface {
	Draw()
	Destroy()
	//pushModelView(modelView mathgl.Mat4f)
	//popModelView()
}

type ModelBase struct {
	mode                  gl.GLenum
	buffer                gl.Buffer
	vertices              []Vertex
	sizeVertices          int
	prg                   gl.Program
	posLoc                gl.AttribLocation
	colLoc                gl.AttribLocation
	vao                   gl.VertexArray
	uniformMVP            gl.UniformLocation
	modelView, projection mathgl.Mat4f
	modelViewBackup       mathgl.Mat4f
	vshader, fshader      gl.Shader
}

type ModelGroup struct {
	models                []*ModelBase
	modelView, projection mathgl.Mat4f
}

func (t *ModelGroup) Add(mode gl.GLenum, vertices []Vertex, vshaderf, fshaderf string) {
	m := &ModelBase{}
	m.Init(mode, vertices, vshaderf, fshaderf)
	t.models = append(t.models, m)
}

//func (t *ModelBase) pushModelView(modelView mathgl.Mat4f) {
//	t.modelViewBackup = t.modelView
//	t.modelView = modelView.Mul4(t.modelView)
//}
//
//func (t *ModelBase) popModelView() {
//	t.modelView = t.modelViewBackup
//}

func (t *ModelBase) Init(mode gl.GLenum, vertices []Vertex, vshaderf, fshaderf string) {
	t.mode = mode
	t.vertices = vertices
	t.sizeVertices = len(t.vertices) * sizeVertex

	// Shaders
	t.vshader = compileShader(gl.VERTEX_SHADER, vshaderf)
	t.fshader = compileShader(gl.FRAGMENT_SHADER, fshaderf)
	t.prg = NewProgram(t.vshader, t.fshader)
	t.posLoc = gl.AttribLocation(0)
	t.colLoc = gl.AttribLocation(1)
	t.uniformMVP = t.prg.GetUniformLocation("modelViewProjection")

	// the projection matrix
	t.projection = mathgl.Ident4f()

	// the model view
	t.modelView = mathgl.Ident4f()

	// Create VBO
	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)
	t.buffer.Unbind(gl.ARRAY_BUFFER)

	// Create VAO
	t.vao = gl.GenVertexArray()
	t.vao.Bind()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	// Attrib vertex data to VAO
	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	t.posLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))
	t.colLoc.EnableArray()

	t.buffer.Unbind(gl.ARRAY_BUFFER)
	t.vao.Unbind()
}

func (t *ModelBase) Draw() {
	t.prg.Use()

	t.vao.Bind()

	mvp := t.modelView.Mul4(t.projection)
	t.uniformMVP.UniformMatrix4f(false, (*[16]float32)(&mvp))

	gl.DrawArrays(t.mode, 0, len(t.vertices))

	t.vao.Unbind()
	t.prg.Unuse()
}

func (t *ModelGroup) Draw() {
	for _, m := range t.models {
		m.modelView = t.modelView
		m.Draw()
	}
}

func (t *ModelBase) Destroy() {
	t.buffer.Delete()
	t.vao.Delete()
	t.vshader.Delete()
	t.fshader.Delete()
	t.prg.Delete()
}

func (t *ModelGroup) Destroy() {
	for _, m := range t.models {
		m.Destroy()
	}
}
