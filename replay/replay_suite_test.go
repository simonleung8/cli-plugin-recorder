package replay_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestReplay(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Replay Suite")
}
