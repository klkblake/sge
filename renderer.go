package sge

type Renderer interface {
	Render(view *View, mvpMatrix *Mat4, pass int)
}
