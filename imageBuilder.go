package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/utilities"
	"flag"
	"github.com/pkg/errors"
	"log"
)

var publishImage = flag.Bool("publishImage", false, "Push image to docker hub")
var codeRunnerConfig = flag.String("codeRunnerConfig", "code-runner.yaml", "Code Runner configuration filename")
var assignmentEnv = flag.String("assignmentEnv", "assignment-env.yaml", "Assignment Environment configuration filename")
var dockerfileName = flag.String("dockerfile", "Dockerfile", "Dockerfile name")

func main() {

	flag.Parse()

	imgBuilder := builder.NewImageBuilder(*dockerfileName)

	err := generateCodeRunnerImage(imgBuilder)
	if err != nil {
		log.Fatalf("error while building code runner image: %v", err)
	}

	err = generateAssignmentEnvImage(imgBuilder)
	if err != nil {
		log.Fatalf("error while building assignment environment image: %v", err)
	}

}

// Generate a dockerfile for code runner server, build its image and push it ot docker hub.
func generateCodeRunnerImage(imgBuilder *builder.ImageBuilder) error {

	// Unmarshal the yaml configuration file and generate a dockerfileName.
	err, _, _ := utilities.WriteDockerfile(*codeRunnerConfig, *dockerfileName)
	if err != nil {
		return errors.Wrap(err, "error in writing dockerfile for code runner")
	}

	err = imgBuilder.BuildImage(false)
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
func generateAssignmentEnvImage(imgBuilder *builder.ImageBuilder) error {
	err, language, version := utilities.WriteDockerfile(*assignmentEnv, *dockerfileName)
	if err != nil {
		return errors.Wrap(err, "error in writing dockerfile for assignment environment")
	}

	imgBuilder.Language = language
	imgBuilder.Version = version
	err = imgBuilder.BuildImage(true)
	if err != nil {
		return err
	}

	return nil
}
