package main

import (
<<<<<<< HEAD
	"assignment-exec/image-builder/dockerConfig"
	"fmt"
=======
	"assignment-exec/image-builder/configurations"
>>>>>>> 500abb0ba41560bd9316ac21f9083176deb1b722
	"log"
)

func main() {
<<<<<<< HEAD
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
=======
	log.Println("Creating Dockerfile...")

	// Unmarshal the yaml configuration file and generate a dockerfile.
	err := configurations.WriteDockerfile()
	if err != nil {
		log.Fatalf("error while writing dockerfile %v", err)
	}
}
>>>>>>> 500abb0ba41560bd9316ac21f9083176deb1b722
