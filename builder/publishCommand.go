package builder

import (
	"assignment-exec/image-builder/constants"
	"fmt"
	"github.com/pkg/errors"
)

type publishCommand struct {
	assgnEnv *assignmentEnvironment
}

func (cmd *publishCommand) execute() error {

	if err := cmd.assgnEnv.publish(); err != nil {
		return err
	}
	dockerRunCmd := fmt.Sprintf("%s %s %s", constants.DockerRunCommand,
		cmd.assgnEnv.ImgBuildConfig.imageTag, constants.PortCmdArg)
	fmt.Printf("\nFollowing is the command for starting %s\n\n", cmd.assgnEnv.ImgBuildConfig.imageTag)
	fmt.Println(dockerRunCmd)
	return nil
}

func (cmd *publishCommand) undo() error {
	err := cmd.assgnEnv.undoBuild()
	if err != nil {
		return errors.Wrap(err, "error in undo publish operation")
	}
	return nil
}
