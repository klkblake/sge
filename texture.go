package sge

import (
	"unsafe"
)

import (
	"atom/sdl"
	"github.com/chsc/gogl/gl33"
)

var activeTexture int
var boundTexture map[gl33.Enum]*Texture

func init() {
	boundTexture = make(map[gl33.Enum]*Texture)
}

type Texture struct {
	Id     gl33.Uint
	Type   gl33.Enum
	Width  int
	Height int
}

func NewTexture2D() *Texture {
	tex := new(Texture)
	GL <- func() {
		gl33.GenTextures(1, &tex.Id)
	}
	tex.Type = gl33.TEXTURE_2D
	return tex
}

func NewTextureArray() *Texture {
	tex := new(Texture)
	GL <- func() {
		gl33.GenTextures(1, &tex.Id)
	}
	tex.Type = gl33.TEXTURE_2D_ARRAY
	return tex
}

func NewTextureCubeMap() *Texture {
	tex := new(Texture)
	GL <- func() {
		gl33.GenTextures(1, &tex.Id)
	}
	tex.Type = gl33.TEXTURE_CUBE_MAP
	return tex
}

func LoadTexture2D(filename string, minFilter int, magFilter int) *Texture {
	tex := NewTexture2D()
	tex.SetFilters(minFilter, magFilter)
	tex.Bind(0)
	GL <- func() {
		surface := sdl.Load(filename)
		if surface == nil {
			panic(sdl.GetError())
		}
		defer surface.Free()
		tex.Width = int(surface.W)
		tex.Height = int(surface.H)
		uploadSurface(gl33.TEXTURE_2D, surface)
		gl33.GenerateMipmap(gl33.TEXTURE_2D)
	}
	return tex
}

func LoadTextureArray(filenames []string, minFilter int, magFilter int) *Texture {
	tex := NewTextureArray()
	tex.SetFilters(minFilter, magFilter)
	tex.Bind(0)
	GL <- func() {
		surfaces := make([]*sdl.Surface, len(filenames))
		for i, filename := range filenames {
			surfaces[i] = sdl.Load(filename)
			if surfaces[i] == nil {
				panic(sdl.GetError())
			}
			defer surfaces[i].Free()
		}
		tex.Width = int(surfaces[0].W)
		tex.Height = int(surfaces[0].H)
		var internalFormat gl33.Int
		var format gl33.Enum
		var size int
		if surfaces[0].Format.BitsPerPixel == 32 {
			internalFormat = gl33.RGBA8
			format = gl33.RGBA
			size = 4
		} else {
			internalFormat = gl33.RGB8
			format = gl33.RGB
			size = 3
		}
		pixels := make([]byte, tex.Width*tex.Height*len(surfaces)*size)
		for i, surface := range surfaces {
			p := uintptr(surface.Pixels)
			for j := 0; j < tex.Width*tex.Height*size; j++ {
				pixels[i*tex.Width*tex.Height*size+j] = *(*byte)(unsafe.Pointer(p + uintptr(j)))
			}
		}
		gl33.TexImage3D(gl33.TEXTURE_2D_ARRAY, 0, internalFormat, gl33.Sizei(tex.Width), gl33.Sizei(tex.Height), gl33.Sizei(len(surfaces)), 0, format, gl33.UNSIGNED_BYTE, gl33.Pointer(&pixels[0]))
		gl33.GenerateMipmap(gl33.TEXTURE_2D_ARRAY)
	}
	return tex
}

func LoadTextureCubeMap(filenames *[6]string, minFilter int, magFilter int) *Texture {
	tex := NewTextureCubeMap()
	tex.SetFilters(minFilter, magFilter)
	tex.Bind(0)
	GL <- func() {
		var surfaces [6]*sdl.Surface
		for i, filename := range filenames {
			surfaces[i] = sdl.Load(filename)
			if surfaces[i] == nil {
				panic(sdl.GetError())
			}
			defer surfaces[i].Free()
		}
		tex.Width = int(surfaces[0].W)
		tex.Height = int(surfaces[0].H)
		uploadSurface(gl33.TEXTURE_CUBE_MAP_POSITIVE_X, surfaces[0])
		uploadSurface(gl33.TEXTURE_CUBE_MAP_NEGATIVE_X, surfaces[1])
		uploadSurface(gl33.TEXTURE_CUBE_MAP_POSITIVE_Y, surfaces[2])
		uploadSurface(gl33.TEXTURE_CUBE_MAP_NEGATIVE_Y, surfaces[3])
		uploadSurface(gl33.TEXTURE_CUBE_MAP_POSITIVE_Z, surfaces[4])
		uploadSurface(gl33.TEXTURE_CUBE_MAP_NEGATIVE_Z, surfaces[5])
		gl33.GenerateMipmap(gl33.TEXTURE_CUBE_MAP)
	}
	return tex
}

func uploadSurface(target gl33.Enum, surface *sdl.Surface) {
	if target == gl33.TEXTURE_CUBE_MAP && surface.W != surface.H {
		panic("Non-square texture in cube map")
	}
	var internalFormat gl33.Int
	var format gl33.Enum
	if surface.Format.BitsPerPixel == 32 {
		internalFormat = gl33.RGBA8
		format = gl33.RGBA
	} else {
		internalFormat = gl33.RGB8
		format = gl33.RGB
	}
	gl33.TexImage2D(target, 0, internalFormat, gl33.Sizei(surface.W), gl33.Sizei(surface.H), 0, format, gl33.UNSIGNED_BYTE, gl33.Pointer(surface.Pixels))
}

func UnbindTexture2D() {
	if boundTexture[gl33.TEXTURE_2D] != nil {
		boundTexture[gl33.TEXTURE_2D].Unbind()
	}
}

func UnbindTextureCubeMap() {
	if boundTexture[gl33.TEXTURE_CUBE_MAP] != nil {
		boundTexture[gl33.TEXTURE_CUBE_MAP].Unbind()
	}
}

func (tex *Texture) Bind(textureUnit int) {
	if activeTexture != textureUnit {
		GL <- func() {
			gl33.ActiveTexture(gl33.Enum(gl33.TEXTURE0 + textureUnit))
		}
		activeTexture = textureUnit
	}
	if boundTexture[tex.Type] != tex {
		GL <- func() {
			gl33.BindTexture(tex.Type, tex.Id)
		}
		boundTexture[tex.Type] = tex
	}
}

func (tex *Texture) Unbind() {
	if boundTexture[tex.Type] == tex {
		GL <- func() {
			gl33.BindTexture(tex.Type, 0)
		}
		boundTexture[tex.Type] = nil
	}
}

func (tex *Texture) Delete() {
	GL <- func() {
		gl33.DeleteTextures(1, &tex.Id)
	}
}

func (tex *Texture) SetFilters(minFilter int, magFilter int) {
	tex.Bind(0)
	GL <- func() {
		gl33.TexParameteri(tex.Type, gl33.TEXTURE_MIN_FILTER, gl33.Int(minFilter))
		gl33.TexParameteri(tex.Type, gl33.TEXTURE_MAG_FILTER, gl33.Int(magFilter))
	}
}
