package domain_error

import (
	"errors"

	"cnores-skeleton-golang-app/app/shared/utils/common"
)

type Error struct {
	Data       map[string]interface{}
	Info       common.ErrorInformation
	Msg        string
	InnerError error
}

var (
	UnexpectedError  = common.ErrorInformation{Name: "UNEXPECTED_ERROR", Notify: false}
	StructValidation = common.ErrorInformation{Name: "STRUCT_VALIDATION", Notify: true}
)

func New(Msg string, errorInformation common.ErrorInformation, Data ...map[string]interface{}) *Error {
	if len(Data) > 0 {
		return &Error{Data: Data[0], Msg: Msg, Info: errorInformation}
	}
	return &Error{Msg: Msg, Info: errorInformation}
}

func NewWithError(Data map[string]interface{}, Msg string, errorInformation common.ErrorInformation, InnerError error) *Error {
	return &Error{Data: Data, Msg: Msg, Info: errorInformation, InnerError: InnerError}
}

func (e *Error) CanNotify() bool {
	return e.Info.Notify
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

func GetErrorInformation(e error) common.ErrorInformation {
	var domainError *Error
	switch {
	case errors.As(e, &domainError):
		return domainError.Info
	default:
		return UnexpectedError
	}
}
