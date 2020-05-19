package builder

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/environment"
	"github.com/jhoonb/archivex"
	"github.com/pkg/errors"
	"log"
	"os"
)

type dockerAuthData struct {
	Username string
	Password string
}

type imageBuildConfig struct {
	authData      *dockerAuthData
	imageTag      string
	dockerfileLoc string
	isPublish     bool
}

type imageBuildConfigOption func(*imageBuildConfig) error

func newImageBuildConfig(options ...imageBuildConfigOption) (*imageBuildConfig, error) {
	imgBuildCfg := &imageBuildConfig{}
	for _, opt := range options {
		if err := opt(imgBuildCfg); err != nil {
			return nil, errors.Wrap(err, "failed to create imageBuildConfig instance")
		}
	}
	return imgBuildCfg, nil
}

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

func withPublishImageFlag(isPublish bool) imageBuildConfigOption {
	return func(imgBuildCfg *imageBuildConfig) error {
		imgBuildCfg.isPublish = isPublish
		return nil
	}
}

// Get the docker authentication details.
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

// Get Docker build Context Tar Reader for building image.
func (imgBuildCfg imageBuildConfig) getDockerBuildContextTar() (*os.File, error) {
	dockerFileReader, err := os.Open(imgBuildCfg.dockerfileLoc)
	if err != nil {
		return nil, errors.Wrap(err, "error in opening dockerfile for build context")
	}
	fileInfo, err := os.Stat(imgBuildCfg.dockerfileLoc)
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
