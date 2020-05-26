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

type AssignmentEnvConfig struct {
	BaseImage string       `yaml:"baseImage"`
	Deps      Dependencies `yaml:"dependencies"`
}

func (config *AssignmentEnvConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type tempAssignmentEnvConfig struct {
		BaseImage string       `yaml:"baseImage"`
		Deps      Dependencies `yaml:"dependencies"`
	}
	temp := &tempAssignmentEnvConfig{}

	if err := unmarshal(temp); err != nil {
		return errors.Wrap(err, "error in unmarshaling assignment environment configuration")
	}

	// Validate the configuration data.
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

func (config AssignmentEnvConfig) GetInstruction() string {
	buf := &bytes.Buffer{}
	buf.WriteString("FROM " + config.BaseImage)
	buf.WriteString("\n")
	buf.WriteString("COPY . /" + constants.CodeRunnerDir)
	buf.WriteString("\n")
	buf.WriteString(config.Deps.GetInstruction() + "\n")
	return buf.String()
}

type Dependencies struct {
	Language  LanguageInfo                  `yaml:",inline"`
	Libraries map[string]LibInstallationCmd `yaml:"lib"`
}

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

type LibInstallationCmd struct {
	Cmd string `yaml:"cmd"`
}

func (libCmd LibInstallationCmd) GetInstruction() string {
	return libCmd.Cmd
}

type LanguageInfo struct {
	Name    string `yaml:"lang"`
	Version string `yaml:"langVersion"`
}

func (langInfo LanguageInfo) GetInstruction() string {
	return fmt.Sprintf("RUN ./%s/%s_%s.sh", constants.InstallationScriptsDir, langInfo.Name, langInfo.Version)
}

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
