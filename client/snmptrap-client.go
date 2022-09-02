package main

import (
	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

func main(){
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = "127.0.0.1"
	g.Default.Port = 162
	g.Default.Version = g.Version2c
	g.Default.Community = "public"
	g.Default.Logger = g.NewLogger(log.New())
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	pdu := g.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0",
		Type:  g.ObjectIdentifier,
		Value: ".1.3.6.1.6.3.1.1.5.1",
	}
	pdustr := g.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0",
		Type:  g.OctetString,
		Value: "hello world",
	}

	trap := g.SnmpTrap{
		Variables: []g.SnmpPDU{pdu,pdustr},
	}

	_, err = g.Default.SendTrap(trap)
	if err != nil {
		log.Fatalf("SendTrap() err: %v", err)
	}
}