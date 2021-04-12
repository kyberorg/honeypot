package config

import (
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
	hostKey = kingpin.Flag("hostkey", "File with private id_rsa key that is used to identify server").
		Envar("HOSTKEY").String()
	skipHostKeyGeneration = kingpin.Flag("skip-hostkey-generation",
		"If set, app won't generate hostkey at start-up").Bool()
)

//logger for access log
var accessLogger *log.Logger

//are params already parsed
var alreadyParsed = false

//AppConfig application configuration values
type AppConfig struct {
	//SSH Port
	Port uint16
	//Access Log filename
	AccessLog string
	//HostKey filename
	HostKey string
	//Generate key, if HostKey absent
	GenerateHostKey bool
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
		HostKey:         *hostKey,
		GenerateHostKey: !*skipHostKeyGeneration,
	}
}

//GetAccessLogger logger for access log
func GetAccessLogger(accessLog string) *log.Logger {
	var logLocation *os.File
	if accessLog == "" {
		logLocation = os.Stderr
	} else {
		var logOpenError error
		logLocation, logOpenError = os.OpenFile(accessLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if logOpenError != nil {
			log.Fatalln("Unable to open access log file" + logOpenError.Error())
		}
	}
	if accessLogger == nil {
		w := io.MultiWriter(logLocation)
		accessLogger = log.New(w, "", 0)
	}
	return accessLogger

}
