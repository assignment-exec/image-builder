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

	serverConfigData, err := getConfigsDataFromNode(node.Content[0])
	if err != nil {
		return nil, fmt.Errorf("can't extract server config from node: %v", err)
	}
	return serverConfigData, nil
}

// Returns the serverConfig instructions provided in yaml.
func getConfigsDataFromNode(node *yaml.Node) (serverConfig, error) {
	var data dockerfileDataYaml

	err := verifyConfigYamlNode(node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	if err := node.Decode(&data); err != nil {
		return nil, err
	}
	var serverConfigData serverConfig

	serverConfigData = data.ServerConfig
	return serverConfigData, nil
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
