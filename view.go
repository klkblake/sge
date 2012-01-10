package sge

import "math"
import "runtime"
import "os"

import "s3dm"
import "atom/sdl"
import "gl"

type View struct {
	Screen *sdl.Surface
	Camera *s3dm.Xform
	PerspectiveMatrix *Mat4
	ViewMatrix *Mat4
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
	gl.Init()
	//gl.Enable(gl.BLEND)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Viewport(0, 0, width, height)
	view.Camera = s3dm.NewXform()
	view.PerspectiveMatrix = calcPerspectiveMatrix(math.Pi/4, float64(width)/float64(height), near, far)
	return view
}

func calcPerspectiveMatrix(fovy float64, aspect float64, near float64, far float64) *Mat4 {
	top := near*math.Tan(fovy)
	right := aspect*top
	return &Mat4{
		near/right, 0, 0, 0,
		0, near/top, 0, 0,
		0, 0, -(far+near)/(far-near), -1,
		0, 0, -2*far*near/(far-near), 0}
}

func (view *View) SetBackgroundColor(red, green, blue float32) {
	gl.ClearColor(gl.GLclampf(red), gl.GLclampf(green), gl.GLclampf(blue), 1.0)
}

func (view *View) UpdateViewMatrix() {
	m := view.Camera.GetMatrix4()
	// Take the inverse of the camera matrix
	view.ViewMatrix = &Mat4{
		m[0], m[4], m[8], 0,
		m[1], m[5], m[9], 0,
		m[2], m[6], m[10], 0,
		-(m[0]*m[12]+m[1]*m[13]+m[2]*m[14]),
		-(m[4]*m[12]+m[5]*m[13]+m[6]*m[14]),
		-(m[8]*m[12]+m[9]*m[13]+m[10]*m[14]),
		1}
}

func (view *View) Close() os.Error {
	sdl.Quit()
	return nil
}
