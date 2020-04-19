package configurations

import (
	"assignment-exec/image-builder/constants"
	"log"
	"os"
)

// Creates a template and writes it to a new Dockerfile.
func WriteDockerfile() error {
	data, err := newDockerFileDataFromYamlFile("code-runner.yaml")
	if err != nil {
		return err
	}

	tmpl := newDockerfileTemplate(data)

	file, err := os.Create(constants.DockerFilepath)
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("error while closing the created Dockerfile: %v", err)
			return
		}
	}()
	if err != nil {
		return err
	}

	err = tmpl.generateDockerfileFromTemplate(file)

	return err
}
