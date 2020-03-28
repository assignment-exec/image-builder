package dockerConfig

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type DockerConfig struct {
	From     string `yaml:"operatingSystem"`
	Compiler string `yaml:"compiler"`
	ServerPort string `yaml:"serverPort"`
	WorkDir string `yaml:"workdir"`
	Command string `yaml:"command"`
}

func GetConfig(configFilename string) (*DockerConfig, error) {
	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {

	}

	c := &DockerConfig{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Error in unmarshalling yaml: %v", err)
	}
	return c, nil
}

