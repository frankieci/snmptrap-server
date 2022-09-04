package main

import (
	"github.com/snmptrap-server/trap"

	log "github.com/sirupsen/logrus"
)

func main() {
	trapConf := trap.TrapServerConf{}

	if trapserver, err := trap.NewTrapServer(&trapConf, nil); err != nil {
		log.WithField("err", err).Fatalf("config TrapServer err")
	} else {
		log.Info("start running trap server")
		trapserver.Run()
	}
}
