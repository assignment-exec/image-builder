package configurations

import (
	"os"
)

func WriteDockerfile() error {
	data, err := newDockerFileDataFromYamlFile("code-runner.yaml")
	if err != nil {
		return err
	}

	tmpl := newDockerfileTemplate(data)

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	err = tmpl.generateDockerfileFromTemplate(file)
	return err
}
