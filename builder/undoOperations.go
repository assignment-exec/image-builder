package builder

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"os"
)

func (assgnEnv *assignmentEnv) undoWrite() {
	// Clear the dockerfile data bytes.
	assgnEnv.DockerfileData.Reset()
}

func (assgnEnv *assignmentEnv) undoBuild() error {
	// Delete the created Dockerfile.
	return os.Remove(assgnEnv.ImgBuildConfig.dockerfileName)
}

func (assgnEnv *assignmentEnv) undoPublish() error {
	// Delete the built image.
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	_, err = dockerClient.ImageRemove(backgroundContext, assgnEnv.ImgBuildConfig.imageTag, types.ImageRemoveOptions{
		Force: true})

	if err != nil {
		return err
	}

	return nil
}
