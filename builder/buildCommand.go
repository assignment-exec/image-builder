package builder

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

type buildCommand struct {
	assgnEnv *assignmentEnv
}

func (cmd *buildCommand) execute() error {
	return cmd.assgnEnv.build()
}

func (cmd *buildCommand) undo() error {
	if err := cmd.assgnEnv.undoBuild(); err != nil {
		return errors.Wrap(err, "error in undo build operation")
	}
	return nil
}

func (assgnEnv *assignmentEnv) build() error {

	if assgnEnv.DockerfileData.Len() > 0 {
		backgroundContext := context.Background()
		dockerClient, err := client.NewEnvClient()
		if err != nil {
			return errors.Wrap(err, "error in creating a docker client")
		}

		// Create a build context tar for the image.
		// build Context is the current working directory and where the Dockerfile is assumed to be located.
		// [cite: https://docs.docker.com/develop/develop-images/dockerfile_best-practices/].
		dockerFilename := assgnEnv.ImgBuildConfig.dockerfileName

		dockerBuildContext, err := assgnEnv.ImgBuildConfig.getDockerBuildContextTar()
		if err != nil {
			return err
		}

		response, err := dockerClient.ImageBuild(
			backgroundContext,
			dockerBuildContext,
			types.ImageBuildOptions{
				Dockerfile: dockerFilename,
				Tags:       []string{assgnEnv.ImgBuildConfig.imageTag}})
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
	}
	return nil
}
