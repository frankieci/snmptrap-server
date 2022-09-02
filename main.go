package main

import (
	"fmt"

	"github.com/snmptrap-server/trap"

	log "github.com/sirupsen/logrus"
)

func main() {
	trapConf := trap.TrapServerConf{}
	// print snmp trap message
	handler := func(snmp *trap.SNMPTrapMessage) error {
		header := trap.GetSNMPTrapMessageHeader(snmp.Header)
		body := trap.GetSNMPTrapMessageBody(snmp.Body)
		fmt.Println(header+body)
		return nil
	}
	
	if trapserver, err := trap.NewTrapServer(&trapConf, handler); err != nil {
		log.WithField("err", err).Fatalf("config TrapServer err")
	} else {
		log.Info("start running trap server")
		trapserver.Run()
	}
}
