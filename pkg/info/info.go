package info

import (
	"os"
	"time"
)

const (
	ENV_KEY          = "TALKY_INFO_ENV"
	DEFAULT_INFO_ENV = "local"
)

var (
	// version & build date gets defined by the build system
	Version   = "dirty"
	BuildDate string

	Env       = DEFAULT_INFO_ENV
	StartDate = time.Now()
)

func getEnv() {
	val, ok := os.LookupEnv(ENV_KEY)
	if !ok {
		Env = DEFAULT_INFO_ENV
	} else {
		Env = val
	}
}

func InitInfo() {
	getEnv()
}
