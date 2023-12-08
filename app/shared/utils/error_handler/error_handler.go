package error_handler

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"cnores-skeleton-golang-app/app/application/domain_error"
	"cnores-skeleton-golang-app/app/infrastructure/infrastructure_errors"
	"cnores-skeleton-golang-app/app/infrastructure/metrics"
	"cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service"
	notification_service_request "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/request"
	"cnores-skeleton-golang-app/app/shared/utils/config"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"strings"
)

type NotifiableError interface {
	CanNotify() bool
	GetErrorType() string
	Error() string
}

type ErrorHandlerInterface interface {
	Handle(ctx context.Context, e error, metricClient metrics.MetricInterface)
	HandleAndNotify(ctx context.Context, e NotifiableError, metricClient metrics.MetricInterface)
}

type errorHandler struct {
	config              config.Config
	notificationService notification_service.NotificationServiceInterface
}

func NewErrorHandler(notificationService notification_service.NotificationServiceInterface, config config.Config) ErrorHandlerInterface {
	return &errorHandler{notificationService: notificationService, config: config}
}

func (eh *errorHandler) Handle(ctx context.Context, e error, metricClient metrics.MetricInterface) {
	log := utils_context.GetLogFromContext(ctx, "infrastructure", "error_handler.Handle")
	log.Info("handling message error")

	errorName := GetErrorName(e)

	log.Info("Increment metric on error handler")
	metricClient.IncrementErrorMetric(ctx, errorName)
	log.Info("Increment metric success on error handler")

	log.Info("error handled: %s", e.Error())
}

func (eh *errorHandler) HandleAndNotify(ctx context.Context, e NotifiableError, metricClient metrics.MetricInterface) {
	eh.notify(ctx, e)
	eh.Handle(ctx, e, metricClient)
}

func (eh *errorHandler) notify(ctx context.Context, e NotifiableError) {
	log := utils_context.GetLogFromContext(ctx, "infrastructure", "error_handler.notify")
	log.Info("handle and notify")

	msg := fmt.Sprintf(
		`[Error handler message::%s] Error Type: %s, Error Message: "%s"`,
		strings.ToUpper(eh.config.CurrentStage),
		e.GetErrorType(),
		e.Error(),
	)

	log.Info("Message built to be notified: \"%s\"", msg)

	if e.CanNotify() {
		eh.notificationService.Notify(
			ctx,
			&notification_service_request.SlackBody{
				Text: msg,
			})
		log.Info("Error message notified")
		return
	}

	log.Info("Error is not notifiable")
}

func GetErrorName(e error) string {
	var infrastructureError *infrastructure_errors.Error
	var domainError *domain_error.Error

	switch {
	case errors.As(e, &infrastructureError):
		return infrastructureError.Info.Name
	case errors.As(e, &domainError):
		return domainError.Info.Name

	default:
		return infrastructure_errors.UnexpectedError.Name
	}
}

func GetErrorCode(e error) string {
	var infrastructureError *infrastructure_errors.Error
	var domainError *domain_error.Error

	switch {
	case errors.As(e, &infrastructureError):
		return infrastructureError.Info.Code
	case errors.As(e, &domainError):
		return domainError.Info.Code

	default:
		return infrastructure_errors.UnexpectedError.Code
	}
}
