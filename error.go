package sge

import (
	"strconv"
)

import (
	"github.com/chsc/gogl/gl33"
)

func PanicOnError() {
	GL <- func() {
		err := gl33.GetError()
		if err != 0 {
			panic("OpenGL error occured. Error code: " + strconv.Itoa(int(err)))
		}
	}
}
