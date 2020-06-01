// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub.
package builder

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/environment"
	"github.com/jhoonb/archivex"
	"github.com/pkg/errors"
	"log"
	"os"
)

// dockerAuthData struct type holds username
// and password for docker hub authentication.
type dockerAuthData struct {
	Username string
	Password string
}

// imageBuildConfig struct type holds docker authentication data,
// image tag, dockerfile location to be created, publishImage image flag.
// All required to build assignment environment image.
type imageBuildConfig struct {
	authData      *dockerAuthData
	imageTag      string
	dockerfileLoc string
	publishImage  bool
}

// imageBuildConfigOption represents options that can be used to help initialize
// an instance of imageBuildConfig.
// Each option is a closure that is responsible for initializing one or more members
// while instantiating imageBuildConfig.
type imageBuildConfigOption func(*imageBuildConfig) error

// newImageBuildConfig constructs an instance of imageBuildConfig
// by applying each of the provided options.
// The construction of the object fails upon the failure of at least one of the given options.
func newImageBuildConfig(options ...imageBuildConfigOption) (*imageBuildConfig, error) {
	imgBuildCfg := &imageBuildConfig{}
	for _, opt := range options {
		if err := opt(imgBuildCfg); err != nil {
			return nil, errors.Wrap(err, "failed to create imageBuildConfig instance")
		}
	}
	return imgBuildCfg, nil
}

// withDockerfileLocation returns an imageBuildConfigOption for initializing the
// dockerfile location.
func withDockerfileLocation(fileLoc string) imageBuildConfigOption {
	return func(imgBuildCfg *imageBuildConfig) error {
		// Validate fileLoc and raise error if validation fails.
		if fileLoc == "" {
			return errors.New("dockerfile name not provided")
		}
		imgBuildCfg.dockerfileLoc = fileLoc
		return nil

	}
}

// withDockerAuthData returns an imageBuildConfigOption for initializing docker
// authentication data.
func withDockerAuthData(authData *dockerAuthData) imageBuildConfigOption {
	return func(imgBuildCfg *imageBuildConfig) error {
		// Validate authData and raise error if validation fails.
		if authData == nil {
			return errors.New("docker authentication data not provided")
		}
		imgBuildCfg.authData = authData
		return nil
	}
}

// withImageTag returns an imageBuildConfigOption for initializing image tag.
func withImageTag(tag string) imageBuildConfigOption {
	return func(imgBuildCfg *imageBuildConfig) error {
		// Validate tag and raise error if validation fails.
		if tag == "" {
			return errors.New("image tag is empty")
		}
		imgBuildCfg.imageTag = tag
		return nil
	}
}

// withPublishImageFlag returns an imageBuildConfigOption for initializing publishImage flag.
func withPublishImageFlag(publishImage bool) imageBuildConfigOption {
	return func(imgBuildCfg *imageBuildConfig) error {
		imgBuildCfg.publishImage = publishImage
		return nil
	}
}

// getAuthData reads the docker authentication data, i.e username
// and password from environment variables.
func getAuthData() (*dockerAuthData, error) {
	username, hasFound := os.LookupEnv(environment.DockerAuthUsername)
	if !hasFound {
		return nil, errors.New("environment variable for username not set")
	}
	password, hasFound := os.LookupEnv(environment.DockerAuthPassword)
	if !hasFound {
		return nil, errors.New("environment variable for password not set")
	}

	c := &dockerAuthData{Username: username, Password: password}
	return c, nil
}

// getDockerBuildContextTar creates a tar file for docker build context.
// The tar holds Dockerfile and installation scripts that are required for
// building the assignment environment image.
func (imgBuildCfg imageBuildConfig) getDockerBuildContextTar() (*os.File, error) {

	// Gets the reader for the Dockerfile.
	dockerFileReader, err := os.Open(imgBuildCfg.dockerfileLoc)
	if err != nil {
		return nil, errors.Wrap(err, "error in opening dockerfile for build context")
	}
	fileInfo, err := os.Stat(imgBuildCfg.dockerfileLoc)
	if err != nil {
		return nil, errors.Wrap(err, "error in verifying dockerfile path")
	}

	// Creates a new tar file named as `buildContext.tar`, add the installation scripts stored in
	// `scripts` directory and add the Dockerfile.
	buildContextTar := new(archivex.TarFile)
	err = buildContextTar.Create(constants.BuildContextTar)
	if err != nil {
		return nil, errors.Wrap(err, "error in creating build context tar")
	}
	err = buildContextTar.AddAll(constants.InstallationScriptsDir, true)
	if err != nil {
		return nil, errors.Wrap(err, "error in adding installation script to build context tar")
	}

	err = buildContextTar.Add(imgBuildCfg.dockerfileLoc, dockerFileReader, fileInfo)
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

	// Reads the build context tar.
	dockerBuildContext, err := os.Open(constants.BuildContextTar)
	if err != nil {
		return nil, errors.Wrap(err, "error in reading build context tar")
	}

	// Removes tha tar file as it is no longer needed.
	err = os.Remove(constants.BuildContextTar)
	if err != nil {
		return nil, errors.Wrap(err, "error in removing the build context tar file")
	}

	return dockerBuildContext, nil
}
