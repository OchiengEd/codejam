package main

import (
	"fmt"
	"github.com/codingconcepts/env"
)

type SystemConf struct {
	User string `env:"USER"`
	HomeDir string `env:"HOME"`
}

func main() {
	var sys SystemConf
	if err := env.Set(&sys); err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Println("Hello world")
	fmt.Printf("%+v\n",  sys)
}
