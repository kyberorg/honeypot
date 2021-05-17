package config

import "github.com/kyberorg/honeypot/cmd/honeypot/util"

var (
	loginAttemptBroadcaster *util.LoginAttemptBroadcaster
)

func init() {
	//init broadcaster
	initLoginAttemptBroadcaster()
}

func GetLoginAttemptBroadcaster() *util.LoginAttemptBroadcaster {
	return loginAttemptBroadcaster
}

func initLoginAttemptBroadcaster() {
	loginAttemptBroadcaster = util.NewLoginAttemptBroadcaster()
	go loginAttemptBroadcaster.Start()
}
