package sge

import "math"

import "s3dm"

type Node interface {
	Xform() *s3dm.Xform
	AABB() s3dm.AABB
	Parent() Node
	Children() []Node
	Passes() int
	SetParent(parent Node)
	Add(child Node)
	Remove(child Node)
	UpdateAABB()
	Update(delta float64)
	Render(view *View, mvpMatrix *s3dm.Mat4, pass int)
}

type BasicNode struct {
	xform    *s3dm.Xform
	aabb     s3dm.AABB
	parent   Node
	children []Node
	passes   int
}

func NewBasicNode(parent Node) *BasicNode {
	node := new(BasicNode)
	node.xform = s3dm.NewXform()
	node.parent = parent
	node.children = make([]Node, 0)
	return node
}

func (node *BasicNode) Xform() *s3dm.Xform {
	return node.xform
}

func (node *BasicNode) AABB() s3dm.AABB {
	return node.aabb
}

func (node *BasicNode) SetAABB(aabb s3dm.AABB) {
	node.aabb = aabb
}

func (node *BasicNode) Parent() Node {
	return node.parent
}

func (node *BasicNode) Children() []Node {
	return node.children
}

func (node *BasicNode) Passes() int {
	return node.passes
}

func (node *BasicNode) SetPasses(passes int) {
	node.passes = passes
}

func (node *BasicNode) SetParent(parent Node) {
	node.parent = parent
}

func (node *BasicNode) Add(child Node) {
	node.children = append(node.children, child)
}

func (node *BasicNode) Remove(child Node) {
	for i, n := range node.children {
		if n == child {
			if i < len(node.children) - 2 {
				node.children = append(node.children[:i], node.children[i+1:]...)
			} else {
				node.children = node.children[:i]
			}
		}
	}
}

func (node *BasicNode) UpdateAABB() {
	if len(node.children) == 0 {
		node.aabb.Min = s3dm.V3{}
		node.aabb.Max = s3dm.V3{}
	} else {
		pos := node.children[0].Xform().Position
		min := pos.Add(node.children[0].AABB().Min)
		max := pos.Add(node.children[0].AABB().Max)
		for i := 1; i < len(node.children); i++ {
			child := node.children[i].AABB()
			pos = node.children[i].Xform().Position
			min.X = math.Min(min.X, child.Min.X+pos.X)
			min.Y = math.Min(min.Y, child.Min.Y+pos.Y)
			min.Z = math.Min(min.Z, child.Min.Z+pos.Z)
			max.X = math.Max(max.X, child.Max.X+pos.X)
			max.Y = math.Max(max.Y, child.Max.Y+pos.Y)
			max.Z = math.Max(max.Z, child.Max.Z+pos.Z)
		}
		node.aabb.Min = min
		node.aabb.Max = max
	}
	if node.parent != nil {
		node.parent.UpdateAABB()
	}
}

func (node *BasicNode) Update(delta float64)                              {}
func (node *BasicNode) Render(view *View, mvpMatrix *s3dm.Mat4, pass int) {}
