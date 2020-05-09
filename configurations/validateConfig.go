package configurations

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/environment"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
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

	username, found := os.LookupEnv(environment.DockerAuthUsername)
	if !found {
		return errors.New("environment variable for username not set")
	}
	password, found := os.LookupEnv(environment.DockerAuthPassword)
	if !found {
		return errors.New("environment variable for password not set")
	}
	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return errors.Wrap(err, "error in encoding authConfig")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	_, err = dockerClient.ImageSearch(backgroundContext, baseImage, types.ImageSearchOptions{
		RegistryAuth: authStr,
		Limit:        10})

	if err != nil {
		return err
	}

	return nil
}
