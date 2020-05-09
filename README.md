# Image Builder
Image builder is an application to build docker image for running assignments. The image is built using user specific configurations.

## Configurations
There are two types of configurations provided in yaml format - Code runner and Assignment Environment.
### Code Runner Configurations
- These configurations are required for running the [code-runner](https://github.com/assignment-exec/code-runner) server.
- Code runner configurations include instructions that are needed in Dockerfile for running the server.
Following is a sample of the configuration yaml
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
- These configurations are provided by professors or teaching staffs. 
- Professor provides assignment specific configurations for creating suitable docker environment.
- Assignment Environment configurations include base [code-runner](https://github.com/assignment-exec/code-runner) image, programming language and its version.
Following is a sample of the configuration yaml
```commandline
from:
  image: assignmentexec/code-runner:1.0
programmingLanguage:
  name: gcc
  version: 7
```

## Build and Publish Image
- Docker images for above mentioned configurations are built locally and published to the docker hub.
- Prerequisite for building an image is that docker engine should be installed.
### Docker Setup
See [instructions](https://docs.docker.com/engine/installation/) for installing docker engine on different supported platforms.

## Run Docker image for Assignment Environment
Following is the command used to run the docker image for assignment environment
```commandline
docker run --publish <code-runner_server_port> <assignment_environment_image>
```



