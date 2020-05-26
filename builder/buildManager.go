package builder

import (
	"assignment-exec/image-builder/configurations"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

type BuildManager struct {
	commands     []command
	undoCommands *stack
}

type BuildManagerOption func(*BuildManager) error

func NewBuildManager(options ...BuildManagerOption) (*BuildManager, error) {
	b := &BuildManager{undoCommands: newStack()}
	for _, opt := range options {
		if err := opt(b); err != nil {
			return nil, errors.Wrap(err, "failed to create build manager instance")
		}
	}
	return b, nil
}

func WithCommands(assgnEnv *assignmentEnvironment) BuildManagerOption {
	return func(b *BuildManager) error {

		var commandList []command
		commandList = append(commandList,
			&verifyCommand{assgnEnv: assgnEnv},
			&writeDockerfileCommand{assgnEnv: assgnEnv},
			&buildCommand{assgnEnv: assgnEnv},
			&publishCommand{assgnEnv: assgnEnv})

		b.commands = commandList
		return nil
	}
}

func (builder *BuildManager) ExecuteCommands() error {
	for _, cmd := range builder.commands {
		builder.undoCommands.push(cmd)

		if err := cmd.execute(); err != nil {
			if undoErr := builder.UndoCommands(); undoErr != nil {
				log.Printf("error in undoing operations: %v", undoErr)
			}
			return err
		}
	}
	return nil
}

func (builder *BuildManager) UndoCommands() error {

	for !builder.undoCommands.isEmpty() {
		undoCmd := builder.undoCommands.pop()
		if err := undoCmd.undo(); err != nil {
			return err
		}
	}
	return nil
}

func GetConfigurations(publishImage bool, configFilepath string, dockerfileLoc string) (*assignmentEnvironment, error) {
	config, err := configurations.GetAssignmentEnvConfig(configFilepath)
	if err != nil {
		return nil, err
	}
	authData, err := getAuthData()
	if err != nil {
		return nil, errors.Wrap(err, "error while getting docker authentication data")
	}

	imageTag := fmt.Sprintf("%s/%s%s", authData.Username,
		config.Deps.Language.Name, config.Deps.Language.Version)

	imgBuilder, err := newImageBuildConfig(
		withDockerAuthData(authData),
		withImageTag(imageTag),
		withDockerfileLocation(dockerfileLoc),
		withPublishImageFlag(publishImage))

	if err != nil {
		return nil, errors.Wrap(err, "error in creating image builder instance for assignment env")
	}

	assgnEnv, err := newAssignmentEnvironment(
		withImageBuildCfg(imgBuilder),
		withAssgnEnvConfig(config))
	if err != nil {
		return nil, errors.Wrap(err, "error in creating assignment env instance")
	}

	return assgnEnv, nil

}
