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
	*geoip.GeoInfo `json:"geoip,omitempty"`
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
	for loginAttempt := range w.loginAttempts {
		accessJson := AccessJson{
			LoginAttempt: loginAttempt,
		}
		if geoip.Enabled && geoip.ReadyToWork {

			geoInfo, err := geoip.LookupIP(loginAttempt.IP)
			if !geoip.IsEmptyGeoInfo(geoInfo) {
				accessJson.GeoInfo = geoInfo
			}
			if err != nil {
				log.Println("GeoIP error:", err)
			}
		}
		jsonObject, _ := json.Marshal(accessJson)

		w.accessLogger.Println(string(jsonObject))
	}
}
