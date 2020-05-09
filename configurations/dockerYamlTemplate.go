package configurations

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type dockerConfig struct {
	From struct {
		Image string `yaml:"image"`
		As    string `yaml:"as"`
	} `yaml:"from"`

	EnvParams map[string]string `yaml:"env"`

	CopyCommand struct {
		BaseDir string `yaml:"basedir"`
		DestDir string `yaml:"destdir"`
	} `yaml:"copy"`

	WorkDir struct {
		BaseDir string `yaml:"dir"`
	} `yaml:"workdir"`

	RunCommand struct {
		Param string `yaml:"param"`
	} `yaml:"runCommand"`

	Port struct {
		Number string `yaml:"number"`
	} `yaml:"port"`

	Cmd struct {
		Params []string `yaml:"params"`
	} `yaml:"cmd"`

	ProgrammingLanguage struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"programmingLanguage"`
}

func getConfig(logConfigFilename string) (*dockerConfig, error) {

	yamlFile, err := ioutil.ReadFile(logConfigFilename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read log dockerConfig file")
	}

	c := &dockerConfig{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, errors.Wrap(err, "error in unmarshaling yaml")
	}

	return c, nil
}
