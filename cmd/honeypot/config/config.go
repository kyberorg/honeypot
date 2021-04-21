package config

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/util"
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
)

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

//prom metrics params
var (
	promMetricsEnabled = kingpin.Flag("prom-metrics-enable", "Enables Prometheus Metrics Module").Bool()
	promMetricsPort    = kingpin.Flag("prom-metrics-port", "Port for serving metrics").Default("2112").
				Uint16()
	promMetricsPath = kingpin.Flag("prom-metrics-path", "Custom path where metrics are served").
			Default("/metrics").String()
	promMetricsPrefix = kingpin.Flag("prom-metrics-prefix", "Custom metrics prefix").String()
)

//LoginAttemptChannel for sending and receiving dto.LoginAttempt objects
var LoginAttemptChannel = getBroadcaster()

//logger for access log
var accessLogger *log.Logger

//application logger
var applicationLogger *logrus.Logger

//are params already parsed
var alreadyParsed = false

//singleton keeper
var broadcasterObject *util.Broadcaster

//AppConfig application configuration values
type AppConfig struct {
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

//GetAppConfig parses args and converts 'em to AppConfig
func GetAppConfig() AppConfig {
	if !alreadyParsed {
		kingpin.Parse()
		alreadyParsed = true
	}

	return AppConfig{
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

//GetAccessLogger logger for access log
func GetAccessLogger() *log.Logger {
	accessLog := GetAppConfig().AccessLog
	var logDestination = getLogDestination(accessLog)

	if accessLogger == nil {
		writer := io.MultiWriter(logDestination)
		accessLogger = log.New(writer, "", 0)
	}
	return accessLogger
}

//GetApplicationLogger main app logger
func GetApplicationLogger() *logrus.Logger {
	applicationLog := GetAppConfig().ApplicationLog
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
	return GetAppConfig().PromMetrics.Enabled
}

func getBroadcaster() *util.Broadcaster {
	if broadcasterObject == nil {
		broadcasterObject = util.NewBroadcaster()
		go broadcasterObject.Start()
	}
	return broadcasterObject
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
