package trap

import (
	"fmt"
	"strings"
	"time"

	"github.com/snmptrap-server/mibtree"
	"github.com/snmptrap-server/stringsx"

	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

var globalMibtree = mibtree.NewMibTree()

type TrapPDU struct {
	Id    int64       `json:"id"`
	OID   string      `json:"oid"`
	Type  interface{} `json:"type"`
	Value interface{} `json:"value"`
	Ts    string      `json:"ts"`
}

type TrapServer struct {
	listener *g.TrapListener
	ip       string
	port     string
}

type TrapServerConf struct {
	Ip         string `mapstructure:"ip" json:"ip" yaml:"ip"`
	Port       int64  `mapstructure:"port" json:"port" yaml:"port"`
	Version    string `mapstructure:"version" json:"version" yaml:"version"`
	Community  string `mapstructure:"community" json:"community" yaml:"community"`
	Timeout    int64  `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Maxoids    int64  `mapstructure:"maxoids" json:"maxoids" yaml:"maxoids"`
	MibMapFile string `mapstructure:"mib_map_file" json:"mib_map_file" yaml:"mib_map_file"`
	Print      bool   `mapstructure:"print" json:"print" yaml:"print"`
}

func NewTrapServer(conf *TrapServerConf, handlers ...TrapHandleFunc) (*TrapServer, error) {
	setTrapServerConf(conf)
	// load mib map file in mibtree
	if err := globalMibtree.LoadFile(conf.MibMapFile); err != nil {
		return nil, err
	}

	tl := g.NewTrapListener()

	if conf.Print {
		printHandler := func(snmp *SNMPTrapMessage) error {
			header := GetSNMPTrapMessageHeader(snmp.Header)
			body := GetSNMPTrapMessageBody(snmp.Body)
			fmt.Println(header + body)
			return nil
		}
		if len(handlers) <= 0 {
			handlers = []TrapHandleFunc{}
		}

		handlers = append(handlers, printHandler)
	}

	trapHandler := &TrapHnadler{Handlers: handlers}
	tl.OnNewTrap = trapHandler.TrapHandlerFunc()
	var version g.SnmpVersion = g.Version2c
	switch conf.Version {
	case "v1":
		version = g.Version1
	case "v2c":
		version = g.Version2c
	case "v3":
		version = g.Version3
	}

	snmpConfig := &g.GoSNMP{
		Port:               uint16(conf.Port),
		Transport:          "udp",
		Community:          conf.Community,
		Version:            version,
		Timeout:            time.Duration(conf.Timeout) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            g.MaxOids,
	}
	tl.Params = snmpConfig
	tl.Params.Logger = g.NewLogger(log.New())

	server := &TrapServer{listener: tl,
		ip:   conf.Ip,
		port: fmt.Sprintf("%d", conf.Port),
	}

	return server, nil
}

func (t *TrapServer) Run() (err error) {
	listenaddr := t.ip + ":" + t.port
	log.WithField("address", listenaddr).Info("set trapserver address ")
	return t.listener.Listen(listenaddr)
}

func GetSNMPTrapMessageHeader(header *Header) string {
	var msg string
	msg = msg + fmt.Sprintf("%s CQRCB_SYSTEM %s\n", time.Now().Format("2006-01-02 15:04:05"), header.IP)
	msg = msg + "PDU INFO:\n"
	msg = msg + "  MESSAGE TYPE: SNMP TRAP \n"
	msg = msg + fmt.Sprintf("  VERSION:  %s\n", header.Version)
	msg = msg + fmt.Sprintf("  FROM: [%s:%d]\n", header.IP, header.Port)
	msg = msg + fmt.Sprintf("  STATUS: [%s]\n", header.Error)
	msg = msg + fmt.Sprintf("  MESSAGE ID: [%v]\n", header.MsgID)
	msg = msg + fmt.Sprintf("  COMMUNITY: [%s]\n", header.Community)
	msg = msg + fmt.Sprintf("  INDEX: [%v]\n", header.ErrorIndex)
	msg = msg + fmt.Sprintf("  REQUEST ID: [%v]\n", header.RequestID)
	msg = msg + "TRAP VARIABLES:\n"
	return msg
}

func GetSNMPTrapMessageBody(pdus []*TrapPDU) string {
	body := &strings.Builder{}
	for _, v := range pdus {
		body.WriteString(fmt.Sprintf("  SNMP-MIB::%v  value=%v: [%v]\n", v.OID, v.Type, v.Value))
	}
	return body.String()
}

func setTrapServerConf(conf *TrapServerConf) {
	if stringsx.IsEmptyWithTrim(conf.Ip) {
		conf.Ip = "0.0.0.0"
	}

	if conf.Port <= 0 {
		conf.Port = 162
	}

	if stringsx.IsEmptyWithTrim(conf.Version) {
		conf.Version = "v2c"
	}

	if stringsx.IsEmptyWithTrim(conf.Community) {
		conf.Community = "public"
	}

	if conf.Timeout <= 0 {
		conf.Timeout = 2
	}

	if conf.Maxoids <= 0 {
		conf.Maxoids = 60
	}

	if stringsx.IsEmptyWithTrim(conf.MibMapFile) {
		conf.MibMapFile = "miblist.txt"
	}
}
