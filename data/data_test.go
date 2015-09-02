package data_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/simonleung8/cli-plugin-recorder/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data", func() {
	var (
		dataStore data.CmdSetData
	)

	BeforeEach(func() {
		wd, err := os.Getwd()
		Ω(err).ToNot(HaveOccurred())
		os.Setenv("CF_PLUGIN_HOME", filepath.Join(wd, "test_fixtures"))

		dataStore = data.NewCmdSetData()
	})

	Describe("GetCmdSet()", func() {
		It("returns a list of recorded commands in the command set", func() {
			Ω(dataStore.GetCmdSet("cmdSet1")).To(Equal([]string{"api api.sample.com", "auth admin admin", "create-org test-org"}))
		})

		It("returns empty []string when data file is not found", func() {
			os.Setenv("CF_PLUGIN_HOME", "path_to_no_where/")

			cmdSet := dataStore.GetCmdSet("cmdSet1")
			Ω(cmdSet).To(Equal([]string{}))
		})
	})

	Describe("ListCmdSetNames()", func() {
		It("returns a list of recorded commands set names", func() {
			Ω(dataStore.ListCmdSetNames()).To(Equal([]string{"cmdSet1", "cmdSet2"}))
		})

		It("returns empty []string when data file is not found", func() {
			os.Setenv("CF_PLUGIN_HOME", "path_to_no_where/")

			cmdSet := dataStore.ListCmdSetNames()
			Ω(cmdSet).To(Equal([]string{}))
		})
	})

	Describe("IsCmdExist()", func() {
		It("returns true if a command set exists", func() {
			Ω(dataStore.IsCmdExist("cmdSet1")).To(BeTrue())
		})

		It("returns false if a command set exists", func() {
			Ω(dataStore.IsCmdExist("no_cmd_set")).To(BeFalse())
		})
	})

	Describe("SaveCmdSet()", func() {
		var (
			wd             string
			err            error
			cmdSetName     string
			cmds           []string
			CF_PLUGIN_HOME string
		)

		BeforeEach(func() {
			wd, err = os.Getwd()
			Ω(err).ToNot(HaveOccurred())

			CF_PLUGIN_HOME = filepath.Join(wd, "test_fixtures", "tmp")
			os.Setenv("CF_PLUGIN_HOME", CF_PLUGIN_HOME)

			err = os.MkdirAll(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins"), 0700)
			Ω(err).ToNot(HaveOccurred())

			cmdSetName = "sample_set"
			cmds = []string{"cmd1", "cmd2"}
		})

		AfterEach(func() {
			err = os.RemoveAll(os.Getenv("CF_PLUGIN_HOME"))
			Ω(err).ToNot(HaveOccurred())
		})

		Context("When data file does not exist", func() {
			It("creates a file called CmdSetRecords", func() {
				_, err = os.Stat(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins", "CmdSetRecords"))
				Ω(os.IsExist(err)).To(BeFalse())

				dataStore.SaveCmdSet(cmdSetName, cmds)
				_, err = os.Stat(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins", "CmdSetRecords"))
				Ω(err).ToNot(HaveOccurred())
			})

			It("writes the recorded command set to file", func() {
				dataStore.SaveCmdSet(cmdSetName, cmds)

				f, err := os.Open(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins", "CmdSetRecords"))
				Ω(err).ToNot(HaveOccurred())

				contents, err := ioutil.ReadAll(f)
				Ω(err).ToNot(HaveOccurred())

				Ω(string(contents)).To(ContainSubstring("--sample_set"))
				Ω(string(contents)).To(ContainSubstring("  cmd1"))
				Ω(string(contents)).To(ContainSubstring("  cmd2"))
			})
		})

		Context("When data file already exists", func() {
			BeforeEach(func() {
				dataStore.SaveCmdSet(cmdSetName, cmds)
			})

			It("appends the recorded command set to file", func() {
				dataStore.SaveCmdSet("sample_set2", []string{"cmda", "cmdb"})

				f, err := os.Open(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins", "CmdSetRecords"))
				Ω(err).ToNot(HaveOccurred())

				contents, err := ioutil.ReadAll(f)
				Ω(err).ToNot(HaveOccurred())

				Ω(string(contents)).To(ContainSubstring("--sample_set"))
				Ω(string(contents)).To(ContainSubstring("  cmd1"))
				Ω(string(contents)).To(ContainSubstring("  cmd2"))
				Ω(string(contents)).To(ContainSubstring("--sample_set2"))
				Ω(string(contents)).To(ContainSubstring("  cmda"))
				Ω(string(contents)).To(ContainSubstring("  cmdb"))
			})
		})
	})

	Describe("Removing data", func() {
		var (
			wd             string
			err            error
			CF_PLUGIN_HOME string
		)

		BeforeEach(func() {
			wd, err = os.Getwd()
			Ω(err).ToNot(HaveOccurred())

			CF_PLUGIN_HOME = filepath.Join(wd, "test_fixtures", "tmp")
			os.Setenv("CF_PLUGIN_HOME", CF_PLUGIN_HOME)

			err = os.MkdirAll(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins"), 0700)
			Ω(err).ToNot(HaveOccurred())

			dataStore.SaveCmdSet("sample_set", []string{"cmd1", "cmd2"})
		})

		AfterEach(func() {
			os.RemoveAll(os.Getenv("CF_PLUGIN_HOME"))
		})

		Describe("ClearCmdSet()", func() {
			It("deletes the entire data file", func() {
				_, err = os.Stat(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins", "CmdSetRecords"))
				Ω(err).ToNot(HaveOccurred())

				dataStore.ClearCmdSets()

				_, err = os.Stat(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins", "CmdSetRecords"))
				Ω(os.IsExist(err)).To(BeFalse())
			})
		})

		Describe("DeleteCmdSet()", func() {
			BeforeEach(func() {
				dataStore.SaveCmdSet("sample_set2", []string{"cmda", "cmdb"})
				dataStore.SaveCmdSet("sample_set3", []string{"cmdI", "cmdII"})
			})

			It("removes the recorded command set from data file", func() {
				dataStore.DeleteCmdSet("sample_set2")

				f, err := os.Open(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins", "CmdSetRecords"))
				Ω(err).ToNot(HaveOccurred())

				contents, err := ioutil.ReadAll(f)
				Ω(err).ToNot(HaveOccurred())

				Ω(string(contents)).To(ContainSubstring("--sample_set"))
				Ω(string(contents)).To(ContainSubstring("  cmd1"))
				Ω(string(contents)).To(ContainSubstring("  cmd2"))
				Ω(string(contents)).ToNot(ContainSubstring("--sample_set2"))
				Ω(string(contents)).ToNot(ContainSubstring("  cmda"))
				Ω(string(contents)).ToNot(ContainSubstring("  cmdb"))
				Ω(string(contents)).To(ContainSubstring("--sample_set3"))
				Ω(string(contents)).To(ContainSubstring("  cmdI"))
				Ω(string(contents)).To(ContainSubstring("  cmdII"))
			})

		})
	})
})
