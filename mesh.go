package sge

import (
	"reflect"
	"unsafe"
)

import "gl"

type Mesh struct {
	Attrs interface{}
	Indicies []uint32
	vao gl.VertexArray
	vertexBO gl.Buffer
	indexBO gl.Buffer
}

func NewMesh(attrs interface{}, indicies []uint32) *Mesh {
	mesh := new(Mesh)
	mesh.Attrs = attrs
	mesh.Indicies = indicies
	mesh.vao = gl.GenVertexArray()
	mesh.vao.Bind()
	attrsValue := reflect.ValueOf(attrs)
	if attrsValue.Kind() != reflect.Slice {
		panic("attrs is not a slice")
	}
	mesh.vertexBO = createBuffer(gl.ARRAY_BUFFER, reflect.ValueOf(attrs))
	vertexSpec := attrsValue.Type().Elem()
	if vertexSpec.Kind() != reflect.Struct && vertexSpec.Kind() != reflect.Array {
		panic("attrs is not a slice of structs or arrays")
	}
	var num int
	if vertexSpec.Kind() == reflect.Struct {
		num = vertexSpec.NumField()
	} else {
		num = vertexSpec.Len()
	}
	for i, offset := 0, uintptr(0); i < num; i++ {
		var field reflect.Type
		var type_ gl.GLenum
		var dimensions uint
		if vertexSpec.Kind() == reflect.Struct {
			field = vertexSpec.Field(i).Type
		} else {
			field = vertexSpec.Elem()
		}
		if field.Kind() == reflect.Array {
			type_ = glType(field.Elem().Kind())
			dimensions = uint(field.Len())
		} else {
			type_ = glType(field.Kind())
			dimensions = 1
		}
		setupAttrib(i, type_, dimensions, offset, int(vertexSpec.Size()))
		offset += field.Size()
	}
	mesh.indexBO = createBuffer(gl.ELEMENT_ARRAY_BUFFER, reflect.ValueOf(indicies))
	return mesh
}

func createBuffer(target gl.GLenum, data reflect.Value) gl.Buffer {
	buf := gl.GenBuffer()
	buf.Bind(target)
	gl.BufferData(target, data.Len()*int(data.Type().Elem().Size()), (*byte)(unsafe.Pointer(data.Pointer())), gl.DYNAMIC_DRAW)
	return buf
}

func setupAttrib(location int, type_ gl.GLenum, dimensions uint, offset uintptr, vertexSize int) {
	attrib := gl.AttribLocation(location)
	attrib.EnableArray()
	if type_ == gl.FLOAT || type_ == gl.DOUBLE {
		attrib.AttribPointerOffset(dimensions, type_, false, vertexSize, offset)
	} else {
		attrib.AttribIPointerOffset(dimensions, type_, vertexSize, offset)
	}
}

func glType(data reflect.Kind) gl.GLenum {
	switch data {
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
	mesh.indexBO.Delete()
}

func (mesh *Mesh) Render() {
	mesh.vao.Bind()
	gl.DrawElements(gl.TRIANGLES, len(mesh.Indicies), gl.UNSIGNED_INT, 0)
}
