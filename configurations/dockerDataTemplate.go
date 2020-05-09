package configurations

import (
	"assignment-exec/image-builder/constants"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

var instructions []string

func ParseInstructions(configFilename string) (*dockerConfig, []string, error) {
	instructions = instructions[:0]
	config, err := getConfig(configFilename)
	if err != nil {
		return nil, nil, err
	}

	err = config.parseFrom()
	if err != nil {
		return nil, nil, err
	}
	config.parseEnv()
	config.parseCopy()
	config.parseWorkDir()
	config.parseRun()
	config.parsePort()
	config.parseCmd()
	err = config.parseLang()
	if err != nil {
		return nil, nil, err
	}

	return config, instructions, err
}

// Parses From instruction.
func (d dockerConfig) parseFrom() error {
	if len(d.From.Image) <= 0 {
		return errors.New("from instruction not provided")
	}
	result := fmt.Sprintf("FROM %s", d.From.Image)
	instructions = append(instructions, result)

	if d.From.As != "" {
		result = fmt.Sprintf("%s As %s", result, d.From.As)
	}

	return nil
}

// Parses Env instruction.
func (d dockerConfig) parseEnv() {
	var result string
	for key, value := range d.EnvParams {
		result = fmt.Sprintf("%s\nENV %s=%s", result, key, value)
	}

	if len(result) > 0 {
		instructions = append(instructions, result)
	}
}

// Parses Copy instruction.
func (d dockerConfig) parseCopy() {
	if len(d.CopyCommand.BaseDir) > 0 && len(d.CopyCommand.DestDir) > 0 {
		result := fmt.Sprintf("COPY %s %s", d.CopyCommand.BaseDir, d.CopyCommand.DestDir)
		instructions = append(instructions, result)
	}
}

// Parses Workdir instruction.
func (d dockerConfig) parseWorkDir() {
	if len(d.CopyCommand.BaseDir) > 0 && len(d.CopyCommand.DestDir) > 0 {
		result := fmt.Sprintf("WORKDIR %s", d.WorkDir.BaseDir)
		instructions = append(instructions, result)
	}
}

// Parses Run instruction.
func (d dockerConfig) parseRun() {
	if len(d.RunCommand.Param) > 0 {
		result := fmt.Sprintf("RUN %s", d.RunCommand.Param)
		instructions = append(instructions, result)
	}
}

// Parses Port instruction.
func (d dockerConfig) parsePort() {
	if len(d.Port.Number) > 0 {
		result := fmt.Sprintf("EXPOSE %s", d.Port.Number)
		instructions = append(instructions, result)
	}
}

// Parses Cmd instruction.
func (d dockerConfig) parseCmd() {
	if len(d.Cmd.Params) > 0 {
		var params []string
		for _, p := range d.Cmd.Params {
			params = append(params, fmt.Sprintf("\"%s\"", p))
		}

		paramsString := strings.Join(params, ", ")
		execFormString := fmt.Sprintf("[%s]", paramsString)
		result := fmt.Sprintf("CMD %s", execFormString)
		instructions = append(instructions, result)
	}
}

// Parses Programming Language instruction.
func (d dockerConfig) parseLang() error {

	if len(d.ProgrammingLanguage.Name) > 0 && len(d.ProgrammingLanguage.Version) > 0 {
		scriptName := fmt.Sprintf("%s_%s.sh", d.ProgrammingLanguage.Name, d.ProgrammingLanguage.Version)

		// Check whether the given language and version are available in the installation scripts.
		currentDir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "error in getting current directory")
		}
		scriptPath := filepath.Join(currentDir, constants.InstallationScriptsDir, scriptName)
		_, err = os.Stat(scriptPath)
		if err == nil {
			result := fmt.Sprintf("RUN ./%s/%s ", constants.InstallationScriptsDir, scriptName)
			instructions = append(instructions, result)
		} else if os.IsNotExist(err) {
			return errors.New("installation scripts for given language and version doesn't exists")
		}
	}
	return nil
}
