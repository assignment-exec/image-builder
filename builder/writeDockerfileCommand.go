// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub.
package builder

// writeDockerfileCommand struct type holds assignmentEnvironmentImageBuilder instance
// which is required to write the dockerfile instructions from bytes buffer
// to an actual Dockerfile.
type writeDockerfileCommand struct {
	asgmtEnv *assignmentEnvironmentImageBuilder
}

// execute invokes writeToDockerfile function to write the stored instructions
// to the Dockerfile.
func (cmd *writeDockerfileCommand) execute() error {
	return cmd.asgmtEnv.writeToDockerfile()
}

// undo invokes functions to reset and clear the dockerfile instructions bytes
// buffer if any error is encountered while writing to Dockerfile.
func (cmd *writeDockerfileCommand) undo() error {
	cmd.asgmtEnv.resetDockerfileData()
	return nil
}
