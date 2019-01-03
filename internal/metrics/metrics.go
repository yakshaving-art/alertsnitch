package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab.com/yakshaving.art/alertsnitch/version"
)

var (
	namespace = "alertsnitch"

	bootTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "boot_time_seconds",
		Help:      "unix timestamp of when the service was started",
	})

	buildInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "build_info",
		Help:      "Build information",
	}, []string{"version", "commit", "date"})
)

// Exported metrics
var (
	AlertsReceivedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "alerts",
		Name:      "received_total",
		Help:      "total number of valid alerts received",
	}, []string{"receiver", "status"})
	WebhooksReceivedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "webhooks",
		Name:      "received_total",
		Help:      "total number of webhooks posts received",
	})
	InvalidWebhooksTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "webhooks",
		Name:      "invalid_total",
		Help:      "total number of invalid webhooks received",
	})

	DatabaseUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "up",
		Help:      "wether the database is accessible or not",
	})

	AlertsSavedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "alerts",
		Name:      "saved_total",
		Help:      "total number of alerts saved in the database",
	}, []string{"receiver", "status"})
	AlertsSavingFailuresTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "alerts",
		Name:      "saving_failures_total",
		Help:      "total number of alerts that failed to be saved in the database",
	}, []string{"receiver", "status"})
)

func init() {
	bootTime.Set(float64(time.Now().Unix()))

	buildInfo.WithLabelValues(version.Version, version.Commit, version.Date).Set(1)

	prometheus.MustRegister(bootTime)
	prometheus.MustRegister(buildInfo)
	prometheus.MustRegister(DatabaseUp)

	prometheus.MustRegister(AlertsReceivedTotal)
	prometheus.MustRegister(AlertsSavedTotal)
	prometheus.MustRegister(AlertsSavingFailuresTotal)

	prometheus.MustRegister(WebhooksReceivedTotal)
	prometheus.MustRegister(InvalidWebhooksTotal)

}
