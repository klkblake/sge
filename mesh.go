package sge

import "reflect"

import "gl"

type Mesh struct {
	Verticies interface{}
	Texcoords interface{}
	Indicies []uint16
	verticiesValue reflect.Value
	texcoordsValue reflect.Value
	vertexDimensions uint
	texDimensions uint
	vao gl.VertexArray
	vertexBO gl.Buffer
	texcoordBO gl.Buffer
	indexBO gl.Buffer
}

func NewMesh(verticies interface{}, texcoords interface{}, indicies []uint16, vertexDimensions uint) *Mesh {
	mesh := new(Mesh)
	mesh.Verticies = verticies
	mesh.Texcoords = texcoords
	mesh.Indicies = indicies
	mesh.verticiesValue = reflect.ValueOf(verticies)
	if mesh.verticiesValue.Kind() != reflect.Slice {
		panic("verticies is not a slice")
	}
	mesh.texcoordsValue = reflect.ValueOf(texcoords)
	if mesh.texcoordsValue.Kind() != reflect.Slice {
		panic("texcoords is not a slice")
	}
	mesh.vertexDimensions = vertexDimensions
	mesh.texDimensions = uint(mesh.texcoordsValue.Len() / (mesh.verticiesValue.Len() / int(vertexDimensions)))
	mesh.vao = gl.GenVertexArray()
	mesh.vao.Bind()
	mesh.vertexBO = setupVBO(0, mesh.verticiesValue, mesh.vertexDimensions)
	mesh.texcoordBO = setupVBO(1, mesh.texcoordsValue, mesh.texDimensions)
	mesh.indexBO = createBuffer(gl.ELEMENT_ARRAY_BUFFER, reflect.ValueOf(indicies))
	return mesh
}

func createBuffer(target gl.GLenum, data reflect.Value) gl.Buffer {
	buf := gl.GenBuffer()
	buf.Bind(target)
	gl.BufferData(target, data.Len()*int(data.Type().Elem().Size()), data.Interface(), gl.DYNAMIC_DRAW)
	return buf
}

func setupVBO(location int, data reflect.Value, dimensions uint) gl.Buffer {
	buf := createBuffer(gl.ARRAY_BUFFER, data)
	attrib := gl.AttribLocation(location)
	attrib.EnableArray()
	attrib.AttribPointerOffset(dimensions, glType(data), false, 0, 0)
	return buf
}

func glType(data reflect.Value) gl.GLenum {
	switch data.Type().Elem().Kind() {
	case reflect.Int8:
		return gl.BYTE
	case reflect.Int16:
		return gl.SHORT
	case reflect.Int32:
		return gl.INT
	case reflect.Uint8:
		return gl.UNSIGNED_BYTE
	case reflect.Uint16:
		return gl.UNSIGNED_SHORT
	case reflect.Uint32:
		return gl.UNSIGNED_INT
	case reflect.Float32:
		return gl.FLOAT
	case reflect.Float64:
		return gl.DOUBLE
	}
	panic("Bad element type")
}

func (mesh *Mesh) Delete() {
	mesh.vao.Delete()
	mesh.vertexBO.Delete()
	mesh.texcoordBO.Delete()
	mesh.indexBO.Delete()
}

func (mesh *Mesh) Render() {
	mesh.vao.Bind()
	gl.DrawElements(gl.TRIANGLES, len(mesh.Indicies), gl.UNSIGNED_SHORT, 0)
}
