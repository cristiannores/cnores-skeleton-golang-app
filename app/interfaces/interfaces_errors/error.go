package interfaces_errors

type Error struct {
	DataDetails map[string]interface{}
	Kind        ErrorTypes
	Msg         string
	InnerError  error
}

type ErrorTypes string

const (
	UnexpectedError ErrorTypes = "UNEXPECTED_ERROR"
)

func New(data map[string]interface{}, Msg string, Kind ErrorTypes) *Error {
	return &Error{DataDetails: data, Msg: Msg, Kind: Kind}
}

func NewWithError(data map[string]interface{}, Msg string, Kind ErrorTypes, InnerError error) *Error {
	return &Error{DataDetails: data, Msg: Msg, Kind: Kind, InnerError: InnerError}
}

func (e *Error) CanNotify() bool {
	return true
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
	return string(e.Kind)
}

func (e *Error) GetChannel() string {
	return ""
}

func (e *Error) Error() string {
	return e.Msg
}
