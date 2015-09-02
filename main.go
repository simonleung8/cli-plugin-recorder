package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/cli-plugin-recorder/data"
	"github.com/simonleung8/cli-plugin-recorder/record"
	"github.com/simonleung8/cli-plugin-recorder/replay"

	"github.com/simonleung8/flags"
)

type CLI_Recorder struct{}

func (c *CLI_Recorder) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "CLI-Recorder",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 2,
		},
		Commands: []plugin.Command{
			{
				Name:     "record",
				HelpText: "record a set of CLI commands for playback",
				UsageDetails: plugin.Usage{
					Usage: `record COMMAND_SET_NAME | OPTIONS 

OPTIONS:
  -n            list all commands within a set (e.g. -n COMMAND_SET_NAME)
  -d            delete a command set (e.g. -d COMMAND_SET_NAME)
  --list, -l    list all recorded command sets
  --clear, -c   clear all recorded commands
`,
				},
			},
			{
				Name:     "replay",
				Alias:    "rp",
				HelpText: "replay a set of recorded CLI commands",
				UsageDetails: plugin.Usage{
					Usage: `replay COMMAND_SET_NAME | OPTIONS
					
OPTIONS:
  --list, -l    list all recorded command sets
`,
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(CLI_Recorder))
}

func (c *CLI_Recorder) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {

	} else if args[0] == "record" {
		fc := flags.New()
		fc.NewBoolFlag("list", "l", "list all recorded command sets")
		fc.NewBoolFlag("clear", "c", "clear all recorded command sets")
		fc.NewStringFlag("n", "", "list all commands within a set (e.g. -n COMMAND_SET_NAME)")
		fc.NewStringFlag("d", "", "to delete a command set (e.g. -d COMMAND_SET_NAME)")
		err := fc.Parse(args[1:]...)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		c.runRecord(cliConnection, fc)

	} else if args[0] == "replay" {
		fc := flags.New()
		fc.NewBoolFlag("list", "l", "list all recorded command sets")
		err := fc.Parse(args[1:]...)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if fc.Bool("l") || fc.Bool("list") {
			r := record.NewRecordCmd(cliConnection, data.NewCmdSetData(), os.Stdin)
			r.ListCmdSets()
		} else if len(args) > 1 {
			p := replay.NewReplayCmds(cliConnection, data.NewCmdSetData(), args[1:]...)
			p.Run()
		} else {
			fmt.Println("Provide the recorded command set name to playback")
			fmt.Printf("\nUSAGE:\n  %s\n", c.GetMetadata().Commands[1].UsageDetails.Usage)
		}
	}
}

func (c *CLI_Recorder) runRecord(cliConnection plugin.CliConnection, fc flags.FlagContext) {
	r := record.NewRecordCmd(cliConnection, data.NewCmdSetData(), os.Stdin)
	if fc.Bool("l") || fc.Bool("list") {
		r.ListCmdSets()
	} else if fc.Bool("clear") {
		r.ClearCmdSets()
	} else if fc.IsSet("n") {
		r.ListCmds(fc.String("n"))
	} else if fc.IsSet("d") {
		r.DeleteCmdSet(fc.String("d"))
	} else if len(fc.Args()) == 1 {
		r.Record(fc.Args()[0])
	} else {
		fmt.Println("Provide command set name for recording")
		fmt.Printf("\nUSAGE:\n  %s\n", c.GetMetadata().Commands[0].UsageDetails.Usage)
	}
}
