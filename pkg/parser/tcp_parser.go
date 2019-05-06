package parser

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type TCPSession struct {
	ServerAddr     gopacket.Endpoint
	ClientAddr     gopacket.Endpoint
	ServerPort     uint16
	ClientPort     uint16
	SequenceNumber uint32
	Packets        [][]byte
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
			newSesion := TCPSession{
				ServerAddr:     dst,
				ClientAddr:     src,
				ServerPort:     uint16(tcp.DstPort),
				ClientPort:     uint16(tcp.SrcPort),
				SequenceNumber: tcp.Seq >> 8,
				Packets:        [][]byte{packet.Data()},
			}
			index := findTcpSession(slidingWindow, packet)
			if index != -1 {
				sessions[relativeIndex] = *slidingWindow[index]
				relativeIndex += 1
				slidingWindow = append(slidingWindow[:index], slidingWindow[index+1:]...)
			}
			slidingWindow = append(slidingWindow, &newSesion)
			continue
		}
		index := findTcpSession(slidingWindow, packet)
		if index != -1 {
			slidingWindow[index].Packets = append(slidingWindow[index].Packets, packet.Data())
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
		if src != session.ClientAddr && src != session.ServerAddr {
			continue
		}
		if dst != session.ClientAddr && dst != session.ServerAddr {
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
