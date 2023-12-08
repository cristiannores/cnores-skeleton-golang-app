package dependencies

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"cnores-skeleton-golang-app/app/infrastructure/metrics"
	"cnores-skeleton-golang-app/app/infrastructure/mongo_client"
	http_client "cnores-skeleton-golang-app/app/infrastructure/services/http"
	"cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service"
	notification_service_request "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/request"
	http_clients_faker "cnores-skeleton-golang-app/app/shared/faker/http_clients"
	"cnores-skeleton-golang-app/app/shared/utils/config"
	"cnores-skeleton-golang-app/app/shared/utils/error_handler"
)

type DependencyContainer struct {
	MongoClient  mongo_client.MongoClientInterface
	ErrorHandler error_handler.ErrorHandlerInterface
	Config       config.Config
	MetricClient metrics.MetricInterface
	Echo         *echo.Echo
}

func NewDependencyContainer(config config.Config, echo *echo.Echo) *DependencyContainer {

	mongoClient, _ := mongo_client.NewClient(config.Database.Url)
	errConnecting := mongoClient.Connect()

	if errConnecting != nil {
		log.Fatal("[Dependencies]: Error initializing. database %s", errConnecting.Error())
	}
	log.Info("mongo initializated")

	//setup metrics
	metricClient := metrics.NewMetric()
	metricClient.InitMetrics()

	// Setup Notifier
	notifierClient := buildNotifierHttpClient(config)
	// Setup Notification Services
	notificationService := notification_service.NewNotificationService(
		notifierClient,
		config.Notifier.Token,
	)
	errorHandler := error_handler.NewErrorHandler(notificationService, config)

	return &DependencyContainer{
		MongoClient:  mongoClient,
		ErrorHandler: errorHandler,
		Config:       config,
		MetricClient: metricClient,
		Echo:         echo,
	}
}

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
