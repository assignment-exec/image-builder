package main

import (
	"assignment-exec/image-builder/readConfig"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Creating Dockerfile...")
	config,err := readConfig.GetConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}