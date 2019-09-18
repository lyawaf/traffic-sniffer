package parser

import (
	"context"
	"encoding/base64"
	"fmt"

	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoadToDB(sessionPCAP string, dbClient *mongo.Client, ctx context.Context) {
	handle, err := pcap.OpenOffline(sessionPCAP)
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	session := CreateSession(packetSource)

	collection := dbClient.Database("streams").Collection("tcpStreams")
	_, err = collection.InsertOne(ctx, bson.M{
		"port":        session.ServerPort,
		"session":     session,
		"last_update": time.Now().Unix(),
	})
	if err != nil {
		log.Fatal("failed to insert session", err)
	}
    fmt.Println("Save new session from", sessionPCAP)
}

func CreateSession(source *gopacket.PacketSource) TCPSession {
	var resultS TCPSession
	for rawPacket := range source.Packets() {
		if len(resultS.Packets) == 0 {
			initSession(&resultS, rawPacket)
			continue
		}
		insertPacket(&resultS, rawPacket)
	}
	return resultS
}

func initSession(session *TCPSession, rawPacket gopacket.Packet) {
	tcpLayer := rawPacket.Layer(layers.LayerTypeTCP)
	tcp, _ := tcpLayer.(*layers.TCP)
	net := rawPacket.NetworkLayer()
	src, dst := net.NetworkFlow().Endpoints()

	session.ServerAddr = dst.String()
	session.ClientAddr = src.String()
	session.ServerPort = uint16(tcp.DstPort)
	session.ClientPort = uint16(tcp.SrcPort)
	session.SequenceNumber = tcp.Seq >> 8
	session.Packets = []Packet{
		{
			Owner: Client,
			Data:  base64.StdEncoding.EncodeToString(rawPacket.Data())},
	}
}

func insertPacket(session *TCPSession, rawPacket gopacket.Packet) {
	src, _ := rawPacket.NetworkLayer().NetworkFlow().Endpoints()
	packet := Packet{
		Data: base64.StdEncoding.EncodeToString(rawPacket.Data()),
	}
	switch src.String() {
	case session.ClientAddr:
		packet.Owner = Client
	case session.ServerAddr:
		packet.Owner = Server
	}
	session.Packets = append(session.Packets, packet)
}
