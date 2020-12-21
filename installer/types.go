package installer

import (
	"github.com/urfave/cli/v2"
)

type Installable interface {
	GetCommands(hooks *Hooks) []*cli.Command
}

type Hooks struct {
	PreInstall cli.ActionFunc
	PostInstall cli.ActionFunc
	OverrideConfig interface{}
}