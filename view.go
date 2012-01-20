package sge

import (
	"math"
	"runtime"
)

import (
	"atom/sdl"
	"github.com/chsc/gogl/gl33"
	"s3dm"
)

type View struct {
	Screen             *sdl.Surface
	Camera             *s3dm.Frustum
	Width              int
	Height             int
	PerspectiveMatrix  *s3dm.Mat4
	OrthographicMatrix *s3dm.Mat4
	ViewMatrix         *s3dm.Mat4
}

func NewView(title string, width int, height int, near float64, far float64) *View {
	view := new(View)
	runtime.LockOSThread()
	if sdl.Init(sdl.INIT_VIDEO) < 0 {
		panic(sdl.GetError())
	}
	sdl.GL_SetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GL_SetAttribute(sdl.GL_SWAP_CONTROL, 1)
	view.Screen = sdl.SetVideoMode(width, height, 32, sdl.OPENGL)
	if view.Screen == nil {
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
	gl33.Viewport(0, 0, gl33.Sizei(width), gl33.Sizei(height))
	fovy := math.Pi / 4
	aspect := float64(width) / float64(height)
	view.Camera = s3dm.NewFrustum(near, far, fovy, aspect)
	view.Update()
	view.Width = width
	view.Height = height
	view.PerspectiveMatrix = s3dm.NewPerspectiveMat4(fovy, aspect, near, far)
	view.OrthographicMatrix = s3dm.NewOrthographicMat4(float64(width), float64(height), 0, 1)
	return view
}

func (view *View) SetBackgroundColor(red, green, blue float32) {
	gl33.ClearColor(gl33.Clampf(red), gl33.Clampf(green), gl33.Clampf(blue), 1.0)
}

func (view *View) Update() {
	view.Camera.Update()
	m := view.Camera.GetMatrix4()
	// Take the inverse of the camera matrix
	view.ViewMatrix = &s3dm.Mat4{
		m[0], m[4], m[8], 0,
		m[1], m[5], m[9], 0,
		m[2], m[6], m[10], 0,
		-(m[0]*m[12] + m[1]*m[13] + m[2]*m[14]),
		-(m[4]*m[12] + m[5]*m[13] + m[6]*m[14]),
		-(m[8]*m[12] + m[9]*m[13] + m[10]*m[14]),
		1}
}

func (view *View) Close() error {
	sdl.Quit()
	return nil
}
