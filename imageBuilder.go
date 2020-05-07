package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"flag"
	"log"
)

var publishImage = flag.Bool("publishImage", false, "Push image to docker hub")
var codeRunnerConfig = flag.String("codeRunnerConfig", "code-runner.yaml", "Code Runner configuration filename")
var assignmentEnvConfig = flag.String("assignmentEnvConfig", "assignment-env.yaml", "Assignment Environment configuration filename")
var dockerAuthConfig = flag.String("dockerAuthConfig", "docker-auth.yaml", "Docker hub authentication configuration filename")
var dockerfileName = flag.String("dockerfile", "Dockerfile", "Dockerfile name")

func main() {

	flag.Parse()

	log.Println("Creating Dockerfile...")

	imgBuilder := builder.NewImageBuilder(*dockerAuthConfig, *dockerfileName)

	err := generateCodeRunnerImage(imgBuilder)
	if err == nil {
		err = generateAssignmentEnvImage(imgBuilder)
	}
}

// Generate a dockerfile for code runner server, build its image and push it ot docker hub.
func generateCodeRunnerImage(imgBuilder *builder.ImageBuilder) error {

	// Unmarshal the yaml configuration file and generate a dockerfileName.
	err := configurations.WriteDockerfile(*codeRunnerConfig, *dockerfileName)
	if err != nil {
		log.Fatalf("error while writing dockerfileName: %v", err)
		return err
	}

	err = imgBuilder.BuildImage(false)
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
		return err
	}

	if *publishImage {
		err = imgBuilder.PublishImage()
		if err != nil {
			log.Fatalf("error while pushing image to docker hub: %v", err)
			return err
		}
	}
	return nil
}

// Generate a dockerfile for assignment environment and build its image.
func generateAssignmentEnvImage(imgBuilder *builder.ImageBuilder) error {
	err := configurations.WriteDockerfile(*assignmentEnvConfig,*dockerfileName)
	if err != nil {
		log.Fatalf("error while writing dockerfile for assignment environment: %v", err)
		return err
	}

	err = imgBuilder.BuildImage(true)
	if err != nil {
		log.Fatalf("error while building image for assignment environment: %v", err)
		return err
	}

	return nil
}
