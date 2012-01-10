package sge

import "gl"
import "atom/sdl"

const (
	PassOpaque = 1 << iota
	PassTranslucent
)

type World struct {
	Root *Node
	Skybox Renderer
	Gui *Node
	matrixStack *Mat4Stack
}

func NewWorld() *World {
	world := new(World)
	world.Root = NewNode()
	world.Gui = NewNode()
	world.matrixStack = NewMat4Stack()
	return world
}

func (world *World) Update(deltaNs int64) {
	world.Root.UpdateAll(deltaNs)
}

func (world *World) Render(view *View) {
	view.UpdateViewMatrix()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	world.matrixStack.Push(view.PerspectiveMatrix)
	world.matrixStack.Push(view.ViewMatrix)
	world.Root.RenderAll(view, world.matrixStack, PassOpaque)
	if world.Skybox != nil {
		world.Skybox.Render(view, world.matrixStack.Top(), PassOpaque)
	}
	world.Root.RenderAll(view, world.matrixStack, PassTranslucent)
	world.matrixStack.Pop()
	world.matrixStack.Pop()
	// XXX Setup orthographic projection.
	world.Gui.RenderAll(view, world.matrixStack, PassOpaque)
	world.Gui.RenderAll(view, world.matrixStack, PassTranslucent)
	sdl.GL_SwapBuffers()
}
