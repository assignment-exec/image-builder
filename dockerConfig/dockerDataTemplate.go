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

func (f from) Render() string {
	result := fmt.Sprintf("FROM %s", f.Image)

	if f.As != "" {
		result = fmt.Sprintf("%s As %s", result, f.As)
	}

	return result
}

type env struct {
	GOMODULE string `yaml:"GOMODULE"`
	GOFLAGS    string `yaml:"GOFLAGS"`
}

func (f env) Render() string {
	result := fmt.Sprintf("ENV GOMODULE=%s", f.GOMODULE)
	result = fmt.Sprintf("%s\nENV GOFLAGS=%s", result, f.GOFLAGS)

	return result
}

type copyCommand struct {
	BaseDir string `yaml:"basedir"`
	DestDir string `yaml:"destdir"`

}

func (f copyCommand) Render() string {
	result := fmt.Sprintf("COPY %s %s", f.BaseDir,f.DestDir)
	return result
}

type workDir struct {
	BaseDir string `yaml:"dir"`

}

func (f workDir) Render() string {
	result := fmt.Sprintf("WORKDIR %s", f.BaseDir)
	return result
}

type runCommand struct {
	Param string `yaml:"param"`

}

func (f runCommand) Render() string {
	result := fmt.Sprintf("RUN %s", f.Param)
	return result
}

type serverPort struct {
	Number string `yaml:"number"`

}

func (f serverPort) Render() string {
	result := fmt.Sprintf("EXPOSE %s", f.Number)
	return result
}

type cmd struct {
	Params []string `yaml:"params"`

}

func (f cmd) Render() string {
	var params []string
	for _,p := range f.Params {
		params = append(params, fmt.Sprintf("\"%s\"", p))
	}

	paramsString := strings.Join(params, ", ")
	execFormString := fmt.Sprintf("[%s]", paramsString)
	result := fmt.Sprintf("CMD %s", execFormString)
	return result
}