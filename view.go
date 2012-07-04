package sge

import (
	"math"
)

import (
	"github.com/klkblake/s3dm"
)

type View struct {
	*s3dm.Frustum
	Width  int
	Height int
}

func NewView(width int, height int, near float64, far float64) *View {
	view := new(View)
	fovy := math.Pi / 4
	aspect := float64(width) / float64(height)
	view.Frustum = s3dm.NewFrustum(near, far, fovy, aspect)
	view.Update()
	view.Width = width
	view.Height = height
	return view
}
