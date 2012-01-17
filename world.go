package sge

import (
	"atom/sdl"
	"github.com/chsc/gogl/gl33"
)

const (
	PassOpaque = 1 << iota
	PassTranslucent
)

type World struct {
	Root        Node
	Skybox      Renderer
	Gui         Node
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
	gl33.Clear(gl33.COLOR_BUFFER_BIT | gl33.DEPTH_BUFFER_BIT)
	vpMatrix := view.PerspectiveMatrix.Mul(view.ViewMatrix)
	render(world.Root, view, vpMatrix, world.matrixStack, PassOpaque)
	if world.Skybox != nil {
		world.Skybox.Render(view, vpMatrix.Mul(world.matrixStack.Top()), PassOpaque)
	}
	render(world.Root, view, vpMatrix, world.matrixStack, PassTranslucent)
	// XXX Setup orthographic projection.
	render(world.Gui, view, vpMatrix, world.matrixStack, PassOpaque)
	render(world.Gui, view, vpMatrix, world.matrixStack, PassTranslucent)
	sdl.GL_SwapBuffers()
}

func render(node Node, view *View, vpMatrix *Mat4, matrixStack *Mat4Stack, pass int) {
	modelMatrix := Mat4(node.Xform().GetMatrix4())
	matrixStack.Push(&modelMatrix)
	defer matrixStack.Pop()
	if view.Camera.IntersectsAABB(node.AABB().MoveGlobal(matrixStack.Top().Position())) < 0 {
		return
	}
	if pass&node.Passes() != 0 {
		node.Render(view, vpMatrix.Mul(matrixStack.Top()), pass)
	}
	for _, child := range node.Children() {
		render(child, view, vpMatrix, matrixStack, pass)
	}
}
