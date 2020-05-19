package configurations

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/environment"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

func validateLang(langName string, langVersion string) error {
	scriptName := fmt.Sprintf("%s_%s.sh", langName, langVersion)

	// Check whether the given language and version are available in the installation scripts.
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "error in getting current directory")
	}
	scriptPath := filepath.Join(currentDir, constants.InstallationScriptsDir, scriptName)
	_, err = os.Stat(scriptPath)
	if os.IsNotExist(err) {
		return errors.New("installation scripts for given language and version doesn't exists")
	}
	return nil
}

func validateBaseImage(baseImage string) error {
	backgroundContext := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	username, hasFound := os.LookupEnv(environment.DockerAuthUsername)
	if !hasFound {
		return errors.New("environment variable for username not set")
	}

	response, err := dockerClient.ImageSearch(backgroundContext, username, types.ImageSearchOptions{
		Limit: 25})
	if err != nil {
		return err
	} else {
		hasFound := false
		for _, result := range response {
			if strings.Contains(baseImage, result.Name) {
				hasFound = true
			}
		}
		if !hasFound {
			return errors.New("code-runner base image not found on docker hub")
		}
	}

	return nil
}
