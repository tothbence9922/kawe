package main

import (
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	simpleService "github.com/tothbence9922/kawe/internal/ping/impl/simple"
	server "github.com/tothbence9922/kawe/internal/server/impl"
)

var (
	wg sync.WaitGroup
)

func main() {

	// Aggregator is started before the pinging service
	aggregator.Start(&wg)

	// The pinging service starts based on the configuration file
	simpleService.Start(&wg)

	server.Start(&wg)

	wg.Wait()
}
