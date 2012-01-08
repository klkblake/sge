package sge

type Mat4Stack struct {
	Stack []*Mat4
}

func NewMat4Stack() *Mat4Stack {
	stack := new(Mat4Stack)
	stack.Stack = make([]*Mat4, 1)
	stack.Stack[0] = NewMat4()
	return stack
}

func (stack *Mat4Stack) Push(matrix *Mat4) {
	stack.Stack = append(stack.Stack, stack.Top().Mul(matrix))
}

func (stack *Mat4Stack) Pop() {
	stack.Stack = stack.Stack[:len(stack.Stack)-1]
}

func (stack *Mat4Stack) Top() *Mat4 {
	return stack.Stack[len(stack.Stack)-1]
}
