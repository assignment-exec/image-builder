package configurations

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var expectedAssgnEnvDockerfileContents = `FROM assignmentexec/code-runner:1.0
RUN ./scripts/gcc_7.sh

`

func TestAssignmentEnvDockerfileTemplate(t *testing.T) {

	err := os.Chdir("..")
	assert.NoError(t, err)

	data, err := GetConfig("assignment-env.yaml")
	assert.NoError(t, err)

	output := &bytes.Buffer{}
	_, err = output.WriteString(data.String())
	assert.NoError(t, err)

	assert.Equal(t, expectedAssgnEnvDockerfileContents, output.String())
}
