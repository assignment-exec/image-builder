package configurations

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"strings"
)

type dockerfileDataYaml struct {
	Config config `yaml:"config"`
}

// Gives all the instructions after parsing them from yaml.
func getInstructions(in []interface{}) []instruction {

	result := make([]instruction, len(in))
	for i, v := range in {
		result[i] = parseAllInstructions(v)
	}
	return result
}

// Parses instructions based on their type.
func parseAllInstructions(v interface{}) instruction {
	switch v := v.(type) {
	case map[string]interface{}:
		return parseInnerInstructions(v)
	}

	log.Fatal("unknown instruction in yaml")
	return nil
}

// Invokes respective functions based on the instruction node being parsed.
func parseSpecificInstruction(instructionName string, value interface{}) instruction {
	v, ok := value.(map[string]interface{})
	if !ok {
		log.Fatal("error in parsing specific instruction node")
		return nil
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
	case "port":
		return parsePort(v)
	case "copy":
		return parseCopy(v)
	case "programminglanguage":
		return parseProgrammingLanguage(v)
	}
	log.Fatal("unknown instruction in yaml")
	return nil
}

// Converts a map of string-interface to a string array.
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

	log.Fatal("convert map tp string: no valid command parameters found")
	return nil

}

// Converts a map of string-interface to a map of string-string.
func convertMapToMap(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)

	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapString[strKey] = strValue
	}

	return mapString
}

// Parses the programmingLanguage instruction from yaml and returns an instance of `programmingLanguage`.
func parseProgrammingLanguage(data map[string]interface{}) instruction {
	convertedData := convertMapToMap(data)
	var langObj programmingLanguage
	if convertedData["name"] != "" {
		langObj.Name = convertedData["name"]
		if strings.Contains(langObj.Name, "+") {
			langObj.Name = strings.ReplaceAll(langObj.Name, "+", "p")
		}
	}

	if convertedData["version"] != "" {
		langObj.Version = convertedData["version"]
	}
	return langObj
}

// Parses the env instruction node from yaml and returns an instance of `env`.
func parseEnv(data map[string]interface{}) instruction {
	convertedData := convertMapToMap(data)
	var envObj env
	envObj.EnvParams = make(map[string]string)
	for key, value := range convertedData {
		envObj.EnvParams[key] = value
	}

	return envObj
}

// Parses the workdir instruction node and returns an instance of `workDir`.
func parseWorkDir(value map[string]interface{}) instruction {
	v := convertMapToMap(value)
	var workDir workDir
	if v["dir"] != "" {
		workDir.BaseDir = v["dir"]
	}
	return workDir
}

// Parses the runCommand instruction node and returns an instance of `runCommand`.
func parseRun(value map[string]interface{}) instruction {
	v := convertMapToMap(value)
	var runCmd runCommand
	if v["param"] != "" {
		runCmd.Param = v["param"]
	}
	return runCmd
}

// Parses the cmd instruction node and returns an instance of `cmd`.
func parseCmd(value map[string]interface{}) instruction {
	var command cmd
	v := convertMapToStringArray(value)
	if v != nil {
		command.Params = v
	}
	return command
}

// Parses the port instruction node and returns an instance of `port`.
func parsePort(value map[string]interface{}) instruction {
	v := convertMapToMap(value)
	var serverPort port
	if v["number"] != "" {
		serverPort.Number = v["number"]
	}
	return serverPort
}

// Parses the copy instruction node and returns an instance of `copyCommand`.
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

// Parses the from instruction node and returns an instance of `from`.
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

// Parses the instruction nodes within every config node.
func parseInnerInstructions(in map[string]interface{}) instruction {
	for key, value := range in {

		switch value.(type) {
		case map[string]interface{}:
			return parseSpecificInstruction(key, value)
		}
	}
	log.Fatal("unknown instruction in yaml")
	return nil
}

// Unmarshal yaml file into a yaml node.
func unmarshalYamlFile(filename string, node *yaml.Node) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "error in reading yaml")
	}
	err = yaml.Unmarshal(yamlFile, node)
	if err != nil {
		return errors.Wrap(err, "error in unmarshaling yaml")
	}
	return nil
}

// Verifies the kind of the yaml node.
func verifyConfigYamlNode(node *yaml.Node) error {

	if node.Kind != yaml.MappingNode {
		return errors.New("yaml should contain a map that contains a valid config name key")
	}

	configKeyNode := node.Content[0]
	if configKeyNode.Kind != yaml.ScalarNode {
		return errors.New("yaml should contain a valid config name key")
	}

	return nil
}
