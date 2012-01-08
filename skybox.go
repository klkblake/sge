package sge

import "math"

import "gl"

type Skybox struct {
	CubeMap *Texture
	Shader *Program
	scaleMatrix *Mat4
	translateMatrix *Mat4
	mvpMatrixUniform gl.UniformLocation
	textureUnitUniform gl.UniformLocation
	mesh *Mesh
}

func NewSkybox(cubeMap *Texture, shader *Program, far float64) *Skybox {
	skybox := new(Skybox)
	skybox.CubeMap = cubeMap
	skybox.Shader = shader
	scale := far / math.Sqrt(3)
	skybox.scaleMatrix = &Mat4{
		scale, 0, 0, 0,
		0, scale, 0 ,0,
		0, 0, scale, 0,
		0, 0, 0, 1,
	}
	skybox.translateMatrix = NewMat4()
	skybox.mvpMatrixUniform = shader.GetUniformLocation("mvpMatrix")
	skybox.textureUnitUniform = shader.GetUniformLocation("textureUnit")
	verticies := []float32{
		// Positive X
		1, 1, -1,
		1, -1, -1,
		1, -1, 1,
		1, 1, 1,
		// Negative X
		-1, 1, 1,
		-1, -1, 1,
		-1, -1, -1,
		-1, 1, -1,
		// Positive Y
		-1, 1, 1,
		-1, 1, -1,
		1, 1, -1,
		1, 1, 1,
		// Negative Y
		-1, -1, -1,
		-1, -1, 1,
		1, -1, 1,
		1, -1, -1,
		// Positive Z
		1, 1, 1,
		1, -1, 1,
		-1, -1, 1,
		-1, 1, 1,
		// Negative Z
		-1, 1, -1,
		-1, -1, -1,
		1, -1, -1,
		1, 1, -1,
	}
	indicies := []uint32 {
		0, 1, 2, 0, 2, 3,
		4, 5, 6, 4, 6, 7,
		8, 9, 10, 8, 10, 11,
		12, 13, 14, 12, 14, 15,
		16, 17, 18, 16, 18, 19,
		20, 21, 22, 20, 22, 23,
	}
	skybox.mesh = NewMesh(verticies, verticies, indicies)
	return skybox
}

func (skybox *Skybox) Render(view *View, mvpMatrix *Mat4, pass int) {
	skybox.CubeMap.Bind()
	skybox.Shader.Use()
	pos := view.Camera.Position()
	skybox.translateMatrix[12] = pos.X
	skybox.translateMatrix[13] = pos.Y
	skybox.translateMatrix[14] = pos.Z
	matrix := mvpMatrix.Mul(skybox.translateMatrix).Mul(skybox.scaleMatrix).GetFloat32Matrix()
	skybox.mvpMatrixUniform.UniformMatrix4fv(false, 1, matrix[:])
	skybox.textureUnitUniform.Uniform1i(0)
	skybox.mesh.Render()
}