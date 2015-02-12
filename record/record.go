package record

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type RecordCmd interface {
	Record(string)
}

type recordCmd struct {
	name   string
	cli    plugin.CliConnection
	cmdset []string
}

func NewRecordCmd(cli plugin.CliConnection) RecordCmd {
	return &recordCmd{
		cli: cli,
	}
}

func (r *recordCmd) Record(cmdsetName string) {
	var cmd string
	var err error

	r.name = strings.TrimSpace(cmdsetName)
	fmt.Printf(`Please start entering CF commands
For example: 'cf api http://api.10.244.0.34.xip.io --skip-ssl-validation'

type 'stop' to stop recording

`)

	for {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		cmd, err = in.ReadString('\n')
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
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

	r.writeFile()
}

func (r *recordCmd) writeFile() {
	f, err := os.OpenFile("./CmdSetRecords", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		quitWithError("Error writing to file: ", err)
	}
	defer f.Close()

	if _, err = f.WriteString("--" + r.name + "\n"); err != nil {
		quitWithError("Error writing to file: ", err)
	}

	for _, str := range r.cmdset {
		if _, err = f.WriteString("  " + str + "\n"); err != nil {
			quitWithError("Error writing to file: ", err)
		}
	}

}

func validCfCmd(cmd string) bool {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(cmd)), "cf ") {
		return true
	}
	return false
}

func quitWithError(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}
