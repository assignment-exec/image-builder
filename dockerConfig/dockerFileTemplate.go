package dockerConfig

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"text/template"
)

type dockerfileTemplate struct {
	Data *dockerfileData
}

func newDockerfileTemplate(data *dockerfileData) *dockerfileTemplate {
	return &dockerfileTemplate{Data: data}
}

func newDockerFileDataFromYamlFile(filename string) (*dockerfileData, error) {
	node := yaml.Node{}

	err := unmarshalYamlFile(filename, &node)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	stages, err := getStagesDataFromNode(node.Content[0])
	if err != nil {
		return nil, fmt.Errorf("Can't extract Stages from node: %v", err)
	}

	return &dockerfileData{Stages: stages}, nil
}

func getStagesDataFromNode(node *yaml.Node) ([]stage, error) {
	var data dockerfileDataYaml

	stagesInOrder, err := getStagesOrderFromYamlNode(node)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
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


func (d *dockerfileTemplate) generateTemplate(writer io.Writer) error {
	templateString := "{{- range .Stages -}}" +
		"{{- range $i, $instruction := . }}" +
		"{{- if gt $i 0 }}\n{{ end }}" +
		"{{ $instruction.Render }}\n" +
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