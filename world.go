package sge

import "gl"
import "atom/sdl"

const (
	PassOpaque = 1 << iota
	PassTranslucent
)

type World struct {
	Root Node
	Skybox Renderer
	Gui Node
	matrixStack *Mat4Stack
}

func NewWorld() *World {
	world := new(World)
	root := NewBasicNode(nil)
	root.SetPasses(-1)
	world.Root = root
	gui := NewBasicNode(nil)
	gui.SetPasses(-1)
	world.Gui = gui
	world.matrixStack = NewMat4Stack()
	return world
}

func (world *World) Update(deltaNs int64) {
	update(world.Root, deltaNs)
}

func update(node Node, deltaNs int64) {
	node.Update(deltaNs)
	for _, child := range node.Children() {
		update(child, deltaNs)
	}
}

func (world *World) Render(view *View) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	world.matrixStack.Push(view.PerspectiveMatrix)
	world.matrixStack.Push(view.ViewMatrix)
	render(world.Root, view, world.matrixStack, PassOpaque)
	if world.Skybox != nil {
		world.Skybox.Render(view, world.matrixStack.Top(), PassOpaque)
	}
	render(world.Root, view, world.matrixStack, PassTranslucent)
	world.matrixStack.Pop()
	world.matrixStack.Pop()
	// XXX Setup orthographic projection.
	render(world.Gui, view, world.matrixStack, PassOpaque)
	render(world.Gui, view, world.matrixStack, PassTranslucent)
	sdl.GL_SwapBuffers()
}

func render(node Node, view *View, matrixStack *Mat4Stack, pass int) {
	modelMatrix := Mat4(node.Xform().GetMatrix4())
	matrixStack.Push(&modelMatrix)
	defer matrixStack.Pop()
	if view.Camera.IntersectsAABB(node.AABB()) < 0 {
		return
	}
	if pass & node.Passes() != 0 {
		node.Render(view, matrixStack.Top(), pass)
	}
	for _, child := range node.Children() {
		render(child, view, matrixStack, pass)
	}
}
