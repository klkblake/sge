package sge

import (
	"github.com/klkblake/s3dm"
)

type Renderer interface {
	Render(view *View, mvpMatrix *s3dm.Mat4, pass int)
}
