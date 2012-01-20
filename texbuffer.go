package sge

import (
	"github.com/chsc/gogl/gl33"
)

type TextureBuffer struct {
	*Texture
	Buffer *Buffer
}

func NewTextureBuffer(data interface{}) *TextureBuffer {
	texbuf := new(TextureBuffer)
	texbuf.Texture = new(Texture)
	gl33.GenTextures(1, &texbuf.Texture.Id)
	texbuf.Texture.Type = gl33.TEXTURE_BUFFER
	texbuf.Buffer = NewBuffer(gl33.TEXTURE_BUFFER, gl33.DYNAMIC_DRAW, data)
	texbuf.Width = texbuf.Buffer.Value.Len()
	var format gl33.Enum
	switch data.(type) {
	case []byte:
		format = gl33.R8
	case []uint16:
		format = gl33.R16
	case [][2]byte:
		format = gl33.RG8
	case [][2]uint16:
		format = gl33.RG16
	case [][4]byte:
		format = gl33.RGBA8
	case [][4]uint16:
		format = gl33.RGBA16
	}
	texbuf.Bind(0)
	gl33.TexBuffer(gl33.TEXTURE_BUFFER, format, texbuf.Buffer.Id)
	return texbuf
}

func (texbuf *TextureBuffer) Update() {
	texbuf.Buffer.Update()
}

func (texbuf *TextureBuffer) Delete() {
	texbuf.Texture.Delete()
	texbuf.Buffer.Delete()
}
