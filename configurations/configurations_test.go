package configurations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedCodeRunnerOutput = `FROM golang:latest

ENV GOMODULE=on
ENV GOFLAGS=-mod=vendor
COPY . /code-runner
WORKDIR /code-runner
RUN git clone https://github.com/assignment-exec/code-runner.git && cd code-runner && make
EXPOSE 8082
CMD ["./code-runner/code-runner-server", "-port", "8082"]
<<<<<<< HEAD

`

var expectedAssgnEnvOutput = `FROM assignmentexec/trial2:2.0
RUN apt-get update && apt-get -y install gcc-7

`

var expectedAssgnEnvOutput = `FROM assignmentexec/trial2:2.0
RUN apt-get update && apt-get -y install gcc-7

`

// Tests dockerfile template generation.
func TestDockerfileTemplate(t *testing.T) {
	data, err := newDockerFileDataFromYamlFile("../code-runner.yaml")
	tmpl := newDockerfileTemplate(data)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.generateDockerfileFromTemplate(output)
	assert.NoError(t, err)

	assert.Equal(t, expectedCodeRunnerOutput, output.String())
}

func TestAssignmentEnvDockerfileTemplate(t *testing.T) {
	data, err := newDockerFileDataFromYamlFile("../assignment-env.yaml")
	tmpl := newDockerfileTemplate(data)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.generateDockerfileFromTemplate(output)
	assert.NoError(t, err)

	assert.Equal(t, expectedAssgnEnvOutput, output.String())
}
