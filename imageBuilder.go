package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"log"
)

func main() {

	// Unmarshal the yaml configuration file and generate a dockerfile.
	err := configurations.WriteDockerfile("code-runner.yaml")
	if err != nil {
		log.Fatalf("error while writing dockerfile: %v", err)
	}

	authData, err := builder.GetAuthData("docker-auth.yaml")
	if err != nil {
		log.Fatalf("error while reading authetication details: %v", err)
	}

	err = builder.BuildImage(*authData)
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
	}

	err = builder.PushImageToHub(*authData)
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
	}

	err = configurations.WriteDockerfile("assignment-env.yaml")
	if err != nil {
		log.Fatalf("error while writing dockerfile: %v", err)
	}
}
