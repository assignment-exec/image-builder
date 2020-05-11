package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/utilities"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

var publishImage = flag.Bool("publishImage", false, "Push image to docker hub")
var assignmentEnv = flag.String("assignmentEnv", "assignment-env.yaml", "Assignment Environment configuration filename")
var dockerfileName = flag.String("dockerfile", "Dockerfile", "Dockerfile name")

func main() {

	flag.Parse()

	authData, err := builder.GetAuthData()
	if err != nil {
		log.Fatalf("error while getting docker authentication data: %v", err)
	}

	err = generateAssignmentEnvImage(authData)
	if err != nil {
		log.Fatalf("error in building assignment environment image: %v", err)
	}

}

// Generate a dockerfile for assignment environment and build its image.
func generateAssignmentEnvImage(authData *builder.DockerAuthData) error {
	err, language, version := utilities.WriteDockerfileForAssignmentEnv(*assignmentEnv, *dockerfileName)
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
