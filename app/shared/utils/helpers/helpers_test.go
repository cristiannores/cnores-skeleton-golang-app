package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"cnores-skeleton-golang-app/app/shared/utils/helpers"
	"time"
)

var _ = Describe("Tests for the helpers", Ordered, func() {
	var helper helpers.HelperInterface
	BeforeAll(func() {
		helper = helpers.NewHelpers()
	})
	type Child struct {
		BoolValue bool `validate:"required"`
	}
	type BasicEntity struct {
		TextValue    string `validate:"required"`
		NumericValue int
		DateValue    time.Time `validate:"required"`
		Child        Child     `validate:"required"`
	}

	When("A entity with all required fields is passed to the ValidateStruct method", func() {
		It("Should not return an error", func() {

			BaseEntity := BasicEntity{
				TextValue: "Text",
				DateValue: time.Date(2022, 12, 01, 10, 30, 0, 0, time.UTC),
				Child:     Child{BoolValue: true},
			}
			err := helper.ValidateStruct(BaseEntity)
			Expect(err).To(BeNil())
		})
	})

	When("A entity without a required value is passed to the ValidateStruct method", func() {
		It("Should return an error", func() {
			const requiredError = "Key: 'BasicEntity.DateValue' Error:Field validation for 'DateValue' failed on the 'required' tag"
			BaseEntity := BasicEntity{
				TextValue:    "Text",
				NumericValue: 33,
				Child:        Child{BoolValue: true},
			}
			err := helper.ValidateStruct(BaseEntity)
			Expect(err.Error()).To(Equal(requiredError))
		})
	})
})
