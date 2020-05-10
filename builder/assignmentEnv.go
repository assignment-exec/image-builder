package builder

import (
	"assignment-exec/image-builder/configurations"
	"bytes"
	"github.com/pkg/errors"
)

type baseEnv interface {
	verifyAndWrite() error
	build() error
	publish() error
	undoWrite()
	undoBuild() error
	undoPublish() error
}

type assignmentEnv struct {
	DockerfileData bytes.Buffer
	ImgBuildConfig *imageBuildConfig
	AssgnEnvConfig *configurations.AssignmentEnvConfig
}

type assignmentImageOption func(*assignmentEnv) error

func newAssignmentImage(options ...assignmentImageOption) (*assignmentEnv, error) {
	assgnEnv := &assignmentEnv{}
	for _, opt := range options {
		if err := opt(assgnEnv); err != nil {
			return nil, errors.Wrap(err, "failed to create assignmentEnv instance")
		}
	}
	return assgnEnv, nil
}

func withImageBuildCfg(imgBuildCfg *imageBuildConfig) assignmentImageOption {
	return func(assgnEnv *assignmentEnv) error {
		if imgBuildCfg == nil {
			return errors.New("image build config instance not provided")
		}
		assgnEnv.ImgBuildConfig = imgBuildCfg
		return nil

	}
}

func withAssgnEnvConfig(configurations *configurations.AssignmentEnvConfig) assignmentImageOption {
	return func(assgnEnv *assignmentEnv) error {
		if configurations == nil {
			return errors.New("assignment configurations not provided")
		}
		assgnEnv.AssgnEnvConfig = configurations
		return nil

	}
}
