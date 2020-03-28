package dockerConfig

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Generates a dockerfile instructions as a string
func render(config dockerfileData) string {
	result := fmt.Sprintf("FROM %s\n", config.ServerConfig.From)

	for goEnv, value := range config.ServerConfig.Env {
		result += fmt.Sprintf("ENV %s=%s\n", goEnv, value)
	}

	result += fmt.Sprintf("COPY %s %s\n", config.ServerConfig.Copy["baseDir"], config.ServerConfig.Copy["destDir"])

	result += fmt.Sprintf("WORKDIR %s\n", config.ServerConfig.WorkDir)

	result += fmt.Sprintf("RUN %s\n", config.ServerConfig.RunCommand)

	result += fmt.Sprintf("EXPOSE %s\n", config.ServerConfig.ServerPort)

	result += fmt.Sprintf("CMD [\"%s\"]\n", config.ServerConfig.FinalCmd)

	fmt.Println(result)
	return result
}

// creates and writes the rendered string into dockerfile
func WriteDockerfile(config dockerfileData) error {
	result := render(config)
	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	_, err = io.WriteString(file, result)
	if err != nil {
		return err
	}
	return file.Sync()
}
