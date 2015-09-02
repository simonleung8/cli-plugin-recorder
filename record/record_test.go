package record_test

import (
	"os"
	"time"

	"github.com/cloudfoundry/cli/plugin/fakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"

	"github.com/simonleung8/cli-plugin-recorder/data/data_fakes"

	"github.com/simonleung8/cli-plugin-recorder/record"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Record Package", func() {
	var (
		fakeData          *data_fakes.FakeCmdSetData
		fakeCliConnection *fakes.FakeCliConnection
		recorder          record.RecordCmd
		rFile             *os.File
		wFile             *os.File
		err               error
	)

	BeforeEach(func() {
		fakeCliConnection = &fakes.FakeCliConnection{}
		fakeData = &data_fakes.FakeCmdSetData{}
		rFile, wFile, err = os.Pipe()
		Ω(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		rFile.Close()
		wFile.Close()
	})

	Describe("Record()", func() {
		It("warns user about an existing command set name", func() {
			fakeData.IsCmdExistReturns(true)

			recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)
			output := io_helpers.CaptureOutput(func() {
				recorder.Record("cmdSet-Name1")
			})

			Ω(output[1]).Should(ContainSubstring("The name cmdSet-Name1 already exists"))
		})

		It("uses the data store to save the recorded command set", func() {
			recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)
			go io_helpers.CaptureOutput(func() {
				recorder.Record("cmdSet-Name2")
			})

			wFile.WriteString("cf api blah\n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("cf target\n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("stop\n")
			time.Sleep(100 * time.Millisecond)
			Ω(fakeData.SaveCmdSetCallCount()).To(Equal(1))

			cmdSetName, cmds := fakeData.SaveCmdSetArgsForCall(0)
			Ω(cmdSetName).To(Equal("cmdSet-Name2"))
			Ω(cmds).To(Equal([]string{"api blah", "target"}))
		})

		It("trims spaces for command set name and commands before saving", func() {
			recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)
			go io_helpers.CaptureOutput(func() {
				recorder.Record("  cmdSet-Name2  ")
			})

			wFile.WriteString("   cf api blah    \n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("stop\n")
			time.Sleep(100 * time.Millisecond)
			Ω(fakeData.SaveCmdSetCallCount()).To(Equal(1))

			cmdSetName, cmds := fakeData.SaveCmdSetArgsForCall(0)
			Ω(cmdSetName).To(Equal("cmdSet-Name2"))
			Ω(cmds).To(Equal([]string{"api blah"}))
		})

		It("abandon recording when user input 'quit'", func() {
			recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)
			go io_helpers.CaptureOutput(func() {
				recorder.Record("cmdSet-Name2")
			})

			wFile.WriteString("cf api blah\n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("quit\n")
			time.Sleep(100 * time.Millisecond)
			Ω(fakeCliConnection.CliCommandCallCount()).To(Equal(1))
			Ω(fakeData.SaveCmdSetCallCount()).To(Equal(0))
		})

		It("does not invoke invalid command without 'cf' prefix", func() {
			recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)
			go io_helpers.CaptureOutput(func() {
				recorder.Record("cmdSet-Name2")
			})

			wFile.WriteString("cf api blah\n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("bad command\n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("stop\n")
			time.Sleep(100 * time.Millisecond)
			Ω(fakeCliConnection.CliCommandCallCount()).To(Equal(1))
		})

		It("calls CLI to execute the cf command", func() {
			recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)
			go io_helpers.CaptureOutput(func() {
				recorder.Record("cmdSet-Name2")
			})

			wFile.WriteString("cf api abc\n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("cf target\n")
			time.Sleep(100 * time.Millisecond)
			wFile.WriteString("stop\n")
			time.Sleep(100 * time.Millisecond)
			Ω(fakeCliConnection.CliCommandCallCount()).To(Equal(2))
			Ω(fakeCliConnection.CliCommandArgsForCall(0)).To(Equal([]string{"api", "abc"}))
			Ω(fakeCliConnection.CliCommandArgsForCall(1)).To(Equal([]string{"target"}))
		})
	})

	Describe("ListCmdSets()", func() {
		It("prints all the available command sets", func() {
			fakeData.ListCmdSetNamesReturns([]string{"cmdSet1", "cmdSet2", "cmdSet3"})
			recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)

			output := io_helpers.CaptureOutput(func() {
				recorder.ListCmdSets()
			})

			Ω(output[0]).Should(ContainSubstring("following command sets are available"))
			Ω(output[1]).Should(ContainSubstring("============"))
			Ω(output[2]).Should(ContainSubstring("cmdSet1"))
			Ω(output[3]).Should(ContainSubstring("cmdSet2"))
			Ω(output[4]).Should(ContainSubstring("cmdSet3"))
		})

		Describe("ListCmds()", func() {
			It("lists the commands within a command set", func() {
				fakeData.GetCmdSetReturns([]string{"cmd1", "cmd2", "cmd3"})

				recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)

				output := io_helpers.CaptureOutput(func() {
					recorder.ListCmds("cmdSet1")
				})

				Ω(output[0]).Should(ContainSubstring("The following commands are in cmdSet1"))
				Ω(output[1]).Should(ContainSubstring("=========="))
				Ω(output[2]).Should(ContainSubstring("cmd1"))
				Ω(output[3]).Should(ContainSubstring("cmd2"))
				Ω(output[4]).Should(ContainSubstring("cmd3"))
			})
		})

		Describe("DeleteCmdSet()", func() {
			It("calls the data store to delete a command set", func() {
				recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)
				recorder.DeleteCmdSet("delete-set")

				Ω(fakeData.DeleteCmdSetCallCount()).To(Equal(1))
				Ω(fakeData.DeleteCmdSetArgsForCall(0)).To(Equal("delete-set"))
			})
		})

		Describe("ClearCmdSet()", func() {
			It("calls the data store to clear all recorded command set", func() {
				recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)

				go io_helpers.CaptureOutput(func() {
					recorder.ClearCmdSets()
				})

				wFile.WriteString("y\n")
				time.Sleep(100 * time.Millisecond)

				Ω(fakeData.ClearCmdSetsCallCount()).To(Equal(1))
			})

			It("does not clear command sets without user confirmation", func() {
				recorder = record.NewRecordCmd(fakeCliConnection, fakeData, rFile)

				go io_helpers.CaptureOutput(func() {
					recorder.ClearCmdSets()
				})

				wFile.WriteString("n\n")
				time.Sleep(100 * time.Millisecond)

				Ω(fakeData.ClearCmdSetsCallCount()).To(Equal(0))
			})
		})
	})

})
