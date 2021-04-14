package conlog

import (
	"encoding/json"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
)

func LogConnection() {
	messageChannel := config.GetBroker().Subscribe()

	for collectedData := range messageChannel {
		collectedDataJson, _ := json.Marshal(collectedData)

		accessLogger := config.GetAccessLogger()
		accessLogger.Println(string(collectedDataJson))
	}
}
