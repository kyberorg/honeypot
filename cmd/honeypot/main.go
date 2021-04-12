package main

import (
	"encoding/json"
	"github.com/gliderlabs/ssh"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/sshutil"
	"github.com/kyberorg/honeypot/cmd/honeypot/util"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var appConfig = config.GetAppConfig()

type collectedData struct {
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

func passwordHandler(ctx ssh.Context, password string) bool {
	ip, ipErr := util.ParseIP(ctx.RemoteAddr().String())
	if ipErr != nil {
		log.Println("new connection")
	} else {
		log.Println("new connection from", ip)
	}

	collectedData := collectedData{
		User:     ctx.User(),
		Password: password,
		IP:       ip,
	}
	collectedDataJson, _ := json.Marshal(collectedData)

	accessLogger := config.GetAccessLogger(appConfig.AccessLog)
	accessLogger.Println(string(collectedDataJson))

	//small delay to emulate "real" SSH
	time.Sleep(1 * time.Second)
	return false
}

func main() {
	appConfig := config.GetAppConfig()

	//getting HostKey
	hostKey, hostKeyErr := sshutil.HostKey(&appConfig)
	if hostKeyErr != nil {
		if hostKeyErr.Error() == sshutil.NoHostKeyMarker {
			hostKey = nil
		} else {
			log.Fatalln(hostKeyErr)
		}
	}

	portString := strconv.Itoa(int(appConfig.Port))

	sshServer := &ssh.Server{
		Addr:            ":" + portString,
		PasswordHandler: passwordHandler,
	}

	if hostKey != nil {
		sshServer.AddHostKey(hostKey)
	}

	log.Println("Starting SSH Server at port", portString)
	log.Println("ready to access connections")

	if appConfig.AccessLog != "" {
		log.Println("Logging connections to ", appConfig.AccessLog)
	}

	log.Fatalln(sshServer.ListenAndServe())
}
