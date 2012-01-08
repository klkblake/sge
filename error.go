package sge

import "strconv"

import "gl"

func PanicOnError() {
	err := gl.GetError()
	if err != 0 {
		panic("OpenGL error occured. Error code: " + strconv.Itoa(int(err)))
	}
}
