package replay

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/cli-plugin-recorder/data"
)

type ReplayCmds interface {
	Run()
}

type replayCmds struct {
	cli        plugin.CliConnection
	cmdsets    []string
	cmdSetData data.CmdSetData
}

func NewReplayCmds(cli plugin.CliConnection, cmdSetData data.CmdSetData, cmdsets ...string) ReplayCmds {
	return &replayCmds{
		cli:        cli,
		cmdsets:    cmdsets,
		cmdSetData: cmdSetData,
	}
}

func (p *replayCmds) Run() {
	for _, cmdset := range p.cmdsets {
		cmds := p.cmdSetData.GetCmdSet(cmdset)

		if len(cmds) == 0 {
			fmt.Printf("Command set %s not found.\n\n", cmdset)
		}

		for _, cmd := range cmds {
			p.cli.CliCommand(strings.Split(cmd, " ")...)
		}
	}
}
