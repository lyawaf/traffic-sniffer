package worker

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/lyawaf/traffic-sniffer/pkg/parser"
	"github.com/lyawaf/traffic-sniffer/pkg/tshark"
	"go.mongodb.org/mongo-driver/mongo"
)

type Worker struct {
	dirPath  string
	nextPath string
	DBClient *mongo.Client
	Watcher  *fsnotify.Watcher
}

func StartWatcher(dirPath string, dbClient *mongo.Client) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = watcher.Close()
	}()

	w := Worker{
		dirPath:  dirPath,
		nextPath: "",
		Watcher:  watcher,
		DBClient: dbClient,
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
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Create {
				go func() {
					files, tempDir := tshark.SeparateSessions(w.nextPath)
					w.nextPath = w.dirPath + event.Name
					for _, s := range files {
						parser.LoadToDB(s, w.DBClient)
					}
					os.RemoveAll(tempDir)
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
