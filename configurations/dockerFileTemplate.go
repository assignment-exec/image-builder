package configurations

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"text/template"
)

type dockerfileTemplate struct {
	Data config
}

// Gives a new instance of dockerfile template.
func NewDockerfileTemplate(data config) *dockerfileTemplate {
	return &dockerfileTemplate{Data: data}
}

// Unmarshal yaml file and gives an instance of config.
func NewDockerFileDataFromYamlFile(filename string) (config, error) {
	node := yaml.Node{}

	err := unmarshalYamlFile(filename, &node)
	if err != nil {
		return nil, err
	}

	configData, err := getConfigsDataFromNode(node.Content[0])
	if err != nil {
		return nil, err
	}
	return configData, nil
}

// Returns the config instructions provided in yaml.
func getConfigsDataFromNode(node *yaml.Node) (config, error) {
	var data dockerfileDataYaml

	err := verifyConfigYamlNode(node)
	if err != nil {
		return nil, err
	}
	if err := node.Decode(&data); err != nil {
		return nil, errors.Wrap(err, "error in decoding yaml data")
	}
	var configData config

	configData = data.Config
	return configData, nil
}

// Generates a new template and writes it to the Dockerfile.
func (d *dockerfileTemplate) GenerateDockerfileFromTemplate(writer io.Writer) error {
	templateString :=
		"{{- range $i, $instruction := . }}" +
			"{{- if gt $i 0 }}\n{{ end }}" +
			"{{ $instruction.WriteInstruction }}\n" +
			"{{- end }}\n\n"

	tmpl, err := template.New("dockerfile.template").Parse(templateString)
	if err != nil {
		return errors.Wrap(err, "error in parsing dockerfile template")
	}

	err = tmpl.Execute(writer, d.Data)
	if err != nil {
		return errors.Wrap(err, "error in executing the template")
	}

	return nil
}

func (d *dockerfileTemplate) GetLanguageFormat() (string, string) {
	for _, inst := range d.Data {
		switch inst := inst.(type) {
		case programmingLanguage:
			return inst.Name, inst.Version
		}
	}

	return "", ""
}
