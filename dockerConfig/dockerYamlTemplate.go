package dockerConfig

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"strings"

	//"log"
)

type dockerfileDataYaml struct {
	Stages map[string]stage `yaml:"stages"`
}

func getInstructions(in []interface{}) []instruction {

	result := make([]instruction, len(in))
	for i, v := range in {
		result[i] = parseAllInstructions(v)
	}
	return result
}

func parseAllInstructions(v interface{}) instruction {
	switch v := v.(type) {
	case map[string]interface{}:
		return parseInnerInstructions(v)
	}

	log.Fatal("unknown instruction in yaml")
	return nil
}

func parseSpecificInstruction(instructionName string, value interface{}) instruction {
	v, ok := value.(map[string]interface{})
	if !ok {
		panic("Error")
	}
	switch strings.ToLower(instructionName) {
	case "from":
		return parseFrom(v)
	case "env":
		return parseEnv(v)
	case "workdir":
		return parseWorkDir(v)
	case "runcommand":
		return parseRun(v)
	case "cmd":
		return parseCmd(v)
	case "serverport":
		return parseServerPort(v)
	case "copy":
		return parseCopy(v)
	}
	log.Fatal("unknown instruction in yaml")
	return nil
}

func convertMapToStringArray(mapInterface map[string]interface{}) []string {
	switch v := mapInterface["params"].(type) {
	case []interface{}:
		mapString := make([]string, len(v))
		for index, vv := range v {
			strValue := fmt.Sprintf("%v", vv)
			mapString[index] = strValue
		}
		return mapString
	}

	log.Fatal("no valid command parameters found")
	return nil

}
func convertMapToMap(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapString[strKey] = strValue
	}

	return mapString
}

func parseEnv(data map[string]interface{}) instruction {
	convertedData := convertMapToMap(data)
	var envObj env
	envObj.EnvParams = make(map[string]string)
	for key, value := range convertedData {
		envObj.EnvParams[key] = value
	}

	return envObj
}

func parseWorkDir(value map[string]interface{}) instruction {
	v := convertMapToMap(value)
	var workDir workDir
	if v["dir"] != "" {
		workDir.BaseDir = v["dir"]
	}
	return workDir
}

func parseRun(value map[string]interface{}) instruction {
	v := convertMapToMap(value)
	var runCmd runCommand
	if v["param"] != "" {
		runCmd.Param = v["param"]
	}
	return runCmd
}

func parseCmd(value map[string]interface{}) instruction {
	var command cmd
	v := convertMapToStringArray(value)
	if v != nil {
		command.Params = v
	}
	return command
}

func parseServerPort(value map[string]interface{}) instruction {
	v := convertMapToMap(value)
	var serverPort serverPort
	if v["number"] != "" {
		serverPort.Number = v["number"]
	}
	return serverPort
}

func parseCopy(value map[string]interface{}) instruction {
	v := convertMapToMap(value)
	var cpy copyCommand
	if v["basedir"] != "" {
		cpy.BaseDir = v["basedir"]
	}

	if v["destdir"] != "" {
		cpy.DestDir = v["destdir"]
	}

	return cpy
}

func parseFrom(value map[string]interface{}) from {
	v := convertMapToMap(value)
	var from from
	if v["image"] != "" {
		from.Image = v["image"]
	}

	if v["as"] != "" {
		from.As = v["as"]
	}

	return from
}

func parseInnerInstructions(in map[string]interface{}) instruction {
	for key, value := range in {

		switch value.(type) {
		case map[string]interface{}:
			return parseSpecificInstruction(key, value)
		}
	}
	log.Fatal("unknown instruction in yaml")
	return  nil
}

func unmarshalYamlFile(filename string, node *yaml.Node) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("yamlFile.Get err #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, node)
	if err != nil {
		return fmt.Errorf("Unmarshal: %v", err)
	}
	return nil
}
func getStagesOrderFromYamlNode(node *yaml.Node) ([]string, error) {
	var stages []string

	if node.Kind != yaml.MappingNode {
		return nil, errors.New("Yaml should contain a map that contains 'Stages' key!")
	}

	stagesKeyNode := node.Content[0]
	if stagesKeyNode.Kind != yaml.ScalarNode {
		return nil, errors.New("Yaml should contain a 'Stages' key!")
	}

	stagesMapNode := node.Content[1]
	if stagesMapNode.Kind != yaml.MappingNode {
		return nil, errors.New("yaml should contain a Stages map that has stage names As keys")
	}

	for i, stage := range stagesMapNode.Content {
		if i%2 == 0 {
			if stage.Kind != yaml.ScalarNode {
				return nil, errors.New("Yaml should contain stage keys in 'staging' map")
			}
			stages = append(stages, stage.Value)
		} else {
			if stage.Kind != yaml.SequenceNode {
				return nil, errors.New("Yaml should contain stage sequences in 'staging' map")
			}
		}
	}

	return stages, nil
}