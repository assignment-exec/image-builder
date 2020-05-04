package configurations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedGenericOutput = `FROM golang:latest

ENV GOMODULE=on
ENV GOFLAGS=-mod=vendor
COPY . /code-runner
WORKDIR /code-runner
RUN git clone https://github.com/assignment-exec/code-runner.git && cd code-runner && make
EXPOSE 8082
<<<<<<< HEAD
<<<<<<< HEAD
CMD ["./code-runner/code-runner-server", "-port", "8082"]
=======
CMD ["./code-runner/code-runner-server -port 8082"]
>>>>>>> 700268d... Changes to configuration test
=======
CMD ["./code-runner/code-runner-server", "-port", "8082"]
>>>>>>> 0eef049661abf8aefcac2783673dba86eee154a9

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
