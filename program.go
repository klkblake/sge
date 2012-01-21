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
	GL <- func() {
		program.Id = gl33.CreateProgram()
	}
	return program
}

func LoadShaderProgram(vertShader *Shader, fragShader *Shader) *Program {
	program := NewShaderProgram()
	program.AttachShader(vertShader)
	program.AttachShader(fragShader)
	program.Link()
	return program
}

func (program *Program) get(param gl33.Enum) int {
	var ret gl33.Int
	gl33.GetProgramiv(program.Id, param, &ret)
	return int(ret)
}

func (program *Program) Get(param gl33.Enum) int {
	ret := make(chan int, 1)
	GL <- func() {
		ret <- program.get(param)
	}
	return <-ret
}

func (program *Program) getInfoLog() string {
	length := program.get(gl33.INFO_LOG_LENGTH)
	log := gl33.GLStringAlloc(gl33.Sizei(length + 1))
	defer gl33.GLStringFree(log)
	gl33.GetProgramInfoLog(program.Id, gl33.Sizei(length), nil, log)
	return gl33.GoString(log)
}

func (program *Program) GetInfoLog() string {
	ret := make(chan string, 1)
	GL <- func() {
		ret <- program.getInfoLog()
	}
	return <-ret
}

func (program *Program) getUniformLocation(name string) int {
	str := gl33.GLString(name)
	defer gl33.GLStringFree(str)
	return int(gl33.GetUniformLocation(program.Id, str))
}

func (program *Program) GetUniformLocation(name string) int {
	ret := make(chan int, 1)
	GL <- func() {
		ret <- program.getUniformLocation(name)
	}
	return <-ret
}

func (program *Program) SetUniform(name string, data interface{}) {
	GL <- func() {
		location := gl33.Int(program.getUniformLocation(name))
		switch value := data.(type) {
		case float32:
			gl33.Uniform1f(location, gl33.Float(value))
		case int:
			gl33.Uniform1i(location, gl33.Int(value))
		case int32:
			gl33.Uniform1i(location, gl33.Int(value))
		case []float32:
			gl33.Uniform1fv(location, gl33.Sizei(len(value)), (*gl33.Float)(&value[0]))
		case []int32:
			gl33.Uniform1iv(location, gl33.Sizei(len(value)), (*gl33.Int)(&value[0]))
		default:
			panic("Unsupported type of uniform")
		}
	}
}

func (program *Program) SetUniformVector(name string, data interface{}, size int) {
	if size < 2 || size > 4 {
		panic("Invalid size")
	}
	GL <- func() {
		location := gl33.Int(program.getUniformLocation(name))
		switch value := data.(type) {
		case []float32:
			switch size {
			case 2:
				gl33.Uniform2fv(location, gl33.Sizei(len(value)), (*gl33.Float)(&value[0]))
			case 3:
				gl33.Uniform3fv(location, gl33.Sizei(len(value)), (*gl33.Float)(&value[0]))
			case 4:
				gl33.Uniform4fv(location, gl33.Sizei(len(value)), (*gl33.Float)(&value[0]))
			}
		case []int32:
			switch size {
			case 2:
				gl33.Uniform2iv(location, gl33.Sizei(len(value)), (*gl33.Int)(&value[0]))
			case 3:
				gl33.Uniform3iv(location, gl33.Sizei(len(value)), (*gl33.Int)(&value[0]))
			case 4:
				gl33.Uniform4iv(location, gl33.Sizei(len(value)), (*gl33.Int)(&value[0]))
			}
		default:
			panic("Unsupported type of uniform")
		}
	}
}

func (program *Program) SetUniformMatrix(name string, data []float32, size int) {
	if size < 2 || size > 4 {
		panic("Invalid size")
	}
	GL <- func() {
		location := gl33.Int(program.getUniformLocation(name))
		switch size {
		case 2:
			gl33.UniformMatrix2fv(location, gl33.Sizei(len(data)/(size*size)), gl33.FALSE, (*gl33.Float)(&data[0]))
		case 3:
			gl33.UniformMatrix3fv(location, gl33.Sizei(len(data)/(size*size)), gl33.FALSE, (*gl33.Float)(&data[0]))
		case 4:
			gl33.UniformMatrix4fv(location, gl33.Sizei(len(data)/(size*size)), gl33.FALSE, (*gl33.Float)(&data[0]))
		}
	}
}

func (program *Program) AttachShader(shader *Shader) {
	GL <- func() {
		gl33.AttachShader(program.Id, shader.Id)
	}
}

func (program *Program) Link() {
	GL <- func() {
		gl33.LinkProgram(program.Id)
		if program.get(gl33.LINK_STATUS) == gl33.FALSE {
			panic(program.getInfoLog())
		}
	}
}

func (program *Program) Use() {
	if usedProgram != program {
		GL <- func() {
			gl33.UseProgram(program.Id)
		}
		usedProgram = program
	}
}
