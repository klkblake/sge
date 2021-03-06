package sge

import (
	"unsafe"
)

import (
	"github.com/chsc/gogl/gl33"
	"github.com/klkblake/Go-SDL/sdl"
)

var activeTexture int
var boundTexture2D *Texture
var boundTexture2DArray *Texture
var boundTextureCubeMap *Texture
var boundTextureBuffer *Texture

type Texture struct {
	Id     gl33.Uint
	Type   gl33.Enum
	Width  int
	Height int
}

func NewTexture2D() *Texture {
	tex := new(Texture)
	gl33.GenTextures(1, &tex.Id)
	tex.Type = gl33.TEXTURE_2D
	return tex
}

func NewTextureArray() *Texture {
	tex := new(Texture)
	gl33.GenTextures(1, &tex.Id)
	tex.Type = gl33.TEXTURE_2D_ARRAY
	return tex
}

func NewTextureCubeMap() *Texture {
	tex := new(Texture)
	gl33.GenTextures(1, &tex.Id)
	tex.Type = gl33.TEXTURE_CUBE_MAP
	return tex
}

func LoadTexture2D(filename string, minFilter int, magFilter int) *Texture {
	ch := make(chan *Texture)
	GL <- func() {
		tex := NewTexture2D()
		tex.SetFilters(minFilter, magFilter)
		tex.Bind(0)
		surface := sdl.Load(filename)
		if surface == nil {
			panic(sdl.GetError())
		}
		defer surface.Free()
		tex.Width = int(surface.W)
		tex.Height = int(surface.H)
		uploadSurface(gl33.TEXTURE_2D, surface)
		gl33.GenerateMipmap(gl33.TEXTURE_2D)
		ch <- tex
	}
	return <-ch
}

func LoadTextureArray(filenames []string, minFilter int, magFilter int) *Texture {
	ch := make(chan *Texture)
	GL <- func() {
		tex := NewTextureArray()
		tex.SetFilters(minFilter, magFilter)
		tex.Bind(0)
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
		ch <- tex
	}
	return <-ch
}

func LoadTextureCubeMap(filenames *[6]string, minFilter int, magFilter int) *Texture {
	ch := make(chan *Texture)
	GL <- func() {
		tex := NewTextureCubeMap()
		tex.SetFilters(minFilter, magFilter)
		tex.Bind(0)
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
		ch <- tex
	}
	return <-ch
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
	if boundTexture2D != nil {
		boundTexture2D.Unbind()
	}
}

func UnbindTextureCubeMap() {
	if boundTextureCubeMap != nil {
		boundTextureCubeMap.Unbind()
	}
}

func (tex *Texture) bound() **Texture {
	switch tex.Type {
	case gl33.TEXTURE_2D:
		return &boundTexture2D
	case gl33.TEXTURE_2D_ARRAY:
		return &boundTexture2DArray
	case gl33.TEXTURE_CUBE_MAP:
		return &boundTextureCubeMap
	case gl33.TEXTURE_BUFFER:
		return &boundTextureBuffer
	}
	panic("Unsupported texture type")
}

func (tex *Texture) Bind(textureUnit int) {
	if activeTexture != textureUnit {
		gl33.ActiveTexture(gl33.Enum(gl33.TEXTURE0 + textureUnit))
		activeTexture = textureUnit
	}
	bound := tex.bound()
	if *bound != tex {
		gl33.BindTexture(tex.Type, tex.Id)
		*bound = tex
	}
}

func (tex *Texture) Unbind() {
	bound := tex.bound()
	if *bound == tex {
		gl33.BindTexture(tex.Type, 0)
		*bound = nil
	}
}

func (tex *Texture) Delete() {
	gl33.DeleteTextures(1, &tex.Id)
}

func (tex *Texture) SetFilters(minFilter int, magFilter int) {
	tex.Bind(0)
	gl33.TexParameteri(tex.Type, gl33.TEXTURE_MIN_FILTER, gl33.Int(minFilter))
	gl33.TexParameteri(tex.Type, gl33.TEXTURE_MAG_FILTER, gl33.Int(magFilter))
}
