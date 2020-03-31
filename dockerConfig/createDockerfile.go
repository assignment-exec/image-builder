package dockerConfig

import (
	"os"
)

func WriteDockerfile() error {
	data, err := newDockerFileDataFromYamlFile("serverConfig.yaml")
	if err != nil {
		return err
	}

	tmpl := newDockerfileTemplate(data)

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	err = tmpl.generateTemplate(file)
	return  err
}
