package use_case_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUseCase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UseCase Suite")
}
