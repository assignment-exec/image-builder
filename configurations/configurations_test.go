package configurations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var expectedCodeRunnerDockerfileContents = `FROM golang:latest

ENV GOMODULE=on
ENV GOFLAGS=-mod=vendor
COPY . /code-runner
WORKDIR /code-runner
RUN git clone https://github.com/assignment-exec/code-runner.git && cd code-runner && make
EXPOSE 8082
CMD ["./code-runner/code-runner-server", "-port", "8082"]

`

var expectedAssgnEnvDockerfileContents = `FROM assignmentexec/code-runner:1.0
RUN ./scripts/gcc_7.sh 

`

// Tests dockerfile template generation.
func TestDockerfileTemplate(t *testing.T) {
	data, instructions, err := ParseInstructions("../code-runner.yaml")
	tmpl := NewDockerConfigTemplate(data, instructions)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.GenerateDockerfileFromTemplate(output)
	assert.NoError(t, err)

	assert.Equal(t, expectedCodeRunnerDockerfileContents, output.String())
}

func TestAssignmentEnvDockerfileTemplate(t *testing.T) {

	os.Chdir("..")
	data, instructions, err := ParseInstructions("assignment-env.yaml")
	tmpl := NewDockerConfigTemplate(data, instructions)
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	err = tmpl.GenerateDockerfileFromTemplate(output)
	assert.NoError(t, err)

	assert.Equal(t, expectedAssgnEnvDockerfileContents, output.String())
}
