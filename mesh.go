package sge

import "reflect"

import "gl"

type Mesh struct {
	Attrs []interface{}
	Indicies []uint32
	attrValues []reflect.Value
	attrDimensions []uint
	vao gl.VertexArray
	attrBOs []gl.Buffer
	indexBO gl.Buffer
}

func NewMesh(numVerticies uint, indicies []uint32, attrs ...interface{}) *Mesh {
	mesh := new(Mesh)
	mesh.Attrs = attrs
	mesh.Indicies = indicies
	mesh.attrValues = make([]reflect.Value, len(attrs))
	mesh.attrDimensions = make([]uint, len(attrs))
	mesh.attrBOs = make([]gl.Buffer, len(attrs))
	mesh.vao = gl.GenVertexArray()
	mesh.vao.Bind()
	for i, attr := range attrs {
		mesh.attrValues[i] = reflect.ValueOf(attr)
		if mesh.attrValues[i].Kind() != reflect.Slice {
			panic("an element of attrs is not a slice")
		}
		mesh.attrDimensions[i] = uint(mesh.attrValues[i].Len()) / numVerticies
		mesh.attrBOs[i] = setupVBO(i, mesh.attrValues[i], mesh.attrDimensions[i])
	}
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
	t := glType(data)
	if t == gl.FLOAT || t == gl.DOUBLE {
		attrib.AttribPointerOffset(dimensions, t, false, 0, 0)
	} else {
		attrib.AttribIPointerOffset(dimensions, t, 0, 0)
	}
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
	for _, bo := range mesh.attrBOs {
		bo.Delete()
	}
	mesh.indexBO.Delete()
}

func (mesh *Mesh) Render() {
	mesh.vao.Bind()
	gl.DrawElements(gl.TRIANGLES, len(mesh.Indicies), gl.UNSIGNED_INT, 0)
}
