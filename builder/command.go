package builder

type command interface {
	execute() error
	undo() error
}

type stack []command

func newStack() *stack {
	var cmdStack []command
	return (*stack)(&cmdStack)
}

func (cmdStack stack) isEmpty() bool {
	return len(cmdStack) == 0
}

func (cmdStack *stack) pop() command {
	var cmd command
	if !cmdStack.isEmpty() {
		cmd = (*cmdStack)[len(*cmdStack)-1]
		*cmdStack = append(stack(nil), (*cmdStack)[:len(*cmdStack)-1]...)
	}
	return cmd
}

func (cmdStack *stack) push(cmd command) {
	*cmdStack = append(*cmdStack, cmd)
}
