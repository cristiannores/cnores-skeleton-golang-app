package common_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Common Suite")
}
