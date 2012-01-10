package sge

import "s3dm"

type Node interface {
	Xform() *s3dm.Xform
	Children() []Node
	Passes() int
	Add(child Node)
	Remove(child Node)
	Update(deltaNs int64)
	Render(view *View, mvpMatrix *Mat4, pass int)
}

type BasicNode struct {
	xform *s3dm.Xform
	children []Node
	passes int
}

func NewBasicNode() *BasicNode {
	node := new(BasicNode)
	node.xform = s3dm.NewXform()
	node.children = make([]Node, 0)
	return node
}

func (node *BasicNode) Xform() *s3dm.Xform {
	return node.xform
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

func (node *BasicNode) Update(deltaNs int64) {}
func (node *BasicNode) Render(view *View, mvpMatrix *Mat4, pass int) {}
