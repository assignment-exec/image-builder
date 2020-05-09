package configurations

import (
	"github.com/pkg/errors"
	"io"
	"text/template"
)

type dockerConfigTemplate struct {
	Data         *dockerConfig
	Instructions []string
}

func NewDockerConfigTemplate(data *dockerConfig, instructions []string) *dockerConfigTemplate {
	return &dockerConfigTemplate{Data: data, Instructions: instructions}
}

// Generates a new template and writes it to the Dockerfile.
func (d *dockerConfigTemplate) GenerateDockerfileFromTemplate(writer io.Writer) error {
	templateString :=
		"{{- range $i, $instruction := . }}" +
			"{{- if gt $i 0 }}\n{{ end }}" +
			"{{ $instruction }}\n" +
			"{{- end }}\n\n"

	tmpl, err := template.New("dockerfile.template").Parse(templateString)
	if err != nil {
		return errors.Wrap(err, "error in parsing dockerfile template")
	}

	err = tmpl.Execute(writer, d.Instructions)
	if err != nil {
		return errors.Wrap(err, "error in executing the template")
	}

	return nil
}

func (d *dockerConfigTemplate) GetLanguageFormat() (string, string) {
	return d.Data.ProgrammingLanguage.Name, d.Data.ProgrammingLanguage.Version
}
