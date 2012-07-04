package sge

const (
	PassOpaque = 1 << iota
	PassTranslucent
)

type World struct {
	Root   *XformNode
	Skybox Leaf
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
	world.Skybox.Update(delta)
	world.Gui.Walk(f)
}
