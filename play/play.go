package play

import (
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/cli-recorder/data"
)

type PlayCmds interface {
	Run()
}

type playCmds struct {
	cli    plugin.CliConnection
	cmdset string
}

func NewPlayCmds(cli plugin.CliConnection, cmdset string) PlayCmds {
	return &playCmds{
		cli:    cli,
		cmdset: cmdset,
	}
}

func (p *playCmds) Run() {
	c := data.NewCmdSetData()
	cmds := c.GetCmdSet(p.cmdset)

	for _, cmd := range cmds {
		p.cli.CliCommand(strings.Split(cmd, " ")...)
	}
}
