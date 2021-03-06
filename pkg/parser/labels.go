package parser

import (
	"encoding/base64"
	"fmt"
)

func (p *Parser) markSession(i int) {
	Labels.Lock()
	for _, label := range Labels.L {
		if label.CheckApply(p.sessions[i]) {
			p.sessions[i].Labels = append(p.sessions[i].Labels, label)
		}
	}
	Labels.Unlock()

}

func (l *Label) CheckApply(session TCPSession) bool {
	labelType := LabelTypeToOwnerType[l.Type]
	for _, packet := range session.Packets {
		if packet.Owner == labelType {
			data, err := base64.StdEncoding.DecodeString(packet.Data)
			if err != nil {
				fmt.Println("покс", err)
			}
			matched := l.Regexp.Match(data)
			if matched {

				return true
			}
		}
	}
	return false
}

var LabelTypeToOwnerType = map[LabelType]OwnerType{
	PacketIN:  Client,
	PacketOUT: Server,
}
