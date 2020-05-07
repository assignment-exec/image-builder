package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"flag"
	"log"
)

var publishImage = flag.Bool("publishImage", false, "Push image to docker hub")
var codeRunnerConfig = flag.String("codeRunnerConfig", "code-runner.yaml", "Code Runner configuration filename")
var dockerAuthConfig = flag.String("dockerAuthConfig", "docker-auth.yaml", "Docker hub authentication configuration filename")
var dockerfileName = flag.String("dockerfile", "Dockerfile", "Dockerfile name")

func main() {
	flag.Parse()

	log.Println("Creating Dockerfile...")

	// Unmarshal the yaml configuration file and generate a dockerfileName.
	err := configurations.WriteDockerfile(*codeRunnerConfig, *dockerfileName)
	if err != nil {
		log.Fatalf("error while writing dockerfileName: %v", err)
		return
	}

	imgBuilder := builder.NewImageBuilder(*dockerAuthConfig, *dockerfileName)
	err = imgBuilder.BuildImage()
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
		return
	}

	if *publishImage {
		err = imgBuilder.PublishImage()
		if err != nil {
			log.Fatalf("error while pushing image to docker hub: %v", err)
		}
	}
}
