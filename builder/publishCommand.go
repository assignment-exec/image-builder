package builder

import (
	"github.com/pkg/errors"
)

type publishCommand struct {
	assgnEnv *assignmentEnv
}

func (cmd *publishCommand) execute() error {
	return cmd.assgnEnv.publish()
}

func (cmd *publishCommand) undo() error {
	err := cmd.assgnEnv.undoPublish()
	if err != nil {
		return errors.Wrap(err, "error in undo publish operation")
	}
	return nil
}
