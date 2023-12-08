package handlers

import (
	"encoding/json"
	"io/ioutil"
	"cnores-skeleton-golang-app/app/shared/utils/error_handler"
	"net/http"

	"github.com/labstack/echo/v4"
	"cnores-skeleton-golang-app/app/infrastructure/infrastructure_errors"
	"cnores-skeleton-golang-app/app/infrastructure/web/models/response"
)

type Source string

var (
	FromApi   Source = "FromApi"
	FromKafka Source = "FromKafka"
)

type Flow string

var (
	Bill        Flow = "Bill"
	Checkout    Flow = "Checkout"
	Delivered   Flow = "Delivered"
	OnRoute     Flow = "OnRoute"
	Picked      Flow = "Picked"
	Pickup      Flow = "Pickup"
	BillRequest Flow = "BillRequest"
)

type FromApiCommon[T any] struct {
}

func (c *FromApiCommon[T]) SetBadRequest(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
		Status:  "error",
		Code:    error_handler.GetErrorName(err),
		Message: err.Error(),
	})
}

func (c *FromApiCommon[T]) ProcessRequestBody(serverContext echo.Context, request *T) (string, error) {
	defer serverContext.Request().Body.Close()
	raw, err := ioutil.ReadAll(serverContext.Request().Body)
	if err != nil {
		return "", infrastructure_errors.New(map[string]interface{}{}, err.Error(), infrastructure_errors.InputValidation)
	}
	err = json.Unmarshal(raw, request)
	if err != nil {
		return "", infrastructure_errors.New(map[string]interface{}{}, err.Error(), infrastructure_errors.InputValidation)
	}
	return string(raw), nil
}
