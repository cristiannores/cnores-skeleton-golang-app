package handlers

import "github.com/labstack/echo/v4"

type FromKafkaInterface interface {
	FromKafka(message []byte) error
}

type FromApiInterface interface {
	FromApi(ctx echo.Context) error
}
 
type IGetServiceByIdRestHandler interface {
	Handle(serverContext echo.Context) error
}
