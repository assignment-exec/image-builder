package builder

import (
	"fmt"
	"os/exec"
)

func BuildCodeRunnerImage() error {

	cmdStr := "sudo docker build -t test_image:1.0 ."
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	fmt.Printf("%s", out)

	return  err
}