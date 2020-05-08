package configurations

import (
	"assignment-exec/image-builder/constants"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type instruction interface {
	WriteInstruction() string
}

type config []instruction

// UnmarshalYAML decodes the yaml configuration to obtain dockerfile instructions.
func (s *config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data []interface{}
	err := unmarshal(&data)
	if err != nil {
		return err
	}
	*s = getInstructions(data)
	return nil
}

// FROM instructions.
type from struct {
	Image string
	As    string
}

// Gives FROM instruction for the values read from yaml.
func (fromObj from) WriteInstruction() string {
	result := fmt.Sprintf("FROM %s", fromObj.Image)

	if fromObj.As != "" {
		result = fmt.Sprintf("%s As %s", result, fromObj.As)
	}

	return result
}

// ENV instruction.
type env struct {
	EnvParams map[string]string
}

// Gives ENV instruction for the values read from yaml.
func (envObj env) WriteInstruction() string {

	var result string
	for key, value := range envObj.EnvParams {
		result = fmt.Sprintf("%s\nENV %s=%s", result, key, value)
	}

	return result
}

// COPY instruction.
type copyCommand struct {
	BaseDir string
	DestDir string
}

// Gives COPY instruction for the values read from yaml.
func (cpyObj copyCommand) WriteInstruction() string {
	result := fmt.Sprintf("COPY %s %s", cpyObj.BaseDir, cpyObj.DestDir)
	return result
}

// WORKDIR instruction.
type workDir struct {
	BaseDir string
}

// Gives WORKDIR instruction for the values read from yaml.
func (wrkObj workDir) WriteInstruction() string {
	result := fmt.Sprintf("WORKDIR %s", wrkObj.BaseDir)
	return result
}

// RUN instruction.
type runCommand struct {
	Param string
}

// Gives RUN instruction for the values read from yaml.
func (runObj runCommand) WriteInstruction() string {
	result := fmt.Sprintf("RUN %s", runObj.Param)
	return result
}

// EXPOSE instruction.
type port struct {
	Number string
}

// Gives EXPOSE instruction for the values read from yaml.
func (portObj port) WriteInstruction() string {
	result := fmt.Sprintf("EXPOSE %s", portObj.Number)
	return result
}

// CMD instruction.
type cmd struct {
	Params []string
}

// Gives CMD instruction for the values read from yaml.
func (cmdObj cmd) WriteInstruction() string {
	var params []string
	for _, p := range cmdObj.Params {
		params = append(params, fmt.Sprintf("\"%s\"", p))
	}

	paramsString := strings.Join(params, ", ")
	execFormString := fmt.Sprintf("[%s]", paramsString)
	result := fmt.Sprintf("CMD %s", execFormString)
	return result
}

// Programming Language installation instruction.
type programmingLanguage struct {
	Name    string
	Version string
}

// Gives RUN instruction for installing the specified programmingLanguage.
func (lang programmingLanguage) WriteInstruction() string {

	scriptName := fmt.Sprintf("%s_%s.sh", lang.Name, lang.Version)

	// Check whether the given language and version are available in the installation scripts.
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("error in getting current directory")
	}
	scriptPath := filepath.Join(currentDir, constants.InstallationScriptsDir, scriptName)
	_, err = os.Stat(scriptPath)
	if err == nil {
		return fmt.Sprintf("RUN ./%s/%s ", constants.InstallationScriptsDir, scriptName)
	} else if os.IsNotExist(err) {
		log.Println("installation scripts for given language and version doesn't exists")
	}
	return ""
}
