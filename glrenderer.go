package sge

import (
	"github.com/chsc/gogl/gl33"
	"github.com/klkblake/Go-SDL/sdl"
	"github.com/klkblake/s3dm"
)

type GLRenderer struct {
	View               *View
	Screen             *sdl.Surface
	PerspectiveMatrix  s3dm.Mat4
	OrthographicMatrix s3dm.Mat4
	ViewMatrix         s3dm.Mat4
}

func NewGLRenderer(view *View, title string) *GLRenderer {
	r := new(GLRenderer)
	r.View = view
	GL <- func() {
		if sdl.Init(sdl.INIT_VIDEO) < 0 {
			panic(sdl.GetError())
		}
		sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
		sdl.GL_SetAttribute(sdl.GL_SWAP_CONTROL, 1)
		r.Screen = sdl.SetVideoMode(view.Width, view.Height, 32, sdl.OPENGL)
		if r.Screen == nil {
			panic(sdl.GetError())
		}
		sdl.WM_SetCaption(title, title)
		err := gl33.Init()
		if err != nil {
			panic(err)
		}
		gl33.Enable(gl33.BLEND)
		gl33.Enable(gl33.CULL_FACE)
		gl33.Enable(gl33.DEPTH_TEST)
		gl33.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl33.BlendFunc(gl33.SRC_ALPHA, gl33.ONE_MINUS_SRC_ALPHA)
		gl33.DepthFunc(gl33.LEQUAL)
		gl33.Viewport(0, 0, gl33.Sizei(view.Width), gl33.Sizei(view.Height))
	}
	r.PerspectiveMatrix = s3dm.PerspectiveMatrix(view.Fovy, view.Aspect, view.Near, view.Far)
	r.OrthographicMatrix = s3dm.OrthographicMatrix(float64(view.Width), float64(view.Height), 0, 1)
	return r
}

func (r *GLRenderer) SetBackgroundColor(red, green, blue float32) {
	GL <- func() {
		gl33.ClearColor(gl33.Clampf(red), gl33.Clampf(green), gl33.Clampf(blue), 1.0)
	}
}

func (r *GLRenderer) UpdatePerspective() {
	r.PerspectiveMatrix = s3dm.PerspectiveMatrix(r.View.Fovy, r.View.Aspect, r.View.Near, r.View.Far)
}

// TODO this should walk a seperate visibility structure.
func (r *GLRenderer) Render(world *World) {
	r.View.Update()
	m := r.View.Matrix(r.View.Position)
	// Take the inverse of the camera matrix
	r.ViewMatrix = s3dm.Mat4{
		m[0], m[4], m[8], 0,
		m[1], m[5], m[9], 0,
		m[2], m[6], m[10], 0,
		-(m[0]*m[12] + m[1]*m[13] + m[2]*m[14]),
		-(m[4]*m[12] + m[5]*m[13] + m[6]*m[14]),
		-(m[8]*m[12] + m[9]*m[13] + m[10]*m[14]),
		1}
	GL <- func() {
		gl33.Clear(gl33.COLOR_BUFFER_BIT | gl33.DEPTH_BUFFER_BIT)
		vpMatrix := r.PerspectiveMatrix.Mul(r.ViewMatrix)
		pass := PassOpaque
		hud := false
		render := func(leaf Leaf) {
			if pass&leaf.Passes() == 0 {
				return
			}
			var modelMatrix s3dm.Mat4
			if hud {
				modelMatrix = leaf.XformNode().WorldXform.Matrix(s3dm.Position{})
			} else {
				modelMatrix = leaf.XformNode().WorldXform.Matrix(r.View.Position)
			}
			mvpMatrix := vpMatrix.Mul(modelMatrix)
			if !hud {
				if leaf.AABB().IntersectsFrustum(r.View.Frustum) < 0 {
					return
				}
			}
			leaf.Render(r.View, mvpMatrix, pass)
		}
		world.Root.Walk(render)
		pass = PassTranslucent
		world.Root.Walk(render)
		if world.Skybox != nil {
			world.Skybox.Render(r.View, vpMatrix, PassOpaque)
		}
		vpMatrix = r.OrthographicMatrix
		hud = true
		pass = PassOpaque
		world.Gui.Walk(render)
		pass = PassTranslucent
		world.Gui.Walk(render)
		sdl.GL_SwapBuffers()
	}
}

func (r *GLRenderer) Close() error {
	GL <- func() {
		sdl.Quit()
	}
	FlushGL()
	return nil
}
