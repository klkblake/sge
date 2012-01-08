package sge

import "gl"

type Mesh struct {
	Verticies []float32
	Texcoords []float32
	Indicies []uint32
	numTexDimensions uint
	vao gl.VertexArray
	vertexBO gl.Buffer
	texcoordBO gl.Buffer
	indexBO gl.Buffer
}

func NewMesh(verticies []float32, texcoords []float32, indicies []uint32) *Mesh {
	mesh := new(Mesh)
	mesh.Verticies = verticies
	mesh.Texcoords = texcoords
	mesh.Indicies = indicies
	mesh.numTexDimensions = uint(len(texcoords) / (len(verticies) / 3))
	mesh.vao = gl.GenVertexArray()
	mesh.vao.Bind()
	mesh.vertexBO = setupVBO(0, verticies, 3)
	mesh.texcoordBO = setupVBO(1, texcoords, mesh.numTexDimensions)
	mesh.indexBO = createBuffer(gl.ELEMENT_ARRAY_BUFFER, indicies, len(indicies)*4)
	return mesh
}

func createBuffer(target gl.GLenum, data interface{}, size int) gl.Buffer {
	buf := gl.GenBuffer()
	buf.Bind(target)
	gl.BufferData(target, size, data, gl.DYNAMIC_DRAW)
	return buf
}

func setupVBO(location int, data []float32, dimensions uint) gl.Buffer {
	buf := createBuffer(gl.ARRAY_BUFFER, data, len(data)*4)
	attrib := gl.AttribLocation(location)
	attrib.EnableArray()
	attrib.AttribPointerOffset(dimensions, gl.FLOAT, false, 0, 0)
	return buf
}

func (mesh *Mesh) Render() {
	mesh.vao.Bind()
	gl.DrawElements(gl.TRIANGLES, len(mesh.Indicies), gl.UNSIGNED_INT, 0)
}
