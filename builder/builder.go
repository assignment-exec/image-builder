package builder

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/environment"
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

type dockerAuthData struct {
	Username   string
	Password   string
	Repository string
	Version    string
}

type ImageBuilder struct {
	authData       *dockerAuthData
	dockerFilename string
	Language       string
	Version        string
}

// New Image Builder instance.
func NewImageBuilder(dockerFilename string) *ImageBuilder {
	dockerAuth := getAuthData()
	imgBuilder := &ImageBuilder{authData: dockerAuth, dockerFilename: dockerFilename}
	return imgBuilder
}

// Get the docker authentication details.
func getAuthData() *dockerAuthData {
	c := &dockerAuthData{}
	c.Username = os.Getenv(environment.DockerAuthUsername)
	c.Password = os.Getenv(environment.DockerAuthPassword)
	c.Repository = os.Getenv(environment.DockerAuthRepository)
	c.Version = os.Getenv(environment.DockerAuthRepositoryVersion)
	return c
}

// Build Docker Image.
func (imgBuilder ImageBuilder) BuildImage(assignmentEnv bool) error {

	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Println("error in creating a docker client: ", err)
		return errors.Wrap(err, "error in creating a docker client")
	}

	// Create a build context tar for the image.
	// Build Context is the current working directory and where the Dockerfile is assumed to be located
	//[cite: https://docs.docker.com/develop/develop-images/dockerfile_best-practices/]..
	dockerFilename := imgBuilder.dockerFilename

	dockerBuildContext, err := imgBuilder.getDockerBuildContextTar()
	if err != nil {
		return err
	}

	tagName := imgBuilder.getTagName(assignmentEnv)
	response, err := dockerClient.ImageBuild(
		backgroundContext,
		dockerBuildContext,
		types.ImageBuildOptions{
			Dockerfile: dockerFilename,
			Tags:       []string{tagName}})
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

// Get Tag name for building image.
func (imgBuilder ImageBuilder) getTagName(assignmentEnv bool) string {
	var tagName string
	if assignmentEnv {
		tagName = fmt.Sprintf("%s/%s%s", imgBuilder.authData.Username, imgBuilder.Language, imgBuilder.Version)
	} else {
		tagName = fmt.Sprintf("%s/%s:%s", imgBuilder.authData.Username,
			imgBuilder.authData.Repository, imgBuilder.authData.Version)
	}

	return tagName
}

// Get Docker Build Context Tar Reader for building image.
func (imgBuilder ImageBuilder) getDockerBuildContextTar() (*os.File, error) {
	dockerFileReader, err := os.Open(imgBuilder.dockerFilename)
	if err != nil {
		return nil, errors.Wrap(err, "error in opening dockerfile for build context")
	}
	fileInfo, err := os.Stat(imgBuilder.dockerFilename)
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

	err = buildContextTar.Add(imgBuilder.dockerFilename, dockerFileReader, fileInfo)
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

	imageString := fmt.Sprintf("%s/%s/%s:%s", constants.DockerIO, imgBuilder.authData.Username,
		imgBuilder.authData.Repository, imgBuilder.authData.Version)

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
