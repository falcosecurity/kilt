package main

import (
	"github.com/falcosecurity/kilt/installer"
	cfn_macro "github.com/falcosecurity/kilt/installer/runtimes/cfn-macro"
	"github.com/markbates/pkger"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	pkger.Include("/artifacts")
	macro := cfn_macro.CfnMacroInstaller{}
	commands := make([]*cli.Command, 0)
	for _, c := range macro.GetCommands(&installer.Hooks{}) {
		commands = append(commands, c)
	}
	app := cli.App{
		Name: "kilt-installer",
		Commands: commands,
	}
	app.Run(os.Args)
}
