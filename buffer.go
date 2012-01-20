package sge

import (
	"reflect"
)

import (
	"github.com/chsc/gogl/gl33"
)

var boundBuffer map[gl33.Enum]*Buffer

func init() {
	boundBuffer = make(map[gl33.Enum]*Buffer)
}

type Buffer struct {
	Id     gl33.Uint
	Target gl33.Enum
	Type   gl33.Enum
	Data   interface{}
	Value  reflect.Value
}

func NewBuffer(target gl33.Enum, type_ gl33.Enum, data interface{}) *Buffer {
	buf := new(Buffer)
	gl33.GenBuffers(1, &buf.Id)
	buf.Target = target
	buf.Type = type_
	buf.Data = data
	buf.Value = reflect.ValueOf(data)
	if buf.Value.Kind() != reflect.Slice {
		panic("data is not a slice")
	}
	buf.Update()
	return buf
}

func (buffer *Buffer) Bind() {
	if boundBuffer[buffer.Target] != buffer {
		gl33.BindBuffer(buffer.Target, buffer.Id)
		boundBuffer[buffer.Target] = buffer
	}
}

func (buffer *Buffer) Update() {
	buffer.Bind()
	data := gl33.Pointer(buffer.Value.Pointer())
	size := gl33.Sizeiptr(buffer.Value.Len()*int(buffer.Value.Type().Elem().Size()))
	gl33.BufferData(buffer.Target, size, data, buffer.Type)
}

func (buffer *Buffer) Delete() {
	gl33.DeleteBuffers(1, &buffer.Id)
}
