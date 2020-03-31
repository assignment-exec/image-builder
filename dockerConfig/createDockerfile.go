package dockerConfig

// Generates a dockerfile instructions As a string
/*func Render(config dockerfileData) string {
	result := fmt.Sprintf("FROM %s\n", config.ServerConfig.from)

	for goEnv, value := range config.ServerConfig.env {
		result += fmt.Sprintf("ENV %s=%s\n", goEnv, value)
	}

	result += fmt.Sprintf("COPY %s %s\n", config.ServerConfig.copyCommand["BaseDir"], config.ServerConfig.copyCommand["DestDir"])

	result += fmt.Sprintf("WORKDIR %s\n", config.ServerConfig.workDir)

	result += fmt.Sprintf("RUN %s\n", config.ServerConfig.runCommand)

	result += fmt.Sprintf("EXPOSE %s\n", config.ServerConfig.serverPort)

	result += fmt.Sprintf("CMD [\"%s\"]\n", config.ServerConfig.FinalCmd)

	fmt.Println(result)
	return result
}

// creates and writes the rendered string into dockerfile
func WriteDockerfile(config dockerfileData) error {
	result := Render(config)
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
}*/
