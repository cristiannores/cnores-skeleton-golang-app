package domain_error_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"cnores-skeleton-golang-app/app/application/domain_error"
)

func TestErrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Domain Errors Suite")
}

var _ = Describe("GetErrorInformation", func() {
	When("function receives a normal error", func() {
		It("should return a domain's UnexpectedError", func() {
			normalError := errors.New("normal error")
			err := domain_error.GetErrorInformation(normalError)
			Expect(err.Name).To(Equal("UNEXPECTED_ERROR"))
		})
	})

	When("function receives a domain error", func() {
		It("should map it correctly and return its information", func() {
			domErr := domain_error.New("dom error", domain_error.StructValidation)
			err := domain_error.GetErrorInformation(domErr)
			Expect(err.Name).To(Equal("STRUCT_VALIDATION"))
		})
	})
})
