// Package builder provides primitives to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub. It uses command pattern to
// perform all operations and perform undo operations when any error is encountered.
package builder

// writeDockerfileCommand struct type holds assignmentEnvironment instance
// which is required to write the dockerfile instructions from bytes buffer
// to an actual Dockerfile.
type writeDockerfileCommand struct {
	assgnEnv *assignmentEnvironment
}

// execute invokes writeToDockerfile function to write the stored instructions
// to the Dockerfile.
func (cmd *writeDockerfileCommand) execute() error {
	return cmd.assgnEnv.writeToDockerfile()
}

// undo invokes functions to reset and clear the dockerfile instructions bytes
// buffer if any error is encountered while writing to Dockerfile.
func (cmd *writeDockerfileCommand) undo() error {
	cmd.assgnEnv.resetDockerfileData()
	return nil
}
