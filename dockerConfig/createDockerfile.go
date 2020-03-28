package dockerConfig

import (
	"fmt"
	"io"
	"os"
)

// Generates a dockerfile instructions as a string
func render(f DockerfileData) string {
	result := fmt.Sprintf("FROM %s\n", f.ServerConfig.From)

	for goEnv, value := range f.ServerConfig.Env {
		result += fmt.Sprintf("ENV %s=%s\n", goEnv,value)
	}

	result += fmt.Sprintf("COPY %s %s\n", f.ServerConfig.Copy["baseDir"], f.ServerConfig.Copy["destDir"])

	result += fmt.Sprintf("WORKDIR %s\n", f.ServerConfig.WorkDir)

	result += fmt.Sprintf("RUN %s\n", f.ServerConfig.RunCommand)

	result += fmt.Sprintf("EXPOSE %s\n", f.ServerConfig.ServerPort)

	result += fmt.Sprintf("CMD [\"%s\"]\n", f.ServerConfig.FinalCmd)

	fmt.Println(result)
	return result
}

// creates and writes the rendered string into dockerfile
func WriteDockerfile(f DockerfileData) error {
	result := render(f)
	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, result)
	if err != nil {
		return err
	}
	return file.Sync()
}