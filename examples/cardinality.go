package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	requestDuration *prometheus.HistogramVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "example",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
			Buckets:   []float64{0.1, 0.2, 0.5, 1, 2, 5, 10},
		}, []string{"path", "method"}),
	}
	reg.MustRegister(m.requestDuration)
	return m
}

func (m *metrics) HelloWorld(resp http.ResponseWriter, req *http.Request) {
	timer := prometheus.NewTimer(m.requestDuration.WithLabelValues(req.URL.Path, req.Method))
	defer timer.ObserveDuration()
	resp.Write([]byte("Hello, world!"))
}

func main() {
	reg := prometheus.NewRegistry()
	m := newMetrics(reg)
	http.HandleFunc("/", m.HelloWorld)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", nil)
}
