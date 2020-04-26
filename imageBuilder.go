package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"log"
)

func main() {

	authData, err := builder.GetAuthData("docker-auth.yaml")

	// Unmarshal the yaml configuration file and generate a dockerfile.
	err = configurations.WriteDockerfile("code-runner.yaml")

	if err != nil {
		log.Fatalf("error while reading docker authetication details: %v", err)
		return
	}

	err = generateCodeRunnerImage(authData)
	if err == nil {
		err = generateAssignmentEnvImage(authData)
	}
}

// Generate a dockerfile for code runner server, build its image and push it ot docker hub.
func generateCodeRunnerImage(authData *builder.DockerAuthData) error {

	err := configurations.WriteDockerfile("code-runner.yaml")
	if err != nil {
		log.Fatalf("error while writing dockerfile for code runner server: %v", err)
		return err
	}

	err = builder.BuildImage(*authData, false)
	if err != nil {
		log.Fatalf("error while building image for code runner server: %v", err)
		return err
	}

	err = builder.PushImageToHub(*authData)
	if err != nil {
		log.Fatalf("error while pushing code runner server image to docker hub: %v", err)
		return err
	}
	return nil
}

// Generate a dockerfile for assignment environment and build its image.
func generateAssignmentEnvImage(authData *builder.DockerAuthData) error {
	err := configurations.WriteDockerfile("assignment-env.yaml")
	if err != nil {
		log.Fatalf("error while writing dockerfile for assignment environment: %v", err)
		return err
	}

	err = builder.BuildImage(*authData, true)
	if err != nil {
		log.Fatalf("error while building image for assignment environment: %v", err)
		return err
	}

	return nil
}
