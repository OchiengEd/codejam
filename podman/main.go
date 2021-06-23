package main

import (
	"fmt"

	"github.com/OchiengEd/codejam/podman/checks"
)

func main() {
	fmt.Println("hello world")

	checks.PodmanSocket()
}
