package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/environment"
	"assignment-exec/image-builder/utilities"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

var publishImage = flag.Bool("publishImage", false, "Push image to docker hub")
var codeRunnerConfig = flag.String("codeRunnerConfig", "code-runner.yaml", "Code Runner configuration filename")
var assignmentEnv = flag.String("assignmentEnv", "assignment-env.yaml", "Assignment Environment configuration filename")
var dockerfileName = flag.String("dockerfile", "Dockerfile", "Dockerfile name")

func main() {

	flag.Parse()

	authData := builder.GetAuthData()

	err := generateCodeRunnerImage(authData)
	if err != nil {
		log.Fatalf("error while building code runner image: %v", err)
	}

	err = generateAssignmentEnvImage(authData)
	if err != nil {
		log.Fatalf("error while building assignment environment image: %v", err)
	}

}

// Generate a dockerfile for code runner server, build its image and push it ot docker hub.
func generateCodeRunnerImage(authData *builder.DockerAuthData) error {

	// Unmarshal the yaml configuration file and generate a dockerfileName.
	err, _, _ := utilities.WriteDockerfile(*codeRunnerConfig, *dockerfileName)
	if err != nil {
		return errors.Wrap(err, "error in writing dockerfile for code runner")
	}

	repositoryName := os.Getenv(environment.CodeRunnerRepository)
	repositoryVersion := os.Getenv(environment.CodeRunnerRepositoryVersion)
	imageTag := fmt.Sprintf("%s/%s:%s", authData.Username, repositoryName, repositoryVersion)

	imgBuilder, err := builder.NewImageBuilder(
		builder.WithDockerAuthData(authData),
		builder.WithImageTag(imageTag),
		builder.WithDockerfileName(*dockerfileName))

	if err != nil {
		return errors.Wrap(err, "error in creating image builder instance for code runner")
	}
	err = imgBuilder.BuildImage()
	if err != nil {
		return err
	}

	if *publishImage {
		err = imgBuilder.PublishImage()
		if err != nil {
			return errors.Wrap(err, "error in pushing code runner image")
		}
	}
	return nil
}

// Generate a dockerfile for assignment environment and build its image.
func generateAssignmentEnvImage(authData *builder.DockerAuthData) error {
	err, language, version := utilities.WriteDockerfile(*assignmentEnv, *dockerfileName)
	if err != nil {
		return errors.Wrap(err, "error in writing dockerfile for assignment environment")
	}

	imageTag := fmt.Sprintf("%s/%s%s", authData.Username, language, version)

	imgBuilder, err := builder.NewImageBuilder(
		builder.WithDockerAuthData(authData),
		builder.WithImageTag(imageTag),
		builder.WithDockerfileName(*dockerfileName))

	if err != nil {
		return errors.Wrap(err, "error in creating image builder instance for assignment env")
	}

	err = imgBuilder.BuildImage()
	if err != nil {
		return err
	}

	if *publishImage {
		err = imgBuilder.PublishImage()
		if err != nil {
			return errors.Wrap(err, "error in pushing code runner image")
		}
	}

	return nil
}
