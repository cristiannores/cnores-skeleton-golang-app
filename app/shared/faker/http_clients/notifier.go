package http_clients_faker

import (
	"fmt"
	"golang.org/x/net/context"
	"cnores-skeleton-golang-app/app/infrastructure/constant"
	http_client "cnores-skeleton-golang-app/app/infrastructure/services/http"
	notification_service_request "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/request"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"net/http"
)

type FakeNotifierHttpClient[Request notification_service_request.SlackBody, Response string] struct {
	http_client.HttpClient[Request, Response]
}

func NewFakeNotifierHttpClient() http_client.HttpClientInterface[notification_service_request.SlackBody, string] {
	return &FakeNotifierHttpClient[notification_service_request.SlackBody, string]{}
}

func (f *FakeNotifierHttpClient[Request, Response]) Post(ctx context.Context, url string, body notification_service_request.SlackBody) (*http_client.HttpResponse[Response], error) {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "http_clients_faker.Post")
	log.Info("init  fake notifier")
	log.Info(fmt.Sprintf("%s %#v ", url, body))

	if body.Text == "FAKE MESSAGE" {
		errorResponse := FailedNotifier()
		log.Info(fmt.Sprintf("Finish fake sync bill with response %v  ", errorResponse))
		return nil, errorResponse
	} else {
		syncBillResponse := SuccessNotifier[Response](body)
		log.Info(fmt.Sprintf("Finish fake notifier with response %v  ", syncBillResponse))
		return syncBillResponse, nil
	}

}

func FailedNotifier() error {
	httpClientError := &http_client.HttpClientError{
		StatusCodeMessage: http.StatusInternalServerError,
		ErrorMessage:      fmt.Sprintf("Error %s", " [ERROR FROM NOTIFIER FAKER] "),
	}
	return httpClientError
}
func SuccessNotifier[Response string](body notification_service_request.SlackBody) *http_client.HttpResponse[Response] {

	response := fmt.Sprintf("Response  %s ", body.Text)
	return &http_client.HttpResponse[Response]{
		StatusCode: http.StatusOK,
		Body:       Response(response),
	}
}
