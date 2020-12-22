package main

import (
	"github.com/markbates/pkger"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	pkger.Include("/artifacts")

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
