// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub.
package builder

import (
	"github.com/pkg/errors"
)

// buildCommand struct type holds assignmentEnvironmentImageBuilder instance
// which is required to perform image build operation.
type buildCommand struct {
	asgmtEnv *assignmentEnvironmentImageBuilder
}

// execute invokes the build function to build the docker image.
func (cmd *buildCommand) execute() error {
	return cmd.asgmtEnv.build()
}

// undo invokes the deleteDockerfile function to delete the created
// dockerfile if any error is encountered while building the image.
func (cmd *buildCommand) undo() error {
	if err := cmd.asgmtEnv.deleteDockerfile(); err != nil {
		return errors.Wrap(err, "error in undo build operation")
	}
	return nil
}
