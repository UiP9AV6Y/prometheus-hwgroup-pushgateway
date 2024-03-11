package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hwg_pushgateway_http_requests_total",
			Help: "Total HTTP requests processed by the HWg-Push protocol gateway, excluding scrapes.",
		},
		[]string{"handler", "code", "method"},
	)
	httpPushSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "hwg_pushgateway_http_push_size_bytes",
			Help:       "HTTP request size for pushes to the HWg-Push protocol gateway.",
			Objectives: map[float64]float64{0.1: 0.01, 0.5: 0.05, 0.9: 0.01},
		},
		[]string{"method"},
	)
	httpPushDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "hwg_pushgateway_http_push_duration_seconds",
			Help:       "HTTP request duration for pushes to the HWg-Push protocol gateway.",
			Objectives: map[float64]float64{0.1: 0.01, 0.5: 0.05, 0.9: 0.01},
		},
		[]string{"method"},
	)
)

func RegisterMetrics(reg prometheus.Registerer) error {
	if err := reg.Register(httpCnt); err != nil {
		return err
	}

	if err := reg.Register(httpPushSize); err != nil {
		return err
	}

	if err := reg.Register(httpPushDuration); err != nil {
		return err
	}

	return nil
}

func InstrumentWithCounter(handlerName string, handler http.Handler) http.HandlerFunc {
	return promhttp.InstrumentHandlerCounter(
		httpCnt.MustCurryWith(prometheus.Labels{"handler": handlerName}),
		handler,
	)
}

func InstrumentWithSummaries(handler http.Handler) http.HandlerFunc {
	return promhttp.InstrumentHandlerRequestSize(
		httpPushSize, promhttp.InstrumentHandlerDuration(httpPushDuration, handler),
	)
}
