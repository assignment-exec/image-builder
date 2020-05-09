package configurations

import (
	"assignment-exec/image-builder/constants"
	"assignment-exec/image-builder/utilities/validation"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	BaseImage string     `yaml:"baseImage"`
	Deps      Dependency `yaml:"dependencies"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tempConfig struct {
		BaseImage string     `yaml:"baseImage"`
		Deps      Dependency `yaml:"dependencies"`
	}
	temp := &tempConfig{}

	if err := unmarshal(temp); err != nil {
		return errors.Wrap(err, "failed to unmarshal assignment environment configuration")
	}

	// Validate the configuration data.
	err := validation.Validate("invalid configuration",
		ValidatorForConfig(Config(*temp),
			withBaseImageValidator(),
			withLanguageValidator(),
			withLibsValidator()))

	if err != nil {
		return err
	}

	// All validations passed.
	c.BaseImage = temp.BaseImage
	c.Deps = temp.Deps
	return nil
}

func (c Config) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("FROM " + c.BaseImage)
	buf.WriteString("\n")
	buf.WriteString(c.Deps.String() + "\n")
	return buf.String()
}

type Dependency struct {
	Language LanguageReq              `yaml:",inline"`
	Libs     map[string]LibInstallCmd `yaml:"lib"`
}

func (d Dependency) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(d.Language.String())
	buf.WriteString("\n")
	for lib, installCmd := range d.Libs {
		buf.WriteString("RUN " + installCmd.String() + " " + lib)
		buf.WriteString("\n")
	}
	return buf.String()
}

type LibInstallCmd struct {
	Cmd string `yaml:"cmd"`
}

func (lic LibInstallCmd) String() string {
	return lic.Cmd
}

type LanguageReq struct {
	Name    string `yaml:"lang"`
	Version string `yaml:"langVersion"`
}

func (lr LanguageReq) String() string {
	return fmt.Sprintf("RUN ./%s/%s_%s.sh", constants.InstallationScriptsDir, lr.Name, lr.Version)
}

func GetConfig(configFilename string) (*Config, error) {

	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log config file")
	}

	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Error in unmarshalling yaml: %v", err)
	}

	return c, nil
}
