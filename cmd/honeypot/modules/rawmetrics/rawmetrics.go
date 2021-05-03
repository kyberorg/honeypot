package rawmetrics

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"sync"
	"sync/atomic"
)

var log = config.GetApplicationLogger()
var singleMetricsWriter *MetricsWriter

const (
	// Prefix for all metrics.
	defaultPrefix = "honeypot"
)

type MetricsWriter struct {
	connectionsCounter uint64
	loginAttempts      chan *dto.LoginAttempt
	uniqueIPs          []string
	prefix             string
	wg                 sync.WaitGroup
}

func init() {
	singleMetricsWriter = &MetricsWriter{
		connectionsCounter: 0,
		loginAttempts:      config.GetLoginAttemptBroadcaster().Subscribe(),
		uniqueIPs:          make([]string, 0),
		prefix:             getPrefix(),
	}
}

func GetMetricsWriter() *MetricsWriter {
	return singleMetricsWriter
}

func (w *MetricsWriter) RecordMetric() {
	for loginAttempt := range w.loginAttempts {
		w.wg.Add(1)
		atomic.AddUint64(&w.connectionsCounter, 1)
		w.wg.Done()

		if w.isNewIPConnected(loginAttempt.IP) {
			w.uniqueIPs = append(w.uniqueIPs, loginAttempt.IP)
		}

		log.Printf("total number of connections: %d (unique sources %d)",
			w.connectionsCounter, len(w.uniqueIPs))

		//TODO map<String(IP), attempts>
	}
}

func (w *MetricsWriter) isNewIPConnected(ip string) bool {
	for i := range w.uniqueIPs {
		if w.uniqueIPs[i] == ip {
			return false
		}
	}
	return true
}

func getPrefix() string {
	metricsPrefix := config.GetAppConfig().RawMetrics.Prefix
	if metricsPrefix == "" {
		metricsPrefix = defaultPrefix
	}
	return metricsPrefix
}
