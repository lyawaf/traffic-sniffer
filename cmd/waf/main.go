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
	result := parser.Parse(packetSource, 0)

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
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	collection := dbClient.Database("streams").Collection("tcpStreams")
	cur, err := collection.Find(ctx, bson.M{"port": 9007})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
}
