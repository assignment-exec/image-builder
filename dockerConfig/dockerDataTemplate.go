package dockerConfig

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type dockerfileData struct {
	ServerConfig struct {
		From       string            `yaml:"from"`
		Env        map[string]string `yaml:"env"`
		Copy       map[string]string `yaml:"copy"`
		WorkDir    string            `yaml:"workdir"`
		RunCommand string            `yaml:"run"`
		ServerPort string            `yaml:"serverPort"`
		FinalCmd   string            `yaml:"cmd"`
	} `yaml:"serverConfig"`
}

func UnmarshalYAML(configFilename string) (*dockerfileData, error) {
	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	config := &dockerfileData{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Printf("Error in unmarshalling yaml: %v", err)
		return nil, err
	}
	return config, nil
}
