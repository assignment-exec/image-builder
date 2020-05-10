package builder

import (
	"assignment-exec/image-builder/constants"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

type publishCommand struct {
	assgnEnv *assignmentEnv
}

func (cmd *publishCommand) execute() error {
	return cmd.assgnEnv.publish()
}

func (cmd *publishCommand) undo() error {
	err := cmd.assgnEnv.undoPublish()
	if err != nil {
		return errors.Wrap(err, "error in undo publish operation")
	}
	return nil
}

func (assgnEnv *assignmentEnv) publish() error {
	if assgnEnv.DockerfileData.Len() <= 0 {
		return assgnEnv.pullImage()
	} else if assgnEnv.ImgBuildConfig.publishImage {
		return assgnEnv.pushImage()
	}
	return nil
}

func (assgnEnv *assignmentEnv) pullImage() error {
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "error in creating new docker client")
	}
	authConfig := types.AuthConfig{
		Username: assgnEnv.ImgBuildConfig.authData.Username,
		Password: assgnEnv.ImgBuildConfig.authData.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return errors.Wrap(err, "error in encoding authConfig")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	resp, err := dockerClient.ImagePull(backgroundContext, assgnEnv.ImgBuildConfig.imageTag, types.ImagePullOptions{
		RegistryAuth: authStr,
	})

	if err != nil {
		return errors.Wrap(err, "error in pulling image from hub")
	}
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return errors.Wrap(err, "error in reading image pull response")
	}
	defer func() {
		err = resp.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	return err

}
func (assgnEnv *assignmentEnv) pushImage() error {
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return errors.Wrap(err, "error in creating new docker client")
	}

	authConfig := types.AuthConfig{
		Username: assgnEnv.ImgBuildConfig.authData.Username,
		Password: assgnEnv.ImgBuildConfig.authData.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return errors.Wrap(err, "error in encoding authConfig")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imageString := fmt.Sprintf("%s/%s", constants.DockerIO, assgnEnv.ImgBuildConfig.imageTag)

	resp, err := dockerClient.ImagePush(backgroundContext, imageString, types.ImagePushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		return errors.Wrap(err, "error in pushing image to hub")
	}
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		return errors.Wrap(err, "error in reading image push response")
	}
	defer func() {
		err = resp.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	return nil
}
