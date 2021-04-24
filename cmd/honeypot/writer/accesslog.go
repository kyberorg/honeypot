package writer

import (
	"encoding/json"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/kyberorg/honeypot/cmd/honeypot/logger"
	"github.com/kyberorg/honeypot/cmd/honeypot/modules/geoip"
	logg "log"
)

var (
	singleAccessLogWriter *AccessLogWriter
)

type AccessJson struct {
	*dto.LoginAttempt
	*geoip.GeoInfo `json:"geoip"`
}

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
	var err error
	for loginAttempt := range w.loginAttempts {
		accessJson := AccessJson{
			LoginAttempt: loginAttempt,
		}
		if geoip.Enabled && geoip.ReadyToWork {
			accessJson.GeoInfo, err = geoip.LookupIP(loginAttempt.IP)
			if err != nil {
				log.Println("GeoIP error:", err)
			}
		}
		jsonObject, _ := json.Marshal(accessJson)

		w.accessLogger.Println(string(jsonObject))
	}
}
