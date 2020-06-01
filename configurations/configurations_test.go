// Package configurations implements routines to read and store the
// assignment environment configuration yaml file, get the docker instructions
// in the specific format for every configuration.
package configurations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var expectedAsgmtEnvDockerfileContents = `FROM assignmentexec/code-runner:1.0
COPY . /code-runner
RUN ./scripts/gcc_7.sh
ENV SUPPORTED_LANGUAGE gcc

`

// TestAssignmentEnvDockerfileTemplate tests the Dockerfile generation
// for the assignment environment config.
func TestAssignmentEnvDockerfileTemplate(t *testing.T) {

	err := os.Chdir("..")
	assert.NoError(t, err)

	data, err := GetAssignmentEnvConfig("assignment-env.yaml")
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	_, err = output.WriteString(data.GetInstruction())
	assert.NoError(t, err)

	assert.Equal(t, expectedAsgmtEnvDockerfileContents, output.String())
}
