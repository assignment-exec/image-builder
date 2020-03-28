package main

import (
	"assignment-exec/image-builder/dockerConfig"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Creating Dockerfile...")
	config,err := dockerConfig.GetConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}