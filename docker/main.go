package main

import (
	"fmt"

	"github.com/OchiengEd/codejam/docker/checks"
)

func main() {
	fmt.Println("hello world")

	cli := checks.DockerClient()
	cli.Operations()
}
