package builder

type verifyCommand struct {
	assgnEnv *assignmentEnv
}

func (cmd *verifyCommand) execute() error {
	return cmd.assgnEnv.verifyAndWrite()
}

func (cmd *verifyCommand) undo() error {
	//No operation.
	return nil
}
