// Package configurations implements routines to read and store the
// assignment environment configuration yaml file, get the docker instructions
// in the specific format for every configuration.
package configurations

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/environment"
	"assignment-exec/image-builder/utilities/validation"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// AssignmentEnvConfig struct type holds the base image and
// dependencies level of the configuration yaml.
type AssignmentEnvConfig struct {
	BaseImage string       `yaml:"baseImage"`
	Deps      Dependencies `yaml:"dependencies"`
}

// UnmarshalYAML unmarshals the config yaml, validates the data
// and stores the configurations to `AssignmentEnvConfig`.
//It returns any error encountered while unmarshaling the file.
func (config *AssignmentEnvConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type tempAssignmentEnvConfig struct {
		BaseImage string       `yaml:"baseImage"`
		Deps      Dependencies `yaml:"dependencies"`
	}
	temp := &tempAssignmentEnvConfig{}

	if err := unmarshal(temp); err != nil {
		return errors.Wrap(err, "error in unmarshaling assignment environment configuration")
	}

	// Validates base image, language and the library dependencies.
	err := validation.Validate("error in configuration",
		ValidatorForConfig(AssignmentEnvConfig(*temp),
			withBaseImageValidator(),
			withLanguageValidator(),
			withLibsValidator()))

	if err != nil {
		return err
	}

	config.BaseImage = temp.BaseImage
	config.Deps = temp.Deps
	return nil
}

// GetInstruction returns the docker instructions for the full configuration
// as a single string.
func (config AssignmentEnvConfig) GetInstruction() string {
	buf := &bytes.Buffer{}
	buf.WriteString("FROM " + config.BaseImage)
	buf.WriteString("\n")
	buf.WriteString("COPY . /" + constants.CodeRunnerDir)
	buf.WriteString("\n")
	buf.WriteString(config.Deps.GetInstruction() + "\n")
	return buf.String()
}

// Dependencies struct type holds the language information
// and library names and their installation command level of the configuration yaml.
type Dependencies struct {
	Language  LanguageInfo                  `yaml:",inline"`
	Libraries map[string]LibInstallationCmd `yaml:"lib"`
}

// GetInstruction returns the docker instructions for the dependencies
// as a single string.
func (langDep Dependencies) GetInstruction() string {
	buf := &bytes.Buffer{}
	buf.WriteString(langDep.Language.GetInstruction())
	buf.WriteString("\n")
	buf.WriteString("ENV " + environment.LanguageEnvKey + " " + langDep.Language.Name)
	buf.WriteString("\n")
	for _, installCmd := range langDep.Libraries {
		buf.WriteString("RUN " + installCmd.GetInstruction())
		buf.WriteString("\n")
	}
	return buf.String()
}

// LibInstallationCmd struct type holds the installation command
// for the respective library name.
type LibInstallationCmd struct {
	Cmd string `yaml:"cmd"`
}

// GetInstruction returns the library installation command.
func (libCmd LibInstallationCmd) GetInstruction() string {
	return libCmd.Cmd
}

// LanguageInfo struct type holds name and version of language.
type LanguageInfo struct {
	Name    string `yaml:"lang"`
	Version string `yaml:"langVersion"`
}

// GetInstruction returns the docker instruction for the language
// as a single string. The instruction includes running the installation
// script for the given language.
func (langInfo LanguageInfo) GetInstruction() string {
	return fmt.Sprintf("RUN ./%s/%s_%s.sh", constants.InstallationScriptsDir, langInfo.Name, langInfo.Version)
}

// GetAssignmentEnvConfig reads the yaml config file and unmarshals it into
// `AssignmentEnvConfig` struct. It returns the `AssignmentEnvConfig` instance
// and any error encountered.
func GetAssignmentEnvConfig(configFilepath string) (*AssignmentEnvConfig, error) {

	yamlFile, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log config file")
	}

	c := &AssignmentEnvConfig{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, errors.Wrap(err, "error in unmarshaling yaml: %v")
	}

	return c, nil
}
