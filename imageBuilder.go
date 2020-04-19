package main

import (
	"assignment-exec/image-builder/builder"
	"assignment-exec/image-builder/configurations"
	"log"
)

func main() {
	log.Println("Creating Dockerfile...")

	// Unmarshal the yaml configuration file and generate a dockerfile.
	err := configurations.WriteDockerfile()
	if err != nil {
		log.Fatalf("error while writing dockerfile: %v", err)
	}

<<<<<<< HEAD
	authData, err := builder.GetAuthData("docker-auth.yaml")
	if err != nil {
		log.Fatalf("error while reading authetication details: %v", err)
	}

	err = builder.BuildImage(*authData)
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
	}

	err = builder.PushImageToHub(*authData)
=======
	err = builder.BuildCodeRunnerImage()
>>>>>>> 9585f69... Initial code for building image
	if err != nil {
		log.Fatalf("error while building image for code runner: %v", err)
	}
}
