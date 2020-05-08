package utilities

import (
	"assignment-exec/image-builder/configurations"
	"github.com/pkg/errors"
	"log"
	"os"
)

// Creates a template and writes it to a new Dockerfile.
func WriteDockerfileForCodeRunner(configFilename string, dockerFilename string) error {
	data, err := configurations.NewDockerFileDataFromYamlFile(configFilename)
	if err != nil {
		return err
	}

	tmpl := configurations.NewDockerfileTemplate(data)

	file, err := os.Create(dockerFilename)
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println("error while closing the created Dockerfile", err)
			return
		}
	}()
	if err != nil {
		return errors.Wrap(err, "error in creating dockerfile")
	}

	err = tmpl.GenerateDockerfileFromTemplate(file)

	return err
}

func WriteDockerfileForAssignmentEnv(configFilename string, dockerFilename string) (error, string, string) {
	data, err := configurations.NewDockerFileDataFromYamlFile(configFilename)
	if err != nil {
		return err, "", ""
	}

	tmpl := configurations.NewDockerfileTemplate(data)

	file, err := os.Create(dockerFilename)
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println("error while closing the created Dockerfile", err)
			return
		}
	}()
	if err != nil {
		return errors.Wrap(err, "error in creating dockerfile"), "", ""
	}

	err = tmpl.GenerateDockerfileFromTemplate(file)

	language, version := tmpl.GetLanguageFormat()

	return err, language, version
}
