package config

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/util"
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"sync"
)

//core flags
var (
	port = kingpin.Flag("port", "Port we start at").Short('p').
		Envar("PORT").Default("22").Uint16()
	accessLog = kingpin.Flag("access-log", "Where to log requests").
			Envar("ACCESS_LOG").String()
	applicationLog = kingpin.Flag("log", "File to sent application logs").
			Envar("LOG_FILE").String()
	hostKey = kingpin.Flag("hostkey", "File with private id_rsa key that is used to identify server").
		Envar("HOSTKEY").String()
	skipHostKeyGeneration = kingpin.Flag("skip-hostkey-generation",
		"If set, app won't generate hostkey at start-up").Bool()
)

//prometheus metrics module flags
var (
	promMetricsEnabled = kingpin.Flag("prom-metrics-enable", "Enables Prometheus Metrics Module").Bool()
	promMetricsPort    = kingpin.Flag("prom-metrics-port", "Port for serving metrics").Default("2112").
				Uint16()
	promMetricsPath = kingpin.Flag("prom-metrics-path", "Custom path where metrics are served").
			Default("/metrics").String()
	promMetricsPrefix = kingpin.Flag("prom-metrics-prefix", "Custom metrics prefix").String()
)

//internal vars
var (
	once              sync.Once
	appConfig         *applicationConfiguration
	broadcasterObject *util.Broadcaster
	accessLogger      *log.Logger
	applicationLogger *logrus.Logger
)

//applicationConfiguration application configuration values
type applicationConfiguration struct {
	//SSH Port
	Port uint16
	//Access Log filename
	AccessLog string
	//Application Log filename
	ApplicationLog string
	//HostKey filename
	HostKey string
	//Generate key, if HostKey absent
	GenerateHostKey bool

	//PromMetrics module flags
	PromMetrics
}

//PromMetrics module flags
type PromMetrics struct {
	//Prom Metrics module
	Enabled bool
	Port    uint16
	Path    string
	Prefix  string
}

func init() {
	once.Do(func() {
		//parse flags
		kingpin.Parse()
		//init broadcaster
		initBroadcaster()
	})

	appConfig = &applicationConfiguration{
		Port:            *port,
		AccessLog:       *accessLog,
		ApplicationLog:  *applicationLog,
		HostKey:         *hostKey,
		GenerateHostKey: !*skipHostKeyGeneration,

		PromMetrics: PromMetrics{
			Enabled: *promMetricsEnabled,
			Port:    *promMetricsPort,
			Path:    *promMetricsPath,
			Prefix:  *promMetricsPrefix,
		},
	}
}

//GetAppConfig returns application configuration object
func GetAppConfig() *applicationConfiguration {
	return appConfig
}

func GetLoginAttemptChannel() *util.Broadcaster {
	return broadcasterObject
}

//GetAccessLogger logger for access log
func GetAccessLogger() *log.Logger {
	accessLog := appConfig.AccessLog
	var logDestination = getLogDestination(accessLog)

	if accessLogger == nil {
		writer := io.MultiWriter(logDestination)
		accessLogger = log.New(writer, "", 0)
	}
	return accessLogger
}

//GetApplicationLogger main app logger
func GetApplicationLogger() *logrus.Logger {
	applicationLog := appConfig.ApplicationLog
	logDestination := getLogDestination(applicationLog)
	writer := io.MultiWriter(logDestination)

	if applicationLogger == nil {
		applicationLogger = logrus.New()
		applicationLogger.SetFormatter(&easy.Formatter{
			TimestampFormat: "02/01/2006 15:04:05-0700",
			LogFormat:       "%time% - %msg%\n",
		})

		applicationLogger.SetOutput(writer)
	}
	return applicationLogger
}

//IsPromMetricsModuleEnabled says if PromMetrics module is enabled or not, based on activation flag.
func IsPromMetricsModuleEnabled() bool {
	return appConfig.PromMetrics.Enabled
}

func initBroadcaster() {
	broadcasterObject = util.NewBroadcaster()
	go broadcasterObject.Start()
}

//log to file or os.Stdout
func getLogDestination(logFile string) *os.File {
	var logLocation *os.File
	if logFile == "" {
		logLocation = os.Stdout
	} else {
		var logOpenError error
		logLocation, logOpenError = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if logOpenError != nil {
			log.Fatalln("Unable to open log file" + logOpenError.Error())
		}
	}
	return logLocation
}
