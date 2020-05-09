package configurations

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/utilities/validation"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type DockerConfig struct {
	BaseImage    string         `yaml:"baseImage"`
	Dependencies LangDependency `yaml:"dependencies"`
}

func (config *DockerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type tempDockerConfig struct {
		BaseImage    string         `yaml:"baseImage"`
		Dependencies LangDependency `yaml:"dependencies"`
	}
	temp := &tempDockerConfig{}

	if err := unmarshal(temp); err != nil {
		return errors.Wrap(err, "error in unmarshaling assignment environment configuration")
	}

	// Validate the configuration data.
	err := validation.Validate("error in configuration",
		ValidatorForConfig(DockerConfig(*temp),
			withBaseImageValidator(),
			withLanguageValidator(),
			withLibsValidator()))

	if err != nil {
		return err
	}

	config.BaseImage = temp.BaseImage
	config.Dependencies = temp.Dependencies
	return nil
}

func (config DockerConfig) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("FROM " + config.BaseImage)
	buf.WriteString("\n")
	buf.WriteString(config.Dependencies.String() + "\n")
	return buf.String()
}

type LangDependency struct {
	Language  LanguageInfo                  `yaml:",inline"`
	Libraries map[string]LibInstallationCmd `yaml:"lib"`
}

func (langDep LangDependency) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(langDep.Language.String())
	buf.WriteString("\n")
	for lib, installCmd := range langDep.Libraries {
		buf.WriteString("RUN " + installCmd.String() + " " + lib)
		buf.WriteString("\n")
	}
	return buf.String()
}

type LibInstallationCmd struct {
	Cmd string `yaml:"cmd"`
}

func (libCmd LibInstallationCmd) String() string {
	return libCmd.Cmd
}

type LanguageInfo struct {
	Name    string `yaml:"lang"`
	Version string `yaml:"langVersion"`
}

func (langInfo LanguageInfo) String() string {
	return fmt.Sprintf("RUN ./%s/%s_%s.sh", constants.InstallationScriptsDir, langInfo.Name, langInfo.Version)
}

func GetConfig(configFilename string) (*DockerConfig, error) {

	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log config file")
	}

	c := &DockerConfig{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, errors.Wrap(err, "error in unmarshaling yaml: %v")
	}

	return c, nil
}
