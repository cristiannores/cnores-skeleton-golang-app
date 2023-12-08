package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricType interface {
	*prometheus.SummaryVec | *prometheus.CounterVec | *prometheus.HistogramVec
}

type MetricPrometheus[T MetricType] struct {
	Metric T
	Name   MetricNames
}

const APP = "medical_connect_users"

type MetricNames string

const (
	ERROR MetricNames = "error_metrics"

	TIME_ELAPSED_API MetricNames = "time_elapsed_api_metrics"

	API_REQUEST MetricNames = "api_request_metrics"
	API_FAILED  MetricNames = "api_failed_metrics"
	API_SUCCESS MetricNames = "api_success_metrics"
)

var MetricsList = []MetricPrometheus[*prometheus.CounterVec]{
	{
		Metric: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s", APP, ERROR),
			Help: "Total of times called this error metric",
		}, []string{"kind"}),
		Name: ERROR,
	},

	{
		Metric: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s", APP, API_REQUEST),
			Help: "Total of api request metrics",
		}, []string{"api_name"}),
		Name: API_REQUEST,
	},
	{
		Metric: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s", APP, API_FAILED),
			Help: "Total of api failed metrics",
		}, []string{"api_name"}),
		Name: API_FAILED,
	},
	{
		Metric: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s", APP, API_SUCCESS),
			Help: "Total of api success metrics",
		}, []string{"api_name"}),
		Name: API_SUCCESS,
	},
}

var HistogramMetricsList = []MetricPrometheus[*prometheus.HistogramVec]{
	{
		Metric: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_%s", APP, TIME_ELAPSED_API),
			Help: "Time elapsed bill order shipping group (by api_name) ",
		}, []string{"api_name"}),
		Name: TIME_ELAPSED_API,
	},
}
