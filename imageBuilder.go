package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"assignment-exec/image-builder/environment"
	"flag"
	"log"
)

var pushImage = flag.Bool("pushImage", false, "Push image to docker hub")

func init() {
	flag.BoolVar(pushImage, "p", false, "Push image to docker hub")
}

func main() {
	flag.Parse()
	log.Println("Creating Dockerfile...")

	// Unmarshal the yaml configuration file and generate a dockerfile.
	err := configurations.WriteDockerfile()
	if err != nil {
		log.Fatalf("error while writing dockerfile: %v", err)
	}

	authData, err := builder.GetAuthData(environment.DockerAuthYaml)
	if err != nil {
		log.Fatalf("error while reading authetication details: %v", err)
	}

	err = builder.BuildImage(*authData)
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
	}

	if *pushImage {
		err = builder.PushImageToHub(*authData)
		if err != nil {
			log.Fatalf("error while pushing image to docker hub: %v", err)
		}
	}
}