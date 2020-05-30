package main

import (
	"assignment-exec/image-builder/builder"
	"flag"
	"log"
)

var publishImage = flag.Bool("publishImage", false, "Publish image to docker hub")
var assignmentEnvConfigFilepath = flag.String("assignmentEnvConfigFilepath", "assignment-env.yaml", "Assignment Environment configuration filepath")
var dockerfileLoc = flag.String("dockerfileLoc", "Dockerfile", "Location for dockerfile to be created")

func main() {

	flag.Parse()

	asgmtEnv, err := builder.GetConfigurations(*publishImage, *assignmentEnvConfigFilepath, *dockerfileLoc)
	if err != nil {
		log.Fatalf("error in getting configurations: %v", err)
	}

	buildManager, err := builder.NewBuildManager(builder.WithCommands(asgmtEnv))
	if err != nil {
		log.Fatalf("error in creating a builder: %v", err)
	}
	if err = buildManager.ExecuteCommands(); err != nil {
		log.Fatalf("error in building assignment environment image: %v", err)
	}
}
