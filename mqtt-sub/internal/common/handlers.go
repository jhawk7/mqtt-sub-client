package common

import log "github.com/sirupsen/logrus"

func LogError(err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatalf("fatal error: %v", err)
		} else {
			log.Errorf("error: %v", err)
		}
	}
}

func LogInfo(msg string) {
	log.Info(msg)
}
