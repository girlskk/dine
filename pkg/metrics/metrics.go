package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// RecoverCounter recover计数
	RecoverCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "recover_counter",
			Help: "统计recover产生情况",
		},
	)
)
