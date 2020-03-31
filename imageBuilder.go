package main

import (
	"assignment-exec/image-builder/dockerConfig"
	"bytes"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Creating Dockerfile...")

	data, err := dockerConfig.NewDockerFileDataFromYamlFile("serverConfig.yaml")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := dockerConfig.NewDockerfileTemplate(data)

	output := &bytes.Buffer{}
	err = tmpl.GenerateTemplate(output)
	fmt.Println(output)

}
