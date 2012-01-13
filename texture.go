package sge

import (
	"unsafe"
)

import (
	"gl"
	"atom/sdl"
)

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

func NewTextureArray() *Texture {
	tex := new(Texture)
	tex.Texture = gl.GenTexture()
	tex.Type = gl.TEXTURE_2D_ARRAY
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
	if surface == nil {
		panic(sdl.GetError())
	}
	defer surface.Free()
	tex := NewTexture2D()
	tex.Width = int(surface.W)
	tex.Height = int(surface.H)
	tex.SetFilters(minFilter, magFilter)
	tex.Bind()
	gl.TexParameteri(gl.TEXTURE_2D, gl.GENERATE_MIPMAP, gl.TRUE)
	uploadSurface(gl.TEXTURE_2D, surface)
	return tex
}

func LoadTextureArray(filenames []string, minFilter int, magFilter int) *Texture {
	surfaces := make([]*sdl.Surface, len(filenames))
	for i, filename := range filenames {
		surfaces[i] = sdl.Load(filename)
		if surfaces[i] == nil {
			panic(sdl.GetError())
		}
		defer surfaces[i].Free()
	}
	tex := NewTextureArray()
	tex.Width = int(surfaces[0].W)
	tex.Height = int(surfaces[0].H)
	tex.SetFilters(minFilter, magFilter)
	tex.Bind()
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.GENERATE_MIPMAP, gl.TRUE)
	var internalFormat int
	var format gl.GLenum
	var size int
	if surfaces[0].Format.BitsPerPixel == 32 {
		internalFormat = gl.RGBA8
		format = gl.RGBA
		size = 4
	} else {
		internalFormat = gl.RGB8
		format = gl.RGB
		size = 3
	}
	pixels := make([]byte, tex.Width*tex.Height*len(surfaces)*size)
	for i, surface := range surfaces {
		p := uintptr(surface.Pixels)
		for j := 0; j < tex.Width*tex.Height*size; j++ {
			pixels[i*tex.Width*tex.Height*size + j] = *(*byte)(unsafe.Pointer(p+uintptr(j)))
		}
	}
	gl.TexImage3D(gl.TEXTURE_2D_ARRAY, 0, internalFormat, tex.Width, tex.Height, len(surfaces), 0, format, gl.UNSIGNED_BYTE, pixels)
	return tex
}

func LoadTextureCubeMap(filenames *[6]string, minFilter int, magFilter int) *Texture {
	var surfaces [6]*sdl.Surface
	for i, filename := range filenames {
		surfaces[i] = sdl.Load(filename)
		if surfaces[i] == nil {
			panic(sdl.GetError())
		}
		defer surfaces[i].Free()
	}
	tex := NewTextureCubeMap()
	tex.Width = int(surfaces[0].W)
	tex.Height = int(surfaces[0].H)
	tex.SetFilters(minFilter, magFilter)
	tex.Bind()
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.GENERATE_MIPMAP, gl.TRUE)
	uploadSurface(gl.TEXTURE_CUBE_MAP_POSITIVE_X, surfaces[0])
	uploadSurface(gl.TEXTURE_CUBE_MAP_NEGATIVE_X, surfaces[1])
	uploadSurface(gl.TEXTURE_CUBE_MAP_POSITIVE_Y, surfaces[2])
	uploadSurface(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, surfaces[3])
	uploadSurface(gl.TEXTURE_CUBE_MAP_POSITIVE_Z, surfaces[4])
	uploadSurface(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, surfaces[5])
	return tex
}

func uploadSurface(target gl.GLenum, surface *sdl.Surface) {
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
	gl.TexImage2D(target, 0, internalFormat, int(surface.W), int(surface.H), 0, format, gl.UNSIGNED_BYTE, (*byte)(surface.Pixels))
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
