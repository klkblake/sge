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
