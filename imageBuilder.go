package main

import (
	"assignment-exec/image-builder/dockerConfig"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Creating Dockerfile...")

	err := dockerConfig.WriteDockerfile()
	if err != nil {
		log.Fatalf("error while writing dockerfile %v", err)
	}
}
