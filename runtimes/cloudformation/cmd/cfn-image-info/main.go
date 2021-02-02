package main

import (
	"fmt"
	"os"

	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
)

func main() {
	if len(os.Args) < 2 {
		println("needs 1 argument: image")
		return
	}

	res, err := cfnpatcher.GetConfigFromRepository(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
}
