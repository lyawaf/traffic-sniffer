package main

import (
	"fmt"
	"os"

	"github.com/lyawaf/traffic-sniffer/pkg/worker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to create db client", err)
		return
	}

	worker.StartWatcher(os.Args[1], dbClient)
}
