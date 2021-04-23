package writer

import (
	"encoding/json"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/kyberorg/honeypot/cmd/honeypot/logger"
	logg "log"
)

var (
	singleAccessLogWriter *AccessLogWriter
)

type AccessLogWriter struct {
	loginAttempts chan *dto.LoginAttempt
	accessLogger  *logg.Logger
}

func init() {
	singleAccessLogWriter = &AccessLogWriter{
		accessLogger:  logger.GetAccessLogger(),
		loginAttempts: config.GetLoginAttemptBroadcaster().Subscribe(),
	}
}

func GetAccessLogWriter() *AccessLogWriter {
	return singleAccessLogWriter
}

func (w *AccessLogWriter) WriteToLog() {
	for loginAttempt := range w.loginAttempts {
		jsonObject, _ := json.Marshal(loginAttempt)

		w.accessLogger.Println(string(jsonObject))
	}
}
