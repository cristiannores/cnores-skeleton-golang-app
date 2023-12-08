package setup_dependencies

import (
	http_client "cnores-skeleton-golang-app/app/infrastructure/services/http"
	notification_service_request "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/request"
	http_clients_faker "cnores-skeleton-golang-app/app/shared/faker/http_clients"
	"cnores-skeleton-golang-app/app/shared/utils/config"
)

func buildNotifierHttpClient(config config.Config) http_client.HttpClientInterface[notification_service_request.SlackBody, string] {
	if config.Notifier.EnableFaker == true {
		return http_clients_faker.NewFakeNotifierHttpClient()
	} else {
		return http_client.NewHttpClient[notification_service_request.SlackBody, string](
			http_client.HttpClientSettings{
				Host:        config.Notifier.Url,
				Headers:     config.Notifier.Headers,
				Timeout:     config.Notifier.Timeout,
				Config:      config,
				ServiceName: "notifier",
				Secure:      true,
			},
		)
	}
}
