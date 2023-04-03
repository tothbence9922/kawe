package main

import (
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	tcpService "github.com/tothbence9922/kawe/internal/ping/tcp"
	server "github.com/tothbence9922/kawe/internal/server"
)

var (
	wg sync.WaitGroup
)

func start(wg *sync.WaitGroup) {
	// Aggregator is started before the pinging service
	aggregator.Start(wg)

	// The pinging service starts based on the configuration file
	tcpService.Start(wg)

	// Servers start based on .env variables
	server.Start(wg)

	wg.Wait()
}

func main() {

	start(&wg)

}
