package configurations

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"text/template"
)

type dockerfileTemplate struct {
	Data *dockerfileData
}

// Gives a new instance of dockerfile template.
func newDockerfileTemplate(data *dockerfileData) *dockerfileTemplate {
	return &dockerfileTemplate{Data: data}
}

// Unmarshal yaml file and gives an instance of dockerfileData.
func newDockerFileDataFromYamlFile(filename string) (*dockerfileData, error) {
	node := yaml.Node{}

	err := unmarshalYamlFile(filename, &node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	stages, err := getStagesDataFromNode(node.Content[0])
	if err != nil {
		return nil, fmt.Errorf("can't extract Stages from node: %v", err)
	}

	return &dockerfileData{Stages: stages}, nil
}

// Gives all the stages and their corresponding data provided in yaml.
func getStagesDataFromNode(node *yaml.Node) ([]stage, error) {
	var data dockerfileDataYaml

	stagesInOrder, err := getStagesOrderFromYamlNode(node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	if err := node.Decode(&data); err != nil {
		return nil, err
	}
	var stages []stage
	for _, stageName := range stagesInOrder {
		stages = append(stages, data.Stages[stageName])
	}
	return stages, nil
}

// Generates a new template and writes it to the Dockerfile.
func (d *dockerfileTemplate) generateDockerfileFromTemplate(writer io.Writer) error {
	templateString := "{{- range .Stages -}}" +
		"{{- range $i, $instruction := . }}" +
		"{{- if gt $i 0 }}\n{{ end }}" +
		"{{ $instruction.WriteInstruction }}\n" +
		"{{- end }}\n\n" +
		"{{ end }}"

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
