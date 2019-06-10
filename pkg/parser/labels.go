package parser

import (
	"encoding/base64"
	"fmt"
)

func markSession(session TCPSession) {
	for _, label := range Labels {
		if label.CheckApply(session) {
			fmt.Println("Add label")
			session.Labels = append(session.Labels, Label{})
		}
	}

}

func (l *Label) CheckApply(session TCPSession) bool {
	labelType := LabelTypeToOwnerType[l.Type]
	for _, packet := range session.Packets {
		if packet.Owner == labelType {
			data, err := base64.URLEncoding.DecodeString(packet.Data)
			if err != nil {
				fmt.Println("покс")
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
