package parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"go.mongodb.org/mongo-driver/bson"
)

func (p *Parser) Parse() {
	go p.saveWorker(WAIT_TIMEOUT * 2 * time.Second)
	for rawPacket := range p.Source.Packets() {
		tcpLayer := rawPacket.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}
		tcp, _ := tcpLayer.(*layers.TCP)
		// new session
		if tcp.SYN && !tcp.ACK {
			newSession := createNewSession(rawPacket)
			oldSessionIndex := p.findTcpSession(rawPacket)
			p.Lock()
			if oldSessionIndex != -1 {
				p.saveSession(oldSessionIndex)
				p.sessions[oldSessionIndex] = newSession
			} else {
				p.sessions = append(p.sessions, newSession)
			}
			p.Unlock()
			continue
		}
		i := p.findTcpSession(rawPacket)
		if i != -1 {
			packet := p.makePacket(i, rawPacket)
			p.addPacketToSession(i, packet)
		}
	}
}

// saveWorker saves sessions which are end by wait timeout.
func (p *Parser) saveWorker(d time.Duration) {
	for x := range time.Tick(d) {
		fmt.Println("[WORKER]", x)
		var sessionsCopy []TCPSession
		p.Lock()
		for i, session := range p.sessions {
			if time.Now().Second()-session.LastUpdate.Second() > WAIT_TIMEOUT {
				fmt.Println("[WORKER] Save session.")
				p.saveSession(i)
				continue
			}
			sessionsCopy = append(sessionsCopy, session)
		}
		p.sessions = sessionsCopy
		p.Unlock()
	}
}

func (p *Parser) findTcpSession(tcpPacket gopacket.Packet) int {
	src, dst := tcpPacket.NetworkLayer().NetworkFlow().Endpoints()
	tcp := tcpPacket.Layer(layers.LayerTypeTCP).(*layers.TCP)
	for i, session := range p.sessions {
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

func createNewSession(rawPacket gopacket.Packet) TCPSession {
	tcpLayer := rawPacket.Layer(layers.LayerTypeTCP)
	tcp, _ := tcpLayer.(*layers.TCP)
	net := rawPacket.NetworkLayer()
	src, dst := net.NetworkFlow().Endpoints()
	newSession := TCPSession{
		ServerAddr:     dst.String(),
		ClientAddr:     src.String(),
		ServerPort:     uint16(tcp.DstPort),
		ClientPort:     uint16(tcp.SrcPort),
		SequenceNumber: tcp.Seq >> 8,
		LastUpdate:     time.Now(),
		Packets: []Packet{
			{
				Owner: Client,
				Data:  base64.URLEncoding.EncodeToString(rawPacket.Data())},
		},
	}
	return newSession
}

func (p *Parser) saveSession(i int) {
	p.markSession(i)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	p.DBClient.Connect(ctx)
	collection := p.DBClient.Database("streams").Collection("tcpStreams")
	_, err := collection.InsertOne(ctx, bson.M{"port": p.sessions[i].ServerPort, "session": p.sessions[i]})
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Parser) makePacket(i int, tcpPacket gopacket.Packet) Packet {
	src, _ := tcpPacket.NetworkLayer().NetworkFlow().Endpoints()
	packet := Packet{
		Owner: Client,
		Data:  base64.URLEncoding.EncodeToString(tcpPacket.Data()),
	}
	switch src.String() {
	case p.sessions[i].ClientAddr:
		packet.Owner = Client
	case p.sessions[i].ServerAddr:
		packet.Owner = Server
	}
	return packet
}

func (p *Parser) addPacketToSession(i int, newPacket Packet) {
	p.Lock()
	p.sessions[i].Packets = append(p.sessions[i].Packets, newPacket)
	p.sessions[i].LastUpdate = time.Now()
	p.Unlock()
}
