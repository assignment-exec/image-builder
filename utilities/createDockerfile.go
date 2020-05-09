package utilities

import (
	"assignment-exec/image-builder/configurations"
	"log"
	"os"
)

func WriteDockerfileForAssignmentEnv(configFilename string, dockerFilename string) (err error, language string, version string) {
	c, err := configurations.GetAssignmentEnvConfig(configFilename)
	if err != nil {
		return err, "", ""
	}
	file, err := os.Create(dockerFilename)
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println("error while closing the created Dockerfile", err)
			return
		}
	}()
	if err != nil {
		return err, "", ""
	}
	_, err = file.WriteString(c.WriteInstruction())

	return err, c.Deps.Language.Name, c.Deps.Language.Version
}
