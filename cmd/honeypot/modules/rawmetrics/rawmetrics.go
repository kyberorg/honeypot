package rawmetrics

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"io/ioutil"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	log                 = config.GetApplicationLogger()
	singleMetricsWriter *MetricsWriter
)

const (
	// Prefix for all metrics.
	defaultPrefix        = "honeypot"
	connectionPostfix    = "_connections"
	uniqueSourcesPostfix = "_unique_sources"
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

		if config.LogRawMetricsToFile() {
			connectionsMetric := w.prefix + connectionPostfix + ": " + strconv.Itoa(int(w.connectionsCounter))
			uniqueIpMetric := w.prefix + uniqueSourcesPostfix + ": " + strconv.Itoa(len(w.uniqueIPs))

			record := []byte(connectionsMetric + "\n" + uniqueIpMetric + "\n")
			err := ioutil.WriteFile(config.GetAppConfig().File, record, 0644)
			if err != nil {
				log.Fatalln("Unable to write raw metrics to " + config.GetAppConfig().File +
					"Since you enabled raw module, this is probably not what you want to expect.")
			}
		} else {
			log.Printf("total number of connections: %d (unique sources %d)",
				w.connectionsCounter, len(w.uniqueIPs))
		}
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
