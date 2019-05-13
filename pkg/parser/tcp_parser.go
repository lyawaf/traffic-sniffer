package parser

import (
	"encoding/base64"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type TCPSession struct {
	ServerAddr     string
	ClientAddr     string
	ServerPort     uint16
	ClientPort     uint16
	SequenceNumber uint32
	Packets        []string
}

func Parse(packetSource *gopacket.PacketSource, relativeIndex int) map[int]TCPSession {
	sessions := make(map[int]TCPSession)
	var slidingWindow []*TCPSession
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
				Packets:        []string{base64.URLEncoding.EncodeToString(packet.Data())},
			}
			println("New server addr", newSession.ServerAddr)
			println("new client addr", newSession.ClientAddr)
			index := findTcpSession(slidingWindow, packet)
			if index != -1 {
				sessions[relativeIndex] = *slidingWindow[index]
				relativeIndex += 1
				slidingWindow = append(slidingWindow[:index], slidingWindow[index+1:]...)
			}
			slidingWindow = append(slidingWindow, &newSession)
			continue
		}
		index := findTcpSession(slidingWindow, packet)
		if index != -1 {
			slidingWindow[index].Packets = append(slidingWindow[index].Packets, base64.URLEncoding.EncodeToString(packet.Data()))
		}
	}
	for _, session := range slidingWindow {
		sessions[relativeIndex] = *session
		relativeIndex += 1
	}
	return sessions
}

func findTcpSession(sessions []*TCPSession, tcpPacket gopacket.Packet) int {
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
		return i
	}
	return -1
}
