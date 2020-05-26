package builder

type writeDockerfileCommand struct {
	assgnEnv *assignmentEnvironment
}

func (cmd *writeDockerfileCommand) execute() error {
	return cmd.assgnEnv.writeToDockerfile()
}

func (cmd *writeDockerfileCommand) undo() error {
	cmd.assgnEnv.resetDockerfileData()
	return nil
}