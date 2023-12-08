package common_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"cnores-skeleton-golang-app/app/shared/utils/common"
)

var _ = Describe("Context", func() {

	When("An empty traceparent is received", func() {
		It("Should fill the trace parent components with zero and generate span id successfully", func() {

			type Object struct {
				Field1 string `json:"field1"`
				Field2 int    `json:"field2"`
			}

			obj := Object{
				Field1: "2323232323",
				Field2: 11121221,
			}

			obj2 := Object{
				Field1: "2323232323",
				Field2: 11121221,
			}

			hash1, _ := common.HashObject(obj)
			hash2, _ := common.HashObject(obj2)
			Expect(hash1).To(Equal(hash2))
		})
	})

})
