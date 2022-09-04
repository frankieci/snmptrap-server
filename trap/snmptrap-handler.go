package trap

import (
	"fmt"
	"net"
	"time"

	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

type TrapHandleFunc func(snmp *SNMPTrapMessage) error

type TrapHnadler struct {
	Handlers []TrapHandleFunc
}

type SNMPTrapMessage struct {
	Header     *Header
	Body      []*TrapPDU
	Timestamp int64
}

type Header struct {
	IP         string
	Port       int
	Version    g.SnmpVersion
	Community  string
	PDUType    g.PDUType
	MsgID      uint32
	RequestID  uint32
	Error      g.SNMPError
	ErrorIndex uint8
}

func (t *TrapHnadler) TrapHandlerFunc() g.TrapHandlerFunc {
	return t.BaseTrapHandler
}

func (t *TrapHnadler) BaseTrapHandler(packet *g.SnmpPacket, addr *net.UDPAddr) {
	log.WithFields(map[string]interface{}{"addr": addr.IP.String(), "port": addr.Port}).Info("got trap package from")

	header := t.header(addr.IP.String(), addr.Port, packet)

	pdus := t.parseSnmpPacket(packet)

	for _, handler := range t.Handlers {
		if handler != nil {
			snmpMessage := &SNMPTrapMessage{
				Header:    header,
				Body:      pdus,
				Timestamp: int64(packet.Timestamp),
			}
			if err := handler(snmpMessage); err != nil {
				log.WithField("err", err).Error("%v HandlerFunc handle SNMPTrapMessage run error", handler)
			}
		}
	}
}

func (t *TrapHnadler) parseSnmpPacket(packet *g.SnmpPacket) []*TrapPDU {
	pdus := make([]*TrapPDU, 0)
	// parse snmp packet
	for _, v := range packet.Variables {
		oidName := ""
		if name, err := globalMibtree.FindNodeName(v.Name); err != nil {
			log.WithField("err", err).Error("trans oid to name error")
			oidName = v.Name
		} else {
			oidName = name
		}

		switch v.Type {
		case g.OctetString:
			b := v.Value.([]byte)
			log.WithField("OID", v.Name).WithField("string", fmt.Sprintf("%s", b)).WithField("Type", v.Type).Info()
			pdu := TrapPDU{
				OID:   oidName,
				Type:  v.Type,
				Value: fmt.Sprintf("%s", b),
				Ts:    time.Now().Format("2006-01-02 15:04:05"),
			}
			pdus = append(pdus, &pdu)
		case g.ObjectIdentifier:
			objid := fmt.Sprintf("%s", v.Value)
			objname := ""
			if name, err := globalMibtree.FindNodeName(objid); err != nil {
				log.WithField("err", err).Error("trans oid to name error")
				objname = objid
			} else {
				objname = name
			}

			pdu := TrapPDU{
				OID:   oidName,
				Type:  v.Type,
				Value: objname,
				Ts:    time.Now().Format("2006-01-02 15:04:05"),
			}
			pdus = append(pdus, &pdu)
		default:
			pdu := TrapPDU{
				OID:   oidName,
				Type:  v.Type,
				Value: v.Value,
				Ts:    time.Now().Format("2006-01-02 15:04:05"),
			}
			pdus = append(pdus, &pdu)
		}
	}
	return pdus
}

func (t *TrapHnadler) header(ip string, port int, packet *g.SnmpPacket) *Header {
	return &Header{
		IP:         ip,
		Port:       port,
		Version:    packet.Version,
		Community:  packet.Community,
		PDUType:    packet.PDUType,
		MsgID:      packet.MsgID,
		RequestID:  packet.RequestID,
		Error:      packet.Error,
		ErrorIndex: packet.ErrorIndex,
	}
}
