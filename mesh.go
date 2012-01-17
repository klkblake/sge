package sge

import (
	"reflect"
)

import (
	"github.com/chsc/gogl/gl33"
)

type Mesh struct {
	Attrs    interface{}
	Indicies []uint32
	vao      gl33.Uint
	vertexBO gl33.Uint
	indexBO  gl33.Uint
}

func NewMesh(attrs interface{}, indicies []uint32) *Mesh {
	mesh := new(Mesh)
	mesh.Attrs = attrs
	mesh.Indicies = indicies
	gl33.GenVertexArrays(1, &mesh.vao)
	gl33.BindVertexArray(mesh.vao)
	attrsValue := reflect.ValueOf(attrs)
	if attrsValue.Kind() != reflect.Slice {
		panic("attrs is not a slice")
	}
	mesh.vertexBO = createBuffer(gl33.ARRAY_BUFFER, reflect.ValueOf(attrs))
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
		var type_ gl33.Enum
		var dimensions int
		if vertexSpec.Kind() == reflect.Struct {
			field = vertexSpec.Field(i).Type
		} else {
			field = vertexSpec.Elem()
		}
		if field.Kind() == reflect.Array {
			type_ = glType(field.Elem().Kind())
			dimensions = field.Len()
		} else {
			type_ = glType(field.Kind())
			dimensions = 1
		}
		setupAttrib(gl33.Uint(i), type_, dimensions, offset, int(vertexSpec.Size()))
		offset += field.Size()
	}
	mesh.indexBO = createBuffer(gl33.ELEMENT_ARRAY_BUFFER, reflect.ValueOf(indicies))
	return mesh
}

func createBuffer(target gl33.Enum, data reflect.Value) gl33.Uint {
	var buf gl33.Uint
	gl33.GenBuffers(1, &buf)
	gl33.BindBuffer(target, buf)
	gl33.BufferData(target, gl33.Sizeiptr(data.Len()*int(data.Type().Elem().Size())), gl33.Pointer(data.Pointer()), gl33.DYNAMIC_DRAW)
	return buf
}

func setupAttrib(location gl33.Uint, type_ gl33.Enum, dimensions int, offset uintptr, vertexSize int) {
	gl33.EnableVertexAttribArray(location)
	if type_ == gl33.FLOAT || type_ == gl33.DOUBLE {
		gl33.VertexAttribPointer(location, gl33.Int(dimensions), type_, gl33.FALSE, gl33.Sizei(vertexSize), gl33.Pointer(offset))
	} else {
		gl33.VertexAttribIPointer(location, gl33.Int(dimensions), type_, gl33.Sizei(vertexSize), gl33.Pointer(offset))
	}
}

func glType(data reflect.Kind) gl33.Enum {
	switch data {
	case reflect.Int8:
		return gl33.BYTE
	case reflect.Int16:
		return gl33.SHORT
	case reflect.Int32:
		return gl33.INT
	case reflect.Uint8:
		return gl33.UNSIGNED_BYTE
	case reflect.Uint16:
		return gl33.UNSIGNED_SHORT
	case reflect.Uint32:
		return gl33.UNSIGNED_INT
	case reflect.Float32:
		return gl33.FLOAT
	case reflect.Float64:
		return gl33.DOUBLE
	}
	panic("Bad element type")
}

func (mesh *Mesh) Delete() {
	gl33.DeleteVertexArrays(1, &mesh.vao)
	gl33.DeleteBuffers(1, &mesh.vertexBO)
	gl33.DeleteBuffers(1, &mesh.indexBO)
}

func (mesh *Mesh) Render() {
	gl33.BindVertexArray(mesh.vao)
	gl33.DrawElements(gl33.TRIANGLES, gl33.Sizei(len(mesh.Indicies)), gl33.UNSIGNED_INT, gl33.Pointer(uintptr(0)))
}
