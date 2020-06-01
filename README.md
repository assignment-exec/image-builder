# Image Builder
Image builder is an application to build docker image for running assignments. The image is built using user specific configurations.

## Configurations
- The configurations required for assignment environment is provided in yaml format.
- A docker image for [code-runner](https://github.com/assignment-exec/code-runner) application is used as base image for any assignment environment.
- The code-runner application is used to build and run assignments.

### Assignment Environment Configurations
- Provides the docker image tag corresponding to a particular version of code-runner. In addition, this configuration file includes the runtime environment requirements for the assignment. This includes the following.
    - Programming language used by students.
    - Additional libraries that are needed, if any.
Following is a sample of the configuration.
```commandline
baseImage: "assignmentexec/code-runner:1.0"
dependencies:
  lang: python
  langVersion: 3.7
  lib:
    numpy:
      cmd: pip3 install numpy
    scipy:
      cmd: pip3 install scipy
```

## Supported Languages
Below is the list of supported languages and the corresponding versions.
- gcc 7
- g++ 7
- python 3.7
- java 8 & 11

### Installation scripts
- Every supported language and its version has an installation script stored in `scripts` directory.
- The scripts are named as `<language_version>.sh`. Example - For language - java and version - 8, script name should be `java_8.sh`.
- The installation scripts are shell scripts that hold commands to install the corresponding language and its version.
- To add support for a new language and version, add a new shell script that holds the appropriate commands for installation.

## Build and Publish Image
- Using above configurations docker images are built locally and published to the docker hub.
- Prerequisite for building an image is that docker engine should be installed.
### Docker Setup
See [instructions](https://docs.docker.com/engine/installation/) for installing docker engine on different supported platforms.
### Compile and Run
Compile the source code using the make tool as shown below.
```commandline
make
```
Use the -h option to get information about the command-line options.
- Use the `-assignmentEnvConfigFilepath` option to specify the path to assignment environment config file.
- Use the `-dockerfileLoc` option to specify the Dockerfile location to be created.
- Use the `-publishImage` option to specify whether to publish image to docker hub.
Below is an example to run the source code.
```commandline
./image-builder -assignmentEnvConfigFilepath <path_to_config_file> -dockerfileLoc <dockerfile_location> -publishImage <true/false>
```
## Run Docker Image for Assignment Environment
Following is the command used to run the docker image for assignment environment.
```commandline
docker run --publish <port_to_expose>:<port_to_run> <assignment_environment_image> -port <port_to_run>
```

