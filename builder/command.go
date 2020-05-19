// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub. It uses command pattern to
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
	var cmdStack []command
	return (*stack)(&cmdStack)
}

// isEmpty checks whether stack is empty.
// It returns a boolean flag accordingly.
func (cmdStack stack) isEmpty() bool {
	return len(cmdStack) == 0
}

// pop removes the last inserted command from the stack.
// It returns the popped command.
func (cmdStack *stack) pop() command {
	var cmd command
	if !cmdStack.isEmpty() {
		cmd = (*cmdStack)[len(*cmdStack)-1]
		*cmdStack = append(stack(nil), (*cmdStack)[:len(*cmdStack)-1]...)
	}
	return cmd
}

// push appends the command to the stack.
func (cmdStack *stack) push(cmd command) {
	*cmdStack = append(*cmdStack, cmd)
}
