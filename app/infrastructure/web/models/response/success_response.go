package response

type SuccessResponse[Data any] struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data,omitempty"`
}
