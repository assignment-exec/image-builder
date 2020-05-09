# Image Builder
Image builder is an application to build docker image for running assignments. The image is built using user specific configurations.

## Configurations
There are two types of configurations provided in yaml format.

### Code Runner Configurations
Contains information to containerize the [code-runner](https://github.com/assignment-exec/code-runner) application and also to run it.
Following is a sample of the configuration.
```commandline
from:
  image: golang:latest
env:
  GOMODULE: on
  GOFLAGS: -mod=vendor
copy:
  basedir: .
  destdir: /code-runner
workdir:
  dir: /code-runner
runCommand:
  param: git clone https://github.com/assignment-exec/code-runner.git
    && cd code-runner
    && make
port:
  number: 8082
cmd:
  params:
    - ./code-runner/code-runner-server
    - -port
    - 8082
```

### Assignment Environment Configurations
- Provides the docker image tag corresponding to a particular version of the code-runner application that can be used to run assignments. In addition, this configuration file includes the runtime environment requirements for the assignment. This includes the following.
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

## Build and Publish Image
- Using above configurations docker images are built locally and published to the docker hub.
- Prerequisite for building an image is that docker engine should be installed.
### Docker Setup
See [instructions](https://docs.docker.com/engine/installation/) for installing docker engine on different supported platforms.

## Run Docker Image for Assignment Environment
Following is the command used to run the docker image for assignment environment.
```commandline
docker run --publish <code-runner_server_port> <assignment_environment_image>
```



