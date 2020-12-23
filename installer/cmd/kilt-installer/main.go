package main

import (
	"os"

	"github.com/markbates/pkger"
	"github.com/urfave/cli/v2"

	"github.com/falcosecurity/kilt/installer"
	cfn_macro "github.com/falcosecurity/kilt/installer/runtimes/cfn-macro"
)

func main() {
	pkger.Include("/artifacts")
	macro := cfn_macro.CfnMacroInstaller{}
	commands := make([]*cli.Command, 0)
	for _, c := range macro.GetCommands(&installer.Hooks{}) {
		commands = append(commands, c)
	}
	app := cli.App{
		Name:     "kilt-installer",
		Commands: commands,
	}
	app.Run(os.Args)
}
