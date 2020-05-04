package builder

import (
	"archive/tar"
	"assignment-exec/image-builder/environment"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type dockerAuthData struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}

// Get the docker authentication details.
func GetAuthData(filename string) (*dockerAuthData, error) {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log config file")
	}

	c := &dockerAuthData{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("error in unmarshalling yaml: %v", err)
	}

	return c, nil
}

// Builds Docker Image for Code Runner.
func BuildImage(authData dockerAuthData) error {

	buildContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("error in creating a docker client: %v", err)
	}

	// Create a tar for build context.
	tarBuffer := new(bytes.Buffer)
	writer := tar.NewWriter(tarBuffer)
	defer func() {
		err = writer.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	dockerFilepath := environment.DockerFilepath
	dockerFileReader, err := os.Open(dockerFilepath)
	if err != nil {
		log.Fatalf(" unable to open dockerfile: %v", err)
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatalf("unable to read dockerfile: %v", err)
	}

	tarHeader := &tar.Header{
		Name: dockerFilepath,
		Size: int64(len(readDockerFile)),
	}
	err = writer.WriteHeader(tarHeader)
	if err != nil {
		log.Fatalf("unable to write tar header: %v", err)
	}
	_, err = writer.Write(readDockerFile)
	if err != nil {
		log.Fatalf("unable to write tar body: %v", err)
	}

	// Use the tar of the Dockerfile while building image.
	dockerFileTarReader := bytes.NewReader(tarBuffer.Bytes())

	response, err := dockerClient.ImageBuild(
		buildContext,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Dockerfile: dockerFilepath,
			Tags:       []string{fmt.Sprintf("%s/%s:%s", authData.Username, authData.Repository, authData.Version)}})
	if err != nil {
		log.Fatalf("unable to build docker image: %v", err)
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		log.Fatalf("unable to read image build response: %v", err)
	}

	return err
}

func PushImageToHub(authData dockerAuthData) error {

	// TODO: setup ssh keys for logging into docker hub

	buildContext := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	authConfig := types.AuthConfig{
		Username: authData.Username,
		Password: authData.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("error when encoding authConfig. err: %v", err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imageString := fmt.Sprintf("%s/%s/%s:%s", environment.DockerIOPath, authData.Username, authData.Repository, authData.Version)
	resp, err := cli.ImagePush(buildContext, imageString, types.ImagePushOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		log.Fatalf("unable to push the code runner image to hub: %v", err)
	}
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		log.Fatalf("unable to read image push response: %v", err)
	}
	defer func() {
		err = resp.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	return err

}