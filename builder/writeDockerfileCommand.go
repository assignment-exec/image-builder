package builder

import (
	"log"
	"os"
)

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

func (assgnEnv *assignmentEnv) writeDockerfile() error {
	if assgnEnv.DockerfileData.Len() > 0 {
		file, err := os.Create(assgnEnv.ImgBuildConfig.dockerfileName)
		defer func() {
			err = file.Close()
			if err != nil {
				log.Println("error while closing the created Dockerfile", err)
				return
			}
		}()
		if err != nil {
			return err
		}
		_, err = file.WriteString(assgnEnv.DockerfileData.String())
		return err
	}
	return nil
}
