package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"assignment-exec/image-builder/environment"
	"flag"
	"log"
	"os"
)

var pushImage = flag.Bool("publishImage", false, "Push image to docker hub")

func main() {
	flag.Parse()
	log.Println("Creating Dockerfile...")

	// Setting environment variables
	err := os.Setenv("CODE_RUNNER_YAML", "code-runner.yaml")
	if err != nil {
		log.Fatalf("error while setting environment variables: %v", err)
		return
	}
	err = os.Setenv("DOCKER_AUTH_YAML", "docker-auth.yaml")
	if err != nil {
		log.Fatalf("error while setting environment variables: %v", err)
		return
	}
	err = os.Setenv("DOCKERFILE_PATH", "Dockerfile")
	if err != nil {
		log.Fatalf("error while setting environment variables: %v", err)
		return
	}
	err = os.Setenv("DOCKER_IO_PATH", "docker.io")
	if err != nil {
		log.Fatalf("error while setting environment variables: %v", err)
		return
	}

	// Unmarshal the yaml configuration file and generate a dockerfile.
	err = configurations.WriteDockerfile()
	if err != nil {
		log.Fatalf("error while writing dockerfile: %v", err)
		return
	}

	imgBuilder := builder.NewImageBuilder(os.Getenv(environment.DockerAuthYaml))

	err = imgBuilder.BuildImage()
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
		return
	}

	if *pushImage {
		err = imgBuilder.PublishImage()
		if err != nil {
			log.Fatalf("error while pushing image to docker hub: %v", err)
		}
	}
}
