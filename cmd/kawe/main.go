package main

import (
	"fmt"
	"log"
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	simpleService "github.com/tothbence9922/kawe/internal/ping/impl/simple"
	server "github.com/tothbence9922/kawe/internal/server/impl"

	"github.com/fsnotify/fsnotify"
)

var (
	wg sync.WaitGroup
)

func start(wg *sync.WaitGroup) {
	// Aggregator is started before the pinging service
	aggregator.Start(wg)

	// The pinging service starts based on the configuration file
	simpleService.Start(wg)

	server.Start(wg)

	wg.Wait()
}

func main() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("fasz")
					return
				}
				fmt.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("fasz2")

					return
				}
				fmt.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("/config.json")

	if err != nil {
		log.Fatal(err)
	}

	<-done

}
