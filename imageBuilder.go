package main

import (
	"assignment-exec/image-builder/dockerConfig"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Creating Dockerfile...")

	// read the yaml config
	config, err := dockerConfig.UnmarshalYAML("config.yaml")
	if err != nil {
		log.Fatalf("Error while reading yaml config %v", err)
	}

	// create a dockerfile with the above read configs
	err = dockerConfig.WriteDockerfile(*config)
	if err != nil {
		log.Fatalf("Error occurred while writing Dockerfile %v", err)
	}
}
