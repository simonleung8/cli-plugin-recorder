package replay

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/cli-recorder/data"
)

type ReplayCmds interface {
	Run()
}

type replayCmds struct {
	cli     plugin.CliConnection
	cmdsets []string
}

func NewReplayCmds(cli plugin.CliConnection, cmdsets ...string) ReplayCmds {
	return &replayCmds{
		cli:     cli,
		cmdsets: cmdsets,
	}
}

func (p *replayCmds) Run() {
	c := data.NewCmdSetData()

	for _, cmdset := range p.cmdsets {
		cmds := c.GetCmdSet(cmdset)

		if len(cmds) == 0 {
			fmt.Printf("Command set %s not found.\n\n", cmdset)
		}

		for _, cmd := range cmds {
			p.cli.CliCommand(strings.Split(cmd, " ")...)
		}
	}
}
