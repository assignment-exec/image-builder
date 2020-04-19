package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"log"
)

func main() {
	log.Println("Creating Dockerfile...")

	// Unmarshal the yaml configuration file and generate a dockerfile.
	err := configurations.WriteDockerfile()
	if err != nil {
		log.Fatalf("error while writing dockerfile: %v", err)
	}

	err = builder.BuildImage()
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
	}
}
