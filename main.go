package main

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/cli-plugin-recorder/record"
	"github.com/simonleung8/cli-plugin-recorder/replay"
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
					Usage: `record [COMMAND SET NAME] | -l | -n [COMMAND SET NAME] | -d [COMMAND SET NAME] | -clear

Options:
-l : to list all the record command sets
-n : to list all commands withint a set
-d : to delete a command set
-clear : clear all record commands
`,
				},
			},
			{
				Name:     "replay",
				HelpText: "replay a set of recorded CLI commands",
				UsageDetails: plugin.Usage{
					Usage: "replay [COMMAND SET NAME]",
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
		runRecord(cliConnection, args)
	} else if args[0] == "replay" {
		if len(args) > 1 {
			p := replay.NewReplayCmds(cliConnection, args[1:]...)
			p.Run()
		} else {
			fmt.Println("Provide the recorded command set name to playback")
		}
	}
}

func runRecord(cliConnection plugin.CliConnection, args []string) {
	r := record.NewRecordCmd(cliConnection)
	if len(args) == 2 {
		if args[1] == "-l" {
			r.ListCmdSets()
		} else if args[1] == "-clear" {
			r.ClearCmdSets()
		} else {
			r.Record(args[1])
		}
	} else if len(args) == 3 {
		if args[1] == "-n" {
			r.ListCmds(args[2])
		} else if args[1] == "-d" {
			r.DeleteCmdSet(args[2])
		}
	} else {
		fmt.Println("Provide the recorded command set name to playback")
	}
}
