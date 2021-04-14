package metrics

import (
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"sync"
	"sync/atomic"
)

//TODO replace it with prometheus metrics

var log = config.GetApplicationLogger()

var connections uint64

var ips = make([]string, 0)

var wg sync.WaitGroup

func RecordMetric() {
	messageChannel := config.GetBroker().Subscribe()

	for collectedData := range messageChannel {
		wg.Add(1)
		atomic.AddUint64(&connections, 1)
		wg.Done()

		if isNewIPConnected(collectedData.IP) {
			ips = append(ips, collectedData.IP)
		}

		log.Println("number of connections: ", connections)
		log.Println("unique connections: ", len(ips))

		//TODO map<String(IP), attempts>
	}
}

func isNewIPConnected(ip string) bool {
	for i := range ips {
		if ips[i] == ip {
			return false
		}
	}
	return true
}
