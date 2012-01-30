package sge

import (
	"atom/sdl"
	"github.com/chsc/gogl/gl33"
	"s3dm"
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

func (world *World) Update(delta float64) {
	update(world.Root, delta)
	update(world.Gui, delta)
}

func update(node Node, delta float64) {
	node.Update(delta)
	for _, child := range node.Children() {
		update(child, delta)
	}
}

func (world *World) Render(view *View) {
	GL <- func() {
		gl33.Clear(gl33.COLOR_BUFFER_BIT | gl33.DEPTH_BUFFER_BIT)
	}
	vpMatrix := view.PerspectiveMatrix.Mul(view.ViewMatrix)
	render(world.Root, view, vpMatrix, world.matrixStack, PassOpaque, true)
	if world.Skybox != nil {
		world.Skybox.Render(view, vpMatrix.Mul(world.matrixStack.Top()), PassOpaque)
	}
	render(world.Root, view, vpMatrix, world.matrixStack, PassTranslucent, true)
	vpMatrix = view.OrthographicMatrix
	render(world.Gui, view, vpMatrix, world.matrixStack, PassOpaque, false)
	render(world.Gui, view, vpMatrix, world.matrixStack, PassTranslucent, false)
	GL <- func() {
		sdl.GL_SwapBuffers()
	}
}

func render(node Node, view *View, vpMatrix *s3dm.Mat4, matrixStack *Mat4Stack, pass int, frustumCull bool) {
	modelMatrix := s3dm.Mat4(node.Xform().GetMatrix4())
	matrixStack.Push(&modelMatrix)
	defer matrixStack.Pop()
	if frustumCull {
		aabb := node.AABB()
		aabb.MoveGlobal(matrixStack.Top().Position())
		if view.Camera.IntersectsAABB(aabb) < 0 {
			return
		}
	}
	if pass&node.Passes() != 0 {
		node.Render(view, vpMatrix.Mul(matrixStack.Top()), pass)
	}
	for _, child := range node.Children() {
		render(child, view, vpMatrix, matrixStack, pass, frustumCull)
	}
}
