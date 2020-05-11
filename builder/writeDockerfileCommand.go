package builder

type writeDockerfileCommand struct {
	assgnEnv *assignmentEnv
}

func (cmd *writeDockerfileCommand) execute() error {
	return cmd.assgnEnv.writeDockerfile()
}

func (cmd *writeDockerfileCommand) undo() error {
	cmd.assgnEnv.undoWrite()
	return nil
}
