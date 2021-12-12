package main

import (
	"sync"

	aggregator "github.com/tothbence9922/kawe/internal/aggregator"
	simpleService "github.com/tothbence9922/kawe/internal/ping/impl/simple"
	prometheusServer "github.com/tothbence9922/kawe/internal/server/impl/prometheus"
	httpServer "github.com/tothbence9922/kawe/internal/server/impl/simple"
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

	httpServer := httpServer.HttpServer{Port: 8080}
	httpServer.Serve(wgPtr)

	prometheusServer := prometheusServer.PrometheusServer{Port: 8080}
	prometheusServer.Serve(wgPtr)

	wg.Wait()
}
