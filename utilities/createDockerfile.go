package utilities

import (
	"assignment-exec/image-builder/configurations"
	"log"
	"os"
)

func WriteDockerfileForAssignmentEnv(configFilename string, dockerFilename string) (error, string, string) {
	c, _ := configurations.GetConfig(configFilename)
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
	_, err = file.WriteString(c.String())

	return err, c.Dependencies.Language.Name, c.Dependencies.Language.Version
}
