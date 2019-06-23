package parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

func UpdateLabels(label Label) {

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	DBClientForUpdater.Connect(ctx)
	collection := DBClientForUpdater.Database("streams").Collection("tcpStreams")

	cur, err := collection.Find(ctx,
		bson.M{"last_update": bson.M{
			"$lt": time.Now().Unix(),
		},
		})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		result := bson.M{
			"_id":     0,
			"session": TCPSession{},
		}
		id := result["id"]
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		var tcpSession TCPSession
		bsonBytes, _ := bson.Marshal(result["session"])
		err = bson.Unmarshal(bsonBytes, &tcpSession)
		for _, p := range tcpSession.Packets {
			data, _ := base64.StdEncoding.DecodeString(p.Data)
			matched := label.Regexp.Match(data)
			if matched {
				fmt.Println("matched")
				tcpSession.Labels = append(tcpSession.Labels, label)
				_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"session": tcpSession}})
				continue
			}
		}
	}
}
