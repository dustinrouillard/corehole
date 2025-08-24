package corehole

import (
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	blockedRequestsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "corehole",
		Name:      "blocked_requests_total",
		Help:      "Counter of the total DNS requests that were blocked by corehole",
	}, []string{"server", "domain"})
)

func reportBlockedRequest(server, domain string) {
	blockedRequestsCount.WithLabelValues(server, domain).Inc()
}
