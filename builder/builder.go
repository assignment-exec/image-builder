package builder

import (
	"archive/tar"
	"assignment-exec/image-builder/constants"
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Builds Docker Image for Code Runner.
func BuildImage() error {

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

	dockerFilepath := constants.DockerFilepath
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
		log.Fatalf("unable to write dockerfile into tar body: %v", err)
	}

	// Use the tar of the Dockerfile while building image.
	dockerFileTarReader := bytes.NewReader(tarBuffer.Bytes())

	response, err := dockerClient.ImageBuild(
		buildContext,
		dockerFileTarReader,
		types.ImageBuildOptions{Dockerfile: dockerFilepath})

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
