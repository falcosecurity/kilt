package main

import (
	"fmt"
	"github.com/markbates/pkger"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	pkger.Include("/artifacts/")
	c,e := pkger.Current()
	fmt.Printf("%v - %v\n", c , e)
	commands := make([]*cli.Command, 0)
	for _, c := range registerCfnMacro() {
		commands = append(commands, c)
	}
	app := cli.App{
		Name: "kilt",
		Commands: commands,
	}
	app.Run(os.Args)
}
