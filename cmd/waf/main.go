package main

import (
	"os"

	"github.com/lyawaf/traffic-sniffer/pkg/worker"
)

func main() {
	worker.StartWatcher(os.Args[1])
}
