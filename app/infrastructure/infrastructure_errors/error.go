package infrastructure_errors

import (
	"cnores-skeleton-golang-app/app/shared/utils/common"
)

type Error struct {
	DataDetails map[string]interface{}
	Info        common.ErrorInformation
	Msg         string
	InnerError  error
}

var (
	DatabaseException  = common.ErrorInformation{Name: "DATABASE_EXCEPTION", Notify: true}
	UnexpectedError    = common.ErrorInformation{Name: "UNEXPECTED_ERROR", Notify: false}
	InputValidation    = common.ErrorInformation{Name: "INPUT_VALIDATION", Notify: true}
	UserNotFound       = common.ErrorInformation{Name: "USER_NOT_FOUND", Notify: true}
	SpecialityNotFound = common.ErrorInformation{Name: "SPECIALITY_NOT_FOUND", Notify: true}
	ServiceNotFound    = common.ErrorInformation{Name: "SERVICE_NOT_FOUND", Notify: true}
)

func New(data map[string]interface{}, Msg string, errorInformation common.ErrorInformation) *Error {
	return &Error{DataDetails: data, Msg: Msg, Info: errorInformation}
}

func NewWithError(data map[string]interface{}, Msg string, errorInformation common.ErrorInformation, InnerError error) *Error {
	return &Error{DataDetails: data, Msg: Msg, Info: errorInformation, InnerError: InnerError}
}

func (e *Error) CanNotify() bool {
	return e.Info.Notify
}

func (e *Error) GetShippingGroupId() string {
	sgInterface, ok := e.DataDetails["ShippingGroupId"]
	if ok {
		if sg, ok := sgInterface.(string); ok {
			return sg
		}
	}
	return ""
}

func (e *Error) GetErrorType() string {
	return e.Info.Name
}

func (e *Error) GetChannel() string {
	return e.Info.SalesChannel
}

func (e *Error) Error() string {
	return e.Msg
}
