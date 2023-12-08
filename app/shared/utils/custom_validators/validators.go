package custom_validators

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"cnores-skeleton-golang-app/app/infrastructure/constant"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"strconv"
	"strings"
)

type CustomValidatorInterface interface {
	Validate(ctx context.Context, value interface{}) error
}

type CustomValidator struct {
	validator   *validator.Validate
	customFuncs map[string]validator.Func
}

func NewCustomValidator() CustomValidatorInterface {
	validate := validator.New()
	// Registrar todas las funciones de validación

	validators := map[string]validator.Func{}
	validators["minlength"] = isMinLengthCharValidator
	validators["datetime"] = isValidDateTimeValidator
	validators["invalues"] = isInValuesValidator
	validators["objectid"] = isValidObjectIDValidator
	validators["validurl"] = isValidUrlValidator
	validators["validphone"] = isValidPhoneValidator
	validators["validemail"] = isValidEmailValidator
	validators["validrut"] = isValidRUTValidator
	validators["gtEqZero"] = greaterThanOrEqualZeroValidator

	customValidator := &CustomValidator{validator: validate, customFuncs: validators}
	customValidator.initValidator()

	return customValidator
}

func (c *CustomValidator) initValidator() {

	for key, custom := range c.customFuncs {
		c.validator.RegisterValidation(key, custom)
	}

}

func (c *CustomValidator) Validate(ctx context.Context, value interface{}) error {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "custom_validators.Validate")
	err := c.validator.Struct(value)
	if err != nil {
		log.Error(fmt.Sprintf("error validating %s", err))
		return err
	}
	return nil
}

func greaterThanOrEqualZeroValidator(fl validator.FieldLevel) bool {
	return GreaterThanOrEqualZero(fl.Field().Int())
}

func isMinLengthCharValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	minLength, err := strconv.Atoi(fl.Param())
	if err != nil {
		// Manejar error si es necesario, o asumir un valor predeterminado para minLength
		return false
	}
	return IsMinLengthChar(value, minLength)

}

func isValidDateTimeValidator(fl validator.FieldLevel) bool {
	stringDate := fl.Field().String()
	return IsValidDateTime(stringDate)
}

func isInValuesValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param() // Obtener los valores permitidos de la etiqueta
	allowedValues := strings.Split(param, ",")
	return IsInValues(value, allowedValues)
}

func isValidObjectIDValidator(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	return IsValidObjectID(id)
}

func isValidUrlValidator(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	return IsValidUrl(url)
}
func isValidPhoneValidator(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	return IsValidPhone(phoneNumber)
}

func isValidEmailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return IsValidEmail(email)
}

func isValidRUTValidator(fl validator.FieldLevel) bool {
	rut := fl.Field().String()
	return IsValidRUT(rut)
}
