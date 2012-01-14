package sge

import "gl"

type Assets struct {
	textures2D map[string]*Texture
	texturesArray map[string]*Texture
	texturesCubeMap map[string]*Texture
	MinFilter int
	MagFilter int
	vertShaders map[string]*Shader
	fragShaders map[string]*Shader
	UseDefaultShaders bool
	shaderPrograms map[string]*Program
}

func NewAssets() *Assets {
	assets := new(Assets)
	assets.textures2D = make(map[string]*Texture)
	assets.texturesArray = make(map[string]*Texture)
	assets.texturesCubeMap = make(map[string]*Texture)
	assets.MinFilter = gl.LINEAR_MIPMAP_LINEAR
	assets.MagFilter = gl.LINEAR
	assets.vertShaders = make(map[string]*Shader)
	assets.fragShaders = make(map[string]*Shader)
	assets.UseDefaultShaders = true
	assets.shaderPrograms = make(map[string]*Program)
	return assets
}

func (assets *Assets) Texture2D(name string) *Texture {
	if assets.textures2D[name] == nil {
		assets.LoadTexture2D(name)
	}
	return assets.textures2D[name]
}

func (assets *Assets) LoadTexture2D(name string) {
	assets.textures2D[name] = LoadTexture2D(name, assets.MinFilter, assets.MagFilter)
}

func (assets *Assets) TextureArray(names []string) *Texture {
	name := names[0]
	for i := 1; i < len(names); i++ {
		name += "\x00" + names[i]
	}
	if assets.texturesArray[name] == nil {
		assets.LoadTextureArray(names)
	}
	return assets.texturesArray[name]
}

func (assets *Assets) LoadTextureArray(names []string) {
	name := names[0]
	for i := 1; i < len(names); i++ {
		name += "\x00" + names[i]
	}
	assets.texturesArray[name] = LoadTextureArray(names, assets.MinFilter, assets.MagFilter)
}

func (assets *Assets) TextureCubeMap(names *[6]string) *Texture {
	name := names[0]
	for i := 1; i < 6; i++ {
		name += "\x00" + names[i]
	}
	if assets.texturesCubeMap[name] == nil {
		assets.LoadTextureCubeMap(names)
	}
	return assets.texturesCubeMap[name]
}

func (assets *Assets) LoadTextureCubeMap(names *[6]string) {
	name := names[0]
	for i := 1; i < 6; i++ {
		name += "\x00" + names[i]
	}
	assets.texturesCubeMap[name] = LoadTextureCubeMap(names, assets.MinFilter, assets.MagFilter)
}

func (assets *Assets) VertexShader(name string) *Shader {
	if assets.vertShaders[name] == nil {
		assets.LoadVertexShader(name)
	}
	return assets.vertShaders[name]
}

func (assets *Assets) LoadVertexShader(name string) {
	if name == "default" && assets.UseDefaultShaders {
		assets.vertShaders[name] = DefaultVertexShader()
	} else {
		assets.vertShaders[name] = LoadVertexShader(name)
	}
}

func (assets *Assets) FragmentShader(name string) *Shader {
	if assets.fragShaders[name] == nil {
		assets.LoadFragmentShader(name)
	}
	return assets.fragShaders[name]
}

func (assets *Assets) LoadFragmentShader(name string) {
	if name == "default" && assets.UseDefaultShaders {
		assets.fragShaders[name] = DefaultFragmentShader()
	} else if name == "defaultCube" && assets.UseDefaultShaders {
		assets.fragShaders[name] = DefaultCubeFragmentShader()
	} else {
		assets.fragShaders[name] = LoadFragmentShader(name)
	}
}

func (assets *Assets) ShaderProgram(vertexName string, fragmentName string) *Program {
	name := vertexName + "\x00" + fragmentName
	if assets.shaderPrograms[name] == nil {
		assets.LoadShaderProgram(vertexName, fragmentName)
	}
	return assets.shaderPrograms[name]
}

func (assets *Assets) LoadShaderProgram(vertexName string, fragmentName string) {
	name := vertexName + "\x00" + fragmentName
	vertShader := assets.VertexShader(vertexName)
	fragShader := assets.FragmentShader(fragmentName)
	assets.shaderPrograms[name] = LoadShaderProgram(vertShader, fragShader)
}
