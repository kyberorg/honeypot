package main

import (
	"github.com/gliderlabs/ssh"
	"github.com/kyberorg/honeypot/cmd/honeypot/config"
	"github.com/kyberorg/honeypot/cmd/honeypot/dto"
	"github.com/kyberorg/honeypot/cmd/honeypot/logger"
	"github.com/kyberorg/honeypot/cmd/honeypot/modules/prom"
	"github.com/kyberorg/honeypot/cmd/honeypot/sshutil"
	"github.com/kyberorg/honeypot/cmd/honeypot/util"
	"github.com/kyberorg/honeypot/cmd/honeypot/writer"
	"strconv"
	"time"
)

var log = logger.GetApplicationLogger()

func main() {

	//register writers (functions receiving published by passwordHandler object)
	registerWriters()

	//getting HostKey
	hostKey, hostKeyErr := sshutil.HostKey()
	if hostKeyErr != nil {
		if hostKeyErr.Error() == sshutil.NoHostKeyMarker {
			hostKey = nil
		} else {
			log.Fatalln(hostKeyErr)
		}
	}

	portString := strconv.Itoa(int(config.GetAppConfig().Port))

	sshServer := &ssh.Server{
		Addr:            ":" + portString,
		PasswordHandler: passwordHandler,
	}

	if hostKey != nil {
		sshServer.AddHostKey(hostKey)
	}

	log.Println("Starting SSH Server at port", portString)
	log.Println("ready to access connections")

	if config.GetAppConfig().AccessLog != "" {
		log.Println("Logging connections to", config.GetAppConfig().AccessLog)
	}

	if config.IsPromMetricsModuleEnabled() {
		go prom.GetPrometheusMetricsHandler().StartMetricsServer()
	}

	log.Fatalln(sshServer.ListenAndServe())
}

func registerWriters() {
	go writer.NewAccessLogWriter().WriteToLog()
	go writer.NewMetricsWriter().RecordMetric()
	if config.IsPromMetricsModuleEnabled() {
		go prom.GetPrometheusMetricsHandler().RecordMetrics()
	}
}

func passwordHandler(ctx ssh.Context, password string) bool {
	ip, ipErr := util.ParseIP(ctx.RemoteAddr().String())
	if ipErr != nil {
		log.Println("new connection")
	} else {
		log.Println("new connection from", ip)
	}

	loginAttempt := dto.LoginAttempt{
		Time:     time.Now().Format("02/01/2006 15:04:05-0700"),
		User:     ctx.User(),
		Password: password,
		IP:       ip,
	}

	config.GetLoginAttemptChannel().Send(&loginAttempt)

	//small delay to emulate "real" SSH
	time.Sleep(1 * time.Second)
	return false
}
