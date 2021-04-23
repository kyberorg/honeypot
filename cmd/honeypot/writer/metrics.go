package writer

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/kyberorg/honeypot/cmd/honeypot/logger"
	"sync"
	"sync/atomic"
)

//TODO replace it with raw metrics module

var log = logger.GetApplicationLogger()
var singleMetricsWriter *MetricsWriter

type MetricsWriter struct {
	connectionsCounter uint64
	loginAttempts      chan *dto.LoginAttempt
	uniqueIPs          []string
	wg                 sync.WaitGroup
}

func init() {
	singleMetricsWriter = &MetricsWriter{
		connectionsCounter: 0,
		loginAttempts:      config.GetLoginAttemptBroadcaster().Subscribe(),
		uniqueIPs:          make([]string, 0),
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
