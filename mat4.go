package sge

import (
	"math"
	"strconv"
)

import (
	"s3dm"
)

// 4x4 matrix in column-major order.
type Mat4 [4 * 4]float64

func NewMat4() *Mat4 {
	matrix := new(Mat4)
	matrix.SetIdentity()
	return matrix
}

func NewPerspectiveMat4(fovy, aspect, near, far float64) *Mat4 {
	top := near * math.Tan(fovy*0.5)
	right := aspect * top
	return &Mat4{
		near / right, 0, 0, 0,
		0, near / top, 0, 0,
		0, 0, -(far + near) / (far - near), -1,
		0, 0, -2 * far * near / (far - near), 0}
}

func NewOrthographicMat4(width, height, near, far float64) *Mat4 {
	return &Mat4{
		2/width, 0, 0, 0,
		0, 2/height, 0, 0,
		0, 0, 2/(near-far), 0,
		-1, -1, (near+far)/(near-far), 1}
}

func (m *Mat4) SetIdentity() {
	*m = [4 * 4]float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func (m *Mat4) GetMatrix() [4 * 4]float64 {
	return [4 * 4]float64(*m)
}

func (m *Mat4) GetFloat32Matrix() [4 * 4]float32 {
	matrix := new([4 * 4]float32)
	for i, v := range m {
		matrix[i] = float32(v)
	}
	return *matrix
}

func (m *Mat4) Position() *s3dm.V3 {
	return s3dm.NewV3(m[12], m[13], m[14])
}

func (m *Mat4) Mul(o *Mat4) *Mat4 {
	res := new(Mat4)
	for row := 0; row < 4; row++ {
		col1 := row
		col2 := col1 + 4
		col3 := col1 + 8
		col4 := col1 + 12
		res[col1] =
			m[col1]*o[0] +
				m[col2]*o[1] +
				m[col3]*o[2] +
				m[col4]*o[3]
		res[col2] =
			m[col1]*o[4] +
				m[col2]*o[5] +
				m[col3]*o[6] +
				m[col4]*o[7]
		res[col3] =
			m[col1]*o[8] +
				m[col2]*o[9] +
				m[col3]*o[10] +
				m[col4]*o[11]
		res[col4] =
			m[col1]*o[12] +
				m[col2]*o[13] +
				m[col3]*o[14] +
				m[col4]*o[15]
	}
	return res
}

func (m *Mat4) Transpose() *Mat4 {
	return &Mat4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15]}
}

func (m *Mat4) String() string {
	s := "["
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			s += strconv.FormatFloat(m[col+row*4], 'g', 2, 64)
			if col != 3 || row != 3 {
				s += ", "
			}
		}
	}
	s += "]"
	return s
}
