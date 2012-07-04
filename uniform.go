package sge

import (
	"github.com/chsc/gogl/gl33"
)

type Uniform struct {
	Location gl33.Int
}

func (uniform *Uniform) Set(data interface{}) {
	loc := uniform.Location
	switch value := data.(type) {
	case float32:
		gl33.Uniform1f(loc, gl33.Float(value))
	case int:
		gl33.Uniform1i(loc, gl33.Int(value))
	case int32:
		gl33.Uniform1i(loc, gl33.Int(value))
	case []float32:
		gl33.Uniform1fv(loc, gl33.Sizei(len(value)), (*gl33.Float)(&value[0]))
	case []int32:
		gl33.Uniform1iv(loc, gl33.Sizei(len(value)), (*gl33.Int)(&value[0]))
	default:
		panic("Unsupported type of uniform")
	}
}

func (uniform *Uniform) SetVector(data interface{}, size int) {
	if size < 2 || size > 4 {
		panic("Invalid size")
	}
	loc := uniform.Location
	switch value := data.(type) {
	case []float32:
		length := gl33.Sizei(len(value))
		ptr := (*gl33.Float)(&value[0])
		switch size {
		case 2:
			gl33.Uniform2fv(loc, length, ptr)
		case 3:
			gl33.Uniform3fv(loc, length, ptr)
		case 4:
			gl33.Uniform4fv(loc, length, ptr)
		}
	case []int32:
		length := gl33.Sizei(len(value))
		ptr := (*gl33.Int)(&value[0])
		switch size {
		case 2:
			gl33.Uniform2iv(loc, length, ptr)
		case 3:
			gl33.Uniform3iv(loc, length, ptr)
		case 4:
			gl33.Uniform4iv(loc, length, ptr)
		}
	default:
		panic("Unsupported type of uniform")
	}
}

func (uniform *Uniform) SetMatrix(data []float32, size int) {
	if size < 2 || size > 4 {
		panic("Invalid size")
	}
	loc := uniform.Location
	length := gl33.Sizei(len(data) / (size * size))
	ptr := (*gl33.Float)(&data[0])
	switch size {
	case 2:
		gl33.UniformMatrix2fv(loc, length, gl33.FALSE, ptr)
	case 3:
		gl33.UniformMatrix3fv(loc, length, gl33.FALSE, ptr)
	case 4:
		gl33.UniformMatrix4fv(loc, length, gl33.FALSE, ptr)
	}
}
