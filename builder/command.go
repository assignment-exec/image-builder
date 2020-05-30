// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publishImage it to docker hub. It uses command pattern to
// perform all operations and perform undo operations when any error is encountered.
package builder

// command interface type represents the execute
// and undo function required by different commands
// to perform respective operations.
type command interface {
	execute() error
	undo() error
}

// stack type for holding the commands in the order of their execution.
type stack []command

// newStack creates and returns a new instance of stack.
func newStack() *stack {
	var commands []command
	return (*stack)(&commands)
}

// isEmpty checks whether stack is empty.
// It returns a boolean flag accordingly.
func (commands stack) isEmpty() bool {
	return len(commands) == 0
}

// pop removes the last inserted command from the stack.
// It returns the popped command.
func (commands *stack) pop() command {
	var cmd command
	if !commands.isEmpty() {
		cmd = (*commands)[len(*commands)-1]
		*commands = append(stack(nil), (*commands)[:len(*commands)-1]...)
	}
	return cmd
}

// push appends the command to the stack.
func (commands *stack) push(cmd command) {
	*commands = append(*commands, cmd)
}
