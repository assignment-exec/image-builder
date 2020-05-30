// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publishImage it to docker hub. It uses command pattern to
// perform all operations and perform undo operations when any error is encountered.
package builder

import (
	"assignment-exec/image-builder/constants"
	"fmt"
	"github.com/pkg/errors"
)

// publishCommand struct type holds assignmentEnvironmentImageBuilder instance
// which is required to perform image publishImage operation.
type publishCommand struct {
	asgmtEnv *assignmentEnvironmentImageBuilder
}

// execute invokes the publishImage function to push the image to docker hub.
func (cmd *publishCommand) execute() error {

	if err := cmd.asgmtEnv.publishImage(); err != nil {
		return err
	}
	dockerRunCmd := fmt.Sprintf("%s %s %s", constants.DockerRunCommand,
		cmd.asgmtEnv.ImgBuildConfig.imageTag, constants.PortCmdArg)
	fmt.Printf("\nFollowing is the command for starting %s\n\n", cmd.asgmtEnv.ImgBuildConfig.imageTag)
	fmt.Println(dockerRunCmd)
	return nil
}

// undo invokes undoBuild function to remove the locally built image
// if any error is encountered while publishing the image.
func (cmd *publishCommand) undo() error {
	err := cmd.asgmtEnv.undoBuild()
	if err != nil {
		return errors.Wrap(err, "error in undo publishImage operation")
	}
	return nil
}
