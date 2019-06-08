package main

// Use tcpdump to create a test file
// tcpdump -w test.pcap
// or use the example above for writing pcap files

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/lyawaf/traffic-sniffer/pkg/parser"
	"github.com/lyawaf/traffic-sniffer/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	pcapFile = "test.pcap"
	handle   *pcap.Handle
	err      error
)

func main() {
	// Open file instead of device
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	result := parser.Parse(packetSource)

	dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to create db client", err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbClient.Connect(ctx)
	for _, session := range result {
		collection := dbClient.Database("streams").Collection("tcpStreams")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := collection.InsertOne(ctx, bson.M{"port": session.ServerPort, "session": session})
		fmt.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	service.Start()
}
