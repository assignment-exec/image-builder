package configurations

import (
	"fmt"
	"log"
	"os"
)

// Creates a template and writes it to a new Dockerfile.
func WriteDockerfile(configFilename string, dockerFilename string) (error, string) {
	data, err := newDockerFileDataFromYamlFile(configFilename)
	if err != nil {
		return err, ""
	}

	tmpl := newDockerfileTemplate(data)

	file, err := os.Create(dockerFilename)
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("error while closing the created Dockerfile: %v", err)
			return
		}
	}()
	if err != nil {
		return err, ""
	}

	err = tmpl.generateDockerfileFromTemplate(file)

	var languageImageFormat string
	for _, inst := range tmpl.Data {
		switch inst := inst.(type) {
		case programmingLanguage:
			languageImageFormat = fmt.Sprintf("%s%s", inst.Name, inst.Version)
			break
		}
	}

	return err, languageImageFormat
}
