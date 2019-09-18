package worker

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/lyawaf/traffic-sniffer/pkg/parser"
	"github.com/lyawaf/traffic-sniffer/pkg/tshark"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Worker struct {
	dirPath  string
	nextPath string
	Watcher  *fsnotify.Watcher
	DbClient *mongo.Client
	Ctx      context.Context
}

func StartWatcher(dirPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = watcher.Close()
	}()

	dbClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("Failed to create db client", err)
		return
	}
	ctx, _ := context.WithCancel(context.Background())
	err = dbClient.Connect(ctx)

	w := Worker{
		dirPath:  dirPath,
		nextPath: "",
		Watcher:  watcher,
		DbClient: dbClient,
		Ctx:      ctx,
	}

	done := make(chan bool)
	go w.handlePath()

	fmt.Println("Adding path", w.Watcher.Add(dirPath))
	<-done
}

func (w *Worker) handlePath() {
	for {
		select {
		case event, ok := <-w.Watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				go func() {
					if w.nextPath != "" {
						files, tempDir := tshark.SeparateSessions(w.nextPath)
						for _, s := range files {
							if s == "" {
								continue
							}
							fmt.Println("trying to load to db with", s)
							parser.LoadToDB(s, w.DbClient, w.Ctx)
						}
						os.RemoveAll(tempDir)
					}
					w.nextPath = event.Name
				}()
			}
		case err, ok := <-w.Watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
