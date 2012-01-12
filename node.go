package sge

import "math"

import "s3dm"

type Node interface {
	Xform() *s3dm.Xform
	AABB() *s3dm.AABB
	Parent() Node
	Children() []Node
	Passes() int
	SetParent(parent Node)
	Add(child Node)
	Remove(child Node)
	UpdateAABB()
	Update(deltaNs int64)
	Render(view *View, mvpMatrix *Mat4, pass int)
}

type BasicNode struct {
	xform *s3dm.Xform
	aabb *s3dm.AABB
	parent Node
	children []Node
	passes int
}

func NewBasicNode(parent Node) *BasicNode {
	node := new(BasicNode)
	node.xform = s3dm.NewXform()
	node.aabb = s3dm.NewAABB(s3dm.NewV3(0, 0, 0), s3dm.NewV3(0, 0, 0))
	node.parent = parent
	node.children = make([]Node, 0)
	return node
}

func (node *BasicNode) Xform() *s3dm.Xform {
	return node.xform
}

func (node *BasicNode) AABB() *s3dm.AABB {
	return node.aabb
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
			node.children = append(node.children[:i], node.children[i+1:]...)
		}
	}
}

func (node *BasicNode) UpdateAABB() {
	if len(node.children) == 0 {
		node.aabb.Min = s3dm.NewV3(0, 0, 0)
		node.aabb.Max = s3dm.NewV3(0, 0, 0)
	} else {
		min := node.children[0].AABB().Min.Copy()
		max := node.children[0].AABB().Max.Copy()
		for i := 1; i < len(node.children); i++ {
			child := node.children[i].AABB()
			min.X = math.Fmin(min.X, child.Min.X)
			min.Y = math.Fmin(min.Y, child.Min.Y)
			min.Z = math.Fmin(min.Z, child.Min.Z)
			max.X = math.Fmax(max.X, child.Max.X)
			max.Y = math.Fmax(max.Y, child.Max.Y)
			max.Z = math.Fmax(max.Z, child.Max.Z)
		}
		node.aabb.Min = min
		node.aabb.Max = max
	}
	if node.parent != nil {
		node.parent.UpdateAABB()
	}
}

func (node *BasicNode) Update(deltaNs int64) {}
func (node *BasicNode) Render(view *View, mvpMatrix *Mat4, pass int) {}
