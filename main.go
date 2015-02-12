package main

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/cli-recorder/play"
	"github.com/simonleung8/cli-recorder/record"
)

type CLI_Recorder struct{}

func (c *CLI_Recorder) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "CLI-Recorder",
		Commands: []plugin.Command{
			{
				Name:     "record",
				HelpText: "record a set of CLI commands",
				UsageDetails: plugin.Usage{
					Usage: "record [COMMAND SET NAME]",
				},
			},
			{
				Name:     "play",
				HelpText: "play back a set of recorded CLI commands",
				UsageDetails: plugin.Usage{
					Usage: "play [COMMAND SET NAME]",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(CLI_Recorder))
}

func (c *CLI_Recorder) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "record" {
		if len(args) > 1 {
			r := record.NewRecordCmd(cliConnection)
			r.Record(args[1])
		} else {
			fmt.Println("Provide the recorded command set name to playback")
		}
	} else if args[0] == "play" {
		if len(args) > 1 {
			p := play.NewPlayCmds(cliConnection, args[1])
			p.Run()
		} else {
			fmt.Println("Provide the recorded command set name to playback")
		}
	}
}
