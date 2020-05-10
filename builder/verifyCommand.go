package builder

import (
	"assignment-exec/image-builder/environment"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type verifyCommand struct {
	assgnEnv *assignmentEnv
}

func (cmd *verifyCommand) execute() error {
	return cmd.assgnEnv.verifyAndWrite()
}

func (cmd *verifyCommand) undo() error {
	//No operation.
	return nil
}

func (assgnEnv *assignmentEnv) verifyAndWrite() error {

	// Verify whether language image is present in registry.
	if err := assgnEnv.verifyLanguage(); err != nil {
		// If no then write the instructions from base image.
		assgnEnv.writeFromBaseImage()
	} else {
		if len(assgnEnv.AssgnEnvConfig.Deps.Libraries) > 0 {
			// Else write the instructions from dependencies.
			assgnEnv.writeFromDependencies()
		}
	}
	return nil
}

func (assgnEnv *assignmentEnv) verifyLanguage() error {
	// Check whether the language image is available on docker hub.
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	username, found := os.LookupEnv(environment.DockerAuthUsername)
	if !found {
		return errors.New("environment variable for username not set")
	}

	resp, err := dockerClient.ImageSearch(backgroundContext, username, types.ImageSearchOptions{
		Limit: 25})

	if err != nil {
		return err
	} else {
		found := false
		for _, result := range resp {
			if strings.Contains(assgnEnv.ImgBuildConfig.imageTag, result.Name) {
				found = true
			}
		}
		if !found {
			return errors.New("code-runner base image not found on docker registry")
		}
	}

	return nil
}

func (assgnEnv *assignmentEnv) writeFromBaseImage() {
	assgnEnv.DockerfileData.WriteString(assgnEnv.AssgnEnvConfig.WriteInstruction())
	// Append library names to image tag.
	for lib := range assgnEnv.AssgnEnvConfig.Deps.Libraries {
		assgnEnv.ImgBuildConfig.imageTag = strings.Join([]string{assgnEnv.ImgBuildConfig.imageTag, lib}, "-")
	}
}

func (assgnEnv *assignmentEnv) writeFromDependencies() {
	buf := &bytes.Buffer{}
	from := fmt.Sprintf("FROM %s", assgnEnv.ImgBuildConfig.imageTag)
	buf.WriteString(from)
	buf.WriteString("\n")

	for lib, installCmd := range assgnEnv.AssgnEnvConfig.Deps.Libraries {
		// Append library names to image tag.
		assgnEnv.ImgBuildConfig.imageTag = strings.Join([]string{assgnEnv.ImgBuildConfig.imageTag, lib}, "-")

		buf.WriteString("RUN " + installCmd.WriteInstruction() + " " + lib)
		buf.WriteString("\n")
	}

	assgnEnv.DockerfileData.WriteString(buf.String())
}
