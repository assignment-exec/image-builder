package main

import (
	"assignment-exec/image-builder/dockerConfig"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Creating Dockerfile...")
<<<<<<< HEAD
	config,err := dockerConfig.GetConfig("config.yaml")
=======

	err := dockerConfig.WriteDockerfile()
>>>>>>> 5917893... Minor changes and fixes
	if err != nil {
		log.Fatalf("error while writing dockerfile %v", err)
	}
<<<<<<< HEAD

	fmt.Println(config)
}
=======
}
>>>>>>> 5917893... Minor changes and fixes
