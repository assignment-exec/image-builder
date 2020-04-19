package configurations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedGenericOutput = `FROM golang:latest

ENV GOFLAGS=-mod=vendor
ENV GOMODULE=on
COPY . /code-runner
WORKDIR /code-runner
RUN git clone https://github.com/assignment-exec/code-runner.git && cd code-runner && make
EXPOSE 8082
CMD ["./code-runner/code-runner-server -port 8082"]

`

// Tests dockerfile template generation.
func TestDockerfileTemplate(t *testing.T) {
	data, err := newDockerFileDataFromYamlFile("../code-runner.yaml")
	tmpl := newDockerfileTemplate(data)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.generateDockerfileFromTemplate(output)
	assert.NoError(t, err)

	assert.Equal(t, expectedGenericOutput, output.String())
}
