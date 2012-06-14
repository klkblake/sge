package sge

import "github.com/klkblake/s3dm"

type XformNode struct {
	Parent     *XformNode
	Children   []*XformNode
	Xform      *s3dm.Xform
	WorldXform *s3dm.Xform
	Leaf       Leaf
}

func NewXformNode() *XformNode {
	node := new(XformNode)
	node.Xform = s3dm.NewXform()
	node.WorldXform = s3dm.NewXform()
	return node
}

func (node *XformNode) Add(child *XformNode) {
	if child.Parent != nil {
		child.Parent.Remove(child)
	}
	node.Children = append(node.Children, child)
	child.Parent = node
}

func (node *XformNode) Remove(child *XformNode) {
	for i, n := range node.Children {
		if n == child {
			node.Children = append(node.Children[:i], node.Children[i+1:]...)
			break
		}
	}
	if len(node.Children) == 0 {
		node.Children = nil
	}
	if child.Parent == node {
		child.Parent = nil
	}
}

func (node *XformNode) Attach(leaf Leaf) {
	node.Add(leaf.XformNode())
}

func (node *XformNode) Update() {
	if node.Parent != nil {
		pxf := node.Parent.WorldXform
		xf := node.Xform
		wxf := node.WorldXform
		wxf.Position = pxf.Position.Add(pxf.Mulv(xf.Position).Mul(pxf.Scale))
		wxf.Mat3 = pxf.Mul(xf.Mat3)
		wxf.Scale = pxf.Scale.Mul(pxf.Mulv(xf.Scale))
	} else {
		*node.WorldXform = *node.Xform
	}
	for _, child := range node.Children {
		child.Update()
	}
}

func (node *XformNode) Walk(f func(Leaf)) {
	if node.Leaf != nil {
		f(node.Leaf)
	}
	for _, child := range node.Children {
		child.Walk(f)
	}
}
