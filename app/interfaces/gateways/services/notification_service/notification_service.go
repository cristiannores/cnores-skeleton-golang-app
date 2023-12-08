package notification_service

import (
	"context"
	"fmt"

	http_client "cnores-skeleton-golang-app/app/infrastructure/services/http"
	"cnores-skeleton-golang-app/app/interfaces/constant"
	notification_service_request "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/request"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
)

type NotificationServiceInterface interface {
	Notify(ctx context.Context, data *notification_service_request.SlackBody)
}

type NotificationService struct {
	httpClient http_client.HttpClientInterface[notification_service_request.SlackBody, string]
	token      string
}

func NewNotificationService(httpClient http_client.HttpClientInterface[notification_service_request.SlackBody, string], token string) NotificationServiceInterface {
	return &NotificationService{httpClient, token}
}

func (n *NotificationService) Notify(ctx context.Context, data *notification_service_request.SlackBody) {
	log := utils_context.GetLogFromContext(ctx, constant.InterfaceLayer, "notification_service.Notify")
	url := fmt.Sprintf("/services/%s", n.token)
	log.Info("calling to %s ", url)
	log.Info("text to %s ", data.Text)
	res, err := n.httpClient.Post(ctx, url, *data)

	if err == nil {
		log.Info("response %v", res)
	} else {
		log.Info("error  %s", err.Error())
	}

}
