package replay_test

import (
	"github.com/cloudfoundry/cli/plugin/fakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"

	"github.com/simonleung8/cli-plugin-recorder/data/data_fakes"
	"github.com/simonleung8/cli-plugin-recorder/replay"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Replay", func() {
	var (
		fakeData          *data_fakes.FakeCmdSetData
		fakeCliConnection *fakes.FakeCliConnection
		cmds              replay.ReplayCmds
	)

	BeforeEach(func() {
		fakeCliConnection = &fakes.FakeCliConnection{}
		fakeData = &data_fakes.FakeCmdSetData{}

		fakeData.GetCmdSetStub = func(cmd string) []string {
			if cmd == "set1" {
				return []string{"cmd1 arg1", "cmd2"}
			} else if cmd == "set2" {
				return []string{"cmda arg1 arg2", "cmdb arg1"}
			}
			return []string{}
		}

	})

	Describe("Run()", func() {
		It("retreive the command set and pass all commands to CLI for execution", func() {
			cmds = replay.NewReplayCmds(fakeCliConnection, fakeData, "set1")
			cmds.Run()

			Ω(fakeCliConnection.CliCommandCallCount()).To(Equal(2))

			args1 := fakeCliConnection.CliCommandArgsForCall(0)
			args2 := fakeCliConnection.CliCommandArgsForCall(1)

			Ω(args1).To(Equal([]string{"cmd1", "arg1"}))
			Ω(args2).To(Equal([]string{"cmd2"}))
		})

		It("accept multiple command sets", func() {
			cmds = replay.NewReplayCmds(fakeCliConnection, fakeData, "set1", "set2")
			cmds.Run()

			Ω(fakeCliConnection.CliCommandCallCount()).To(Equal(4))

			args1 := fakeCliConnection.CliCommandArgsForCall(0)
			args2 := fakeCliConnection.CliCommandArgsForCall(1)
			args3 := fakeCliConnection.CliCommandArgsForCall(2)
			args4 := fakeCliConnection.CliCommandArgsForCall(3)

			Ω(args1).To(Equal([]string{"cmd1", "arg1"}))
			Ω(args2).To(Equal([]string{"cmd2"}))
			Ω(args3).To(Equal([]string{"cmda", "arg1", "arg2"}))
			Ω(args4).To(Equal([]string{"cmdb", "arg1"}))
		})

		It("reports when a command set is not found", func() {
			cmds = replay.NewReplayCmds(fakeCliConnection, fakeData, "bad-set", "set1")
			output := io_helpers.CaptureOutput(func() {
				cmds.Run()
			})

			Ω(output[0]).To(ContainSubstring("Command set bad-set not found"))

			Ω(fakeCliConnection.CliCommandCallCount()).To(Equal(2))

			args1 := fakeCliConnection.CliCommandArgsForCall(0)
			args2 := fakeCliConnection.CliCommandArgsForCall(1)

			Ω(args1).To(Equal([]string{"cmd1", "arg1"}))
			Ω(args2).To(Equal([]string{"cmd2"}))
		})

	})

})
