package main

import (
	"github.com/google/go-containerregistry/pkg/crane"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		println("needs 1 argument: image")
		return
	}
	res, err := crane.Config(os.Args[1])
	if err != nil {
		panic(err)
	}
	println(string(res))
}
