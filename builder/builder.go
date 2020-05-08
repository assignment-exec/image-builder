package builder

import (
	"assignment-exec/image-builder/constants"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

// Build Docker Image.
func (imgBuilder ImageBuilder) BuildImage() error {

	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Println("error in creating a docker client: ", err)
		return errors.Wrap(err, "error in creating a docker client")
	}

	// Create a build context tar for the image.
	// Build Context is the current working directory and where the Dockerfile is assumed to be located
	//[cite: https://docs.docker.com/develop/develop-images/dockerfile_best-practices/]..
	dockerFilename := imgBuilder.dockerfileName

	dockerBuildContext, err := imgBuilder.getDockerBuildContextTar()
	if err != nil {
		return err
	}

	response, err := dockerClient.ImageBuild(
		backgroundContext,
		dockerBuildContext,
		types.ImageBuildOptions{
			Dockerfile: dockerFilename,
			Tags:       []string{imgBuilder.imageTag}})
	if err != nil {
		return errors.Wrap(err, "error in building docker image")
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		return errors.Wrap(err, "error in reading image build response")
	}

	return nil
}

// Get Docker Build Context Tar Reader for building image.
func (imgBuilder ImageBuilder) getDockerBuildContextTar() (*os.File, error) {
	dockerFileReader, err := os.Open(imgBuilder.dockerfileName)
	if err != nil {
		return nil, errors.Wrap(err, "error in opening dockerfile for build context")
	}
	fileInfo, err := os.Stat(imgBuilder.dockerfileName)
	if err != nil {
		return nil, errors.Wrap(err, "error in verifying dockerfile path")
	}

	buildContextTar := new(archivex.TarFile)
	err = buildContextTar.Create(constants.BuildContextTar)
	if err != nil {
		return nil, errors.Wrap(err, "error in creating build context tar")
	}
	err = buildContextTar.AddAll(constants.InstallationScriptsDir, true)
	if err != nil {
		return nil, errors.Wrap(err, "error in adding installation script to build context tar")
	}

	err = buildContextTar.Add(imgBuilder.dockerfileName, dockerFileReader, fileInfo)
	if err != nil {
		return nil, errors.Wrap(err, "error in adding dockerfile to build context tar")
	}

	defer func() {
		err = buildContextTar.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	dockerBuildContext, err := os.Open(constants.BuildContextTar)
	if err != nil {
		return nil, errors.Wrap(err, "error in reading build context tar")
	}

	err = os.Remove(constants.BuildContextTar)
	if err != nil {
		return nil, errors.Wrap(err, "error in removing the build context tar file")
	}

	return dockerBuildContext, nil
}

func (imgBuilder ImageBuilder) PublishImage() error {
	// TODO: setup ssh keys for logging into docker hub

	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "error in creating new docker client")
	}

	authConfig := types.AuthConfig{
		Username: imgBuilder.authData.Username,
		Password: imgBuilder.authData.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return errors.Wrap(err, "error in encoding authConfig")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imageString := fmt.Sprintf("%s/%s", constants.DockerIO, imgBuilder.imageTag)

	resp, err := dockerClient.ImagePush(backgroundContext, imageString, types.ImagePushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return errors.Wrap(err, "error in pushing image to hub")
	}
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return errors.Wrap(err, "error in reading image push response")
	}
	defer func() {
		err = resp.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	return err

}
