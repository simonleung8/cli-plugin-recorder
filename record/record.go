package record

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/cli-plugin-recorder/data"
)

type RecordCmd interface {
	Record(string)
	ListCmdSets()
	ClearCmdSets()
	ListCmds(string)
	DeleteCmdSet(string)
}

type recordCmd struct {
	name        string
	cli         plugin.CliConnection
	cmdset      []string
	cmdSetData  data.CmdSetData
	inputStream *os.File
}

func NewRecordCmd(cli plugin.CliConnection, cmdSetData data.CmdSetData, inputStream *os.File) RecordCmd {
	return &recordCmd{
		cli:         cli,
		cmdSetData:  cmdSetData,
		inputStream: inputStream,
	}
}

func (r *recordCmd) ListCmdSets() {
	fmt.Println("The following command sets are available:")
	fmt.Println("=========================================")

	cmdsets := r.cmdSetData.ListCmdSetNames()

	for _, name := range cmdsets {
		fmt.Println(name)
	}
	fmt.Println()
}

func (r *recordCmd) ListCmds(cmdset string) {
	fmt.Printf("The following commands are in %s:\n", cmdset)
	fmt.Println("=========================================")

	cmds := r.cmdSetData.GetCmdSet(cmdset)

	for _, cmd := range cmds {
		fmt.Println(cmd)
	}
	fmt.Println()
}

func (r *recordCmd) Record(cmdsetName string) {
	var cmd string
	var err error

	r.name = strings.TrimSpace(cmdsetName)

	if r.cmdSetData.IsCmdExist(r.name) {
		fmt.Printf("\nThe name %s already exists, please use another ...\n\n", r.name)
		return
	}

	fmt.Printf(`Please start entering CF commands
For example: 'cf api http://api.10.244.0.34.xip.io --skip-ssl-validation'

type 'stop' to stop recording and save
type 'quit' to quit recording without saving

`)

	for {
		fmt.Printf("\n(recording) >> ")
		in := bufio.NewReader(r.inputStream)
		cmd, err = in.ReadString('\n')
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if strings.TrimSpace(cmd) == "quit" {
			return
		}

		if strings.TrimSpace(cmd) == "stop" {
			break
		}

		if validCfCmd(cmd) {
			cmd = strings.TrimSpace(cmd)[3:]
			r.cmdset = append(r.cmdset, cmd)
			r.cli.CliCommand(strings.Split(cmd, " ")...)
		} else {
			fmt.Printf("Invalid CF command\n\n")
		}
	}

	r.cmdSetData.SaveCmdSet(r.name, r.cmdset)
}

func (r *recordCmd) DeleteCmdSet(cmdset string) {
	r.cmdSetData.DeleteCmdSet(cmdset)
}

func (r *recordCmd) ClearCmdSets() {
	fmt.Print("WARNING: You are about to delete all recorded command sets (y or n): ")
	in := bufio.NewReader(r.inputStream)
	response, err := in.ReadString('\n')
	if err != nil {
		fmt.Println("Error: ", err)
	}

	if strings.TrimSpace(response) == "y" {
		r.cmdSetData.ClearCmdSets()
		fmt.Println("All recorded command sets has been removed")
	} else {
		fmt.Println("Action aborted.")
	}
}

func validCfCmd(cmd string) bool {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(cmd)), "cf ") {
		return true
	}
	return false
}
