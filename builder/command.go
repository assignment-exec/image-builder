package builder

type command interface {
	execute() error
	undo() error
}

type stack []command

func newStack() *stack {
	var s []command
	return (*stack)(&s)
}

func (s *stack) pop() command {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *stack) push(h command) {
	*s = append(*s, h)
}
