// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub. It uses command pattern to
// perform all operations and perform undo operations when any error is encountered.
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
// image tag, dockerfile location to be created, publish image flag.
// All required to build assignment environment image.
type imageBuildConfig struct {
	authData      *dockerAuthData
	imageTag      string
	dockerfileLoc string
	isPublish     bool
}

// imageBuildConfigOption is a function interface that
// is supplied as different options while creating new instance of
// 'imageBuildConfig' type. This function returns any error encountered.
type imageBuildConfigOption func(*imageBuildConfig) error

// newImageBuildConfig takes one or more options and
// returns new instance of imageBuildConfig.
func newImageBuildConfig(options ...imageBuildConfigOption) (*imageBuildConfig, error) {
	imgBuildCfg := &imageBuildConfig{}
	for _, opt := range options {
		if err := opt(imgBuildCfg); err != nil {
			return nil, errors.Wrap(err, "failed to create imageBuildConfig instance")
		}
	}
	return imgBuildCfg, nil
}

// withDockerfileLocation is used as an option while creating imageBuildConfig instance. It takes
// docker file location as a parameter and returns 'imageBuildConfigOption' function.
// This returned function in turn sets the docker file location within imageBuildConfig instance.
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

// withDockerAuthData is used as an option while creating imageBuildConfig instance. It takes
// 'dockerAuthData' instance as a parameter and returns 'imageBuildConfigOption' function.
// This returned function in turn sets the 'dockerAuthData' within imageBuildConfig instance.
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

// withImageTag is used as an option while creating imageBuildConfig instance. It takes
// image tag string as a parameter and returns 'imageBuildConfigOption' function.
// This returned function in turn sets the 'imageTag' within imageBuildConfig instance.
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

// withPublishImageFlag is used as an option while creating imageBuildConfig instance. It takes
// publish image flag as a parameter and returns 'imageBuildConfigOption' function.
// This returned function in turn sets the 'isPublish' flag within imageBuildConfig instance.
func withPublishImageFlag(isPublish bool) imageBuildConfigOption {
	return func(imgBuildCfg *imageBuildConfig) error {
		imgBuildCfg.isPublish = isPublish
		return nil
	}
}

// getAuthData reads the docker authentication data, i.e username
// and password from environment variables. It returns the 'dockerAuthData'
// instance and any error encountered while reading from environment variables.
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
// It returns the read tar file and any error encountered while creating and reading tar.
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
