package remote

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// NewPrometheusMiddleware returns Middleware for Prometheus metrics
func NewPrometheusMiddleware() Middleware {
	return func(next RemoteCaller) RemoteCaller {

		var totalRequests = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rpc_requests_total",
				Help: "Number of RPC requests.",
			},
			[]string{"method"},
		)

		var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name: "rpc_response_time_seconds",
			Help: "Duration of RPC requests to hidden services.",
		}, []string{"method"})

		return CallerFunc(func(req Request) []byte {
			timer := prometheus.NewTimer(httpDuration.With(prometheus.Labels{"method": req.Method}))
			totalRequests.With(prometheus.Labels{"method": req.Method}).Inc()

			resp := next.Call(req)

			timer.ObserveDuration()
			return resp
		})
	}
}
