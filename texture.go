package sge

import "gl"
import "atom/sdl"

var boundTexture map[gl.GLenum]*Texture

func init() {
	boundTexture = make(map[gl.GLenum]*Texture)
}

type Texture struct {
	gl.Texture
	Type gl.GLenum
	Width int
	Height int
}

func NewTexture2D() *Texture {
	tex := new(Texture)
	tex.Texture = gl.GenTexture()
	tex.Type = gl.TEXTURE_2D
	return tex
}

func NewTextureCubeMap() *Texture {
	tex := new(Texture)
	tex.Texture = gl.GenTexture()
	tex.Type = gl.TEXTURE_CUBE_MAP
	return tex
}

func LoadTexture2D(filename string, minFilter int, magFilter int) *Texture {
	surface := sdl.Load(filename)
	defer surface.Free()
	if surface == nil {
		panic(sdl.GetError())
	}
	tex := NewTexture2D()
	tex.Width = int(surface.W)
	tex.Height = int(surface.H)
	tex.SetFilters(minFilter, magFilter)
	tex.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP, gl.TRUE)
	uploadSurface(gl.TEXTURE_2D, surface, minFilter, magFilter)
	return tex
}

func LoadTextureCubeMap(filenames *[6]string, minFilter int, magFilter int) *Texture {
	var surfaces [6]*sdl.Surface
	for i, _ := range surfaces {
		surfaces[i] = sdl.Load(filenames[i])
		defer surfaces[i].Free()
		if surfaces[i] == nil {
			panic(sdl.GetError())
		}
	}
	tex := NewTextureCubeMap()
	tex.Width = int(surfaces[0].W)
	tex.Height = int(surfaces[0].H)
	PanicOnError()
	tex.SetFilters(minFilter, magFilter)
	PanicOnError()
	tex.Bind()
	PanicOnError()
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.GENERATE_MIPMAP, gl.TRUE)
	PanicOnError()
	uploadSurface(gl.TEXTURE_CUBE_MAP_POSITIVE_X, surfaces[0], minFilter, magFilter)
	PanicOnError()
	uploadSurface(gl.TEXTURE_CUBE_MAP_NEGATIVE_X, surfaces[1], minFilter, magFilter)
	PanicOnError()
	uploadSurface(gl.TEXTURE_CUBE_MAP_POSITIVE_Y, surfaces[2], minFilter, magFilter)
	PanicOnError()
	uploadSurface(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, surfaces[3], minFilter, magFilter)
	PanicOnError()
	uploadSurface(gl.TEXTURE_CUBE_MAP_POSITIVE_Z, surfaces[4], minFilter, magFilter)
	PanicOnError()
	uploadSurface(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, surfaces[5], minFilter, magFilter)
	PanicOnError()
	return tex
}

func uploadSurface(target gl.GLenum, surface *sdl.Surface, minFilter int, magFilter int) {
	if target == gl.TEXTURE_CUBE_MAP && surface.W != surface.H {
		panic("Non-square texture in cube map")
	}
	var internalFormat int
	var format gl.GLenum
	if surface.Format.BitsPerPixel == 32 {
		internalFormat = gl.RGBA8
		format = gl.RGBA
	} else {
		internalFormat = gl.RGB8
		format = gl.RGB
	}
	PanicOnError()
	gl.TexImage2D(target, 0, internalFormat, int(surface.W), int(surface.H), 0, format, gl.UNSIGNED_BYTE, (*byte)(surface.Pixels))
	PanicOnError()
	// Work around a driver bug in the AMD proprietary drivers.
	//gl.Enable(gl.TEXTURE_2D)
	// XXX Go-OpenGL currently doesn't support glGenerateMipmap
	//gl.GenerateMipmap(gl.TEXTURE_2D)
}

func UnbindTexture2D() {
	boundTexture[gl.TEXTURE_2D].Unbind()
}

func UnbindTextureCubeMap() {
	boundTexture[gl.TEXTURE_CUBE_MAP].Unbind()
}

func (tex *Texture) Bind() {
	if boundTexture[tex.Type] != tex {
		tex.Texture.Bind(tex.Type)
		boundTexture[tex.Type] = tex
	}
}

func (tex *Texture) Unbind() {
	if boundTexture[tex.Type] == tex {
		tex.Texture.Unbind(tex.Type)
		boundTexture[tex.Type] = nil
	}
}

func (tex *Texture) SetFilters(minFilter int, magFilter int) {
	tex.Bind()
	PanicOnError()
	gl.TexParameteri(tex.Type, gl.TEXTURE_MIN_FILTER, minFilter)
	PanicOnError()
	gl.TexParameteri(tex.Type, gl.TEXTURE_MAG_FILTER, magFilter)
	PanicOnError()
}
