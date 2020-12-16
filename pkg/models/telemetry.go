package models

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var (
	TransactionDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "datastore_tx_duration_second",
			Help:       "the duration of the backend datastore by operation",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.05, 0.95: 0.05, 0.99: 0.05},
		},
		[]string{"store", "op"},
	)

	TransactionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "datastore_tx_errors_total",
			Help: "total number of transaction errors by operation",
		},
		[]string{"store", "op"},
	)

	RecordCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "datastore_records_current",
			Help: "number of unique records currently in the datastore",
		},
		[]string{"store"},
	)
)

func GetCounterValue(metric *prometheus.CounterVec, labels ...string) (float64, error) {
	var m = &dto.Metric{}
	if err := metric.WithLabelValues(labels...).Write(m); err != nil {
		return 0, err
	}
	return m.Counter.GetValue(), nil
}

func init() {
	prometheus.MustRegister(TransactionDuration)
	prometheus.MustRegister(TransactionErrors)
	prometheus.MustRegister(RecordCount)
}
