// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub.
package builder

import (
	"assignment-exec/image-builder/configurations"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

// BuildManager struct type holds array of commands to execute
// and a stack for commands to perform the corresponding undo.
type BuildManager struct {
	commands     []command
	undoCommands *stack
}

// BuildManagerOption represents options that can be used to help initialize
// an instance of BuildManager.
// Each option is a closure that is responsible for initializing one or more members
// while instantiating BuildManager.
type BuildManagerOption func(*BuildManager) error

// NewBuildManager constructs an instance of BuildManager
// by applying each of the provided options.
// The construction of the object fails upon the failure of at least one of the given options.
func NewBuildManager(options ...BuildManagerOption) (*BuildManager, error) {
	b := &BuildManager{undoCommands: newStack()}
	for _, opt := range options {
		if err := opt(b); err != nil {
			return nil, errors.Wrap(err, "failed to create build manager instance")
		}
	}
	return b, nil
}

// WithCommands returns a BuildManagerOption for initializing the commands.
func WithCommands(asgmtEnv *assignmentEnvironmentImageBuilder) BuildManagerOption {
	return func(b *BuildManager) error {

		var commandList []command
		commandList = append(commandList,
			&verifyCommand{asgmtEnv: asgmtEnv},
			&writeDockerfileCommand{asgmtEnv: asgmtEnv},
			&buildCommand{asgmtEnv: asgmtEnv},
			&publishCommand{asgmtEnv: asgmtEnv})

		b.commands = commandList
		return nil
	}
}

// ExecuteCommands invokes execute function for all commands sequentially.
// If error is encountered in any command execution then perform undo operations in
// the reverse order of execution.
func (builder *BuildManager) ExecuteCommands() error {
	for _, cmd := range builder.commands {
		builder.undoCommands.push(cmd)

		if err := cmd.execute(); err != nil {
			if undoErr := builder.UndoCommands(); undoErr != nil {
				// Logs the error encountered during undoing command execution.
				log.Printf("error in undoing operations: %v", undoErr)
			}
			// Returns the error encountered while executing commands.
			return err
		}
	}
	return nil
}

// UndoCommands pops the all commands from stack and invokes its
// respective undo function.
func (builder *BuildManager) UndoCommands() error {

	for !builder.undoCommands.isEmpty() {
		undoCmd := builder.undoCommands.pop()
		if err := undoCmd.undo(); err != nil {
			return err
		}
	}
	return nil
}

// GetConfigurations takes image publishImage flag, assignment environment configuration file path
// and dockerfile location, reads the config file, sets the imageBuildConfig instance,
// sets the assignmentEnvironmentImageBuilder instance.
// It returns the assignmentEnvironmentImageBuilder instance and any error encountered.
func GetConfigurations(publishImage bool, configFilepath string, dockerfileLoc string) (*assignmentEnvironmentImageBuilder, error) {
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

	asgmtEnv, err := newAssignmentEnvironmentImageBuilder(
		withImageBuildCfg(imgBuilder),
		withAsgmtEnvConfig(config))
	if err != nil {
		return nil, errors.Wrap(err, "error in creating assignment env instance")
	}

	return asgmtEnv, nil

}
