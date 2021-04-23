package logger

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"log"
	"os"
	"sync"
)

var (
	accessLogger      *log.Logger
	applicationLogger *logrus.Logger
)

var once sync.Once

func init() {
	once.Do(func() {
		initApplicationLogger()
		initAccessLogger()
	})
}

//GetAccessLogger logger for access log
func GetAccessLogger() *log.Logger {
	return accessLogger
}

//GetApplicationLogger main app logger
func GetApplicationLogger() *logrus.Logger {
	return applicationLogger
}

func initAccessLogger() {
	accessLog := config.GetAppConfig().AccessLog
	var logDestination = getLogDestination(accessLog)

	if accessLogger == nil {
		writer := io.MultiWriter(logDestination)
		accessLogger = log.New(writer, "", 0)
	}
}

func initApplicationLogger() {
	applicationLog := config.GetAppConfig().ApplicationLog
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
