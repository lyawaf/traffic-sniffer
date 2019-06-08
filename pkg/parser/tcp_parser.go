package parser

import (
	"encoding/base64"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func Parse(packetSource *gopacket.PacketSource) []TCPSession {
	var sessions []TCPSession
	for packet := range packetSource.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		// new session
		if tcp.SYN && !tcp.ACK {
			net := packet.NetworkLayer()
			src, dst := net.NetworkFlow().Endpoints()
			newSession := TCPSession{
				ServerAddr:     dst.String(),
				ClientAddr:     src.String(),
				ServerPort:     uint16(tcp.DstPort),
				ClientPort:     uint16(tcp.SrcPort),
				SequenceNumber: tcp.Seq >> 8,
				Packets: []Packet{
					{
						Owner: Client,
						Data:  base64.URLEncoding.EncodeToString(packet.Data())},
				},
			}
			sessions = append(sessions, newSession)
			continue
		}
		findTcpSession(sessions, packet)
	}
	markSessions(sessions)
	return sessions
}

func findTcpSession(sessions []TCPSession, tcpPacket gopacket.Packet) {
	src, dst := tcpPacket.NetworkLayer().NetworkFlow().Endpoints()
	tcp := tcpPacket.Layer(layers.LayerTypeTCP).(*layers.TCP)
	for i, session := range sessions {
		if src.String() != session.ClientAddr && src.String() != session.ServerAddr {
			continue
		}
		if dst.String() != session.ClientAddr && dst.String() != session.ServerAddr {
			continue
		}
		if uint16(tcp.SrcPort) != session.ClientPort && uint16(tcp.SrcPort) != session.ServerPort {
			continue
		}
		if uint16(tcp.DstPort) != session.ClientPort && uint16(tcp.DstPort) != session.ServerPort {
			continue
		}
		packet := Packet{
			Owner: Client,
			Data: base64.URLEncoding.EncodeToString(tcpPacket.Data()),
		}
		switch src.String() {
		case session.ClientAddr:
			packet.Owner = Client
		case session.ServerAddr:
			packet.Owner = Server
		}
		sessions[i].Packets = append(sessions[i].Packets, packet)
	}
}

func markSessions(sessions []TCPSession) {
	for i, session := range sessions {
		for _, label := range Labels {
			if label.CheckApply(session) {
				fmt.Println("Add label")
				sessions[i].Labels = append(sessions[i].Labels, Label{})
			}
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
