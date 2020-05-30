// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub.
package builder

import (
	"assignment-exec/image-builder/configurations"
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/environment"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// assignmentEnvironmentImageBuilder is used to build an image with an environment
// that satisfies all the language and library dependencies of an assignment.
//
// The image is built in phases to allow for rollback in situations where the execution of a phase fails.
// Once the configuration passes a verification phase, a Dockerfile is created and
// the docker client (here - https://github.com/moby/moby/tree/master/client) is used to generate the image.
// The image is then pushed to docker hub repository, if required.
type assignmentEnvironmentImageBuilder struct {
	DockerfileInstructions bytes.Buffer
	ImgBuildConfig         *imageBuildConfig
	AsgmtEnvConfig         *configurations.AssignmentEnvConfig
	ImageExists            bool
}

// assignmentEnvironmentImageBuilderOption represents options that can be used to help initialize
// an instance of assignmentEnvironmentImageBuilder.
// Each option is a closure that is responsible for initializing one or more members
// while instantiating assignmentEnvironmentImageBuilder.
type assignmentEnvironmentImageBuilderOption func(*assignmentEnvironmentImageBuilder) error

// newAssignmentEnvironmentImageBuilder constructs an instance of assignmentEnvironmentImageBuilder
// by applying each of the provided options.
// The construction of the object fails upon the failure of at least one of the given options.
func newAssignmentEnvironmentImageBuilder(options ...assignmentEnvironmentImageBuilderOption) (*assignmentEnvironmentImageBuilder, error) {
	asgmtEnv := &assignmentEnvironmentImageBuilder{}
	for _, opt := range options {
		if err := opt(asgmtEnv); err != nil {
			return nil, errors.Wrap(err, "failed to create assignmentEnvironmentImageBuilder instance")
		}
	}
	return asgmtEnv, nil
}

// withImageBuildCfg returns an assignmentEnvironmentImageBuilderOption for
// initializing the image build config.
func withImageBuildCfg(imgBuildCfg *imageBuildConfig) assignmentEnvironmentImageBuilderOption {
	return func(asgmtEnv *assignmentEnvironmentImageBuilder) error {
		if imgBuildCfg == nil {
			return errors.New("image build config instance not provided")
		}
		asgmtEnv.ImgBuildConfig = imgBuildCfg
		return nil

	}
}

// withAsgmtEnvConfig returns an assignmentEnvironmentImageBuilderOption for initializing
// the assignment environment configuration.
func withAsgmtEnvConfig(assignCfgs *configurations.AssignmentEnvConfig) assignmentEnvironmentImageBuilderOption {
	return func(asgmtEnv *assignmentEnvironmentImageBuilder) error {
		if assignCfgs == nil {
			return errors.New("assignment environment configurations not provided")
		}
		asgmtEnv.AsgmtEnvConfig = assignCfgs
		return nil

	}
}

// verifyAndWriteInstructions checks whether a docker image for the language given in
// configuration file is already present in docker hub and accordingly writes the Dockerfile
// from either the base image or from the existing language image.
func (asgmtEnv *assignmentEnvironmentImageBuilder) verifyAndWriteInstructions() error {

	// Verify whether language image is present in registry.
	if err := asgmtEnv.verifyLanguage(); err != nil {
		// If no then write the instructions from base image.
		asgmtEnv.writeInstructionsLayerOnBaseImage()
	} else {
		if len(asgmtEnv.AsgmtEnvConfig.Deps.Libraries) > 0 {
			// Else write the instructions from dependencies.
			asgmtEnv.writeInstructionsLayerOnLanguageImage()
		}
	}

	if asgmtEnv.DockerfileInstructions.Len() <= 0 {
		asgmtEnv.ImageExists = true
	}
	return nil
}

// verifyLanguage searches the docker image for the given language on docker hub.
func (asgmtEnv *assignmentEnvironmentImageBuilder) verifyLanguage() error {
	// Check whether the language image is available on docker hub.
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	username, usernameFound := os.LookupEnv(environment.DockerAuthUsername)
	if !usernameFound {
		return errors.New("environment variable for username not set")
	}

	response, err := dockerClient.ImageSearch(backgroundContext, username, types.ImageSearchOptions{
		Limit: 25})

	if err != nil {
		return err
	}
	langImageFound := false
	for _, result := range response {
		if strings.Contains(asgmtEnv.ImgBuildConfig.imageTag, result.Name) {
			langImageFound = true
		}
	}
	if !langImageFound {
		return errors.New("language image not found on docker hub")
	}

	return nil
}

// writeInstructionsLayerOnBaseImage writes the docker instructions to starting from the
// base code runner image. Which is then followed by the required language and its dependencies.
func (asgmtEnv *assignmentEnvironmentImageBuilder) writeInstructionsLayerOnBaseImage() {
	var instructions []string
	instructions = append(instructions, asgmtEnv.AsgmtEnvConfig.GetInstruction())

	var libraryNames []string
	for lib := range asgmtEnv.AsgmtEnvConfig.Deps.Libraries {
		libraryNames = append(libraryNames, lib)
	}

	// Generate the image tag.
	asgmtEnv.ImgBuildConfig.imageTag = strings.Join([]string{asgmtEnv.ImgBuildConfig.imageTag,
		strings.Join(libraryNames, "-")}, "-")
	asgmtEnv.DockerfileInstructions.WriteString(strings.Join(instructions, "\n"))
}

// writeInstructionsLayerOnLanguageImage writes the docker instructions starting
// from the respective language image. Which is then followed by the language dependencies.
func (asgmtEnv *assignmentEnvironmentImageBuilder) writeInstructionsLayerOnLanguageImage() {
	var instructions []string

	// FROM instruction.
	instructions = append(instructions, fmt.Sprintf("FROM %s", asgmtEnv.ImgBuildConfig.imageTag))
	// COPY instruction.
	instructions = append(instructions, fmt.Sprintf("COPY . /"+constants.CodeRunnerDir))

	var libraryNames []string
	for lib, installCmd := range asgmtEnv.AsgmtEnvConfig.Deps.Libraries {
		libraryNames = append(libraryNames, lib)

		// RUN instruction.
		instructions = append(instructions, "RUN "+installCmd.GetInstruction()+" "+lib)
	}

	// Generate the image tag.
	asgmtEnv.ImgBuildConfig.imageTag = strings.Join([]string{asgmtEnv.ImgBuildConfig.imageTag,
		strings.Join(libraryNames, "-")}, "-")
	asgmtEnv.DockerfileInstructions.WriteString(strings.Join(instructions, "\n"))
}

// writeToDockerfile creates a Dockerfile at the specified location and writes
// the docker instructions to it.
func (asgmtEnv *assignmentEnvironmentImageBuilder) writeToDockerfile() error {
	if !asgmtEnv.ImageExists {

		file, err := os.Create(asgmtEnv.ImgBuildConfig.dockerfileLoc)
		defer func() {
			err = file.Close()
			if err != nil {
				log.Println("error while closing the created Dockerfile", err)
				return
			}
		}()

		if err != nil {
			return err
		}
		_, err = file.WriteString(asgmtEnv.DockerfileInstructions.String())
		return err
	}
	return nil
}

// build a docker image for the given language. If the image is already present,
// then it simply pull the image.
func (asgmtEnv *assignmentEnvironmentImageBuilder) build() error {

	if !asgmtEnv.ImageExists {
		backgroundContext := context.Background()
		dockerClient, err := client.NewEnvClient()
		if err != nil {
			return errors.Wrap(err, "error in creating a docker client")
		}

		// Create a build context tar for the image.
		// build Context is the current working directory and where the Dockerfile is assumed to be located.
		// [cite: https://docs.docker.com/develop/develop-images/dockerfile_best-practices/].
		dockerfileLoc := filepath.Base(asgmtEnv.ImgBuildConfig.dockerfileLoc)

		dockerBuildContext, err := asgmtEnv.ImgBuildConfig.getDockerBuildContextTar()
		if err != nil {
			return err
		}

		response, err := dockerClient.ImageBuild(
			backgroundContext,
			dockerBuildContext,
			types.ImageBuildOptions{
				Dockerfile: dockerfileLoc,
				Tags:       []string{asgmtEnv.ImgBuildConfig.imageTag}})
		if err != nil {
			return errors.Wrap(err, "error in building docker image")
		}
		defer func() {
			err := response.Body.Close()
			if err != nil {
				log.Println(err)
				return
			}
		}()

		_, err = io.Copy(os.Stdout, response.Body)
		if err != nil {
			return errors.Wrap(err, "error in reading image build response")
		}

		return err
	} else {
		return asgmtEnv.pullImage()
	}
}

// publishImage pushes the built image to docker hub, if required, if it is not already present.
func (asgmtEnv *assignmentEnvironmentImageBuilder) publishImage() error {
	if !asgmtEnv.ImageExists && asgmtEnv.ImgBuildConfig.publishImage {
		backgroundContext := context.Background()
		dockerClient, err := client.NewEnvClient()
		if err != nil {
			return errors.Wrap(err, "error in creating new docker client")
		}

		authConfig := types.AuthConfig{
			Username: asgmtEnv.ImgBuildConfig.authData.Username,
			Password: asgmtEnv.ImgBuildConfig.authData.Password,
		}
		authJson, err := json.Marshal(authConfig)
		if err != nil {
			return errors.Wrap(err, "error in encoding authConfig")
		}

		authString := base64.URLEncoding.EncodeToString(authJson)

		imageString := fmt.Sprintf("%s/%s", constants.DockerIO, asgmtEnv.ImgBuildConfig.imageTag)

		response, err := dockerClient.ImagePush(backgroundContext, imageString, types.ImagePushOptions{
			RegistryAuth: authString,
		})
		if err != nil {
			return errors.Wrap(err, "error in pushing image to hub")
		}
		_, err = io.Copy(os.Stdout, response)
		if err != nil {
			return errors.Wrap(err, "error in reading image push response")
		}
		defer func() {
			err = response.Close()
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
	return nil
}

// pullImage pulls the required docker image for given language from docker hub.
func (asgmtEnv *assignmentEnvironmentImageBuilder) pullImage() error {
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "error in creating new docker client")
	}
	authConfig := types.AuthConfig{
		Username: asgmtEnv.ImgBuildConfig.authData.Username,
		Password: asgmtEnv.ImgBuildConfig.authData.Password,
	}
	authJson, err := json.Marshal(authConfig)
	if err != nil {
		return errors.Wrap(err, "error in encoding authConfig")
	}

	authString := base64.URLEncoding.EncodeToString(authJson)
	imageString := fmt.Sprintf("%s/%s", constants.DockerIO, asgmtEnv.ImgBuildConfig.imageTag)
	response, err := dockerClient.ImagePull(backgroundContext, imageString, types.ImagePullOptions{
		RegistryAuth: authString,
	})

	if err != nil {
		return errors.Wrap(err, "error in pulling image from hub")
	}
	_, err = io.Copy(os.Stdout, response)
	if err != nil {
		return errors.Wrap(err, "error in reading image pull response")
	}
	defer func() {
		err = response.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	return nil

}

// resetDockerfileData resets the Dockerfile instructions buffer.
func (asgmtEnv *assignmentEnvironmentImageBuilder) resetDockerfileData() {
	// Clear the Dockerfile data.
	asgmtEnv.DockerfileInstructions.Reset()
}

// deleteDockerfile deletes the dockerfile that was created
// while building the assignment environment image.
func (asgmtEnv *assignmentEnvironmentImageBuilder) deleteDockerfile() error {
	// Delete the created Dockerfile.
	_, err := os.Stat(asgmtEnv.ImgBuildConfig.dockerfileLoc)
	if err == nil {
		return os.Remove(asgmtEnv.ImgBuildConfig.dockerfileLoc)
	}
	return nil
}

// undoBuild removes the assignment environment image that was built locally.
func (asgmtEnv *assignmentEnvironmentImageBuilder) undoBuild() error {
	// Delete the built image.
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	_, err = dockerClient.ImageRemove(backgroundContext, asgmtEnv.ImgBuildConfig.imageTag, types.ImageRemoveOptions{
		Force: true})

	return err
}
