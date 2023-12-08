package helpers

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
	"cnores-skeleton-golang-app/app/shared/utils/custom_validators"
)

type Helpers struct {
}

type HelperInterface interface {
	GetFileName() string
	ValidateStruct(o interface{}) error
	JsonParser(data []byte, v interface{}) (err error)
}

func NewHelpers() HelperInterface {
	return &Helpers{}
}

func (i *Helpers) JsonParser(data []byte, v interface{}) (err error) {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
func (i *Helpers) GetFileName() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		err := errors.New("failed to get filename")
		panic(err)
	}
	filename := filepath.Base(filePath)
	filename = strings.Replace(filename, ".go", "", 1)
	return filename
}

func (i *Helpers) ValidateStruct(entity interface{}) error {
	var validate *validator.Validate
	validate = validator.New()
	validate.RegisterValidation("gtezero", custom_validators.GreaterThanOrEqualZero)
	errValidate := validate.Struct(entity)
	if errValidate != nil {
		return errValidate
	}
	return nil
}
