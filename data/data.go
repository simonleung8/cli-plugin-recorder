package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type CmdSetData interface {
	GetCmdSet(string) []string
	SaveCmdSet(string, []string)
	DeleteCmdSet(string)
	IsCmdExist(t string) bool
	ListCmdSetNames() []string
	ClearCmdSets()
}

type cmdSetData struct {
}

func NewCmdSetData() CmdSetData {
	return &cmdSetData{}
}

func (c *cmdSetData) GetCmdSet(cmdset string) []string {
	b, err := getAllData()
	if err != nil && !os.IsExist(err) {
		return []string{}
	} else if err != nil {
		fmt.Println("Error reading from file: ", err)
	}

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

func (c *cmdSetData) ListCmdSetNames() []string {
	b, err := getAllData()
	if err != nil && !os.IsExist(err) {
		return []string{}
	} else if err != nil {
		fmt.Println("Error reading from file: ", err)
	}

	contents := strings.Split(string(b), "\n")

	var result []string
	for _, s := range contents {
		if strings.HasPrefix(s, "--") {
			result = append(result, s[2:])
		}
	}

	return result
}

func (c *cmdSetData) IsCmdExist(cmdset string) bool {
	b, err := getAllData()
	if err != nil && !os.IsExist(err) {
		return false
	} else if err != nil {
		fmt.Println("Error reading from file: ", err)
	}

	contents := strings.Split(string(b), "\n")

	for _, s := range contents {
		if strings.ToLower(s) == "--"+strings.ToLower(cmdset) {
			return true
		}

	}
	return false
}

func (c *cmdSetData) SaveCmdSet(cmdsetName string, cmdset []string) {
	var f *os.File

	fp := getFilePath()
	_, err := os.Stat(fp)
	if err != nil && !os.IsExist(err) {
		f, err = os.Create(fp)
	} else {
		f, err = os.OpenFile(fp, os.O_RDWR|os.O_APPEND, 0660)
	}
	if err != nil {
		quitWithError("Error writing to file: ", err)
	}
	defer f.Close()

	if _, err = f.WriteString("--" + cmdsetName + "\n"); err != nil {
		quitWithError("Error writing to file: ", err)
	}

	for _, str := range cmdset {
		if _, err = f.WriteString("  " + str + "\n"); err != nil {
			quitWithError("Error writing to file: ", err)
		}
	}

}

func (c *cmdSetData) ClearCmdSets() {
	err := os.Remove(getFilePath())
	if err != nil {
		fmt.Println("Error removing file: ", err)
		os.Exit(1)
	}

	fmt.Println("All record command sets has be removed")
}

func (c *cmdSetData) DeleteCmdSet(cmdsetName string) {
	cmds := c.GetCmdSet(cmdsetName)

	inputf, err := os.Open(getFilePath())
	if err != nil && !os.IsExist(err) {
		fmt.Printf("\n%s does not exist in the command sets", cmdsetName)
		return
	}
	defer inputf.Close()

	b, err := ioutil.ReadAll(inputf)
	b = bytes.Replace(b, []byte("--"+cmdsetName+"\n"), []byte(""), 1)
	for _, cmd := range cmds {
		b = bytes.Replace(b, []byte("  "+cmd+"\n"), []byte(""), 1)
	}

	fp := getFilePath()
	f, err := os.Create(fp)
	if err != nil {
		quitWithError("Error writing to file: ", err)
	}
	defer f.Close()

	if _, err = f.Write(b); err != nil {
		quitWithError("Error writing to file: ", err)
	}
}

func getFilePath() string {
	fp := filepath.Join(userHomeDir(), ".cf", "plugins")
	if os.Getenv("CF_PLUGIN_HOME") != "" {
		fp = filepath.Join(os.Getenv("CF_PLUGIN_HOME"), ".cf", "plugins")
	}

	return filepath.Join(fp, "CmdSetRecords")
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDIRVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

func getAllData() ([]byte, error) {
	f, err := os.Open(getFilePath())
	if err != nil {
		return []byte{}, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func quitWithError(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}
