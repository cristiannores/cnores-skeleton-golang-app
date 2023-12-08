package metrics

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"cnores-skeleton-golang-app/app/infrastructure/constant"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
)

type MetricInterface interface {
	InitMetrics()
	GetMetricByName(ctx context.Context, name MetricNames) *prometheus.CounterVec
	IncrementErrorMetric(ctx context.Context, labelValue string)
	ObserveTimeElapsedApiMetric(ctx context.Context, labels map[string]string, secondsToAdd float64)
	IncrementMetric(ctx context.Context, labels map[string]string, metricName MetricNames)
}

type MetricClient struct {
	metrics          []MetricPrometheus[*prometheus.CounterVec]
	histogramMetrics []MetricPrometheus[*prometheus.HistogramVec]
}

func NewMetric() MetricInterface {
	metrics := MetricsList
	histogramMetrics := HistogramMetricsList
	return &MetricClient{metrics, histogramMetrics}
}

func (m *MetricClient) InitMetrics() {
	log := utils_context.GetLogFromContext(context.Background(), constant.InfrastructureLayer, "metrics_client.InitMetrics")
	for _, metric := range m.metrics {
		err := m.addMetric(metric.Metric)
		if err != nil {
			log.Error("Error initializing metrics %s", metric.Name)
			return
		}
		log.Debug("metric initialized %s", metric.Name)
	}

}

func (m *MetricClient) addMetric(metric *prometheus.CounterVec) error {
	return prometheus.Register(metric)
}

func (m *MetricClient) GetMetricByName(ctx context.Context, name MetricNames) *prometheus.CounterVec {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "metrics_client.GetMetricByName")
	log.Info("Getting metric by name %s", name)
	for _, metric := range m.metrics {
		if metric.Name == name {
			return metric.Metric
		}
	}

	log.Info("Error getting metric by name %s", name)
	return nil
}

func (m *MetricClient) addHistogramMetric(metric *prometheus.HistogramVec) error {
	return prometheus.Register(metric)
}

func (m *MetricClient) addSummaryMetric(metric *prometheus.SummaryVec) error {
	return prometheus.Register(metric)
}

func (m *MetricClient) IncrementErrorMetric(ctx context.Context, labelValue string) {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "metrics_client.IncrementErrorMetric")
	log.Info("Increment metric error label: %s", labelValue)
	m.GetMetricByName(ctx, ERROR).With(prometheus.Labels{
		"kind": labelValue,
	}).Inc()
}

func (m *MetricClient) IncrementMetric(ctx context.Context, labels map[string]string, metricName MetricNames) {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "metrics_client.IncrementReceivedBillMetric")
	log.Info("Increment metric %s label: %s", metricName, labels)
	m.GetMetricByName(ctx, metricName).With(labels).Inc()
}

func (m *MetricClient) GetHistogramMetricByName(ctx context.Context, name MetricNames) *prometheus.HistogramVec {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "metrics_client.GetHistogramMetricByName")
	log.Info("Getting histogram metric by name %s", name)
	for _, metric := range m.histogramMetrics {
		if metric.Name == name {
			return metric.Metric
		}
	}

	log.Info("Error getting histogram metric by name %s", name)
	return nil
}

func (m *MetricClient) ObserveTimeElapsedApiMetric(ctx context.Context, labels map[string]string, secondsToAdd float64) {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "metrics_client.ObserveTimeElapsedApiMetric")
	log.Info(fmt.Sprintf("Increment metric elapsed time bill order labels %#v processed in %f seconds", labels, secondsToAdd))
	m.GetHistogramMetricByName(ctx, TIME_ELAPSED_API).With(labels).Observe(secondsToAdd)
}
