package main_test

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cloudfoundry/cli/testhelpers/rpc_server"
	fake_rpc_handlers "github.com/cloudfoundry/cli/testhelpers/rpc_server/fakes"
	"github.com/simonleung8/cli-plugin-recorder/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("CliPluginRecorder", func() {
	var (
		validPluginPath = "./main.exe"
		CF_PLUGIN_HOME  string
		rpcHandlers     *fake_rpc_handlers.FakeHandlers
		ts              *test_rpc_server.TestServer
	)

	BeforeEach(func() {
		wd, err := os.Getwd()
		Ω(err).ToNot(HaveOccurred())

		CF_PLUGIN_HOME = filepath.Join(wd, "data", "test_fixtures", "tmp")
		os.Setenv("CF_PLUGIN_HOME", CF_PLUGIN_HOME)

		err = os.MkdirAll(filepath.Join(CF_PLUGIN_HOME, ".cf", "plugins"), 0700)
		Ω(err).ToNot(HaveOccurred())

		dataStore := data.NewCmdSetData()
		dataStore.SaveCmdSet("test_set1", []string{"cmd1", "cmd2"})
		dataStore.SaveCmdSet("test_set2", []string{"cmda", "cmdb"})

		rpcHandlers = &fake_rpc_handlers.FakeHandlers{}
		ts, err = test_rpc_server.NewTestRpcServer(rpcHandlers)
		Expect(err).NotTo(HaveOccurred())

		err = ts.Start()
		Expect(err).NotTo(HaveOccurred())

		//set rpc.CallCoreCommand to a successful call
		rpcHandlers.CallCoreCommandStub = func(_ []string, retVal *bool) error {
			*retVal = true
			return nil
		}

		//set rpc.GetOutputAndReset to return empty string; this is used by CliCommand()/CliWithoutTerminalOutput()
		rpcHandlers.GetOutputAndResetStub = func(_ bool, retVal *[]string) error {
			*retVal = []string{"{}"}
			return nil
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(CF_PLUGIN_HOME)
		Ω(err).ToNot(HaveOccurred())

		ts.Stop()
	})

	AfterSuite(func() {
		err := os.RemoveAll("./main.exe")
		Ω(err).ToNot(HaveOccurred())
	})

	Describe("Command Record", func() {
		It("complains about invalid flag", func() {
			args := []string{ts.Port(), "record", "--bad-flag"}
			session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait()
			Ω(session).To(gbytes.Say("Invalid flag: --bad-flag"))
		})

		It("needs at least 1 argument", func() {
			args := []string{ts.Port(), "record"}
			session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait()
			Ω(session).To(gbytes.Say("Provide command set name for recording"))
		})

		It("calls data.Record() to record commands", func() {
			args := []string{ts.Port(), "record", "test_command"}
			session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait()
			Ω(session).To(gbytes.Say("Please start entering CF commands"))
		})

		Context("--list or -l flag", func() {
			It("lists all recorded command sets", func() {
				wd, _ := os.Getwd()
				os.Setenv("CF_PLUGIN_HOME", filepath.Join(wd, "data", "test_fixtures"))

				args := []string{ts.Port(), "record", "--list"}
				session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				session.Wait()
				Eventually(session).Should(gbytes.Say("The following command sets are available"))
				Eventually(session).Should(gbytes.Say("============================="))
				Eventually(session).Should(gbytes.Say("cmdSet1"))
				Eventually(session).Should(gbytes.Say("cmdSet2"))
			})
		})

		Context("-n flag", func() {
			It("lists all commands within a command sets", func() {
				args := []string{ts.Port(), "record", "-n", "test_set1"}
				session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				session.Wait()
				Eventually(session).Should(gbytes.Say("The following commands are in test_set1"))
				Eventually(session).Should(gbytes.Say("cmd1"))
				Eventually(session).Should(gbytes.Say("cmd2"))
			})
		})

		Context("--clear or -c flag", func() {
			It("clears the entire data file by calling data.ClearCmdSet()", func() {
				args := []string{ts.Port(), "record", "--clear"}
				session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				session.Wait()
				Eventually(session).Should(gbytes.Say("WARNING: You are about to delete all recorded command sets"))
			})
		})
	})

	Describe("Command Replay", func() {
		It("needs at least 1 argument", func() {
			args := []string{ts.Port(), "replay"}
			session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait()
			Ω(session).To(gbytes.Say("Provide the recorded command set name to playback"))
		})

		It("calls replay.Run() to replay commands", func() {
			args := []string{ts.Port(), "replay", "bad_command_set"}
			session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait()
			Eventually(session).Should(gbytes.Say("Command set bad_command_set not found"))
		})
	})

	Context("--list or -l flag", func() {
		It("lists all recorded command sets", func() {
			wd, _ := os.Getwd()
			os.Setenv("CF_PLUGIN_HOME", filepath.Join(wd, "data", "test_fixtures"))

			args := []string{ts.Port(), "replay", "--list"}
			session, err := gexec.Start(exec.Command(validPluginPath, args...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			session.Wait()
			Eventually(session).Should(gbytes.Say("The following command sets are available"))
			Eventually(session).Should(gbytes.Say("============================="))
			Eventually(session).Should(gbytes.Say("cmdSet1"))
			Eventually(session).Should(gbytes.Say("cmdSet2"))
		})
	})
})
