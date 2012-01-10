package sge

import "s3dm"

type Node struct {
	s3dm.Xform
	Children []*Node
	Passes int
}

func NewNode() *Node {
	node := new(Node)
	node.ResetXform()
	node.Children = make([]*Node, 0)
	return node
}

func (node *Node) Add(child *Node) {
	node.Children = append(node.Children, child)
}

func (node *Node) Remove(child *Node) {
	for i, n := range node.Children {
		if n == child {
			node.Children = append(node.Children[:i], node.Children[i+1:]...)
		}
	}
}

func (node *Node) UpdateAll(deltaNs int64) {
	node.Update(deltaNs)
	for _, child := range node.Children {
		child.UpdateAll(deltaNs)
	}
}

func (node *Node) RenderAll(view *View, matrixStack *Mat4Stack, pass int) {
	modelMatrix := Mat4(node.Xform.GetMatrix4())
	matrixStack.Push(&modelMatrix)
	if pass & node.Passes != 0 {
		node.Render(view, matrixStack.Top(), pass)
	}
	for _, child := range node.Children {
		child.RenderAll(view, matrixStack, pass)
	}
	matrixStack.Pop()
}

func (node *Node) Update(deltaNs int64) {}
func (node *Node) Render(view *View, mvpMatrix *Mat4, pass int) {}
