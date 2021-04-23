package writer

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/kyberorg/honeypot/cmd/honeypot/logger"
	"sync"
	"sync/atomic"
)

//TODO replace it with prometheus metrics

var log = logger.GetApplicationLogger()

type MetricsWriter struct {
	connectionsCounter uint64
	messageChannel     chan *dto.LoginAttempt
	uniqueIPs          []string
	wg                 sync.WaitGroup
}

func NewMetricsWriter() *MetricsWriter {
	return &MetricsWriter{
		connectionsCounter: 0,
		messageChannel:     config.GetLoginAttemptChannel().Subscribe(),
		uniqueIPs:          make([]string, 0),
	}
}

func (w *MetricsWriter) RecordMetric() {
	for collectedData := range w.messageChannel {
		w.wg.Add(1)
		atomic.AddUint64(&w.connectionsCounter, 1)
		w.wg.Done()

		if w.isNewIPConnected(collectedData.IP) {
			w.uniqueIPs = append(w.uniqueIPs, collectedData.IP)
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
