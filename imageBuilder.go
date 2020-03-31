package main

import (
	"assignment-exec/image-builder/configurations"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Creating Dockerfile...")

	err := configurations.WriteDockerfile()
	if err != nil {
		log.Fatalf("error while writing dockerfile %v", err)
	}
}
