package error_handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestErrorHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ErrorHandler Suite")
}
