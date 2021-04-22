package writer

import (
	"encoding/json"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/kyberorg/honeypot/cmd/honeypot/logger"
)

type AccessLogWriter struct {
	messageChannel chan *dto.LoginAttempt
}

func NewAccessLogWriter() *AccessLogWriter {
	return &AccessLogWriter{
		messageChannel: config.GetLoginAttemptChannel().Subscribe(),
	}
}

func (alw *AccessLogWriter) WriteToLog() {
	for collectedData := range alw.messageChannel {
		collectedDataJson, _ := json.Marshal(collectedData)

		accessLogger := logger.GetAccessLogger()
		accessLogger.Println(string(collectedDataJson))
	}
}
