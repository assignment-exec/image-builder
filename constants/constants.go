// Package constants provide list of constant values
// required for building and publishing the assignment
// environment image.
package constants

const InstallationScriptsDir = "scripts"
const DockerIO = "docker.io"
const BuildContextTar = "buildContext.tar"

const CodeRunnerDir = "code-runner"
const DockerRunCommand = "docker run --publish 52453:52453"
const PortCmdArg = "-port 52453"
