package sge

import (
	"math"
)

import (
	"github.com/klkblake/s3dm"
)

type Skybox struct {
	*BasicLeaf
	CubeMap         *Texture
	Shader          *Program
	cachedMatrix    s3dm.Mat4
	mvpMatrix       *Uniform
	textureSampler  *Uniform
	mesh            *Mesh
}

func NewSkybox(cubeMap *Texture, shader *Program, far float64) *Skybox {
	skybox := new(Skybox)
	skybox.BasicLeaf = NewBasicLeaf()
	skybox.Link(skybox)
	scale := far / math.Sqrt(3)
	*skybox.Xform() = s3dm.XformIdentity
	skybox.Xform().Scale = s3dm.V3{scale, scale, scale}
	skybox.cachedMatrix = skybox.Xform().Matrix()
	skybox.CubeMap = cubeMap
	skybox.Shader = shader
	skybox.mvpMatrix = shader.Uniform("mvpMatrix")
	skybox.textureSampler = shader.Uniform("textureSampler")
	type skyboxVertex [2][3]float32
	verticies := []skyboxVertex{
		// Positive X
		{{1, 1, -1}, {1, 1, -1}},
		{{1, -1, -1}, {1, -1, -1}},
		{{1, -1, 1}, {1, -1, 1}},
		{{1, 1, 1}, {1, 1, 1}},
		// Negative X
		{{-1, 1, 1}, {-1, 1, 1}},
		{{-1, -1, 1}, {-1, -1, 1}},
		{{-1, -1, -1}, {-1, -1, -1}},
		{{-1, 1, -1}, {-1, 1, -1}},
		// Positive Y
		{{-1, 1, 1}, {-1, 1, 1}},
		{{-1, 1, -1}, {-1, 1, -1}},
		{{1, 1, -1}, {1, 1, -1}},
		{{1, 1, 1}, {1, 1, 1}},
		// Negative Y
		{{-1, -1, -1}, {-1, -1, -1}},
		{{-1, -1, 1}, {-1, -1, 1}},
		{{1, -1, 1}, {1, -1, 1}},
		{{1, -1, -1}, {1, -1, -1}},
		// Positive Z
		{{1, 1, 1}, {1, 1, 1}},
		{{1, -1, 1}, {1, -1, 1}},
		{{-1, -1, 1}, {-1, -1, 1}},
		{{-1, 1, 1}, {-1, 1, 1}},
		// Negative Z
		{{-1, 1, -1}, {-1, 1, -1}},
		{{-1, -1, -1}, {-1, -1, -1}},
		{{1, -1, -1}, {1, -1, -1}},
		{{1, 1, -1}, {1, 1, -1}},
	}
	indicies := []uint32{
		0, 1, 2, 0, 2, 3,
		4, 5, 6, 4, 6, 7,
		8, 9, 10, 8, 10, 11,
		12, 13, 14, 12, 14, 15,
		16, 17, 18, 16, 18, 19,
		20, 21, 22, 20, 22, 23,
	}
	skybox.mesh = NewMesh(verticies, indicies)
	return skybox
}

func (skybox *Skybox) Render(view *View, mvpMatrix s3dm.Mat4, pass int) {
	if skybox.CubeMap != nil {
		skybox.CubeMap.Bind(0)
	}
	skybox.Shader.Use()
	pos := view.Camera.Position
	if !pos.Equals(skybox.Xform().Position) {
		skybox.Xform().Position = pos
		skybox.cachedMatrix = skybox.Xform().Matrix()
	}
	matrix := mvpMatrix.Mul(skybox.cachedMatrix).RawMatrix32()
	skybox.mvpMatrix.SetMatrix(matrix[:], 4)
	if skybox.CubeMap != nil {
		skybox.textureSampler.Set(0)
	}
	skybox.mesh.Render()
}
