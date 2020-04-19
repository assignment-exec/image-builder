# Image Builder
Image builder is an application to build docker image for running assignments. The image is built using user specific configurations.

## Configurations
There are two types of configurations provided - Code runner and Assignment Environment.

### Code Runner Configurations
- These configurations are required for running the [code-runner](https://github.com/assignment-exec/code-runner) server.
- Code runner configurations are specified in yaml format and include instructions that are needed in Dockerfile for running the server.
Following is a sample of the configuration yaml
```commandline
stages:
  serverConfig:
    - from:
        image: golang:latest
        as: golang
    - env:
        GOMODULE: on
        GOFLAGS: -mod=vendor
    - copy:
        basedir: .
        destdir: /code-runner
    - workdir:
        dir: /code-runner
    - runCommand:
        param: make
    - serverPort:
        number: 8082
    - cmd:
        params:
          - ./code-runner-server
```

### Assignment Environment Configurations
- These configurations are provided by professors or teaching staffs. 
- Professor provides assignment specific configurations for creating suitable docker environment.
- These configurations include specifying operating system and its version, compiler and commands required to run the specific assignment.

## Build Image
- Code Runner configurations are used to build a docker image locally.
- Prerequisite for building an image is that docker engine should be installed.
### Docker Setup
See [instructions](https://docs.docker.com/engine/installation/) for installing docker engine on different supported platforms.


