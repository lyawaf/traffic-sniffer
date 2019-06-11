package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/lyawaf/traffic-sniffer/pkg/parser"
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
	http.HandleFunc("/getLabels", s.GetLabels)
	http.HandleFunc("/addLabel", s.AddLabel)
	fmt.Println("[SERVICE] Start.")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func (s *Service) GetSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("failed to get body")
		}
		var timeStamp struct {
			LastUpdate int64 `json:"lastUpdate"`
		}
		err = json.Unmarshal(body, &timeStamp)
		if err != nil {
			fmt.Println("failed to parse body")
		}

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		s.dbClient.Connect(ctx)
		collection := s.dbClient.Database("streams").Collection("tcpStreams")
		cur, err := collection.Find(ctx,
			bson.M{"last_update": bson.M{
				"$gt": timeStamp.LastUpdate,
			},
			})
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
		writeAnswer(w, bytes)
	}
}

func (s *Service) GetLabels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		bytes, err := json.Marshal(parser.Labels)
		if err != nil {
			log.Fatal("marshal failed")
		}
		w.Header().Set("Content-Type", "application/json")
		writeAnswer(w, bytes)
	case "POST":
	}
}

func (s *Service) AddLabel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("failed to get body")
		}
		var rawLabel RawLabel
		err = json.Unmarshal(body, &rawLabel)
		if err != nil {
			fmt.Println("failed to parse body")
		}
		if ok := validateLabel(w, rawLabel); !ok {
			return
		}
		decodedRegexp, err := base64.URLEncoding.DecodeString(rawLabel.Regexp)
		newRegexp := regexp.MustCompile(string(decodedRegexp))
		newLabel := parser.Label{
			Name:      rawLabel.Name,
			Color:     rawLabel.Color,
			Type:      parser.LabelType(rawLabel.Type),
			Regexp:    newRegexp,
			RawRegexp: rawLabel.Regexp,
		}
		parser.Labels.Lock()
		parser.Labels.L = append(parser.Labels.L, newLabel)
		parser.Labels.Unlock()
		fmt.Println("[SERVICE] Add new label", newLabel)

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		s.dbClient.Connect(ctx)
		collection := s.dbClient.Database("streams").Collection("labels")
		_, err = collection.InsertOne(ctx, newLabel)
		if err != nil {
			writeAnswer(w, []byte("ERROR: Failed to insert to database"))
			return
		}

		go parser.UpdateLabels(newLabel)

		writeAnswer(w, []byte("SUCCESS"))
	}
}

func validateLabel(w http.ResponseWriter, label RawLabel) bool {
	decodedRegexp, err := base64.URLEncoding.DecodeString(label.Regexp)
	fmt.Println("Decoded regexp", string(decodedRegexp))
	if err != nil {
		writeAnswer(w, []byte("ERROR: Failed to decode regexp base64"))
		return false
	}
	_, err = regexp.Compile(string(decodedRegexp))
	if err != nil {
		writeAnswer(w, []byte("ERROR: Failed to compile regexp"))
		return false
	}
	switch label.Type {
	case parser.PacketIN, parser.PacketOUT:
	default:
		writeAnswer(w, []byte("ERROR: Unknown label type"))
		return false
	}
	return true
}

func writeAnswer(w http.ResponseWriter, data []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err := w.Write(data)
	if err != nil {
		fmt.Println("Failed to send to client", err)
	}
}
