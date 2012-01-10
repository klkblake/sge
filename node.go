package sge

import "s3dm"

type Node interface {
	Xform() *s3dm.Xform
	Children() []Node
	Passes() int
	Add(child Node)
	Remove(child Node)
	UpdateAll(deltaNs int64)
	RenderAll(view *View, matrixStack *Mat4Stack, pass int)
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
	node.Init()
	return node
}

func (node *BasicNode) Init() {
	node.xform = s3dm.NewXform()
	node.children = make([]Node, 0)
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

func (node *BasicNode) UpdateAll(deltaNs int64) {
	node.Update(deltaNs)
	for _, child := range node.children {
		child.UpdateAll(deltaNs)
	}
}

func (node *BasicNode) RenderAll(view *View, matrixStack *Mat4Stack, pass int) {
	modelMatrix := Mat4(node.Xform().GetMatrix4())
	matrixStack.Push(&modelMatrix)
	if pass & node.Passes() != 0 {
		node.Render(view, matrixStack.Top(), pass)
	}
	for _, child := range node.children {
		child.RenderAll(view, matrixStack, pass)
	}
	matrixStack.Pop()
}

func (node *BasicNode) Update(deltaNs int64) {}
func (node *BasicNode) Render(view *View, mvpMatrix *Mat4, pass int) {}
