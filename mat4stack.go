package sge

import (
	"s3dm"
)

type Mat4Stack struct {
	Stack []*s3dm.Mat4
}

func NewMat4Stack() *Mat4Stack {
	stack := new(Mat4Stack)
	stack.Stack = make([]*s3dm.Mat4, 1)
	stack.Stack[0] = s3dm.NewMat4()
	return stack
}

func (stack *Mat4Stack) Push(matrix *s3dm.Mat4) {
	stack.Stack = append(stack.Stack, stack.Top().Mul(matrix))
}

func (stack *Mat4Stack) Pop() {
	stack.Stack = stack.Stack[:len(stack.Stack)-1]
}

func (stack *Mat4Stack) Top() *s3dm.Mat4 {
	return stack.Stack[len(stack.Stack)-1]
}
