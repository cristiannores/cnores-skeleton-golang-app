package controllers

import (
	"context"
	"errors"

	"cnores-skeleton-golang-app/app/infrastructure/metrics"
	"cnores-skeleton-golang-app/app/shared/utils/error_handler"
)

func HandleNotifiableError(ctx context.Context, err error, errorHandler error_handler.ErrorHandlerInterface, metricsClient metrics.MetricInterface) {
	var notifiableError error_handler.NotifiableError
	if ok := errors.As(err, &notifiableError); ok {
		errorHandler.HandleAndNotify(ctx, notifiableError, metricsClient)
	} else {
		errorHandler.Handle(ctx, err, metricsClient)
	}
}
