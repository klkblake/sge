package sge

import "github.com/klkblake/s3dm"

type XformNode struct {
	Parent      *XformNode
	Children    []*XformNode
	Xform       s3dm.XformScale
	WorldXform  s3dm.XformScale
	Leaf        Leaf
}

func NewXformNode() *XformNode {
	node := new(XformNode)
	node.Xform = s3dm.XformScaleIdentity
	node.WorldXform = s3dm.XformScaleIdentity
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
		// XXX this should be a function in s3dm.
		pxf := node.Parent.WorldXform
		xf := node.Xform
		var wxf s3dm.XformScale
		// XXX Dodgy hax, will be going away soon.
		wxf.Position = pxf.Position.Add(xf.Position.Sub(s3dm.Position{}).Rotate(pxf.Rotation).Mul(pxf.Scale))
		wxf.Rotation = pxf.Rotation.Mul(xf.Rotation)
		wxf.Scale = pxf.Scale.Mul(xf.Scale.Rotate(pxf.Rotation))
		node.WorldXform = wxf
	} else {
		node.WorldXform = node.Xform
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
