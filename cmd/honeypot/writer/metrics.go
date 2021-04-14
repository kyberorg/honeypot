package writer

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

//TODO replace it with prometheus metrics

type MetricsWriter struct {
	connectionsCounter uint64
	messageChannel     chan *dto.LoginAttempt
	log                *logrus.Logger
	uniqueIPs          []string
	wg                 sync.WaitGroup
}

func NewMetricsWriter() *MetricsWriter {
	return &MetricsWriter{
		connectionsCounter: 0,
		messageChannel:     config.LoginAttemptChannel.Subscribe(),
		log:                config.GetApplicationLogger(),
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

		w.log.Printf("total number of connections: %d (unique sources %d)",
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
