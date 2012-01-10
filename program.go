package sge

import "gl"

var usedProgram *Program

type Program struct {
	gl.Program
}

func NewShaderProgram() *Program {
	return &Program{gl.CreateProgram()}
}

func LoadShaderProgram(vertShader *Shader, fragShader *Shader) *Program {
	program := NewShaderProgram()
	program.AttachShader(vertShader.Shader)
	program.AttachShader(fragShader.Shader)
	program.Link()
	return program
}

func (program *Program) SetUniform(name string, data interface{}) {
	location := program.GetUniformLocation(name)
	switch value := data.(type) {
	case float32:
		location.Uniform1f(value)
	case int:
		location.Uniform1i(value)
	case []float32:
		location.Uniform1fv(value)
	case []int:
		location.Uniform1iv(value)
	default:
		panic("Unsupported type of uniform")
	}
}

func (program *Program) SetUniformVector(name string, data interface{}, size int) {
	if size < 2 || size > 4 {
		panic("Invalid size")
	}
	location := program.GetUniformLocation(name)
	switch value := data.(type) {
	case []float32:
		switch size {
		case 2:
			location.Uniform2fv(value)
		case 3:
			location.Uniform3fv(value)
		case 4:
			location.Uniform4fv(value)
		}
	case []int32:
		switch size {
		case 2:
			location.Uniform2iv(value)
		case 3:
			location.Uniform3iv(value)
		case 4:
			location.Uniform4iv(value)
		}
	default:
		panic("Unsupported type of uniform")
	}
}

func (program *Program) SetUniformMatrix(name string, data []float32, size int) {
	if size != 4 {
		panic ("Invalid size")
	}
	location := program.GetUniformLocation(name)
	location.UniformMatrix4fv(false, len(data) / (size*size), data)
}

func (program *Program) Link() {
	program.Program.Link()
	if program.Get(gl.LINK_STATUS) == gl.FALSE {
		panic(program.GetInfoLog())
	}
}

func (program *Program) Use() {
	if usedProgram != program {
		program.Program.Use()
		usedProgram = program
	}
}
