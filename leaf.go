package sge

import "github.com/klkblake/s3dm"

type Leaf interface {
	Xform() *s3dm.XformScale
	Update(delta float64)
	Passes() int
	AABB() s3dm.AABB
	Render(view *View, mvpMatrix s3dm.Mat4, pass int)
	XformNode() *XformNode
}

type BasicLeaf struct {
	xformNode *XformNode
}

func NewBasicLeaf() *BasicLeaf {
	return &BasicLeaf{xformNode: NewXformNode()}
}

func (leaf *BasicLeaf) Link(self Leaf) {
	leaf.xformNode.Leaf = self
}

func (leaf *BasicLeaf) Xform() *s3dm.XformScale {
	return &leaf.xformNode.Xform
}

func (leaf *BasicLeaf) Passes() int {
	return PassOpaque
}

func (leaf *BasicLeaf) AABB() s3dm.AABB {
	pos := leaf.XformNode().WorldXform.Position
	return s3dm.AABB{pos, pos}
}

func (leaf *BasicLeaf) XformNode() *XformNode {
	return leaf.xformNode
}

func (leaf *BasicLeaf) Update(delta float64)                             {}
func (leaf *BasicLeaf) Render(view *View, mvpMatrix s3dm.Mat4, pass int) {}
