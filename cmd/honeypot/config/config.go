package config

import (
	"gopkg.in/alecthomas/kingpin.v2"
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
	promMetricsEnabled = kingpin.Flag("with-prom-metrics", "Enables Prometheus Metrics Module").Bool()
	promMetricsPort    = kingpin.Flag("prom-metrics-port", "Port for serving metrics").Default("2112").
				Uint16()
	promMetricsPath = kingpin.Flag("prom-metrics-path", "Custom path where metrics are served").
			Default("/metrics").String()
	promMetricsPrefix = kingpin.Flag("prom-metrics-prefix", "Custom metrics prefix").String()
)

//geoip module flags
var (
	geoIpDatabaseFile = kingpin.Flag("geoip-mmdb-file", "Location of MaxMind City MMDB file").String()
)

//raw metrics module flags
var (
	rawMetricsEnabled = kingpin.Flag("with-raw-metrics", "Enables Raw Metrics Module").Bool()
	rawMetricsFile    = kingpin.Flag("raw-metrics-file", "File to write metrics to").String()
	rawMetricsPrefix  = kingpin.Flag("raw-metrics-prefix", "Custom metrics prefix").String()
)

//internal vars
var (
	appConfig *applicationConfiguration
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

	//GeoIP module flags
	GeoIP

	//RawMetrics module flags
	RawMetrics
}

//PromMetrics module flags
type PromMetrics struct {
	//Prom Metrics module
	Enabled bool
	Port    uint16
	Path    string
	Prefix  string
}

type GeoIP struct {
	DatabaseFile string
}

//RawMetrics module flags
type RawMetrics struct {
	Enabled bool
	File    string
	Prefix  string
}

func init() {
	//parse flags
	kingpin.Parse()

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
		GeoIP: GeoIP{
			DatabaseFile: *geoIpDatabaseFile,
		},
		RawMetrics: RawMetrics{
			Enabled: *rawMetricsEnabled,
			File:    *rawMetricsFile,
			Prefix:  *rawMetricsPrefix,
		},
	}
}

//GetAppConfig returns application configuration object
func GetAppConfig() *applicationConfiguration {
	return appConfig
}

//IsPromMetricsModuleEnabled says if PromMetrics module is enabled or not, based on activation flag.
func IsPromMetricsModuleEnabled() bool {
	return appConfig.PromMetrics.Enabled
}

//IsRawMetricsModuleEnabled reports if RawMetrics module is enabled or not, based on activation flag.
func IsRawMetricsModuleEnabled() bool {
	return appConfig.RawMetrics.Enabled
}
