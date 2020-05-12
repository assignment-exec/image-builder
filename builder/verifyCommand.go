package builder

type verifyCommand struct {
	assgnEnv *assignmentEnvironment
}

func (cmd *verifyCommand) execute() error {
	return cmd.assgnEnv.verifyAndWriteInstructions()
}

func (cmd *verifyCommand) undo() error {
	// No operation.
	return nil
}
