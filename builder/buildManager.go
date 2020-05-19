// Package builder implements routines to write dockerfile for assignment environment,
// build its docker image and publish it to docker hub. It uses command pattern to
// perform all operations and perform undo operations when any error is encountered.
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

// BuildManagerOption is a function interface that
// is supplied as different options while creating new instance of
// 'BuildManager' type. This function returns any error encountered.
type BuildManagerOption func(*BuildManager) error

// NewBuildManager takes one or more options and
// returns new instance of BuildManager.
func NewBuildManager(options ...BuildManagerOption) (*BuildManager, error) {
	b := &BuildManager{undoCommands: newStack()}
	for _, opt := range options {
		if err := opt(b); err != nil {
			return nil, errors.Wrap(err, "failed to create build manager instance")
		}
	}
	return b, nil
}

// WithCommands is used as an option while creating BuildManager instance. It takes
// assignmentEnvironment as a parameter and returns 'BuildManagerOption' function.
// This returned function in turn creates a new commands array, sets the assignmentEnvironment
// instance for every command and assigns this command array to BuildManager.
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

// ExecuteCommands invokes execute function for all commands sequentially.
// Every command is pushed to the 'undoCommands' stack to keep track of order of execution.
// If error is encountered in any command execution then perform all undo operations from the stack.
// It returns error encountered in execution of commands.
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
// respective undo function. It returns any error encountered
// while undoing command execution.
func (builder *BuildManager) UndoCommands() error {

	for !builder.undoCommands.isEmpty() {
		undoCmd := builder.undoCommands.pop()
		if err := undoCmd.undo(); err != nil {
			return err
		}
	}
	return nil
}

// GetConfigurations takes image publish flag, assignment environment configuration file path
// and dockerfile location, reads the config file, sets the imageBuildConfig instance,
// sets the assignmentEnvironment instance.
// It returns the assignmentEnvironment instance and any error encountered.
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
