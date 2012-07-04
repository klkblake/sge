package sge

import (
	"io/ioutil"
)

import (
	"github.com/chsc/gogl/gl33"
)

const defaultVertexSource = "#version 330\n" +
	"layout(location=0) in vec3 position;\n" +
	"layout(location=1) in vec3 texcoord;\n" +
	"out vec3 fragTexcoord;\n" +
	"uniform mat4 mvpMatrix;\n" +
	"void main(void) {\n" +
	"	gl_Position = mvpMatrix * vec4(position, 1.0);\n" +
	"	fragTexcoord = texcoord;\n" +
	"}\n"

const defaultFragmentSource = "#version 330\n" +
	"uniform sampler2D textureSampler;\n" +
	"in vec3 fragTexcoord;\n" +
	"out vec4 color;\n" +
	"void main(void) {\n" +
	"	color = texture(textureSampler, fragTexcoord.xy);\n" +
	"}\n"

const defaultCubeFragmentSource = "#version 330\n" +
	"uniform samplerCube textureSampler;\n" +
	"in vec3 fragTexcoord;\n" +
	"out vec4 color;\n" +
	"void main(void) {\n" +
	"	color = texture(textureSampler, fragTexcoord);\n" +
	"}\n"

type Shader struct {
	Id gl33.Uint
}

func NewShader(type_ gl33.Enum) *Shader {
	shader := new(Shader)
	shader.Id = gl33.CreateShader(type_)
	if shader.Id == 0 {
		panic(gl33.GetError())
	}
	return shader
}

func NewVertexShader() *Shader {
	return NewShader(gl33.VERTEX_SHADER)
}

func NewFragmentShader() *Shader {
	return NewShader(gl33.FRAGMENT_SHADER)
}

func LoadShader(filename string, type_ gl33.Enum) *Shader {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ch := make(chan *Shader)
	GL <- func() {
		shader := NewShader(type_)
		shader.Source(string(data))
		shader.Compile()
		ch <- shader
	}
	return <-ch
}

func LoadVertexShader(filename string) *Shader {
	return LoadShader(filename, gl33.VERTEX_SHADER)
}

func LoadFragmentShader(filename string) *Shader {
	return LoadShader(filename, gl33.FRAGMENT_SHADER)
}

func DefaultVertexShader() *Shader {
	ch := make(chan *Shader)
	GL <- func() {
		shader := NewVertexShader()
		shader.Source(defaultVertexSource)
		shader.Compile()
		ch <- shader
	}
	return <-ch
}

func DefaultFragmentShader() *Shader {
	ch := make(chan *Shader)
	GL <- func() {
		shader := NewFragmentShader()
		shader.Source(defaultFragmentSource)
		shader.Compile()
		ch <- shader
	}
	return <-ch
}

func DefaultCubeFragmentShader() *Shader {
	ch := make(chan *Shader)
	GL <- func() {
		shader := NewFragmentShader()
		shader.Source(defaultCubeFragmentSource)
		shader.Compile()
		ch <- shader
	}
	return <-ch
}

func (shader *Shader) Get(param gl33.Enum) int {
	var ret gl33.Int
	gl33.GetShaderiv(shader.Id, param, &ret)
	return int(ret)
}

func (shader *Shader) GetInfoLog() string {
	length := shader.Get(gl33.INFO_LOG_LENGTH)
	log := gl33.GLStringAlloc(gl33.Sizei(length + 1))
	defer gl33.GLStringFree(log)
	gl33.GetShaderInfoLog(shader.Id, gl33.Sizei(length), nil, log)
	return gl33.GoString(log)
}

func (shader *Shader) Source(source string) {
	str := gl33.GLString(source)
	defer gl33.GLStringFree(str)
	length := gl33.Int(len(source))
	gl33.ShaderSource(shader.Id, 1, &str, &length)
}

func (shader *Shader) Compile() {
	gl33.CompileShader(shader.Id)
	if shader.Get(gl33.COMPILE_STATUS) == gl33.FALSE {
		panic(shader.GetInfoLog())
	}
}
