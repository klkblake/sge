package sge

import (
	"github.com/chsc/gogl/gl33"
	"github.com/klkblake/Go-SDL/sdl"
	"github.com/klkblake/s3dm"
)

const (
	PassOpaque = 1 << iota
	PassTranslucent
)

type World struct {
	Root   *XformNode
	Skybox Renderer
	Gui    *XformNode
}

func NewWorld() *World {
	world := new(World)
	world.Root = NewXformNode()
	world.Gui = NewXformNode()
	return world
}

func (world *World) Update(delta float64) {
	f := func(leaf Leaf) {
		leaf.Update(delta)
	}
	world.Root.Walk(f)
	world.Gui.Walk(f)
}

// TODO this should walk a seperate visibility structure.
func (world *World) Render(view *View) {
	GL <- func() {
		gl33.Clear(gl33.COLOR_BUFFER_BIT | gl33.DEPTH_BUFFER_BIT)
	}
	vpMatrix := view.PerspectiveMatrix.Mul(view.ViewMatrix)
	pass := PassOpaque
	frustumCull := true
	render := func(leaf Leaf) {
		if pass&leaf.Passes() == 0 {
			return
		}
		modelMatrix := s3dm.Mat4(leaf.XformNode().WorldXform.GetMatrix4())
		mvpMatrix := vpMatrix.Mul(&modelMatrix)
		if frustumCull {
			aabb := leaf.AABB().MoveGlobal(modelMatrix.Position())
			if view.Camera.IntersectsAABB(aabb) < 0 {
				return
			}
		}
		leaf.Render(view, mvpMatrix, pass)
	}
	world.Root.Walk(render)
	pass = PassTranslucent
	world.Root.Walk(render)
	if world.Skybox != nil {
		world.Skybox.Render(view, vpMatrix, PassOpaque)
	}
	vpMatrix = view.OrthographicMatrix
	frustumCull = false
	pass = PassOpaque
	world.Gui.Walk(render)
	pass = PassTranslucent
	world.Gui.Walk(render)
	GL <- func() {
		sdl.GL_SwapBuffers()
	}
}
