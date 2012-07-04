package sge

import (
	"github.com/chsc/gogl/gl33"
)

var usedProgram *Program

type Program struct {
	Id gl33.Uint
}

func NewShaderProgram() *Program {
	program := new(Program)
	program.Id = gl33.CreateProgram()
	return program
}

func LoadShaderProgram(vertShader *Shader, fragShader *Shader) *Program {
	ch := make(chan *Program)
	GL <- func() {
		program := NewShaderProgram()
		program.AttachShader(vertShader)
		program.AttachShader(fragShader)
		program.Link()
		ch <- program
	}
	return <-ch
}

func (program *Program) Get(param gl33.Enum) int {
	var ret gl33.Int
	gl33.GetProgramiv(program.Id, param, &ret)
	return int(ret)
}

func (program *Program) GetInfoLog() string {
	length := program.Get(gl33.INFO_LOG_LENGTH)
	log := gl33.GLStringAlloc(gl33.Sizei(length + 1))
	defer gl33.GLStringFree(log)
	gl33.GetProgramInfoLog(program.Id, gl33.Sizei(length), nil, log)
	return gl33.GoString(log)
}

func (program *Program) Uniform(name string) *Uniform {
	str := gl33.GLString(name)
	defer gl33.GLStringFree(str)
	return &Uniform{gl33.GetUniformLocation(program.Id, str)}
}

func (program *Program) AttachShader(shader *Shader) {
	gl33.AttachShader(program.Id, shader.Id)
}

func (program *Program) Link() {
	gl33.LinkProgram(program.Id)
	if program.Get(gl33.LINK_STATUS) == gl33.FALSE {
		panic(program.GetInfoLog())
	}
}

func (program *Program) Use() {
	if usedProgram != program {
		gl33.UseProgram(program.Id)
		usedProgram = program
	}
}
