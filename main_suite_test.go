package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cli/testhelpers/plugin_builder"

	"testing"
)

func TestCliPluginRecorder(t *testing.T) {
	RegisterFailHandler(Fail)

	plugin_builder.BuildTestBinary("", "main")

	RunSpecs(t, "CliPluginRecorder Suite")
}
