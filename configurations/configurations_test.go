package configurations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedGenericOutput = `FROM golang:latest As golang

ENV GOMODULE=on
ENV GOFLAGS=-mod=vendor
COPY . /code-runner
WORKDIR /code-runner
RUN go build -o code-runner-server
EXPOSE 8082
CMD ["./code-runner-server"]

`

// Tests dockerfile template generation
func TestDockerfileTemplate(t *testing.T) {
	data, err := newDockerFileDataFromYamlFile("../code-runner.yaml")
	tmpl := newDockerfileTemplate(data)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.generateDockerfileFromTemplate(output)
	assert.NoError(t, err)

	assert.Equal(t, expectedGenericOutput, output.String())
}