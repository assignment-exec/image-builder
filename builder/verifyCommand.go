// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub.
package builder

// verifyCommand struct type holds assignmentEnvironmentImageBuilder instance
// which is required to verify language image and write the dockerfile instructions.
type verifyCommand struct {
	asgmtEnv *assignmentEnvironmentImageBuilder
}

// execute invokes the verifyAndWriteInstructions function to verify whether
// a docker image for given language is already present on docker hub and accordingly
// write dockerfile instructions for the provided assignment environment configurations.
func (cmd *verifyCommand) execute() error {
	return cmd.asgmtEnv.verifyAndWriteInstructions()
}

// undo is a No operation function as there is no possible undo to be performed
// if the verification fails.
func (cmd *verifyCommand) undo() error {
	// No operation.
	return nil
}
