package builder

import (
	"assignment-exec/image-builder/constants"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
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

type ImageBuilder struct {
	authData            *dockerAuthData
	dockerFilename      string
	LanguageImageFormat string
}

// New Image Builder instance.
func NewImageBuilder(dockerAuthConfig string, dockerFilename string) *ImageBuilder {
	dockerAuth, err := getAuthData(dockerAuthConfig)
	imgBuilder := &ImageBuilder{authData: dockerAuth, dockerFilename: dockerFilename}
	if err != nil {
		log.Fatalf("error in reading docker authentication data: %v", err)
	}
	return imgBuilder
}

// Get the docker authentication details.
func getAuthData(filename string) (*dockerAuthData, error) {

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

// Build Docker Image.
func (imgBuilder ImageBuilder) BuildImage(assignmentEnv bool) error {

	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("error in creating a docker client: %v", err)
		return err
	}

	// Create a build context tar for the image.
	dockerFilepath := imgBuilder.dockerFilename

	dockerBuildContext, err := getDockerBuildContextTar(dockerFilepath)
	if err != nil {
		log.Fatalf("error in creating a docker client: %v", err)
		return err
	}

	tagName := getTagName(imgBuilder, assignmentEnv)

	response, err := dockerClient.ImageBuild(
		backgroundContext,
		dockerBuildContext,
		types.ImageBuildOptions{
			Dockerfile: dockerFilepath,
			Tags:       []string{tagName}})
	if err != nil {
		log.Fatalf("unable to build docker image: %v", err)
		return err
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
		return err
	}

	return nil
}

// Get Tag name for building image.
func getTagName(imgBuilder ImageBuilder, assignmentEnv bool) string {
	var tagName string
	if assignmentEnv {
		tagName = fmt.Sprintf("%s/%s", imgBuilder.authData.Username, imgBuilder.LanguageImageFormat)
	} else {
		tagName = fmt.Sprintf("%s/%s:%s", imgBuilder.authData.Username,
			imgBuilder.authData.Repository, imgBuilder.authData.Version)
	}

	return tagName
}

// Get Docker Build Context Tar Reader for building image.
func getDockerBuildContextTar(dockerFilepath string) (*os.File, error) {
	dockerFileReader, err := os.Open(dockerFilepath)
	if err != nil {
		log.Fatalf(" unable to open dockerfile: %v", err)
	}
	fileInfo, err := os.Stat(dockerFilepath)
	if err != nil {
		return nil, err
	}

	buildContextTar := new(archivex.TarFile)
	err = buildContextTar.Create(constants.BuildContextTar)
	if err != nil {
		return nil, err
	}
	err = buildContextTar.AddAll(constants.InstallationScriptsDir, true)
	if err != nil {
		return nil, err
	}
	err = buildContextTar.Add(dockerFilepath, dockerFileReader, fileInfo)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = buildContextTar.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	dockerBuildContext, err := os.Open(constants.BuildContextTar)
	if err != nil {
		return nil, err
	}

	err = os.Remove(constants.BuildContextTar)
	if err != nil {
		return nil, err
	}

	return dockerBuildContext, nil
}

func (imgBuilder ImageBuilder) PublishImage() error {
	// TODO: setup ssh keys for logging into docker hub

	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	authConfig := types.AuthConfig{
		Username: imgBuilder.authData.Username,
		Password: imgBuilder.authData.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return fmt.Errorf("error when encoding authConfig. err: %v", err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	imageString := fmt.Sprintf("%s/%s/%s:%s", constants.DockerIO, imgBuilder.authData.Username,
		imgBuilder.authData.Repository, imgBuilder.authData.Version)

	resp, err := dockerClient.ImagePush(backgroundContext, imageString, types.ImagePushOptions{
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
