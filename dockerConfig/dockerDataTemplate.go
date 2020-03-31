package dockerConfig

import (
	"fmt"
	"strings"
)

type instruction interface {
	Render() string
}

type dockerfileData struct {
	Stages []stage `yaml:"stages"`
}

type stage []instruction

func (s *stage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data []interface{}
	var result []instruction
	err := unmarshal(&data)
	if err != nil {
		return err
	}
	*s = append(result, getInstructions(data)...)
	return nil
}

type from struct {
	Image string `yaml:"image"`
	As    string `yaml:"as"`
}

func (fromObj from) Render() string {
	result := fmt.Sprintf("FROM %s", fromObj.Image)

	if fromObj.As != "" {
		result = fmt.Sprintf("%s As %s", result, fromObj.As)
	}

	return result
}

type env struct {
	EnvParams map[string]string
}

func (envObj env) Render() string {

	var result string
	for key, value := range envObj.EnvParams {
		result = fmt.Sprintf("%s\nENV %s=%s", result, key, value)
	}

	return result
}

type copyCommand struct {
	BaseDir string `yaml:"basedir"`
	DestDir string `yaml:"destdir"`
}

func (cpyObj copyCommand) Render() string {
	result := fmt.Sprintf("COPY %s %s", cpyObj.BaseDir, cpyObj.DestDir)
	return result
}

type workDir struct {
	BaseDir string `yaml:"dir"`
}

func (wrkObj workDir) Render() string {
	result := fmt.Sprintf("WORKDIR %s", wrkObj.BaseDir)
	return result
}

type runCommand struct {
	Param string `yaml:"param"`
}

func (runObj runCommand) Render() string {
	result := fmt.Sprintf("RUN %s", runObj.Param)
	return result
}

type serverPort struct {
	Number string `yaml:"number"`
}

func (portObj serverPort) Render() string {
	result := fmt.Sprintf("EXPOSE %s", portObj.Number)
	return result
}

type cmd struct {
	Params []string `yaml:"params"`
}

func (cmdObj cmd) Render() string {
	var params []string
	for _, p := range cmdObj.Params {
		params = append(params, fmt.Sprintf("\"%s\"", p))
	}

	paramsString := strings.Join(params, ", ")
	execFormString := fmt.Sprintf("[%s]", paramsString)
	result := fmt.Sprintf("CMD %s", execFormString)
	return result
}
