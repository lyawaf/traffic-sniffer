package main

// Use tcpdump to create a test file
// tcpdump -w test.pcap
// or use the example above for writing pcap files

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/lyawaf/traffic-sniffer/pkg/parser"
	"github.com/lyawaf/traffic-sniffer/pkg/service"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"github.com/google/gopacket/pcap"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	device       string = "lo"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
)

var DBClient *mongo.Client

func main() {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	var filter string = "tcp and port 5000"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing TCP port 5000 packets.")

	DBClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	parser.DBClientForUpdater, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to create db client", err)
		return
	}

	newParser := parser.Parser{
		Source: gopacket.NewPacketSource(handle, handle.LinkType()),
	}
	parser.DBClient = DBClient
	go newParser.Parse()
	fmt.Println("[MAIN] Parser start")
	service.Start()
}
