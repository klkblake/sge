package sge

import "io/ioutil"

import "gl"

const defaultVertexSource =
	"#version 330\n" +
	"layout(location=0) in vec3 position;\n" +
	"layout(location=1) in vec3 texcoord;\n" +
	"out vec3 fragTexcoord;\n" +
	"uniform mat4 mvpMatrix;\n" +
	"void main(void) {\n" +
	"	gl_Position = mvpMatrix * vec4(position, 1.0);\n" +
	"	fragTexcoord = texcoord;\n" +
	"}\n"

const defaultFragmentSource =
	"#version 330\n" +
	"uniform sampler2D textureUnit;\n" +
	"in vec3 fragTexcoord;\n" +
	"out vec4 color;\n" +
	"void main(void) {\n" +
	"	color = texture(textureUnit, fragTexcoord.xy);\n" +
	"}\n"

const defaultCubeFragmentSource =
	"#version 330\n" +
	"uniform samplerCube textureUnit;\n" +
	"in vec3 fragTexcoord;\n" +
	"out vec4 color;\n" +
	"void main(void) {\n" +
	"	color = texture(textureUnit, fragTexcoord);\n" +
	"}\n"

type Shader struct {
	gl.Shader
}

func NewShader(type_ gl.GLenum) *Shader {
	shader := &Shader{gl.CreateShader(type_)}
	if shader.Shader == 0 {
		panic(gl.GetError())
	}
	return shader
}

func NewVertexShader() *Shader {
	return NewShader(gl.VERTEX_SHADER)
}

func NewFragmentShader() *Shader {
	return NewShader(gl.FRAGMENT_SHADER)
}

func LoadShader(filename string, type_ gl.GLenum) *Shader {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.String())
	}
	shader := NewShader(type_)
	shader.Source(string(data))
	shader.Compile()
	return shader
}

func LoadVertexShader(filename string) *Shader {
	return LoadShader(filename, gl.VERTEX_SHADER)
}

func LoadFragmentShader(filename string) *Shader {
	return LoadShader(filename, gl.FRAGMENT_SHADER)
}

func DefaultVertexShader() *Shader {
	shader := NewVertexShader()
	shader.Source(defaultVertexSource)
	shader.Compile()
	return shader
}

func DefaultFragmentShader() *Shader {
	shader := NewFragmentShader()
	shader.Source(defaultFragmentSource)
	shader.Compile()
	return shader
}

func DefaultCubeFragmentShader() *Shader {
	shader := NewFragmentShader()
	shader.Source(defaultCubeFragmentSource)
	shader.Compile()
	return shader
}

func (shader *Shader) Compile() {
	shader.Shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) == gl.FALSE {
		panic(shader.GetInfoLog())
	}
}
