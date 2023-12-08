package error_handler_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"cnores-skeleton-golang-app/app/application/domain_error"
	"cnores-skeleton-golang-app/app/infrastructure/infrastructure_errors"
	mock_metrics "cnores-skeleton-golang-app/app/infrastructure/metrics/mock"
	mock_notification_service "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/mock"
	notification_service_request "cnores-skeleton-golang-app/app/interfaces/gateways/services/notification_service/request"
	"cnores-skeleton-golang-app/app/shared/utils/common"
	"cnores-skeleton-golang-app/app/shared/utils/config"
	"cnores-skeleton-golang-app/app/shared/utils/error_handler"
)

var _ = Describe("Error handler", Ordered, func() {
	configMock := config.Config{
		CurrentStage: "",
		Developers:   nil,
		Url:          "",
		Port:         "",
		Services:     nil,
		Notifier: config.Notifier{
			Url:          "",
			Timeout:      0,
			Headers:      nil,
			Token:        "",
			AlertChannel: "test-chan",
			Channels: map[string]string{
				"test-chan": "test-chan",
			},
		},
	}
	When("is called with a notifiable error", func() {
		When("the error can't be notified (CanNotify = false)", func() {
			It("shouldn't call Notify from service", func() {
				ctrl := gomock.NewController(GinkgoT())
				notificationSrvc := mock_notification_service.NewMockNotificationServiceInterface(ctrl)

				metricsClient := mock_metrics.NewMockMetricInterface(ctrl)

				notificationSrvc.EXPECT().Notify(gomock.Any(), gomock.Any()).Times(0)
				metricsClient.EXPECT().IncrementErrorMetric(gomock.Any(), gomock.Any()).Times(1)

				infraError := infrastructure_errors.New(
					map[string]interface{}{
						"ShippingGroupId": "SG",
					},
					"test infra error",
					common.ErrorInformation{Notify: false},
				)

				errorHandler := error_handler.NewErrorHandler(notificationSrvc, configMock)

				errorHandler.HandleAndNotify(context.TODO(), infraError, metricsClient)
			})
		})

		When("the error can be notified (CanNotify = true)", func() {
			It("should call notify service onto salesChannel", func() {
				ctrl := gomock.NewController(GinkgoT())
				notificationSrvc := mock_notification_service.NewMockNotificationServiceInterface(ctrl)

				metricsClient := mock_metrics.NewMockMetricInterface(ctrl)

				sg := "SG"
				source := "Testing"
				store := "99"
				errInfo := common.ErrorInformation{
					Notify:       true,
					SalesChannel: "test-chan",
					Name:         "test-error",
				}

				msgError := "test infra error"

				infraError := infrastructure_errors.New(
					map[string]interface{}{
						"ShippingGroupId": sg,
						"Store":           store,
						"Source":          source,
					},
					msgError,
					errInfo,
				)

				//Error processing shippingGroupID: %s, store: %s Error Type: %s, Error Message: "%s", IN: %s
				expectedMessage := fmt.Sprintf(`Error processing shippingGroupID: %s, store: %s Error Type: %s, Error Message: "%s", IN: %s`, sg, store, errInfo.Name, msgError, source)

				ctx := context.TODO()
				ctx = context.WithValue(ctx, "fieldsToLog", map[string]string{
					"Source":          source,
					"ShippingGroupId": sg,
					"Store":           store,
				})

				notificationSrvc.EXPECT().Notify(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, slackBody *notification_service_request.SlackBody) {
					Expect(slackBody.Text).To(Equal(expectedMessage))
					Expect(slackBody.Channel).To(Equal("test-chan"))
				}).Times(1)

				metricsClient.EXPECT().IncrementErrorMetric(gomock.Any(), gomock.Any()).Times(1)

				errorHandler := error_handler.NewErrorHandler(notificationSrvc, configMock)
				errorHandler.HandleAndNotify(ctx, infraError, metricsClient)
			})

			It("should call notify service onto general channel", func() {
				ctrl := gomock.NewController(GinkgoT())
				notificationSrvc := mock_notification_service.NewMockNotificationServiceInterface(ctrl)

				metricsClient := mock_metrics.NewMockMetricInterface(ctrl)

				sg := "SG"
				source := "Testing"
				store := "99"

				errInfo := common.ErrorInformation{
					Notify: true,
					Name:   "test-error",
				}

				msgError := "test infra error"

				domError := domain_error.New(
					msgError,
					errInfo,
					map[string]interface{}{
						"ShippingGroupId": sg,
						"Store":           store,
						"Source":          source,
					},
				)

				expectedMessage := fmt.Sprintf(`Error processing shippingGroupID: %s, store: %s Error Type: %s, Error Message: "%s", IN: %s`, sg, store, errInfo.Name, msgError, source)

				ctx := context.TODO()
				ctx = context.WithValue(ctx, "fieldsToLog", map[string]string{
					"ShippingGroupId": sg,
					"Store":           store,
					"Source":          source,
				})

				notificationSrvc.EXPECT().Notify(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, slackBody *notification_service_request.SlackBody) {
					Expect(slackBody.Text).To(Equal(expectedMessage))
					Expect(slackBody.Channel).To(Equal("test-chan"))
				}).Times(1)

				metricsClient.EXPECT().IncrementErrorMetric(gomock.Any(), gomock.Any()).Times(1)

				errorHandler := error_handler.NewErrorHandler(notificationSrvc, configMock)
				errorHandler.HandleAndNotify(ctx, domError, metricsClient)
			})

			It("should call notify service onto general channel with no Source and Sg provided in context", func() {
				ctrl := gomock.NewController(GinkgoT())
				notificationSrvc := mock_notification_service.NewMockNotificationServiceInterface(ctrl)

				metricsClient := mock_metrics.NewMockMetricInterface(ctrl)

				sg := "SG"
				source := "Testing"
				store := "99"
				errInfo := common.ErrorInformation{
					Notify: true,
					Name:   "test-error",
				}

				msgError := "test infra error"

				domError := domain_error.New(
					msgError,
					errInfo,
				)

				ctx := context.TODO()
				ctx = context.WithValue(ctx, "fieldsToLog", map[string]string{
					"ShippingGroupId": sg,
					"Store":           store,
					"Source":          source,
				})

				expectedMessage := fmt.Sprintf(`Error processing shippingGroupID: %s, store: %s Error Type: %s, Error Message: "%s", IN: %s`, sg, store, errInfo.Name, msgError, source)

				notificationSrvc.EXPECT().Notify(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, slackBody *notification_service_request.SlackBody) {
					Expect(slackBody.Text).To(Equal(expectedMessage))
					Expect(slackBody.Channel).To(Equal("test-chan"))
				}).Times(1)

				metricsClient.EXPECT().IncrementErrorMetric(gomock.Any(), gomock.Any()).Times(1)

				errorHandler := error_handler.NewErrorHandler(notificationSrvc, configMock)
				errorHandler.HandleAndNotify(ctx, domError, metricsClient)
			})
		})
	})
})
