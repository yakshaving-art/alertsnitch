package metrics_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"gitlab.com/yakshaving.art/alertsnitch/internal/metrics"
)

func TestMetricsAreRegistered(t *testing.T) {
	a := assert.New(t)
	a.True(prometheus.DefaultRegisterer.Unregister(metrics.AlertsReceivedTotal),
		"alerts received total not registered")
	a.True(prometheus.DefaultRegisterer.Unregister(metrics.AlertsSavedTotal),
		"alerts saved total not registered")
	a.True(prometheus.DefaultRegisterer.Unregister(metrics.DatabaseUp),
		"database up")
	a.True(prometheus.DefaultRegisterer.Unregister(metrics.InvalidWebhooksTotal),
		"invalid webhooks total")
	a.True(prometheus.DefaultRegisterer.Unregister(metrics.WebhooksReceivedTotal),
		"webhooks received total")
	a.True(prometheus.DefaultRegisterer.Unregister(metrics.AlertsSavingFailuresTotal),
		"alerts failed to be saved")

}
