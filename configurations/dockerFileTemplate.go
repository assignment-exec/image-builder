package configurations

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"text/template"
)

type dockerfileTemplate struct {
	Data serverConfig
}

// Gives a new instance of dockerfile template.
func newDockerfileTemplate(data serverConfig) *dockerfileTemplate {
	return &dockerfileTemplate{Data: data}
}

// Unmarshal yaml file and gives an instance of serverConfig.
func newDockerFileDataFromYamlFile(filename string) (serverConfig, error) {
	node := yaml.Node{}

	err := unmarshalYamlFile(filename, &node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	stageData, err := getStagesDataFromNode(node.Content[0])
	if err != nil {
		return nil, fmt.Errorf("can't extract Stages from node: %v", err)
	}
	return stageData, nil
}

// Returns the serverConfig instructions provided in yaml.
func getStagesDataFromNode(node *yaml.Node) (serverConfig, error) {
	var data dockerfileDataYaml

	err := verifyStageYamlNode(node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	if err := node.Decode(&data); err != nil {
		return nil, err
	}
	var stageData serverConfig

	stageData = data.ServerConfig
	return stageData, nil
}

// Generates a new template and writes it to the Dockerfile.
func (d *dockerfileTemplate) generateDockerfileFromTemplate(writer io.Writer) error {
	templateString :=
		"{{- range $i, $instruction := . }}" +
			"{{- if gt $i 0 }}\n{{ end }}" +
			"{{ $instruction.WriteInstruction }}\n" +
			"{{- end }}\n\n"

	tmpl, err := template.New("dockerfile.template").Parse(templateString)
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, d.Data)
	if err != nil {
		return err
	}

	return nil
}
