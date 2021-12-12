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
	wgPtr := &wg

	// Aggregator is started before the pinging service
	aggregator.Start(wgPtr)

	// The pinging service starts based on the configuration file
	simpleService.Start(wgPtr)

	server.Start(wgPtr)

	wg.Wait()
}
