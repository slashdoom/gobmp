package message

import (
	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/bmp"
)

// sessionInitData holds router identity fields from the BMP Initiation message.
// Stored in memory on the producer for the lifetime of the BMP session.
type sessionInitData struct {
	SysName  string
	SysDescr string
}

// produceSessionInitMessage stores Init TLV data in memory for later enrichment.
func (p *producer) produceSessionInitMessage(msg bmp.Message) {
	im, ok := msg.Payload.(*bmp.InitiationMessage)
	if !ok {
		glog.Errorf("session init: invalid payload type %T", msg.Payload)
		return
	}
	data := &sessionInitData{}
	for _, tlv := range im.TLV {
		switch tlv.InformationType {
		case 1:
			data.SysDescr = string(tlv.Information)
		case 2:
			data.SysName = string(tlv.Information)
		}
	}
	p.pendingSessionInit = data
}

// sessionEnrich returns sys_name and sys_descr from the BMP Init message when
// session_tracking is enabled and an Init has been received on this connection.
func (p *producer) sessionEnrich() (sysName, sysDescr string) {
	if !p.sessionTracking {
		return
	}
	if p.pendingSessionInit != nil {
		sysName = p.pendingSessionInit.SysName
		sysDescr = p.pendingSessionInit.SysDescr
	}
	return
}
