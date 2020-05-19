// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub. It uses command pattern to
// perform all operations and perform undo operations when any error is encountered.
package builder

// verifyCommand struct type holds assignmentEnvironment instance
// which is required to verify language image and write the dockerfile instructions.
type verifyCommand struct {
	assgnEnv *assignmentEnvironment
}

// execute invokes the verifyAndWriteInstructions function to verify
// whether language image is already present on docker hub and accordingly
// write dockerfile instructions to the bytes buffer for the provided
// assignment environment configurations.
func (cmd *verifyCommand) execute() error {
	return cmd.assgnEnv.verifyAndWriteInstructions()
}

// undo is a No operation function as there is no possible undo to be performed
// if the verification fails.
func (cmd *verifyCommand) undo() error {
	// No operation.
	return nil
}
