package builder

import (
	"assignment-exec/image-builder/environment"
	"github.com/pkg/errors"
	"os"
)

type DockerAuthData struct {
	Username string
	Password string
}

type ImageBuilder struct {
	authData       *DockerAuthData
	imageTag       string
	dockerfileName string
}

type ImageBuilderOption func(*ImageBuilder) error

func NewImageBuilder(options ...ImageBuilderOption) (*ImageBuilder, error) {
	b := &ImageBuilder{}
	for _, opt := range options {
		if err := opt(b); err != nil {
			return nil, errors.Wrap(err, "failed to build ImageBuilder instance")
		}
	}
	return b, nil
}

func WithDockerfileName(filename string) ImageBuilderOption {
	return func(b *ImageBuilder) error {
		// Validate filename and raise error if validation fails.
		_, err := os.Stat(filename)
		if err != nil {
			return errors.Wrap(err, "dockerfile not found")
		}
		b.dockerfileName = filename
		return nil

	}
}

func WithDockerAuthData(authData *DockerAuthData) ImageBuilderOption {
	return func(b *ImageBuilder) error {
		// Validate authData and raise error if validation fails.
		if authData == nil {
			return errors.New("docker authentication data not provided")
		}
		b.authData = authData
		return nil
	}
}

func WithImageTag(tag string) ImageBuilderOption {
	return func(b *ImageBuilder) error {
		// Validate tag and raise error if validation fails.
		if tag == "" {
			return errors.New("image tag is empty")
		}
		b.imageTag = tag
		return nil
	}
}

// Get the docker authentication details.
func GetAuthData() *DockerAuthData {
	c := &DockerAuthData{}
	c.Username = os.Getenv(environment.DockerAuthUsername)
	c.Password = os.Getenv(environment.DockerAuthPassword)
	return c
}
