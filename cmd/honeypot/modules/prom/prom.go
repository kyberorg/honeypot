package prom

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
)

var (
	once          sync.Once
	singleHandler *PrometheusMetricsHandler
)

type PrometheusMetricsHandler struct {
	connectionsCounter  prometheus.Counter
	uniqueSourceCounter prometheus.Counter
	messageChannel      chan *dto.LoginAttempt
	log                 *logrus.Logger
	uniqueIPs           []string
}

// Metric name parts.
const (
	// Prefix for all metrics.
	defaultPrefix = "honeypot"
)

func init() {
	prefix := getPrefix()

	singleHandler = &PrometheusMetricsHandler{
		connectionsCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_connections",
			Help: "Total number of connections",
		}),
		uniqueSourceCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: prefix + "_unique_sources",
			Help: "Number of unique sources",
		}),

		messageChannel: config.GetLoginAttemptChannel().Subscribe(),
		log:            config.GetApplicationLogger(),
		uniqueIPs:      make([]string, 0),
	}
}

func GetPrometheusMetricsHandler() *PrometheusMetricsHandler {
	return singleHandler
}

func (h *PrometheusMetricsHandler) StartMetricsServer() {
	once.Do(func() {
		port := strconv.Itoa(int(config.GetAppConfig().PromMetrics.Port))
		path := config.GetAppConfig().PromMetrics.Path

		h.log.Printf("Starting metrics server at port %s", port)

		http.Handle(path, promhttp.Handler())
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			h.log.Fatalln("Unable to start prometheus metrics server." +
				"Since you enabled prom module, this is probably not what you want to expect")
		}
	})
}

func (h *PrometheusMetricsHandler) RecordMetrics() {
	for loginAttempt := range h.messageChannel {
		h.connectionsCounter.Inc()

		//adding uniq ips list if unique
		if h.isNewIPConnected(loginAttempt.IP) {
			h.uniqueSourceCounter.Inc()
			h.uniqueIPs = append(h.uniqueIPs, loginAttempt.IP)
		}
	}
}

func (h *PrometheusMetricsHandler) isNewIPConnected(ip string) bool {
	for i := range h.uniqueIPs {
		if h.uniqueIPs[i] == ip {
			return false
		}
	}
	return true
}

func getPrefix() string {
	metricsPrefix := config.GetAppConfig().PromMetrics.Prefix
	if metricsPrefix == "" {
		metricsPrefix = defaultPrefix
	}
	return metricsPrefix
}
