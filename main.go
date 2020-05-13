package main

import (
	"assignment-exec/image-builder/builder"
	"flag"
	"log"
)

var publishImage = flag.Bool("publishImage", false, "push image to docker hub")
var assignmentEnvConfigFilename = flag.String("assignmentEnvConfig", "assignment-env.yaml", "Assignment Environment configuration filename")
var dockerfileLoc = flag.String("dockerfileLoc", "Dockerfile", "Dockerfile Location")

func main() {

	flag.Parse()

	assgnEnv, err := builder.GetConfigurations(*publishImage, *assignmentEnvConfigFilename, *dockerfileLoc)
	if err != nil {
		log.Fatalf("error in getting configurations: %v", err)
	}

	buildManager, err := builder.NewBuildManager(builder.WithCommands(assgnEnv))
	if err != nil {
		log.Fatalf("error in creating a builder: %v", err)
	}
	if err = buildManager.ExecuteCommands(); err != nil {
		log.Fatalf("error in building assignment environment image: %v", err)
	}
}
