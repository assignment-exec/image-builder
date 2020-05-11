package builder

import (
	"github.com/pkg/errors"
)

type buildCommand struct {
	assgnEnv *assignmentEnv
}

func (cmd *buildCommand) execute() error {
	return cmd.assgnEnv.build()
}

func (cmd *buildCommand) undo() error {
	if err := cmd.assgnEnv.undoBuild(); err != nil {
		return errors.Wrap(err, "error in undo build operation")
	}
	return nil
}
