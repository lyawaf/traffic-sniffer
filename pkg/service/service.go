package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	port     int
	dbClient *mongo.Client
}

func Start() {
	var s Service
	dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to create db client", err)
		return
	}
	s.dbClient = dbClient
	http.HandleFunc("/", s.GetSessions)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func (s *Service) GetSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "POST":
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		s.dbClient.Connect(ctx)
		collection := s.dbClient.Database("streams").Collection("tcpStreams")
		cur, err := collection.Find(ctx, bson.M{"port": 9007})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)
		var sessions []bson.M
		for cur.Next(ctx) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			sessions = append(sessions, result)
		}
		bytes, err := json.Marshal(sessions)
		if err != nil {
			log.Fatal("marshal failed")
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Println("Failed to send to client", err)
		}
	}
}
