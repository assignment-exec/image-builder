package builder

type command interface {
	execute() error
	undo() error
}

type stack []command

func newStack() *stack {
	var commands []command
	return (*stack)(&commands)
}

func (commands stack) isEmpty() bool {
	return len(commands) == 0
}

func (commands *stack) pop() command {
	var cmd command
	if !commands.isEmpty() {
		cmd = (*commands)[len(*commands)-1]
		*commands = append(stack(nil), (*commands)[:len(*commands)-1]...)
	}
	return cmd
}

func (commands *stack) push(cmd command) {
	*commands = append(*commands, cmd)
}
