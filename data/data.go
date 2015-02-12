package data

import (
	"io/ioutil"
	"os"
	"strings"
)

type CmdSetData interface {
	GetCmdSet(string) []string
}

type cmdSetData struct {
}

func NewCmdSetData() CmdSetData {
	return &cmdSetData{}
}

func (c *cmdSetData) GetCmdSet(cmdset string) []string {
	f, err := os.Open("./CmdSetRecords")
	if err != nil && os.IsExist(err) {
		return []string{}
	}

	b, err := ioutil.ReadAll(f)
	contents := strings.Split(string(b), "\n")
	found := false
	var result []string
	for _, s := range contents {
		if found && strings.HasPrefix(s, "  ") {
			result = append(result, strings.TrimSpace(s))
		} else if found && strings.HasPrefix(s, "--") {
			break
		}

		if strings.ToLower(s) == "--"+strings.ToLower(cmdset) {
			found = true
		}

	}
	return result
}
